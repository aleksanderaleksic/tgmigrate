package e2e

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/command"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplyCommand(t *testing.T) {
	t.Parallel()
	err := RunE2E(t, func(t *testing.T, testId string) error {
		app := command.GetApp()

		args := []string{
			"",
			"-yes",
			"-config=../data/e2e/.tgmigrate.hcl",
			fmt.Sprintf("-config-variables=TEST_ID=%s", testId),
			"apply",
			"test",
		}

		if err := app.Run(args); err != nil {
			return err
		}

		return nil
	}, func(t *testing.T, testId string, beforeState map[string]*tfjson.State, afterState map[string]*tfjson.State) error {

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

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to run e2e test, error: %s", err)
	}
}

func TestApplyCommandWithoutMigrationFiles(t *testing.T) {
	t.Parallel()
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
	}, func(t *testing.T, testId string, beforeState map[string]*tfjson.State, afterState map[string]*tfjson.State) error {
		assert.Equal(t, beforeState, afterState)
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to run e2e test, error: %s", err)
	}
}
