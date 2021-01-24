package state

import (
	"context"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
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

func GetStateInterface(c config.Config, ctx common.Context) (State, error) {
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
			context: ctx,
			State:   c.State,
			Sync: S3Sync{
				config:  conf,
				session: *sess,
			},
			Terraform: nil,
		}, nil
	case "local":
		return &LocalState{
			context:   ctx,
			State:     c.State,
			Terraform: nil,
		}, nil
	default:
		return nil, fmt.Errorf("unknown history storage type: %s", c.History.Storage.Type)
	}
}

func initializeTerraformExec(stateConfig config.State) (*tfexec.Terraform, error) {
	fmt.Println("Initializing terraform")
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

func move(terraform *tfexec.Terraform, stateDir string, stateFileName string, from ResourceContext, to ResourceContext) (bool, error) {
	fmt.Printf("%s %s from %s to %s %s\n",
		common.ColorString(common.Green, "Moving"),
		common.ColorString(common.Blue, from.Resource),
		common.ColorString(common.Gray, from.State),
		common.ColorString(common.Gray, to.State),
		common.ColorString(common.Blue, to.Resource),
	)

	fromStateFilePath, _ := filepath.Abs(filepath.Join(stateDir, from.State, stateFileName))
	toStateFilePath, _ := filepath.Abs(filepath.Join(stateDir, to.State, stateFileName))
	fromBackupStatePath, _ := filepath.Abs(filepath.Join(stateDir, from.State, stateFileName+".backup"))
	toBackupStatePath, _ := filepath.Abs(filepath.Join(stateDir, to.State, stateFileName+".backup"))

	err := terraform.StateMv(
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
	return true, nil
}

func remove(terraform *tfexec.Terraform, stateDir string, stateFileName string, resource ResourceContext) (bool, error) {
	fmt.Printf("%s %s from %s\n",
		common.ColorString(common.Red, "Removing"),
		common.ColorString(common.Blue, resource.Resource),
		common.ColorString(common.Gray, resource.State),
	)

	stateFilePath, _ := filepath.Abs(filepath.Join(stateDir, resource.State, stateFileName))
	backupStatePath, _ := filepath.Abs(filepath.Join(stateDir, resource.State, stateFileName+".backup"))

	err := terraform.StateRm(
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
