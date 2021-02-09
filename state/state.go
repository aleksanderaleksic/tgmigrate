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
	"github.com/seqsense/s3sync"
	"os/exec"
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
	Cleanup()
}

func GetStateInterface(c config.Config, ctx common.Context, cache common.Cache) (State, error) {
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

		var safeSyncManger *s3sync.Manager
		if ctx.DryRun {
			safeSyncManger = s3sync.New(sess, s3sync.WithDryRun())
		} else {
			safeSyncManger = s3sync.New(sess)
		}

		return &S3State{
			context: ctx,
			State:   c.State,
			Sync: S3Sync{
				config:          conf,
				syncManager:     *s3sync.New(sess),
				safeSyncManager: *safeSyncManger,
				cache:           cache,
			},
			Terraform: nil,
			Cache:     cache,
		}, nil
	default:
		return nil, fmt.Errorf("unknown history storage type: %s", c.History.Storage.Type)
	}
}

func initializeTerraformExec(workingDirectory string) (*tfexec.Terraform, error) {
	execPath, err := exec.LookPath("terraform")
	if err != nil {
		return nil, err
	}

	workingDir, _ := filepath.Abs(workingDirectory)
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

	var options []tfexec.StateMvCmdOption
	if from.State == to.State {
		options = []tfexec.StateMvCmdOption{
			tfexec.State(fromStateFilePath),
			tfexec.Backup(fromBackupStatePath),
			tfexec.BackupOut(toBackupStatePath),
		}
	} else {
		options = []tfexec.StateMvCmdOption{
			tfexec.State(fromStateFilePath),
			tfexec.StateOut(toStateFilePath),
			tfexec.Backup(fromBackupStatePath),
			tfexec.BackupOut(toBackupStatePath),
		}
	}

	err := terraform.StateMv(
		context.Background(),
		from.Resource,
		to.Resource,
		options...,
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

func getStateDirPath(c common.Cache) string {
	path, _ := filepath.Abs(filepath.Join(c.GetCacheDirectoryPath(), "state"))
	return path
}
