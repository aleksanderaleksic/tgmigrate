package state

import (
	"context"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"path/filepath"
	"strings"
)

type S3State struct {
	State     config.State
	Sync      S3Sync
	Terraform *tfexec.Terraform
}

func (s *S3State) InitializeState() error {
	tf, err := initializeTerraformExec(s.State)
	s.Terraform = tf
	if err != nil {
		return err
	}

	err = s.Sync.DownSync3State()
	if err != nil {
		return err
	}

	return nil
}

func (s S3State) Complete() error {
	err := s.Sync.UpSync3State()
	os.RemoveAll(s.State.Config.GetStateDirectory())
	os.RemoveAll(s.State.Config.GetBackupStateDirectory())

	if err != nil {
		return err
	}

	return nil
}

func (s S3State) Move(from ResourceContext, to ResourceContext) (bool, error) {
	fromStateFilePath := filepath.Join(s.getAbsoluteStateDirPath(), from.State, s.State.Config.GetStateFileName())
	toStateFilePath := filepath.Join(s.getAbsoluteStateDirPath(), to.State, s.State.Config.GetStateFileName())
	fromBackupStatePath := filepath.Join(s.State.Config.GetBackupStateDirectory(), s.backupStateFileName(from))
	toBackupStatePath := filepath.Join(s.State.Config.GetBackupStateDirectory(), s.backupStateFileName(to))

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
	stateFilePath := filepath.Join(s.getAbsoluteStateDirPath(), resource.State, s.State.Config.GetStateFileName())
	backupStatePath := filepath.Join(s.State.Config.GetBackupStateDirectory(), s.backupStateFileName(resource))

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
	path, _ := filepath.Abs(s.State.Config.GetStateDirectory())
	return path
}

func (s S3State) backupStateFileName(resourceContext ResourceContext) string {
	return fmt.Sprintf("%s-%s", strings.ReplaceAll(resourceContext.State, "/", "-"), s.State.Config.GetStateFileName())
}
