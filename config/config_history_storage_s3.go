package config

import (
	"github.com/hashicorp/hcl/v2/gohcl"
)

type S3HistoryStorageConfig struct {
	Bucket string `hcl:"bucket,optional"`
	Key    string `hcl:"key"`
}

func ParseHistoryS3StorageConfig(configFile File) (*S3HistoryStorageConfig, error) {
	var config S3HistoryStorageConfig
	diags := gohcl.DecodeBody(configFile.Migration.History.Storage.Remain, nil, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}
