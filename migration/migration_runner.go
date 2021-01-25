package migration

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/state"
)

type Runner struct {
	Context          *common.Context
	HistoryInterface history.History
	StateInterface   state.State
	MigrationFiles   []File
}

func (r Runner) Apply(environment *string) error {
	migrationsToBeApplied, err := r.getMigrationsToBeApplied(environment)
	if err != nil {
		return err
	}

	if len(*migrationsToBeApplied) == 0 {
		fmt.Printf("No migrations will be applied")
		return nil
	}

	err = r.StateInterface.InitializeState()
	if err != nil {
		return err
	}

	for _, migrationFile := range *migrationsToBeApplied {
		fmt.Printf("Migrations for %s will be applied\n", migrationFile.Metadata.FileName)

		var isSuccess = true
		var migrationError error
		var failingMigration = ""

		for _, migration := range migrationFile.Migrations {
			if isSuccess == false {
				break
			}
			switch migration.Type {
			case "remove":
				success, removeError := r.StateInterface.Remove(
					state.ResourceContext{
						State:    migration.Remove.State,
						Resource: migration.Remove.Resource,
					},
				)

				if !success {
					isSuccess = false
					migrationError = removeError
					failingMigration = fmt.Sprintf("Remove %s %s",
						migration.Remove.State,
						migration.Remove.Resource,
					)
					break
				}

			case "move":
				success, moveError := r.StateInterface.Move(
					state.ResourceContext{
						State:    migration.Move.From.State,
						Resource: migration.Move.From.Resource,
					}, state.ResourceContext{
						State:    migration.Move.To.State,
						Resource: migration.Move.To.Resource,
					},
				)

				if !success {
					isSuccess = false
					migrationError = moveError
					failingMigration = fmt.Sprintf("Move %s %s -> %s %s",
						migration.Move.From.State,
						migration.Move.From.Resource,
						migration.Move.To.State,
						migration.Move.To.Resource,
					)
					break
				}
			}
		}

		if isSuccess && migrationError == nil {
			r.HistoryInterface.StoreMigrationObject(migrationFile.Metadata.FileName, history.SuccessResult, migrationFile.Metadata.FileHash)
		} else {
			r.HistoryInterface.StoreMigrationObject(migrationFile.Metadata.FileName, history.FailedResult, migrationFile.Metadata.FileHash)
			_ = r.HistoryInterface.WriteToStorage()

			return fmt.Errorf("failed to apply migrtaion '%s' '%s' \n with error: %s", migrationFile.Metadata.FileName, failingMigration, migrationError)
		}
	}

	err = r.HistoryInterface.WriteToStorage()
	if err != nil {
		return err
	}
	return nil
}

func (r Runner) getMigrationsToBeApplied(environment *string) (*[]File, error) {
	var migrationsToBeApplied []File

	for _, migration := range r.MigrationFiles {
		if environment != nil && !common.StringListContains(migration.Config.Environments, *environment) {
			break
		}

		historyResult, err := r.HistoryInterface.IsMigrationApplied(migration.Metadata.FileHash)
		if err != nil {
			return nil, err
		}

		if historyResult.IsUnapplied() || historyResult.IsFailure() {
			migrationsToBeApplied = append(migrationsToBeApplied, migration)
		}
	}

	return &migrationsToBeApplied, nil
}
