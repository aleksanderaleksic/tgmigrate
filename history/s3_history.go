package history

import (
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/seqsense/s3sync"
	"os"
)

type S3History struct {
	S3StorageConfig config.S3HistoryStorageConfig
	session         session.Session
	StorageHistory  *StorageHistory
}

func (h S3History) IsMigrationApplied(hash string) (*Result, error) {
	for _, m := range h.StorageHistory.AppliedMigration {
		if m.Hash == hash {
			return &m.Result, nil
		}
	}
	return &Result{State: ResultStateUnapplied}, nil
}

func (h *S3History) InitializeHistory(skipUserInteraction bool) (*StorageHistory, error) {
	historyPath := h.S3StorageConfig.GetLocalHistoryPath() + "/" + h.S3StorageConfig.Key

	err := s3sync.New(&h.session).Sync("s3://"+h.S3StorageConfig.Bucket+"/"+h.S3StorageConfig.Key, historyPath)
	if err != nil {
		return nil, err
	}

	storageHistory, err := getOrCreateNewHistoryFile(historyPath, skipUserInteraction)
	if err != nil {
		return nil, err
	}

	h.StorageHistory = storageHistory

	return h.StorageHistory, nil
}

func (h *S3History) StoreMigrationObject(migrationName string, result Result, fileHash string) {
	storeMigrationObject(h.StorageHistory, migrationName, result, fileHash)
}

func (h *S3History) WriteToStorage() error {
	historyPath := h.S3StorageConfig.GetLocalHistoryPath() + "/" + h.S3StorageConfig.Key

	err := writeToStorage(historyPath, *h.StorageHistory)
	if err != nil {
		return err
	}

	err = s3sync.New(&h.session, s3sync.WithDryRun()).Sync(historyPath, "s3://"+h.S3StorageConfig.Bucket+"/"+h.S3StorageConfig.Key)
	if err != nil {
		return err
	}

	os.RemoveAll(h.S3StorageConfig.GetLocalHistoryPath())

	return nil
}
