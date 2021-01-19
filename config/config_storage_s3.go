package config

import (
	"github.com/hashicorp/hcl/v2/gohcl"
)

type S3StorageConfig struct {
	Bucket string `hcl:"bucket,optional"`
	Key    string `hcl:"key"`
}

func ParseS3StorageConfig(configFile File) (*S3StorageConfig, error)  {
	var config S3StorageConfig
	diags := gohcl.DecodeBody(configFile.Migration.History.Storage.Remain, nil, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}
