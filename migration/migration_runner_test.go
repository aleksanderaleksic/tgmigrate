package migration

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/state"
	. "github.com/aleksanderaleksic/tgmigrate/test"
	. "github.com/aleksanderaleksic/tgmigrate/test/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func getMocks(t *testing.T) (*gomock.Controller, *MockHistory, *MockState) {
	ctrl := gomock.NewController(t)

	mHistory := NewMockHistory(ctrl)
	mState := NewMockState(ctrl)
	return ctrl, mHistory, mState
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
		StoreMigrationObject("V1__move.hcl", false, gomock.Any()).
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

	mHistory.EXPECT().
		StoreMigrationObject("V1__move.hcl", true, gomock.Any()).
		Times(1)
	mHistory.EXPECT().
		StoreMigrationObject("V2__remove.hcl", false, gomock.Any()).
		Times(1)
	mHistory.EXPECT().
		WriteToStorage().
		Times(1).
		Return(nil)

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

	mHistory.EXPECT().
		StoreMigrationObject("V1__move.hcl", true, gomock.Any()).
		Times(1)
	mHistory.EXPECT().
		StoreMigrationObject("V2__remove.hcl", true, gomock.Any()).
		Times(1)
	mHistory.EXPECT().
		WriteToStorage().
		Times(1).
		Return(nil)
	mState.EXPECT().
		Complete().
		Times(1).
		Return(nil)

	err := run.Apply(nil)
	assert.Nil(t, err)
}
