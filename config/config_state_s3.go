package config

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"os"
)

type S3StateConfig struct {
	Bucket               string  `hcl:"bucket"`
	Region               string  `hcl:"region"`
	StateFileName        *string `hcl:"state_file_name,optional"`
	AssumeRole           *string `hcl:"assume_role,optional"`
	stateDirectory       *string
	backupStateDirectory *string
}

func ParseS3StateConfig(configFile File) (*S3StateConfig, error) {
	var config S3StateConfig
	diags := gohcl.DecodeBody(configFile.Migration.State.Remain, nil, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}

func (s *S3StateConfig) GetStateDirectory() string {
	if s.stateDirectory != nil {
		return *s.stateDirectory
	}
	dirName := s.Bucket + "_state"
	s.stateDirectory = &dirName
	_ = os.Mkdir(dirName, 0777)
	return *s.stateDirectory
}

func (s *S3StateConfig) GetBackupStateDirectory() string {
	return s.GetStateDirectory()
}

func (s S3StateConfig) GetStateFileName() string {
	if s.StateFileName != nil {
		return *s.StateFileName
	}
	return defaultStateFileName
}
