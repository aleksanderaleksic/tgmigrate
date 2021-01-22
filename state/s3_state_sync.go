package state

import (
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/seqsense/s3sync"
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
	syncManager := s3sync.New(&s.session, s3sync.WithDryRun())

	err := syncManager.Sync(s.config.GetStateDirectory(), "s3://"+s.config.Bucket)
	if err != nil {
		return err
	}

	return nil
}
