package history

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/config"
)

type History interface {
	IsMigrationApplied(hash string) (*Result, error)
	InitializeStorage(skipUserInteraction bool) (*StorageHistory, error)
	StoreMigrationObject(migrationName string, result Result, fileHash string)
	WriteToStorage() error
}

func GetHistoryInterface(c config.Config) (History, error) {
	switch c.History.Storage.Type {
	case "s3":
		conf := c.History.Storage.Config.(*config.S3HistoryStorageConfig)
		return &S3History{S3StorageConfig: *conf}, nil
	case "local":
		conf := c.History.Storage.Config.(*config.LocalHistoryStorageConfig)
		return &LocalHistory{LocalStorageConfig: *conf}, nil
	default:
		return nil, fmt.Errorf("unknown history storage type: %s", c.History.Storage.Type)
	}
}
