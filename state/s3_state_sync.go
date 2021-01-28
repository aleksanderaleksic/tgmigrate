package state

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/seqsense/s3sync"
	"os"
	"path/filepath"
	"strings"
)

type S3Sync struct {
	config  config.S3StateConfig
	session session.Session
	cache   common.Cache
}

func (s S3Sync) DownSync3State() error {
	syncManager := s3sync.New(&s.session)
	err := syncManager.Sync("s3://"+s.config.Bucket, getStateDirPath(s.cache))
	if err != nil {
		return err
	}

	return nil
}

func (s S3Sync) UpSync3State(dryRun bool) error {
	stateDirPath := getStateDirPath(s.cache)
	//Remove all backup files, dont want to upload them to s3
	_ = filepath.Walk(stateDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".backup") {
			return nil
		}

		os.RemoveAll(path)

		return nil
	})

	var syncManager *s3sync.Manager
	if dryRun {
		syncManager = s3sync.New(&s.session, s3sync.WithDryRun())
	} else {
		syncManager = s3sync.New(&s.session)
	}

	err := syncManager.Sync(stateDirPath, "s3://"+s.config.Bucket)
	if err != nil {
		return err
	}

	return nil
}
