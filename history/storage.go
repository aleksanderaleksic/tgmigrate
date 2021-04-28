package history

import (
	"encoding/json"
	"github.com/aleksanderaleksic/tgmigrate/common"
)

const StorageHistoryVersion = "v1"
const StorageHistoryObjectVersion = "v1"

type StorageHistory struct {
	SchemaVersion     string                        `json:"schema_version"`
	AppliedMigration  []AppliedStorageHistoryObject `json:"applied_migration"`
	SkippedMigrations []SkippedStorageHistoryObject `json:"skipped_migration"`
	FailedMigrations  []FailedStorageHistoryObject  `json:"failed_migration"`
}

type AppliedStorageHistoryObject struct {
	SchemaVersion string          `json:"schema_version"`
	Applied       common.JSONTime `json:"applied"`
	Hash          string          `json:"hash"`
	Name          string          `json:"name"`
	Metadata      MetadataWrapper `json:"metadata"`
}

type FailedStorageHistoryObject struct {
	SchemaVersion string          `json:"schema_version"`
	Failed        common.JSONTime `json:"applied"`
	Hash          string          `json:"hash"`
	Name          string          `json:"name"`
}

type SkippedStorageHistoryObject struct {
	SchemaVersion string          `json:"schema_version"`
	Skipped       common.JSONTime `json:"skipped"`
	Name          string          `json:"name"`
}

func EmptyStorageHistory() StorageHistory {
	return StorageHistory{
		SchemaVersion:     StorageHistoryVersion,
		AppliedMigration:  []AppliedStorageHistoryObject{},
		SkippedMigrations: []SkippedStorageHistoryObject{},
		FailedMigrations:  []FailedStorageHistoryObject{},
	}
}

func DecodeStorageHistory(source []byte) (*StorageHistory, error) {
	var obj StorageHistory
	err := json.Unmarshal(source, &obj)
	return &obj, err
}

func EncodeStorageHistory(history StorageHistory) ([]byte, error) {
	return json.Marshal(history)
}
