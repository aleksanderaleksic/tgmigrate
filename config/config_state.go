package config

const defaultStateFileName = "terraform.tfstate"

type State struct {
	Type   string
	Config StateConfig
}

type StateConfig interface {
	GetStateDirectory() string
	GetBackupStateDirectory() string
	GetStateFileName() string
}
