package history

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type LocalHistory struct {
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

func (h *LocalHistory) InitializeStorage(skipUserInteraction bool) (*StorageHistory, error) {
	historyStoragePath := h.getHistoryPath()
	if _, err := os.Stat(historyStoragePath); os.IsNotExist(err) {
		fmt.Printf("Migration history not found at '%s', do you want to create a new history on this location?\n", historyStoragePath)
		if !skipUserInteraction && !common.AskUserToConfirm() {
			os.Exit(0)
			return nil, nil
		}
		emptyHistoryObj := EmptyStorageHistory()
		err = h.writeStorageHistory(emptyHistoryObj)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Created a empty history at '%s'\n", historyStoragePath)

		h.StorageHistory = &emptyHistoryObj
		return &emptyHistoryObj, nil
	} else {
		history, err := h.readStorageHistory()
		if err != nil {
			return nil, err
		}
		h.StorageHistory = history
		return history, nil
	}
}

func (h LocalHistory) getHistoryPath() string {
	historyStoragePath, _ := filepath.Abs(h.LocalStorageConfig.Path)
	return historyStoragePath
}

func (h LocalHistory) writeStorageHistory(storageHistory StorageHistory) error {
	historyStoragePath := h.getHistoryPath()
	encodedStorageHistory, err := EncodeStorageHistory(storageHistory)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(historyStoragePath, encodedStorageHistory, 0777)
	if err != nil {
		return fmt.Errorf("failed to write history to '%s' err: %s", historyStoragePath, err)
	}
	return nil
}

func (h LocalHistory) readStorageHistory() (*StorageHistory, error) {
	historyStoragePath := h.getHistoryPath()
	source, err := ioutil.ReadFile(historyStoragePath)
	if err != nil {
		return nil, err
	}
	history, err := DecodeStorageHistory(source)
	if err != nil {
		return nil, fmt.Errorf("failed to read history from '%s' err: %s", historyStoragePath, err)
	}

	return history, nil
}

func (h *LocalHistory) StoreMigrationObject(migrationName string, result Result, fileHash string) {
	h.StorageHistory.AppliedMigration = append(h.StorageHistory.AppliedMigration, StorageHistoryObject{
		SchemaVersion: StorageHistoryObjectVersion,
		Applied:       common.JSONTime(time.Now()),
		Hash:          fileHash,
		Name:          migrationName,
		Result:        result,
	})
}

func (h LocalHistory) WriteToStorage() error {
	return h.writeStorageHistory(*h.StorageHistory)
}
