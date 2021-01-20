package history

import (
	"encoding/json"
)

const StorageHistoryVersion = "v1"
const StorageHistoryObjectVersion = "v1"

type StorageHistory struct {
	SchemaVersion    string                 `json:"schema_version"`
	AppliedMigration []StorageHistoryObject `json:"applied_migration"`
}

type StorageHistoryObject struct {
	SchemaVersion string `json:"schema_version"`
	Hash          string `json:"hash"`
	Name          string `json:"name"`
	Result        Result `json:"result"`
}

func EmptyStorageHistory() StorageHistory {
	return StorageHistory{
		SchemaVersion:    StorageHistoryVersion,
		AppliedMigration: []StorageHistoryObject{},
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
