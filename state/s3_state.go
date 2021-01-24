package state

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"path/filepath"
)

type S3State struct {
	context   common.Context
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
	err := s.Sync.UpSync3State(s.context.DryRun)
	os.RemoveAll(s.State.Config.GetStateDirectory())
	if _, err := os.Stat(filepath.Dir(s.Terraform.ExecPath())); os.IsNotExist(err) {
		os.RemoveAll(filepath.Dir(s.Terraform.ExecPath()))
	}

	if err != nil {
		return err
	}

	return nil
}

func (s S3State) Move(from ResourceContext, to ResourceContext) (bool, error) {
	return move(s.Terraform, s.getAbsoluteStateDirPath(), s.State.Config.GetStateFileName(), from, to)
}

func (s S3State) Remove(resource ResourceContext) (bool, error) {
	return remove(s.Terraform, s.getAbsoluteStateDirPath(), s.State.Config.GetStateFileName(), resource)
}

func (s S3State) getAbsoluteStateDirPath() string {
	path, _ := filepath.Abs(s.State.Config.GetStateDirectory())
	return path
}
