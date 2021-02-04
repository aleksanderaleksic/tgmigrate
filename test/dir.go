package test

import (
	"os"
	"testing"
)

func ChangeWorkingDirectory(t *testing.T, newDir string) (string, string) {
	return changeWorkingDirectory(t, newDir, true)
}

func changeWorkingDirectory(t *testing.T, newDir string, cleanup bool) (string, string) {
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal("Failed to get current working directory")
	}
	err = os.Chdir(newDir)
	if err != nil {
		t.Fatalf("Failed to change working directory from '%s' to '%s'", oldDir, newDir)
	}
	if cleanup {
		t.Cleanup(func() {
			changeWorkingDirectory(t, oldDir, false)
		})
	}
	return newDir, oldDir
}
