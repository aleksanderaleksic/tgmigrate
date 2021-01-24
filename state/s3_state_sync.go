package state

import (
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
}

func (s S3Sync) DownSync3State() error {
	syncManager := s3sync.New(&s.session)
	err := syncManager.Sync("s3://"+s.config.Bucket, s.config.GetStateDirectory())
	if err != nil {
		return err
	}

	return nil
}

func (s S3Sync) UpSync3State() error {
	//Remove all backup files, dont want to upload them to s3
	_ = filepath.Walk(s.config.GetStateDirectory(),
		func(path string, info os.FileInfo, err error) error {
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

	syncManager := s3sync.New(&s.session, s3sync.WithDryRun())

	err := syncManager.Sync(s.config.GetStateDirectory(), "s3://"+s.config.Bucket)
	if err != nil {
		return err
	}

	return nil
}
