package history

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestDecodeEmptyStorageHistory(t *testing.T) {
	source := []byte(`
{
	"schema_version": "v1",
	"applied_migration": []
}
`)

	storageHistory, err := DecodeStorageHistory(source)
	assert.Nil(t, err)
	assert.Equal(t, "v1", storageHistory.SchemaVersion)
	assert.Equal(t, []StorageHistoryObject{}, storageHistory.AppliedMigration)
}

func TestDecodeStorageHistoryWithHistoryObject(t *testing.T) {
	source := []byte(`
{
	"schema_version": "v1",
	"applied_migration": [
		{
			"schema_version": "v1",
			"applied": "2021-01-02T15:04:05Z",
			"hash": "sample_hash",
			"name": "V1__move.hcl",
			"result": {
				"state": "success"
			}
		}
	]
}
`)

	storageHistory, err := DecodeStorageHistory(source)
	assert.Nil(t, err)
	assert.Equal(t, "v1", storageHistory.SchemaVersion)
	assert.Equal(t, []StorageHistoryObject{
		{
			SchemaVersion: "v1",
			Applied:       common.JSONTime(time.Date(2021, 1, 2, 15, 4, 5, 0, time.UTC)),
			Hash:          "sample_hash",
			Name:          "V1__move.hcl",
			Result:        SuccessResult,
		},
	}, storageHistory.AppliedMigration)
}

func TestEncodeEmptyStorageHistory(t *testing.T) {
	obj := EmptyStorageHistory()
	expected := clearWhitespace(`
{
	"schema_version": "v1",
	"applied_migration": []
}`)

	output, err := EncodeStorageHistory(obj)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(output))
}

func TestEncodeStorageHistoryWithHistoryObject(t *testing.T) {
	obj := StorageHistory{
		SchemaVersion: "v1",
		AppliedMigration: []StorageHistoryObject{
			{
				SchemaVersion: "v1",
				Applied:       common.JSONTime(time.Date(2021, 1, 2, 15, 4, 5, 0, time.UTC)),
				Hash:          "sample_hash",
				Name:          "V1__move.hcl",
				Result:        SuccessResult,
			},
		},
	}
	expected := clearWhitespace(`
{
"schema_version": "v1",
	"applied_migration": [
		{
			"schema_version": "v1",
			"applied": "2021-01-02T15:04:05Z",
			"hash": "sample_hash",
			"name": "V1__move.hcl",
			"result": {
				"state": "success"
			}
		}
	]
}
`)

	output, err := EncodeStorageHistory(obj)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(output))
}

func clearWhitespace(source string) string {
	var re = regexp.MustCompile(`[\n\t ]`)
	output := re.ReplaceAllString(source, "")
	return output
}
