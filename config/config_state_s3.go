package config

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
)

type S3StateConfig struct {
	Bucket        string  `hcl:"bucket"`
	Region        string  `hcl:"region"`
	Prefix        *string `hcl:"prefix,optional"`
	StateFileName *string `hcl:"state_file_name,optional"`
	AssumeRole    *string `hcl:"assume_role,optional"`
}

func ParseS3StateConfig(configFile File, ctx *hcl.EvalContext) (*S3StateConfig, error) {
	var config S3StateConfig
	diags := gohcl.DecodeBody(configFile.Migration.State.Remain, ctx, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}

func (s S3StateConfig) GetStateFileName() string {
	if s.StateFileName != nil {
		return *s.StateFileName
	}
	return defaultStateFileName
}
