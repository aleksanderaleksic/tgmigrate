package config

import (
	"github.com/hashicorp/hcl/v2/gohcl"
	"io/ioutil"
	"os"
)

type S3StateConfig struct {
	Bucket string `hcl:"bucket"`
	LocalDirectory *string `hcl:"local_directory,optional"`
	StateFileName *string `hcl:"state_file_name,optional"`
}

func ParseS3StateConfig(configFile File) (*S3StateConfig, error) {
	var config S3StateConfig
	diags := gohcl.DecodeBody(configFile.Migration.State.Remain, nil, &config)

	if diags.HasErrors() {
		return nil, diags
	}

	return &config, nil
}

func (s S3StateConfig) GetStateDirectory() string {
	if s.LocalDirectory != nil {
		return *s.LocalDirectory
	}
	tmpDir, _ := ioutil.TempDir("", "tgmigrate")
	defer os.RemoveAll(tmpDir)
	return tmpDir
}

func (s S3StateConfig) GetStateFileName() string {
	if s.StateFileName != nil {
		return *s.StateFileName
	}
	return defaultStateFileName
}