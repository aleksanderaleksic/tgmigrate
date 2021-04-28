package migration

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Skipper struct {
	Context          *common.Context
	Config           *config.Config
	S3StateClient    *s3.S3
	HistoryInterface history.History
}

func (r Skipper) Skip(migrationName string) error {
	hist, err := r.HistoryInterface.InitializeHistory()
	if err != nil {
		return err
	}

	migrationFiles, err := GetMigrationFiles(r.Config.AbsoluteMigrationDir)
	if err != nil {
		return fmt.Errorf("could not get migration files from '%s', error: %s", r.Config.AbsoluteMigrationDir, err)
	}
	migrationFile := getMigration(migrationName, migrationFiles)
	if migrationFile == nil {
		return fmt.Errorf("migration file with name: '%s' was not found", migrationName)
	}

	for _, m := range hist.AppliedMigration {
		if m.Name == migrationName {
			return fmt.Errorf("migration file with name: '%s' is already applied", migrationName)
		}
	}
	for _, m := range hist.SkippedMigrations {
		if m.Name == migrationName {
			return fmt.Errorf("migration file with name: '%s' is already skipped", migrationName)
		}
	}

	r.HistoryInterface.StoreSkippedMigration(&history.SkippedStorageHistoryObject{
		SchemaVersion: history.StorageHistoryObjectVersion,
		Skipped:       common.JSONTime{},
		Name:          migrationName,
	})

	r.HistoryInterface.WriteToStorage()

	return nil
}

func getMigration(migrationName string, migrations *[]File) *File {
	for _, migration := range *migrations {
		if migration.Metadata.FileName == migrationName {
			return &migration
		}
	}
	return nil
}
