package history

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/seqsense/s3sync"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type History interface {
	IsMigrationApplied(hash string) (bool, error)
	InitializeHistory() (*StorageHistory, error)
	StoreAppliedMigration(migration *AppliedStorageHistoryObject)
	StoreFailedMigration(migration *FailedStorageHistoryObject)
	RemoveAppliedMigration(migrationName string)
	WriteToStorage() error
	Cleanup()
}

func GetHistoryInterface(c config.Config, ctx common.Context, cache common.Cache) (History, error) {
	switch c.History.Storage.Type {
	case "s3":
		conf := c.History.Storage.Config.(*config.S3HistoryStorageConfig)
		sess, err := session.NewSession(&aws.Config{
			Region: &conf.Region,
		})

		if err != nil {
			return nil, err
		}

		if conf.AssumeRole != nil {
			credentials := stscreds.NewCredentials(sess, *conf.AssumeRole)
			sess, err = session.NewSession(&aws.Config{
				Region:      &conf.Region,
				Credentials: credentials,
			})

			if err != nil {
				return nil, err
			}
		}

		var safeSyncManager *s3sync.Manager
		if ctx.DryRun {
			safeSyncManager = s3sync.New(sess, s3sync.WithDryRun())
		} else {
			safeSyncManager = s3sync.New(sess)
		}

		return &S3History{
			context:         ctx,
			S3StorageConfig: *conf,
			safeSyncManager: *safeSyncManager,
			syncManager:     *s3sync.New(sess),
			Cache:           cache,
		}, nil
	default:
		return nil, fmt.Errorf("unknown history storage type: %s", c.History.Storage.Type)
	}
}

func writeStorageHistory(path string, storageHistory StorageHistory) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	encodedStorageHistory, err := EncodeStorageHistory(storageHistory)
	if err != nil {
		return err
	}
	absPath, _ := filepath.Abs(path)
	err = ioutil.WriteFile(absPath, encodedStorageHistory, 0777)
	if err != nil {
		return fmt.Errorf("failed to write history to '%s' err: %s", path, err)
	}
	return nil
}

func readStorageHistory(path string) (*StorageHistory, error) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	history, err := DecodeStorageHistory(source)
	if err != nil {
		return nil, fmt.Errorf("failed to read history from '%s' err: %s", path, err)
	}

	return history, nil
}

func getOrCreateNewHistoryFile(historyPath string, skipUserInteraction bool) (*StorageHistory, error) {
	if _, err := os.Stat(historyPath); os.IsNotExist(err) {
		fmt.Printf("Migration history not found at '%s', do you want to create a new history?\n", historyPath)
		if !skipUserInteraction && !common.AskUserToConfirm() {
			os.Exit(0)
		}
		emptyHistoryObj := EmptyStorageHistory()
		err = writeStorageHistory(historyPath, emptyHistoryObj)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Created a empty history at '%s'\n", historyPath)

		return &emptyHistoryObj, nil
	} else {
		history, err := readStorageHistory(historyPath)
		if err != nil {
			return nil, err
		}
		return history, nil
	}
}

func (s *StorageHistory) storeAppliedMigration(migration *AppliedStorageHistoryObject) {
	migration.Applied = common.JSONTime(time.Now())
	s.AppliedMigration = append(s.AppliedMigration, *migration)
}
func (s *StorageHistory) storeFailedMigration(migration *FailedStorageHistoryObject) {
	migration.Failed = common.JSONTime(time.Now())
	s.FailedMigrations = append(s.FailedMigrations, *migration)
}

func writeToStorage(historyPath string, storageHistory StorageHistory) error {
	return writeStorageHistory(historyPath, storageHistory)
}
