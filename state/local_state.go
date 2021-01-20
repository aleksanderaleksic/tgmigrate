package state

import (
	"context"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/hashicorp/terraform-exec/tfexec"
	"path/filepath"
)

type LocalState struct {
	state     config.State
	Terraform *tfexec.Terraform
}

func (s *LocalState) InitializeState() error {
	tf, err := initializeTerraformExec(s.state)
	s.Terraform = tf
	return err
}

func (s *LocalState) Move(from ResourceContext, to ResourceContext) (bool, error) {
	fromStateFilePath := filepath.Join(s.getAbsoluteStateDirPath(), from.State, s.state.Config.GetStateFileName())
	toStateFilePath := filepath.Join(s.getAbsoluteStateDirPath(), to.State, s.state.Config.GetStateFileName())

	err := s.Terraform.StateMv(
		context.Background(),
		from.Resource,
		to.Resource,
		tfexec.State(fromStateFilePath),
		tfexec.StateOut(toStateFilePath),
	)
	if err != nil {
		return false, err
	}
	return true, err
}

func (s *LocalState) Remove(resource ResourceContext) (bool, error) {
	stateFilePath := filepath.Join(s.getAbsoluteStateDirPath(), resource.State, s.state.Config.GetStateFileName())

	err := s.Terraform.StateRm(
		context.Background(),
		resource.Resource,
		tfexec.State(stateFilePath),
	)
	if err != nil {
		return false, err
	}
	return true, err
}

func (s *LocalState) getAbsoluteStateDirPath() string {
	path, _ := filepath.Abs(s.state.Config.GetStateDirectory())
	return path
}
