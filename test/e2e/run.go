package e2e

import (
	"context"
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/test"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/seqsense/s3sync"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const stateBasePath = "../data/e2e/state"
const bucketName = "tgmigrate-e2e-test-bucket"

func RunE2E(
	t *testing.T,
	run func(t *testing.T, testId string) error,
	after func(t *testing.T, testId string, beforeState map[string]*tfjson.State, afterState map[string]*tfjson.State) error,
) error {
	fmt.Printf("Initializing test\n")
	testId, _ := uuid.GenerateUUID()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1"),
	})

	if err != nil {
		return err
	}

	defer clearTestStateFromBucket(sess, testId)

	fmt.Printf("Uploading test state\n")
	beforeStates, err := uploadTestState(t, sess, testId)
	if err != nil {
		return err
	}

	fmt.Printf("Test initialization complete\n")
	fmt.Printf("Running test\n")
	if err := run(t, testId); err != nil {
		return err
	}

	remoteState, err := getRemoteState(t, sess, testId)
	if err := after(t, testId, beforeStates, remoteState); err != nil {
		return err
	}

	return nil
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

func getStateDirs(basePath string) []string {
	return []string{
		filepath.Join(basePath, "file1"),
		filepath.Join(basePath, "file2"),
		filepath.Join(basePath, "file3"),
	}
}

func uploadTestState(t *testing.T, sess *session.Session, testId string) (map[string]*tfjson.State, error) {
	applyDir := t.TempDir()
	tempDir := t.TempDir()

	if err := test.CopyFilesWithPredicate(stateBasePath, applyDir, func(path string, file os.FileInfo) bool {
		return true
	}); err != nil {
		return nil, err
	}

	var states = make(map[string]*tfjson.State, 0)
	for _, workingDir := range getStateDirs(applyDir) {
		terraform, err := initializeTerraformExec(workingDir)
		if err != nil {
			return nil, err
		}
		if err := terraform.Init(context.Background()); err != nil {
			return nil, err
		}
		if err := terraform.Apply(context.Background()); err != nil {
			return nil, err
		}
		state, err := terraform.ShowStateFile(context.Background(), "terraform.tfstate")
		if err != nil {
			return nil, err
		}
		key := filepath.Base(workingDir)
		states[key] = state
	}

	if err := test.CopyFilesWithPredicate(applyDir, tempDir, func(path string, file os.FileInfo) bool {
		return file.Name() == "terraform.tfstate" && !strings.Contains(path, ".terraform")
	}); err != nil {
		return nil, err
	}

	syncManager := s3sync.New(sess)
	err := syncManager.Sync(tempDir, "s3://"+filepath.Join(bucketName, testId))
	if err != nil {
		return nil, err
	}

	return states, nil
}

func clearTestStateFromBucket(sess *session.Session, testId string) error {
	fmt.Printf("Clearing test state\n")
	s3Client := s3.New(sess)

	listOutput, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(testId),
	})
	if err != nil {
		return err
	}

	var objectsToDelete []*s3.ObjectIdentifier
	for _, obj := range listOutput.Contents {
		fmt.Printf("Deleting %s from: s3://%s \n", *obj.Key, bucketName)
		objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
			Key: obj.Key,
		})
	}

	_, err = s3Client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(bucketName),
		Delete: &s3.Delete{
			Objects: objectsToDelete,
			Quiet:   aws.Bool(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func getRemoteState(t *testing.T, session *session.Session, testId string) (map[string]*tfjson.State, error) {
	dir := t.TempDir()
	destinationPath := filepath.Join(dir, "after")
	if err := test.CopyFilesWithPredicate(stateBasePath, destinationPath, func(path string, file os.FileInfo) bool {
		return true
	}); err != nil {
		return nil, err
	}

	err := filepath.Walk(destinationPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "terraform.tfstate") {
			if err := os.Remove(path); err != nil {
				return err
			}
		}
		return nil
	})

	manager := s3sync.New(session)
	if err := manager.Sync("s3://"+filepath.Join(bucketName, testId), destinationPath); err != nil {
		return nil, err
	}

	var states = make(map[string]*tfjson.State, 0)
	err = filepath.Walk(destinationPath, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, "terraform.tfstate") {
			key := filepath.Base(filepath.Dir(path))
			state, err := getLocalState(filepath.Dir(path))
			if err != nil {
				return err
			}
			states[key] = state
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return states, nil
}

func getLocalState(stateDir string) (*tfjson.State, error) {
	tf, err := initializeTerraformExec(stateDir)
	if err != nil {
		return nil, err
	}
	if err := tf.Init(context.Background()); err != nil {
		return nil, err
	}
	state, err := tf.ShowStateFile(context.Background(), "terraform.tfstate")
	if err != nil {
		return nil, err
	}
	return state, nil
}
