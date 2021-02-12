package migration

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"sort"
)

type Reverter struct {
	Context          *common.Context
	Config           *config.Config
	S3StateClient    *s3.S3
	HistoryInterface history.History
}

func (r Reverter) RevertToIncluding(migrationName string) error {
	hist, err := r.HistoryInterface.InitializeHistory()
	if err != nil {
		return err
	}

	migration := getMigrationWithName(hist.AppliedMigration, migrationName)
	if migration == nil {
		return fmt.Errorf("could not find a applied migration with the name: %s", migrationName)
	}

	migrationsToRevert := migrationsToRevert(hist.AppliedMigration, *migration)

	defer r.HistoryInterface.WriteToStorage()

	for _, migration := range migrationsToRevert {
		err := r.revert(migration)
		if err != nil {
			return err
		}
		r.HistoryInterface.RemoveAppliedMigration(migration.Name)
	}

	return nil
}

func getMigrationWithName(slice []history.AppliedStorageHistoryObject, migrationName string) *history.AppliedStorageHistoryObject {
	for _, migration := range slice {
		if migration.Name == migrationName {
			return &migration
		}
	}
	return nil
}

func migrationsToRevert(all []history.AppliedStorageHistoryObject, inclusive history.AppliedStorageHistoryObject) []history.AppliedStorageHistoryObject {
	migrationsToRevert := []history.AppliedStorageHistoryObject{inclusive}

	for _, migration := range all {
		if migration.Applied.Time().After(inclusive.Applied.Time()) {
			migrationsToRevert = append(migrationsToRevert, migration)
		}
	}

	sort.Slice(migrationsToRevert, func(i, j int) bool {
		return migrationsToRevert[i].Applied.Time().After(migrationsToRevert[j].Applied.Time())
	})

	return migrationsToRevert
}

func (r Reverter) revert(migration history.AppliedStorageHistoryObject) error {
	s3Config := r.Config.State.Config.(*config.S3StateConfig)

	fmt.Printf("The following object version will be deleted:\n")
	objectIdentifiersToDelete := make([]*s3.ObjectIdentifier, 0)
	for _, obj := range migration.Metadata.S3Metadata.ChangedObjects {
		objectIdentifiersToDelete = append(objectIdentifiersToDelete, &s3.ObjectIdentifier{
			Key:       aws.String(obj.Key),
			VersionId: aws.String(obj.ToVersionId),
		})
		fmt.Printf("\t- key: %s, version: %s\n", obj.Key, obj.ToVersionId)
	}
	fmt.Printf("Are you sure you would like to continue?\n")
	if common.AskUserToConfirm() {
		fmt.Printf("Cancelled by user, will not delete objects.\n")
		os.Exit(0)
	}

	_, err := r.S3StateClient.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(s3Config.Bucket),
		Delete: &s3.Delete{
			Objects: objectIdentifiersToDelete,
			Quiet:   aws.Bool(false),
		},
	})
	if err != nil {
		return err
	}
	fmt.Printf("Successfully deleted object versions\n")

	return nil
}
