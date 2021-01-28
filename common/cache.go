package common

import (
	"path/filepath"
)

type Cache struct {
	ConfigFilePath string
}

func (c Cache) GetCacheDirectoryPath() string {
	return filepath.Join(filepath.Dir(c.ConfigFilePath), ".tgmigrate_cache")
}

func (c Cache) filepath(name string) string {
	return filepath.Join(c.GetCacheDirectoryPath(), name)
}
