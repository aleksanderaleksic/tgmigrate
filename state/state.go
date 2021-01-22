package state

import (
	"context"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
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
	switch c.State.Type {
	case "s3":
		conf := *c.State.Config.(*config.S3StateConfig)
		sess, err := session.NewSession(&aws.Config{
			Region: &conf.Region,
		})

		if err != nil {
			return nil, err
		}

		if conf.AssumeRole != nil {
			credentials := stscreds.NewCredentials(sess, *conf.AssumeRole)
			sess, err = session.NewSession(&aws.Config{
				Region:      &conf.Region,
				Credentials: credentials,
			})

			if err != nil {
				return nil, err
			}
		}

		return &S3State{
			State: c.State,
			Sync: S3Sync{
				config:  conf,
				session: *sess,
			},
			Terraform: nil,
		}, nil
	case "local":
		return &LocalState{
			State:     c.State,
			Terraform: nil,
		}, nil
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
