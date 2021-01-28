package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type S3HistoryStorageConfig struct {
	Bucket     string  `hcl:"bucket"`
	Region     string  `hcl:"region"`
	Key        string  `hcl:"key"`
	AssumeRole *string `hcl:"assume_role,optional"`
}

func ParseHistoryS3StorageConfig(configFile File, ctx *hcl.EvalContext) (*S3HistoryStorageConfig, error) {
	var config S3HistoryStorageConfig
	diags := gohcl.DecodeBody(configFile.Migration.History.Storage.Remain, ctx, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}
