package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStateFileName(t *testing.T) {
	stateName := "state.tfstate"
	config := S3StateConfig{
		Bucket:        "",
		Region:        "",
		StateFileName: &stateName,
		AssumeRole:    nil,
	}
	assert.Equal(t, config.GetStateFileName(), stateName)
}

func TestGetDefaultStateFileName(t *testing.T) {
	config := S3StateConfig{
		Bucket:        "",
		Region:        "",
		StateFileName: nil,
		AssumeRole:    nil,
	}
	assert.Equal(t, config.GetStateFileName(), defaultStateFileName)
}
