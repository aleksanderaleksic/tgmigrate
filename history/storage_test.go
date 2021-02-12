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
	"applied_migration": [],
	"failed_migration": []
}
`)

	storageHistory, err := DecodeStorageHistory(source)
	assert.Nil(t, err)
	assert.Equal(t, "v1", storageHistory.SchemaVersion)
	assert.Equal(t, []AppliedStorageHistoryObject{}, storageHistory.AppliedMigration)
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
			"metadata": {
				"schema_version": "v1",
				"type": "s3",
				"changed_objects": [
					{
						"key": "file1/terraform.tfstate",
						"from_version_id": null,
						"to_version_id": "d37ff0e71c144cabf5449ec442580c4a"
					}
				]
			}
		}
	],
	"failed_migration": []
}
`)

	storageHistory, err := DecodeStorageHistory(source)
	assert.Nil(t, err)
	assert.Equal(t, "v1", storageHistory.SchemaVersion)
	assert.Equal(t, []AppliedStorageHistoryObject{
		{
			SchemaVersion: "v1",
			Applied:       common.JSONTime(time.Date(2021, 1, 2, 15, 4, 5, 0, time.UTC)),
			Hash:          "sample_hash",
			Name:          "V1__move.hcl",
			Metadata: *StorageS3Metadata{
				SchemaVersion: "v1",
				Type:          "s3",
				ChangedObjects: []ChangedS3Object{
					{
						Key:           "file1/terraform.tfstate",
						FromVersionId: nil,
						ToVersionId:   "d37ff0e71c144cabf5449ec442580c4a",
					},
				},
			}.Wrap(),
		},
	}, storageHistory.AppliedMigration)
}

func TestEncodeEmptyStorageHistory(t *testing.T) {
	obj := EmptyStorageHistory()
	expected := clearWhitespace(`
{
	"schema_version": "v1",
	"applied_migration": [],
	"failed_migration": []
}`)

	output, err := EncodeStorageHistory(obj)
	assert.Nil(t, err)
	assert.Equal(t, expected, string(output))
}

func TestEncodeStorageHistoryWithHistoryObject(t *testing.T) {
	metadata := StorageS3Metadata{
		SchemaVersion: "v1",
		Type:          "s3",
		ChangedObjects: []ChangedS3Object{
			{
				Key:           "file1/terraform.tfstate",
				FromVersionId: nil,
				ToVersionId:   "d37ff0e71c144cabf5449ec442580c4a",
			},
		},
	}.Wrap()

	obj := StorageHistory{
		SchemaVersion: "v1",
		AppliedMigration: []AppliedStorageHistoryObject{
			{
				SchemaVersion: "v1",
				Applied:       common.JSONTime(time.Date(2021, 1, 2, 15, 4, 5, 0, time.UTC)),
				Hash:          "sample_hash",
				Name:          "V1__move.hcl",
				Metadata:      *metadata,
			},
		},
		FailedMigrations: []FailedStorageHistoryObject{},
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
			"metadata": {
				"schema_version": "v1",
				"type": "s3",
				"changed_objects": [
					{
						"key": "file1/terraform.tfstate",
						"from_version_id": null,
						"to_version_id": "d37ff0e71c144cabf5449ec442580c4a"
					}
				]
			}
		}
	],
	"failed_migration": []
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
