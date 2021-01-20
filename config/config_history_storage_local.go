package config

import (
	"github.com/hashicorp/hcl/v2/gohcl"
)

type LocalHistoryStorageConfig struct {
	Path string `hcl:"path"`
}

func ParseLocalHistoryStorageConfig(configFile File) (*LocalHistoryStorageConfig, error) {
	var config LocalHistoryStorageConfig
	diags := gohcl.DecodeBody(configFile.Migration.History.Storage.Remain, nil, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}
