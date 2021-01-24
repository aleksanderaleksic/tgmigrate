package state

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/hashicorp/terraform-exec/tfexec"
	"os"
	"path/filepath"
	"strings"
)

type LocalState struct {
	context   common.Context
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
	return nil
}

func (s LocalState) Move(from ResourceContext, to ResourceContext) (bool, error) {
	return move(s.Terraform, s.getAbsoluteStateDirPath(), s.State.Config.GetStateFileName(), from, to)
}

func (s LocalState) Remove(resource ResourceContext) (bool, error) {
	return remove(s.Terraform, s.getAbsoluteStateDirPath(), s.State.Config.GetStateFileName(), resource)
}

func (s LocalState) getAbsoluteStateDirPath() string {
	path, _ := filepath.Abs(s.State.Config.GetStateDirectory())
	return path
}

func (s LocalState) backupStateFileName(resourceContext ResourceContext) string {
	return fmt.Sprintf("%s-%s", strings.ReplaceAll(resourceContext.State, "/", "-"), s.State.Config.GetStateFileName())
}
