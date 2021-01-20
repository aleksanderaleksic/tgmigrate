package command

import (
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/aleksanderaleksic/tgmigrate/state"
	"github.com/urfave/cli/v2"
)

func Initialize(c *cli.Context) (*migration.Runner, error) {
	cfg, err := config.GetConfigFile(c)
	if err != nil {
		return nil, err
	}

	migrationFiles, err := migration.GetMigrationFiles(*cfg)
	if err != nil {
		return nil, err
	}

	historyInterface, err := history.GetHistoryInterface(*cfg)
	if err != nil {
		return nil, err
	}
	_, err = historyInterface.InitializeStorage(c.Bool("y"))
	if err != nil {
		return nil, err
	}

	stateInterface, err := state.GetStateInterface(*cfg)
	if err != nil {
		return nil, err
	}
	err = stateInterface.InitializeState()
	if err != nil {
		return nil, err
	}



	runner := migration.Runner{
		HistoryInterface: historyInterface,
		StateInterface:   stateInterface,
		MigrationFiles:   *migrationFiles,
	}

	return &runner, nil
}


