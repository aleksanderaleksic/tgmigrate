package state

import (
	"context"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/hashicorp/terraform-exec/tfexec"
	"path/filepath"
)

type S3State struct {
	state     config.State
	Terraform *tfexec.Terraform
}

func (s *S3State) InitializeState() error {
	//TODO download state from s3 bucket to temporary dir
	return nil
}

func (s S3State) Complete() error {
	//TODO upload state to s3 bucket
	//TODO Remove downloaded state dir from local
	return nil
}

func (s S3State) Move(from ResourceContext, to ResourceContext) (bool, error) {
	fromStateFilePath := filepath.Join(s.getAbsoluteStateDirPath(), from.State, s.state.Config.GetStateFileName())
	toStateFilePath := filepath.Join(s.getAbsoluteStateDirPath(), to.State, s.state.Config.GetStateFileName())
	fromBackupStatePath := filepath.Join(s.state.Config.GetBackupStateDirectory(), from.State, s.backupStateFileName())
	toBackupStatePath := filepath.Join(s.state.Config.GetBackupStateDirectory(), to.State, s.backupStateFileName())

	err := s.Terraform.StateMv(
		context.Background(),
		from.Resource,
		to.Resource,
		tfexec.State(fromStateFilePath),
		tfexec.StateOut(toStateFilePath),
		tfexec.Backup(fromBackupStatePath),
		tfexec.BackupOut(toBackupStatePath),
	)
	if err != nil {
		return false, err
	}
	return true, err
}

func (s S3State) Remove(resource ResourceContext) (bool, error) {
	stateFilePath := filepath.Join(s.getAbsoluteStateDirPath(), resource.State, s.state.Config.GetStateFileName())
	backupStatePath := filepath.Join(s.state.Config.GetBackupStateDirectory(), resource.State, s.backupStateFileName())

	err := s.Terraform.StateRm(
		context.Background(),
		resource.Resource,
		tfexec.State(stateFilePath),
		tfexec.Backup(backupStatePath),
	)
	if err != nil {
		return false, err
	}
	return true, err
}

func (s S3State) getAbsoluteStateDirPath() string {
	path, _ := filepath.Abs(s.state.Config.GetStateDirectory())
	return path
}

func (s S3State) backupStateFileName() string {
	return fmt.Sprintf("%s%s", "backup-", s.state.Config.GetStateFileName())
}
