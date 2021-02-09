package history

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/seqsense/s3sync"
	"os"
	"path/filepath"
)

type S3History struct {
	context         common.Context
	S3StorageConfig config.S3HistoryStorageConfig
	safeSyncManager s3sync.Manager
	syncManager     s3sync.Manager
	StorageHistory  *StorageHistory
	Cache           common.Cache
}

func (h S3History) IsMigrationApplied(hash string) (bool, error) {
	for _, m := range h.StorageHistory.AppliedMigration {
		if m.Hash == hash {
			return true, nil
		}
	}
	return false, nil
}

func (h *S3History) InitializeHistory() (*StorageHistory, error) {
	historyPath := h.getHistoryStoragePath()

	err := h.syncManager.Sync("s3://"+h.S3StorageConfig.Bucket+"/"+h.S3StorageConfig.Key, historyPath)
	if err != nil {
		return nil, err
	}

	storageHistory, err := getOrCreateNewHistoryFile(historyPath, h.context.SkipUserInteraction)
	if err != nil {
		return nil, err
	}

	h.StorageHistory = storageHistory

	return h.StorageHistory, nil
}

func (h *S3History) StoreMigrationObject(migrationName string, success bool, fileHash string) {
	storeMigrationObject(h.StorageHistory, migrationName, success, fileHash)
}

func (h *S3History) WriteToStorage() error {
	historyPath := h.getHistoryStoragePath()

	err := writeToStorage(historyPath, *h.StorageHistory)
	if err != nil {
		return err
	}

	err = h.safeSyncManager.Sync(historyPath, "s3://"+h.S3StorageConfig.Bucket+"/"+h.S3StorageConfig.Key)

	if err != nil {
		return err
	}

	return nil
}

func (h S3History) getHistoryStoragePath() string {
	return filepath.Join(h.Cache.GetCacheDirectoryPath(), "history", "history.json")
}

func (h S3History) Cleanup() {
	os.RemoveAll(h.getHistoryStoragePath())
}
