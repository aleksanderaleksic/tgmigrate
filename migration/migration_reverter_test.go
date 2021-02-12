package migration

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReverter_getMigrationWithName_SingleExpectedMigration(t *testing.T) {

	expected := history.AppliedStorageHistoryObject{
		SchemaVersion: "v1",
		Applied:       common.JSONTime(time.Now()),
		Hash:          "sample_hash",
		Name:          "V1__expected.hcl",
		Metadata: history.MetadataWrapper{
			S3Metadata: nil,
		},
	}

	slice := []history.AppliedStorageHistoryObject{
		expected,
	}

	actual := getMigrationWithName(slice, expected.Name)
	assert.Equal(t, expected, *actual)
}

func TestReverter_getMigrationWithName_SingleUnexpectedMigration(t *testing.T) {
	slice := []history.AppliedStorageHistoryObject{
		{
			SchemaVersion: "v1",
			Applied:       common.JSONTime(time.Now()),
			Hash:          "sample_hash",
			Name:          "V1__unexpected.hcl",
			Metadata: history.MetadataWrapper{
				S3Metadata: nil,
			},
		},
	}

	actual := getMigrationWithName(slice, "V1__expected.hcl")
	assert.Nil(t, actual)
}

func TestReverter_getMigrationWithName_MultipleMigrationsWithExpected(t *testing.T) {

	expected := history.AppliedStorageHistoryObject{
		SchemaVersion: "v1",
		Applied:       common.JSONTime(time.Now()),
		Hash:          "sample_hash",
		Name:          "V1__expected.hcl",
		Metadata: history.MetadataWrapper{
			S3Metadata: nil,
		},
	}

	slice := []history.AppliedStorageHistoryObject{
		expected,
		{
			SchemaVersion: "v1",
			Applied:       common.JSONTime(time.Now()),
			Hash:          "sample_hash",
			Name:          "V2__unexpected.hcl",
			Metadata: history.MetadataWrapper{
				S3Metadata: nil,
			},
		},
	}

	actual := getMigrationWithName(slice, expected.Name)
	assert.Equal(t, expected, *actual)
}

func TestReverter_getMigrationWithName_MultipleMigrationsWithoutExpected(t *testing.T) {
	slice := []history.AppliedStorageHistoryObject{
		{
			SchemaVersion: "v1",
			Applied:       common.JSONTime(time.Now()),
			Hash:          "sample_hash",
			Name:          "V1__unexpected.hcl",
			Metadata: history.MetadataWrapper{
				S3Metadata: nil,
			},
		},
		{
			SchemaVersion: "v1",
			Applied:       common.JSONTime(time.Now()),
			Hash:          "sample_hash",
			Name:          "V2__unexpected.hcl",
			Metadata: history.MetadataWrapper{
				S3Metadata: nil,
			},
		},
	}

	actual := getMigrationWithName(slice, "V1__expected.hcl")
	assert.Nil(t, actual)
}

func TestReverter_migrationsToRevert_(t *testing.T) {
	baseTime := time.Now()
	inclusive := history.AppliedStorageHistoryObject{
		SchemaVersion: "v1",
		Applied:       common.JSONTime(baseTime.Add(1 * time.Hour)),
		Hash:          "sample_hash",
		Name:          "V2__include_second.hcl",
		Metadata: history.MetadataWrapper{
			S3Metadata: nil,
		},
	}
	slice := []history.AppliedStorageHistoryObject{
		{
			SchemaVersion: "v1",
			Applied:       common.JSONTime(baseTime.Add(2 * time.Hour)),
			Hash:          "sample_hash",
			Name:          "V3__include_first.hcl",
			Metadata: history.MetadataWrapper{
				S3Metadata: nil,
			},
		},
		inclusive,
		{
			SchemaVersion: "v1",
			Applied:       common.JSONTime(baseTime),
			Hash:          "sample_hash",
			Name:          "V1__unexpected.hcl",
			Metadata: history.MetadataWrapper{
				S3Metadata: nil,
			},
		},
	}

	expected := []history.AppliedStorageHistoryObject{
		slice[0],
		slice[1],
	}

	actual := migrationsToRevert(slice, inclusive)

	assert.Equal(t, expected, actual)
}
