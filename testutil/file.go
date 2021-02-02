package testutil

import (
	"io/ioutil"
	"testing"
)

func TestFile(t *testing.T, path string, value string) {
	data := []byte(value)
	err := ioutil.WriteFile(path, data, 0700)
	if err != nil {
		t.Fatal("Failed to create test config file")
	}
}
