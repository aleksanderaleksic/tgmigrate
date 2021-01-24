package history

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type History interface {
	IsMigrationApplied(hash string) (*Result, error)
	InitializeHistory(ctx common.Context) (*StorageHistory, error)
	StoreMigrationObject(migrationName string, result Result, fileHash string)
	WriteToStorage() error
}

func GetHistoryInterface(c config.Config, ctx common.Context) (History, error) {
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
		return &S3History{
			context:         ctx,
			S3StorageConfig: *conf,
			session:         *sess,
		}, nil
	default:
		return nil, fmt.Errorf("unknown history storage type: %s", c.History.Storage.Type)
	}
}

func writeStorageHistory(path string, storageHistory StorageHistory) error {
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

func storeMigrationObject(storageHistory *StorageHistory, migrationName string, result Result, fileHash string) {
	storageHistory.AppliedMigration = append(storageHistory.AppliedMigration, StorageHistoryObject{
		SchemaVersion: StorageHistoryObjectVersion,
		Applied:       common.JSONTime(time.Now()),
		Hash:          fileHash,
		Name:          migrationName,
		Result:        result,
	})
}

func writeToStorage(historyPath string, storageHistory StorageHistory) error {
	return writeStorageHistory(historyPath, storageHistory)
}
