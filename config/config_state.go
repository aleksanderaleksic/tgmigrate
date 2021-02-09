package config

var defaultStateFileName = "terraform.tfstate"

type State struct {
	Type   string
	Config StateConfig
}

type StateConfig interface {
	GetStateFileName() string
}
