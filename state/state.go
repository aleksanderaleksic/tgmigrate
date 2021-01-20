package state

import (
	"context"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"io/ioutil"
	"path/filepath"
)

type ResourceContext struct {
	State    string
	Resource string
}

type State interface {
	InitializeState() error
	Complete() error
	Move(from ResourceContext, to ResourceContext) (bool, error)
	Remove(resource ResourceContext) (bool, error)
}

func GetStateInterface(c config.Config) (State, error) {
	switch c.History.Storage.Type {
	case "s3":
		return nil, nil
	case "local":
		return &LocalState{state: c.State}, nil
	default:
		return nil, fmt.Errorf("unknown history storage type: %s", c.History.Storage.Type)
	}
}

func initializeTerraformExec(stateConfig config.State) (*tfexec.Terraform, error) {
	tmpDir, err := ioutil.TempDir("", "tfinstall")
	if err != nil {
		return nil, err
	}
	execPath, err := tfinstall.Find(context.Background(), tfinstall.LatestVersion(tmpDir, false))
	if err != nil {
		return nil, err
	}

	workingDir, _ := filepath.Abs(stateConfig.Config.GetStateDirectory())
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		return nil, err
	}

	return tf, nil
}
