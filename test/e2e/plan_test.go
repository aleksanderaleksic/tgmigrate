package e2e

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/command"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlanCommand(t *testing.T) {
	t.Parallel()
	err := RunE2E(t, func(t *testing.T, testId string) error {
		app := command.GetApp()

		args := []string{
			"",
			"-yes",
			"-config=../data/e2e/.tgmigrate.hcl",
			fmt.Sprintf("-config-variables=TEST_ID=%s", testId),
			"plan",
			"test",
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

func TestPlanCommandWithoutMigrationFiles(t *testing.T) {
	t.Parallel()
	err := RunE2E(t, func(t *testing.T, testId string) error {
		app := command.GetApp()

		args := []string{
			"",
			"-yes",
			"-config=../data/e2e/.tgmigrate.hcl",
			fmt.Sprintf("-config-variables=TEST_ID=%s", testId),
			"plan",
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
