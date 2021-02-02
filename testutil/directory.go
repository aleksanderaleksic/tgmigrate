package testutil

import (
	"io/ioutil"
	"testing"
)

func CreateTempTestDir(t *testing.T) string {
	testDir, err := ioutil.TempDir("", "tgmigrate_test")
	if err != nil {
		t.Fatal("Failed to create tempdir for test")
	}
	return testDir
}
