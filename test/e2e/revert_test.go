package e2e

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/command"
	"github.com/aleksanderaleksic/tgmigrate/history"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RevertCommand_AfterRun1(t *testing.T) {
	err := RunE2E(t, func(t *testing.T, testId string) error {
		app := command.GetApp()

		baseArgs := []string{
			"",
			"-yes",
			"-config=../data/e2e/.tgmigrate.hcl",
			fmt.Sprintf("-config-variables=TEST_ID=%s", testId),
		}

		if err := app.Run(append(baseArgs, "apply", "run1")); err != nil {
			return err
		}
		if err := app.Run(append(baseArgs, "revert", "V1__move.hcl")); err != nil {
			return err
		}

		return nil
	}, func(t *testing.T, testId string, beforeState map[string]*tfjson.State, afterState map[string]*tfjson.State, afterHistory *history.StorageHistory) error {

		for key, value := range beforeState {
			assert.Equal(t, *value, *afterState[key])
		}
		assert.Len(t, afterHistory.AppliedMigration, 0)

		return nil
	})

	if err != nil {
		t.Fatalf("Failed to run e2e test, error: %s", err)
	}
}
