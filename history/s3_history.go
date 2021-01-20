package history

import (
	"github.com/aleksanderaleksic/tgmigrate/config"
)

type S3History struct {
	S3StorageConfig config.S3HistoryStorageConfig
}

func (h S3History) IsMigrationApplied(hash string) (*Result, error) {
	r := Result{State: ResultStateUnapplied}
	return &r, nil
}

func (h *S3History) InitializeStorage(skipUserInteraction bool) (*StorageHistory, error) {
	return nil, nil
}

func (h S3History) StoreMigrationObject(migrationName string, result Result, fileHash string) {
	return
}

func (h S3History) WriteToStorage() error {
	return nil
}
