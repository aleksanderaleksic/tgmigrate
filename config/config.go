package config

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

type File struct {
	Migration struct {
		MigrationDir string `hcl:"migration_dir"`
		History      struct {
			Storage struct {
				Type   string   `hcl:"type,label"`
				Remain hcl.Body `hcl:",remain"`
			} `hcl:"storage,block"`
		} `hcl:"history,block"`
	} `hcl:"migration,block"`
}

type Config struct {
	MigrationDir string
	History      History
}

type History struct {
	Storage HistoryStorage
}

type HistoryStorage struct {
	Type   string
	Config interface{}
}

func ParseConfigFile(filename string, source []byte) (*Config, error) {
	var f File
	err := hclsimple.Decode(filename, source, nil, &f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config file: %s, err: %s", filename, err)
	}

	storageConfig, err := getHistoryStorageConfig(f)
	if err != nil {
		return nil, err
	}

	cfg := Config{
		MigrationDir: f.Migration.MigrationDir,
		History: History{
			Storage: HistoryStorage{
				Type:   f.Migration.History.Storage.Type,
				Config: storageConfig,
			},
		},
	}

	return &cfg, nil
}

func getHistoryStorageConfig(file File) (interface{}, error) {
	t := file.Migration.History.Storage.Type
	switch t {
	case "s3":
		return ParseS3StorageConfig(file)
	default:
		return nil, fmt.Errorf("failed to get storage block from config file, no storage block configuration with type: %s", t)
	}
}
