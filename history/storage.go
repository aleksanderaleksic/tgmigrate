package history

import (
	"encoding/json"
	"fmt"
	"time"
)

const StorageHistoryVersion = "v1"
const StorageHistoryObjectVersion = "v1"

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-02T15:04:05"))
	return []byte(stamp), nil
}

type StorageHistory struct {
	SchemaVersion    string                 `json:"schema_version"`
	AppliedMigration []StorageHistoryObject `json:"applied_migration"`
}

type StorageHistoryObject struct {
	SchemaVersion string   `json:"schema_version"`
	Applied       JSONTime `json:"applied"`
	Hash          string   `json:"hash"`
	Name          string   `json:"name"`
	Result        Result   `json:"result"`
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
