package migration

import (
	"crypto/sha256"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	history "github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/state"
	. "github.com/aleksanderaleksic/tgmigrate/test"
	. "github.com/aleksanderaleksic/tgmigrate/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func getMocks(t *testing.T) (*gomock.Controller, *MockHistory, *MockState) {
	ctrl := gomock.NewController(t)

	mHistory := NewMockHistory(ctrl)
	mState := NewMockState(ctrl)
	return ctrl, mHistory, mState
}

func migrationFileHash(path string) string {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to get file hash for: %s", path)
		return ""
	}
	return fmt.Sprintf("%x", sha256.Sum256(source))
}

func TestShouldHandleInitializeHistoryError(t *testing.T) {
	_ = assert.New(t)
	testDir := t.TempDir()
	CopyTestData(t, "simple", testDir)
	ChangeWorkingDirectory(t, testDir)

	ctrl, mHistory, mState := getMocks(t)
	defer ctrl.Finish()

	ctx := GetContext(filepath.Join(testDir, config.DefaultConfigFile), &map[string]string{
		"ACCOUNT_ID": "12345678",
	})
	conf, _ := config.GetConfigFile(ctx)
	run := Runner{
		Context: &common.Context{
			SkipUserInteraction: false,
			DryRun:              false,
		},
		Config:           conf,
		HistoryInterface: mHistory,
		StateInterface:   mState,
	}

	mHistory.EXPECT().
		Cleanup().
		Times(1)
	mState.EXPECT().
		Cleanup().
		Times(1)
	mHistory.
		EXPECT().
		InitializeHistory().
		Times(1).
		Return(nil, fmt.Errorf("failed to get history from s3"))

	err := run.Apply(nil)
	assert.EqualError(t, err, "unable to initialize history, error: failed to get history from s3")
}

func TestShouldHandleNoMigrationFiles(t *testing.T) {
	_ = assert.New(t)
	testDir := t.TempDir()
	CopyTestData(t, "config_no_migration_files", testDir)
	ChangeWorkingDirectory(t, testDir)

	ctrl, mHistory, mState := getMocks(t)
	defer ctrl.Finish()

	ctx := GetContext(filepath.Join(testDir, config.DefaultConfigFile), &map[string]string{
		"ACCOUNT_ID": "12345678",
	})
	conf, _ := config.GetConfigFile(ctx)

	run := Runner{
		Context: &common.Context{
			SkipUserInteraction: false,
			DryRun:              false,
		},
		Config:           conf,
		HistoryInterface: mHistory,
		StateInterface:   mState,
	}

	mHistory.EXPECT().
		Cleanup().
		Times(1)
	mState.EXPECT().
		Cleanup().
		Times(1)
	mHistory.
		EXPECT().
		InitializeHistory().
		Times(1).
		Return(&history.StorageHistory{
			SchemaVersion:    "v1",
			AppliedMigration: nil,
		}, nil)

	err := run.Apply(nil)
	assert.Nil(t, err)
}

func TestShouldFailToApplyWithMoveCommandFailing(t *testing.T) {
	_ = assert.New(t)
	testDir := t.TempDir()
	CopyTestData(t, "simple", testDir)
	ChangeWorkingDirectory(t, testDir)

	ctrl, mHistory, mState := getMocks(t)
	defer ctrl.Finish()

	ctx := GetContext(filepath.Join(testDir, config.DefaultConfigFile), &map[string]string{
		"ACCOUNT_ID": "12345678",
	})
	conf, _ := config.GetConfigFile(ctx)
	run := Runner{
		Context: &common.Context{
			SkipUserInteraction: false,
			DryRun:              false,
		},
		Config:           conf,
		HistoryInterface: mHistory,
		StateInterface:   mState,
	}

	mHistory.EXPECT().
		Cleanup().
		Times(1)
	mState.EXPECT().
		Cleanup().
		Times(1)

	mHistory.
		EXPECT().
		InitializeHistory().
		Times(1).
		Return(&history.StorageHistory{
			SchemaVersion:    "v1",
			AppliedMigration: nil,
		}, nil)

	mState.EXPECT().
		InitializeState().
		MaxTimes(1).
		Return(nil)

	mHistory.EXPECT().
		IsMigrationApplied(gomock.Any()).
		Times(2).
		Return(false, nil)

	mState.EXPECT().
		Move(state.ResourceContext{
			State:    "us-east-1/apis/rest",
			Resource: "aws_lambda_function.rest_api",
		}, state.ResourceContext{
			State:    "us-east-1/apis/rest_v2",
			Resource: "aws_lambda_function.rest_api",
		}).
		Times(1).
		Return(false, fmt.Errorf("invalid resource address"))

	mHistory.EXPECT().
		StoreFailedMigration(&history.FailedStorageHistoryObject{
			SchemaVersion: history.StorageHistoryObjectVersion,
			Hash:          migrationFileHash(filepath.Join(testDir, "migrations", "V1__move.hcl")),
			Name:          "V1__move.hcl",
		}).
		Times(1)
	mHistory.EXPECT().
		WriteToStorage().
		Times(1).
		Return(nil)

	err := run.Apply(nil)
	assert.EqualError(t, err, fmt.Sprintf("failed to apply migrtaion '%s' '%s' \n with error: %s",
		"V1__move.hcl",
		fmt.Sprintf("Move %s %s -> %s %s",
			"us-east-1/apis/rest",
			"aws_lambda_function.rest_api",
			"us-east-1/apis/rest_v2",
			"aws_lambda_function.rest_api",
		),
		fmt.Errorf("invalid resource address"),
	))
}

func TestShouldFailToApplyWithRemoveCommandFailing(t *testing.T) {
	_ = assert.New(t)
	testDir := t.TempDir()
	CopyTestData(t, "simple", testDir)
	ChangeWorkingDirectory(t, testDir)

	ctrl, mHistory, mState := getMocks(t)
	defer ctrl.Finish()

	ctx := GetContext(filepath.Join(testDir, config.DefaultConfigFile), &map[string]string{
		"ACCOUNT_ID": "12345678",
	})
	conf, _ := config.GetConfigFile(ctx)
	run := Runner{
		Context: &common.Context{
			SkipUserInteraction: false,
			DryRun:              false,
		},
		Config:           conf,
		HistoryInterface: mHistory,
		StateInterface:   mState,
	}

	mHistory.EXPECT().
		Cleanup().
		Times(1)
	mState.EXPECT().
		Cleanup().
		Times(1)

	mHistory.
		EXPECT().
		InitializeHistory().
		Times(1).
		Return(&history.StorageHistory{
			SchemaVersion:    "v1",
			AppliedMigration: nil,
		}, nil)

	mState.EXPECT().
		InitializeState().
		MaxTimes(1).
		Return(nil)

	mHistory.EXPECT().
		IsMigrationApplied(gomock.Any()).
		Times(2).
		Return(false, nil)

	mState.EXPECT().
		Move(state.ResourceContext{
			State:    "us-east-1/apis/rest",
			Resource: "aws_lambda_function.rest_api",
		}, state.ResourceContext{
			State:    "us-east-1/apis/rest_v2",
			Resource: "aws_lambda_function.rest_api",
		}).
		Times(1).
		Return(true, nil)

	mState.EXPECT().
		Remove(state.ResourceContext{
			State:    "us-east-1/files",
			Resource: "file.test_file",
		}).
		Times(1).
		Return(false, fmt.Errorf("invalid resource address"))

	s3HistoryMetadata := *history.StorageS3Metadata{
		SchemaVersion: "v1",
		ChangedObjects: []history.ChangedS3Object{
			{
				Key:           "us-east-1/apis/rest",
				FromVersionId: nil,
				ToVersionId:   "c1ce6dec475dcec2aa8c85ca1397465c",
			},
			{
				Key:           "us-east-1/apis/rest_v2",
				FromVersionId: nil,
				ToVersionId:   "c1ce6dec475dcec2aa8c85ca1397465c",
			},
		},
	}.Wrap()

	mHistory.EXPECT().
		StoreAppliedMigration(&history.AppliedStorageHistoryObject{
			SchemaVersion: history.StorageHistoryObjectVersion,
			Hash:          migrationFileHash(filepath.Join(testDir, "migrations", "V1__move.hcl")),
			Name:          "V1__move.hcl",
			Metadata:      s3HistoryMetadata,
		}).
		Times(1)
	mHistory.EXPECT().
		StoreFailedMigration(&history.FailedStorageHistoryObject{
			SchemaVersion: history.StorageHistoryObjectVersion,
			Hash:          migrationFileHash(filepath.Join(testDir, "migrations", "V2__remove.hcl")),
			Name:          "V2__remove.hcl",
		}).
		Times(1)
	mHistory.EXPECT().
		WriteToStorage().
		Times(2).
		Return(nil)

	mState.EXPECT().
		Complete().
		Times(1).
		Return(&s3HistoryMetadata, nil)

	err := run.Apply(nil)
	assert.EqualError(t, err, fmt.Sprintf("failed to apply migrtaion '%s' '%s' \n with error: %s",
		"V2__remove.hcl",
		fmt.Sprintf("Remove %s %s",
			"us-east-1/files",
			"file.test_file",
		),
		fmt.Errorf("invalid resource address"),
	))
}

func TestShouldApplyMigrations(t *testing.T) {
	_ = assert.New(t)
	testDir := t.TempDir()
	CopyTestData(t, "simple", testDir)
	ChangeWorkingDirectory(t, testDir)

	ctrl, mHistory, mState := getMocks(t)
	defer ctrl.Finish()

	ctx := GetContext(filepath.Join(testDir, config.DefaultConfigFile), &map[string]string{
		"ACCOUNT_ID": "12345678",
	})
	conf, _ := config.GetConfigFile(ctx)
	run := Runner{
		Context: &common.Context{
			SkipUserInteraction: false,
			DryRun:              false,
		},
		Config:           conf,
		HistoryInterface: mHistory,
		StateInterface:   mState,
	}

	mHistory.EXPECT().
		Cleanup().
		Times(1)
	mState.EXPECT().
		Cleanup().
		Times(1)

	mHistory.
		EXPECT().
		InitializeHistory().
		Times(1).
		Return(&history.StorageHistory{
			SchemaVersion:    "v1",
			AppliedMigration: nil,
		}, nil)

	mState.EXPECT().
		InitializeState().
		MaxTimes(1).
		Return(nil)

	mHistory.EXPECT().
		IsMigrationApplied(gomock.Any()).
		Times(2).
		Return(false, nil)

	mState.EXPECT().
		Move(state.ResourceContext{
			State:    "us-east-1/apis/rest",
			Resource: "aws_lambda_function.rest_api",
		}, state.ResourceContext{
			State:    "us-east-1/apis/rest_v2",
			Resource: "aws_lambda_function.rest_api",
		}).
		Times(1).
		Return(true, nil)

	mState.EXPECT().
		Remove(state.ResourceContext{
			State:    "us-east-1/files",
			Resource: "file.test_file",
		}).
		Times(1).
		Return(true, nil)

	s3MoveHistoryMetadata := *history.StorageS3Metadata{
		SchemaVersion: "v1",
		Type:          "s3",
		ChangedObjects: []history.ChangedS3Object{
			{
				Key:           "us-east-1/apis/rest",
				FromVersionId: nil,
				ToVersionId:   "c1ce6dec475dcec2aa8c85ca1397465c",
			},
			{
				Key:           "us-east-1/apis/rest_v2",
				FromVersionId: nil,
				ToVersionId:   "c1ce6dec475dcec2aa8c85ca1397465c",
			},
		},
	}.Wrap()

	s3RemoveHistoryMetadata := *history.StorageS3Metadata{
		SchemaVersion: "v1",
		Type:          "s3",
		ChangedObjects: []history.ChangedS3Object{
			{
				Key:           "us-east-1/files",
				FromVersionId: nil,
				ToVersionId:   "c1ce6dec475dcec2aa8c85ca1397465c",
			},
		},
	}.Wrap()

	mHistory.EXPECT().
		StoreAppliedMigration(&history.AppliedStorageHistoryObject{
			SchemaVersion: history.StorageHistoryObjectVersion,
			Hash:          migrationFileHash(filepath.Join(testDir, "migrations", "V1__move.hcl")),
			Name:          "V1__move.hcl",
			Metadata:      s3MoveHistoryMetadata,
		}).
		Times(1)
	mHistory.EXPECT().
		StoreAppliedMigration(&history.AppliedStorageHistoryObject{
			SchemaVersion: history.StorageHistoryObjectVersion,
			Hash:          migrationFileHash(filepath.Join(testDir, "migrations", "V2__remove.hcl")),
			Name:          "V2__remove.hcl",
			Metadata:      s3RemoveHistoryMetadata,
		}).
		Times(1)
	mHistory.EXPECT().
		WriteToStorage().
		Times(2).
		Return(nil)

	mState.EXPECT().
		Complete().
		Times(1).
		Return(&s3MoveHistoryMetadata, nil)

	mState.EXPECT().
		Complete().
		Times(1).
		Return(&s3RemoveHistoryMetadata, nil)

	err := run.Apply(nil)
	assert.Nil(t, err)
}
