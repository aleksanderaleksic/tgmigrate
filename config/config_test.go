package config

import (
	"flag"
	"github.com/aleksanderaleksic/tgmigrate/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
	"testing"
)

func getContextWithConfigFIle(path string) *cli.Context {
	app := cli.NewApp()

	set := flag.NewFlagSet("apply", 0)
	set.String("config", path, "")

	return cli.NewContext(app, set, nil)
}

func TestReadNonExistentConfigFile(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	defer os.RemoveAll(testDir)
	ctx := getContextWithConfigFIle(filepath.Join(testDir, ".tgmigrate.hcl"))

	c, err := GetConfigFile(ctx)
	ass.Nil(c)
	ass.NotNil(err)
}

func TestHandleErrorParsingConfig(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	defer os.RemoveAll(testDir)

	confFilePath := filepath.Join(testDir, ".tgmigrate.hcl")
	testutil.TestFile(t, confFilePath, `
migration {}
`)
	ctx := getContextWithConfigFIle(confFilePath)

	c, err := GetConfigFile(ctx)
	ass.Nil(c)
	ass.NotNil(err)
}

func TestHandleErrorParsingInvalidHistoryStorageConfig(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	defer os.RemoveAll(testDir)

	confFilePath := filepath.Join(testDir, ".tgmigrate.hcl")
	testutil.TestFile(t, confFilePath, `
migration {
  migration = "./migrations"

  history {
    storage "s3" {
      region = "us-east-1"
      assume_role = "arn:aws:iam::12345678:role/history"
      key = "history.json"
    }
  }

  state "s3" {
    bucket = "test-state-bucket"
    region = "us-east-2"
    assume_role = "arn:aws:iam::12345678:role/state"
  }
}
`)
	ctx := getContextWithConfigFIle(confFilePath)

	c, err := GetConfigFile(ctx)
	ass.Nil(c)
	ass.NotNil(err)
}

func TestHandleErrorParsingInvalidStateStorageConfig(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	defer os.RemoveAll(testDir)

	confFilePath := filepath.Join(testDir, ".tgmigrate.hcl")
	testutil.TestFile(t, confFilePath, `
migration {
  migration = "./migrations"

  history {
    storage "s3" {
      bucket = "test-history-bucket"
      region = "us-east-1"
      assume_role = "arn:aws:iam::12345678:role/history"
      key = "history.json"
    }
  }

  state "s3" {
    region = "us-east-2"
    assume_role = "arn:aws:iam::12345678:role/state"
  }
}
`)
	ctx := getContextWithConfigFIle(confFilePath)

	c, err := GetConfigFile(ctx)
	ass.Nil(c)
	ass.NotNil(err)
}

func TestHandleErrorFindingConfigInParentFolder(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	defer os.RemoveAll(testDir)

	app := cli.NewApp()
	set := flag.NewFlagSet("apply", 0)
	ctx := cli.NewContext(app, set, nil)

	c, err := GetConfigFile(ctx)
	ass.Nil(c)
	ass.NotNil(err)
}

func TestFindingConfigInParentFolder(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	confFilePath := filepath.Join("../", ".tgmigrate.hcl")
	defer os.RemoveAll(testDir)
	defer os.RemoveAll(confFilePath)

	err := os.Chdir(testDir)
	if err != nil {
		t.Fatal("Failed to change directory for test")
	}

	testutil.TestFile(t, confFilePath, `
migration {
  migration = "./migrations"

  history {
    storage "s3" {
      bucket = "test-history-bucket"
      region = "us-east-1"
      assume_role = "arn:aws:iam::12345678:role/history"
      key = "history.json"
    }
  }

  state "s3" {
    bucket = "test-state-bucket"
    region = "us-east-2"
    assume_role = "arn:aws:iam::12345678:role/state"
  }
}
`)

	app := cli.NewApp()
	set := flag.NewFlagSet("apply", 0)
	ctx := cli.NewContext(app, set, nil)

	c, err := GetConfigFile(ctx)
	ass.Nil(err)
	ass.NotNil(c)
}

func TestReadConfig(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	defer os.RemoveAll(testDir)

	confFilePath := filepath.Join(testDir, ".tgmigrate.hcl")
	testutil.TestFile(t, confFilePath, `
migration {
  migration = "./migrations"

  history {
    storage "s3" {
      bucket = "test-history-bucket"
      region = "us-east-1"
      assume_role = "arn:aws:iam::12345678:role/history"
      key = "history.json"
    }
  }

  state "s3" {
    bucket = "test-state-bucket"
    region = "us-east-2"
    assume_role = "arn:aws:iam::12345678:role/state"
  }
}
`)

	ctx := getContextWithConfigFIle(confFilePath)

	c, err := GetConfigFile(ctx)
	if err != nil {
		t.Errorf("error getting config file %s", err)
	}

	ass.Equal(c.Path, confFilePath)
	ass.Equal(c.MigrationDir, "./migrations")
	ass.Equal(c.AbsoluteMigrationDir, filepath.Join(testDir, "migrations"))

	stateAssumeRoleArn := "arn:aws:iam::12345678:role/state"
	ass.Equal(c.State, State{
		Type: "s3",
		Config: &S3StateConfig{
			Bucket:        "test-state-bucket",
			Region:        "us-east-2",
			StateFileName: nil,
			AssumeRole:    &stateAssumeRoleArn,
		},
	})

	historyAssumeRoleArn := "arn:aws:iam::12345678:role/history"
	ass.Equal(c.History, History{
		Storage: HistoryStorage{
			Type: "s3",
			Config: &S3HistoryStorageConfig{
				Bucket:     "test-history-bucket",
				Region:     "us-east-1",
				Key:        "history.json",
				AssumeRole: &historyAssumeRoleArn,
			},
		}})
}

func TestReadConfigWithVariables(t *testing.T) {
	ass := assert.New(t)
	testDir := t.TempDir()
	defer os.RemoveAll(testDir)

	confFilePath := filepath.Join(testDir, ".tgmigrate.hcl")
	testutil.TestFile(t, confFilePath, `
migration {
  migration = "./migrations"

  history {
    storage "s3" {
      bucket = "${BUCKET}"
      region = "${REGION}"
      assume_role = "${ASSUME_ROLE}"
      key = "history.json"
    }
  }

  state "s3" {
    bucket = "${BUCKET}"
    region = "${REGION}"
    assume_role = "${ASSUME_ROLE}"
  }
}
`)
	app := cli.NewApp()

	set := flag.NewFlagSet("apply", 0)
	set.String("config", confFilePath, "")
	set.String("cv", "BUCKET=test-bucket;REGION=us-east-1;ASSUME_ROLE=arn:aws:iam::12345678:role/test", "")

	ctx := cli.NewContext(app, set, nil)

	c, err := GetConfigFile(ctx)
	if err != nil {
		t.Errorf("error getting config file %s", err)
		t.FailNow()
	}

	ass.Equal(c.Path, confFilePath)
	ass.Equal(c.MigrationDir, "./migrations")
	ass.Equal(c.AbsoluteMigrationDir, filepath.Join(testDir, "migrations"))

	assumeRoleArn := "arn:aws:iam::12345678:role/test"
	ass.Equal(c.State, State{
		Type: "s3",
		Config: &S3StateConfig{
			Bucket:        "test-bucket",
			Region:        "us-east-1",
			StateFileName: nil,
			AssumeRole:    &assumeRoleArn,
		},
	})

	ass.Equal(c.History, History{
		Storage: HistoryStorage{
			Type: "s3",
			Config: &S3HistoryStorageConfig{
				Bucket:     "test-bucket",
				Region:     "us-east-1",
				Key:        "history.json",
				AssumeRole: &assumeRoleArn,
			},
		}})
}
