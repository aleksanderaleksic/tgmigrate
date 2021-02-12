package e2e

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/command"
	"github.com/aleksanderaleksic/tgmigrate/history"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlanCommand(t *testing.T) {
	err := RunE2E(t, func(t *testing.T, testId string) error {
		app := command.GetApp()

		args := []string{
			"",
			"-yes",
			"-config=../data/e2e/.tgmigrate.hcl",
			fmt.Sprintf("-config-variables=TEST_ID=%s", testId),
			"plan",
			"run1",
		}

		if err := app.Run(args); err != nil {
			return err
		}

		return nil
	}, func(t *testing.T, testId string, beforeState map[string]*tfjson.State, afterState map[string]*tfjson.State, afterHistory *history.StorageHistory) error {
		assert.Equal(t, beforeState, afterState)
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to run e2e test, error: %s", err)
	}
}

func TestPlanCommandWithoutMigrationFiles(t *testing.T) {
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
	}, func(t *testing.T, testId string, beforeState map[string]*tfjson.State, afterState map[string]*tfjson.State, afterHistory *history.StorageHistory) error {
		assert.Equal(t, beforeState, afterState)
		return nil
	})

	if err != nil {
		t.Fatalf("Failed to run e2e test, error: %s", err)
	}
}
