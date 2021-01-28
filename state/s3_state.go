package state

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
)

type S3State struct {
	context   common.Context
	State     config.State
	Sync      S3Sync
	Terraform *tfexec.Terraform
	Cache     common.Cache
}

func (s *S3State) InitializeState() error {
	err := os.MkdirAll(getStateDirPath(s.Cache), os.ModePerm)
	if err != nil {
		return err
	}

	tf, err := initializeTerraformExec(getStateDirPath(s.Cache))
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
	if err != nil {
		return err
	}

	return nil
}

func (s S3State) Move(from ResourceContext, to ResourceContext) (bool, error) {
	return move(s.Terraform, getStateDirPath(s.Cache), s.State.Config.GetStateFileName(), from, to)
}

func (s S3State) Remove(resource ResourceContext) (bool, error) {
	return remove(s.Terraform, getStateDirPath(s.Cache), s.State.Config.GetStateFileName(), resource)
}

func (s S3State) Cleanup() {
	os.RemoveAll(getStateDirPath(s.Cache))
}
