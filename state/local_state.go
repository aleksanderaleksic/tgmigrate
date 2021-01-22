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

type LocalState struct {
	State     config.State
	Terraform *tfexec.Terraform
}

func (s *LocalState) InitializeState() error {
	tf, err := initializeTerraformExec(s.State)
	s.Terraform = tf
	return err
}

func (s LocalState) Complete() error {
	os.RemoveAll(filepath.Dir(s.Terraform.ExecPath()))
	if strings.HasPrefix(s.State.Config.GetStateDirectory(), "/tmp") {
		os.RemoveAll(filepath.Dir(s.State.Config.GetStateDirectory()))
	}
	return nil
}

func (s LocalState) Move(from ResourceContext, to ResourceContext) (bool, error) {
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

func (s LocalState) Remove(resource ResourceContext) (bool, error) {
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

func (s LocalState) getAbsoluteStateDirPath() string {
	path, _ := filepath.Abs(s.State.Config.GetStateDirectory())
	return path
}

func (s LocalState) backupStateFileName(resourceContext ResourceContext) string {
	return fmt.Sprintf("%s-%s", strings.ReplaceAll(resourceContext.State, "/", "-"), s.State.Config.GetStateFileName())
}
