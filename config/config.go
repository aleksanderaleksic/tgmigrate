package config

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"path/filepath"
)

const defaultConfigFile = ".tgmigrate.hcl"

type File struct {
	Migration struct {
		MigrationDir string `hcl:"migration"`
		History      struct {
			Storage struct {
				Type   string   `hcl:"type,label"`
				Remain hcl.Body `hcl:",remain"`
			} `hcl:"storage,block"`
		} `hcl:"history,block"`
		State struct {
			Type   string   `hcl:"type,label"`
			Remain hcl.Body `hcl:",remain"`
		} `hcl:"state,block"`
	} `hcl:"migration,block"`
}

type Config struct {
	MigrationDir         string
	AbsoluteMigrationDir string
	History              History
	State                State
}

func GetConfigFile(ctx *cli.Context) (*Config, error) {
	confFilePath := getConfigFilePathFromFlags(ctx)
	source, err := ioutil.ReadFile(confFilePath)
	if err != nil {
		return nil, err
	}

	conf, err := parseConfigFile(confFilePath, source)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func parseConfigFile(filename string, source []byte) (*Config, error) {
	var f File
	err := hclsimple.Decode(filename, source, nil, &f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %s, err: %s", filename, err)
	}

	historyStorageConfig, err := getHistoryStorageConfig(f)
	if err != nil {
		return nil, err
	}

	stateConfig, err := getStateConfig(f)
	if err != nil {
		return nil, err
	}

	path, err := getAbsoluteMigrationDirPath(filename, f.Migration.MigrationDir)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		MigrationDir:         f.Migration.MigrationDir,
		AbsoluteMigrationDir: *path,
		History: History{
			Storage: HistoryStorage{
				Type:   f.Migration.History.Storage.Type,
				Config: historyStorageConfig,
			},
		},
		State: State{
			Type:   f.Migration.State.Type,
			Config: stateConfig,
		},
	}

	return &cfg, nil
}

func getHistoryStorageConfig(file File) (interface{}, error) {
	t := file.Migration.History.Storage.Type
	switch t {
	case "local":
		return ParseLocalHistoryStorageConfig(file)
	case "s3":
		return ParseHistoryS3StorageConfig(file)
	default:
		return nil, fmt.Errorf("failed to get storage block from config file, no storage block configuration with type: %s", t)
	}
}

func getStateConfig(file File) (StateConfig, error) {
	t := file.Migration.State.Type
	switch t {
	case "local":
		return ParseLocalStateConfig(file)
	case "s3":
		return ParseS3StateConfig(file)
	default:
		return nil, fmt.Errorf("failed to get storage block from config file, no storage block configuration with type: %s", t)
	}
}

func getAbsoluteMigrationDirPath(configFilePath string, migrationDir string) (*string, error) {
	absoluteConfigFilePath, err := filepath.Abs(configFilePath)
	if err != nil {
		return nil, err
	}
	absoluteConfigFileDirPath := filepath.Dir(absoluteConfigFilePath)

	var migrationDirectory = fmt.Sprintf("%s/%s", absoluteConfigFileDirPath, migrationDir)
	migrationDirectory = filepath.Clean(migrationDirectory)

	return &migrationDirectory, nil
}

func getConfigFilePathFromFlags(c *cli.Context) string {
	configFlagValue := c.String("config")

	if configFlagValue != "" {
		path, _ := filepath.Abs(configFlagValue)
		return path
	}
	return defaultConfigFile
}
