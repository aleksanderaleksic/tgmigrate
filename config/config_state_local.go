package config

import (
	"github.com/hashicorp/hcl/v2/gohcl"
)

type LocalStateConfig struct {
	Directory     string `hcl:"directory"`
	StateFileName *string `hcl:"state_file_name,optional"`
}

func ParseLocalStateConfig(configFile File) (*LocalStateConfig, error) {
	var config LocalStateConfig
	diags := gohcl.DecodeBody(configFile.Migration.State.Remain, nil, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}

func (l LocalStateConfig) GetStateDirectory() string {
	return l.Directory
}

func (l LocalStateConfig) GetStateFileName() string {
	if l.StateFileName != nil {
		return *l.StateFileName
	}
	return defaultStateFileName
}
