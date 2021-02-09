package state

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/seqsense/s3sync"
	"os"
	"path/filepath"
	"strings"
)

type S3Sync struct {
	config          config.S3StateConfig
	safeSyncManager s3sync.Manager
	syncManager     s3sync.Manager
	cache           common.Cache
}

func (s S3Sync) DownSync3State() error {
	var s3Prefix = s.config.Prefix
	if s.config.Prefix == nil {
		s3Prefix = aws.String("")
	}
	err := s.syncManager.Sync("s3://"+filepath.Join(s.config.Bucket, *s3Prefix), getStateDirPath(s.cache))
	if err != nil {
		return err
	}

	return nil
}

func (s S3Sync) UpSync3State() error {
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

	var s3Prefix = s.config.Prefix
	if s.config.Prefix == nil {
		s3Prefix = aws.String("")
	}

	err := s.safeSyncManager.Sync(stateDirPath, "s3://"+filepath.Join(s.config.Bucket, *s3Prefix))
	if err != nil {
		return err
	}

	return nil
}
