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
	s3Path := "s3://" + filepath.Join(h.S3StorageConfig.Bucket, filepath.Dir(h.S3StorageConfig.Key))
	err := h.syncManager.Sync(s3Path, filepath.Dir(historyPath))
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

func (h S3History) StoreAppliedMigration(migration *AppliedStorageHistoryObject) {
	h.StorageHistory.storeAppliedMigration(migration)
}
func (h S3History) StoreFailedMigration(migration *FailedStorageHistoryObject) {
	h.StorageHistory.storeFailedMigration(migration)
}

func (h S3History) RemoveAppliedMigration(migrationName string) {
	for index, migration := range h.StorageHistory.AppliedMigration {
		if migration.Name == migrationName {
			h.StorageHistory.AppliedMigration = append(h.StorageHistory.AppliedMigration[:index], h.StorageHistory.AppliedMigration[index+1:]...)
		}
	}
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
