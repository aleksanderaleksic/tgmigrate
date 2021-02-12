package e2e

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/command"
	"github.com/aleksanderaleksic/tgmigrate/history"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyFirstRunCommand(t *testing.T) {
	err := RunE2E(t, func(t *testing.T, testId string) error {
		app := command.GetApp()

		args := []string{
			"",
			"-yes",
			"-config=../data/e2e/.tgmigrate.hcl",
			fmt.Sprintf("-config-variables=TEST_ID=%s", testId),
			"apply",
			"run1",
		}

		if err := app.Run(args); err != nil {
			return err
		}

		return nil
	}, func(t *testing.T, testId string, beforeState map[string]*tfjson.State, afterState map[string]*tfjson.State, afterHistory *history.StorageHistory) error {

		assert.Nil(t, afterState["file1"].Values)
		assert.Contains(t, afterState["file2"].Values.RootModule.Resources, beforeState["file1"].Values.RootModule.Resources[0])

		var removedResource *tfjson.StateResource
		for _, resource := range beforeState["file3"].Values.RootModule.Resources {
			if resource.Name == "test_file1" {
				removedResource = resource
				break
			}
		}
		assert.NotNil(t, removedResource)
		assert.NotContains(t, afterState["file3"].Values.RootModule.Resources, removedResource)
		assert.NotNil(t, afterHistory)
		assert.Len(t, afterHistory.AppliedMigration, 2)
		assert.Len(t, afterHistory.FailedMigrations, 0)

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to run e2e test, error: %s", err)
	}
}

func TestApplySecondRunCommand(t *testing.T) {
	err := RunE2E(t, func(t *testing.T, testId string) error {
		app := command.GetApp()

		args1 := []string{
			"",
			"-yes",
			"-config=../data/e2e/.tgmigrate.hcl",
			fmt.Sprintf("-config-variables=TEST_ID=%s", testId),
			"apply",
			"run1",
		}

		args2 := []string{
			"",
			"-yes",
			"-config=../data/e2e/.tgmigrate.hcl",
			fmt.Sprintf("-config-variables=TEST_ID=%s", testId),
			"apply",
			"run2",
		}

		if err := app.Run(args1); err != nil {
			return err
		}
		if err := app.Run(args2); err != nil {
			return err
		}

		return nil
	}, func(t *testing.T, testId string, beforeState map[string]*tfjson.State, afterState map[string]*tfjson.State, afterHistory *history.StorageHistory) error {

		assert.Equal(t, afterState["file1"].Values.RootModule.Resources, beforeState["file1"].Values.RootModule.Resources)

		assert.NotEqual(t, nil, afterHistory)
		assert.Equal(t,
			afterHistory.AppliedMigration[0].Metadata.S3Metadata.ChangedObjects[0].ToVersionId,
			*afterHistory.AppliedMigration[2].Metadata.S3Metadata.ChangedObjects[0].FromVersionId,
		)

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to run e2e test, error: %s", err)
	}
}

func TestApplyCommandWithoutMigrationFiles(t *testing.T) {
	err := RunE2E(t, func(t *testing.T, testId string) error {
		app := command.GetApp()

		args := []string{
			"",
			"-yes",
			"-config=../data/e2e/.tgmigrate.hcl",
			fmt.Sprintf("-config-variables=TEST_ID=%s", testId),
			"apply",
			"no_env",
		}

		if err := app.Run(args); err != nil {
			return err
		}

		return nil
	}, func(t *testing.T, testId string, beforeState map[string]*tfjson.State, afterState map[string]*tfjson.State, afterHistory *history.StorageHistory) error {
		assert.Equal(t, beforeState, afterState)
		assert.Nil(t, afterHistory)
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to run e2e test, error: %s", err)
	}
}
