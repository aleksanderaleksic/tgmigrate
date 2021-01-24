package history

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"path/filepath"
)

type LocalHistory struct {
	context            common.Context
	LocalStorageConfig config.LocalHistoryStorageConfig
	StorageHistory     *StorageHistory
}

func (h LocalHistory) IsMigrationApplied(hash string) (*Result, error) {
	for _, m := range h.StorageHistory.AppliedMigration {
		if m.Hash == hash {
			return &m.Result, nil
		}
	}
	return &Result{State: ResultStateUnapplied}, nil
}

func (h *LocalHistory) InitializeHistory(ctx common.Context) (*StorageHistory, error) {
	historyPath := h.getHistoryPath()
	storageHistory, err := getOrCreateNewHistoryFile(historyPath, ctx.SkipUserInteraction)
	if err != nil {
		return nil, err
	}

	h.StorageHistory = storageHistory
	return h.StorageHistory, nil
}

func (h LocalHistory) getHistoryPath() string {
	historyStoragePath, _ := filepath.Abs(h.LocalStorageConfig.Path)
	return historyStoragePath
}

func (h *LocalHistory) StoreMigrationObject(migrationName string, result Result, fileHash string) {
	storeMigrationObject(h.StorageHistory, migrationName, result, fileHash)
}

func (h LocalHistory) WriteToStorage() error {
	return writeToStorage(h.getHistoryPath(), *h.StorageHistory)
}
