package test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const testDataDir = "../test/data"

func CopyTestData(t *testing.T, dataSetName string, destination string) {
	copyTestDataWithMode(t, dataSetName, destination, nil)
}

func CopyTestDataWithMode(t *testing.T, dataSetName string, destination string, mode os.FileMode) {
	copyTestDataWithMode(t, dataSetName, destination, &mode)
}

func copyTestDataWithMode(t *testing.T, dataSetName string, destination string, mode *os.FileMode) {
	datasetDir, err := filepath.Abs(filepath.Join(testDataDir, dataSetName))
	if err != nil {
		t.Fatalf("failed to get path for dataset '%s'", dataSetName)
	}

	err = copyFiles(datasetDir, destination, mode)
	if err != nil {
		t.Fatalf("failed to copy dataset '%s' from '%s' to '%s', err: %s", dataSetName, datasetDir, destination, err)
	}
}

func copyFiles(path string, dstPath string, fileMode *os.FileMode) error {
	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, info := range infos {
		var mode os.FileMode
		if fileMode == nil {
			mode = info.Mode()
		} else {
			mode = *fileMode
		}
		srcPath := filepath.Join(path, info.Name())
		if info.IsDir() {
			newDir := filepath.Join(dstPath, info.Name())
			err = os.MkdirAll(newDir, mode)
			if err != nil {
				return err
			}
			err = copyFiles(srcPath, newDir, fileMode)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath, mode)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func copyFile(path string, dstPath string, filemode os.FileMode) error {
	srcF, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcF.Close()

	di, err := os.Stat(dstPath)
	if err != nil {
		return err
	}
	if di.IsDir() {
		_, file := filepath.Split(path)
		dstPath = filepath.Join(dstPath, file)
	}

	dstF, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, filemode)
	if err != nil {
		return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
}
