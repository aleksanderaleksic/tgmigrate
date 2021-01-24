package config

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"os"
)

type S3HistoryStorageConfig struct {
	Bucket      string  `hcl:"bucket"`
	Region      string  `hcl:"region"`
	Key         string  `hcl:"key"`
	AssumeRole  *string `hcl:"assume_role,optional"`
	historyPath *string
}

func ParseHistoryS3StorageConfig(configFile File) (*S3HistoryStorageConfig, error) {
	var config S3HistoryStorageConfig
	diags := gohcl.DecodeBody(configFile.Migration.History.Storage.Remain, nil, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}

func (s *S3HistoryStorageConfig) GetLocalHistoryPath() string {
	if s.historyPath != nil {
		return *s.historyPath
	}
	dirName := s.Bucket + "_history"
	s.historyPath = &dirName
	_ = os.Mkdir(dirName, 0777)
	return *s.historyPath
}
