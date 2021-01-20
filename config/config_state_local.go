package config

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"io/ioutil"
)

type LocalStateConfig struct {
	Directory            string  `hcl:"directory"`
	BackupStateDirectory *string `hcl:"backup_state_directory,optional"`
	StateFileName        *string `hcl:"state_file_name,optional"`
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

func (l LocalStateConfig) GetBackupStateDirectory() string {
	if l.BackupStateDirectory != nil {
		return *l.BackupStateDirectory
	}
	path, _ := ioutil.TempDir("", "tgmigrate-state-backup")
	l.BackupStateDirectory = &path
	return path
}

func (l LocalStateConfig) GetStateFileName() string {
	if l.StateFileName != nil {
		return *l.StateFileName
	}
	return defaultStateFileName
}
