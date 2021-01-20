package main

import (
	"github.com/aleksanderaleksic/tgmigrate/command"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/aleksanderaleksic/tgmigrate/state"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

// Version is a version number.
var version = "0.0.0"

var globalRunner = migration.Runner{
	HistoryInterface: nil,
	StateInterface:   nil,
	MigrationFiles:   nil,
}
var applyCommand = &command.ApplyCommand{Runner: &globalRunner}
var planCommand = &command.PlanCommand{}

func main() {
	app := &cli.App{
		Version: version,
		Name:    "tgmigrate",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
			&cli.BoolFlag{
				Name:    "yes",
				Aliases: []string{"y"},
				Usage:   "Skip all yes confirm steps",
			},
		},
		Before: func(context *cli.Context) error {
			runner, err := Initialize(context)
			if err != nil {
				return err
			}
			globalRunner = *runner
			return nil
		},
		After: func(context *cli.Context) error {
			return globalRunner.StateInterface.Complete()
		},
		Commands: []*cli.Command{
			applyCommand.GetCLICommand(),
			planCommand.GetCLICommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

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
