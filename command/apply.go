package command

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/aleksanderaleksic/tgmigrate/state"
	"github.com/urfave/cli/v2"
)

type ApplyCommand struct {
	Runner *migration.Runner
}

func (command *ApplyCommand) GetCLICommand() *cli.Command {
	cmd := cli.Command{
		Name:         "apply",
		Aliases:      nil,
		Usage:        "",
		UsageText:    "",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       command.initialize,
		After:        nil,
		Action:       command.runAll,
		OnUsageError: nil,
		Subcommands:  nil,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "dryrun",
				Aliases: []string{"d"},
				Usage:   "Dont do any permanent changes, like actually uploading the state changes",
			},
		},
		SkipFlagParsing:        false,
		HideHelp:               false,
		HideHelpCommand:        false,
		Hidden:                 false,
		UseShortOptionHandling: false,
		HelpName:               "",
		CustomHelpTemplate:     "",
	}
	return &cmd
}

func (command *ApplyCommand) runAll(c *cli.Context) error {
	environment := c.Args().First()
	if environment == "" {
		return command.Runner.Apply(nil)
	}
	return command.Runner.Apply(&environment)
}

func (command *ApplyCommand) initialize(c *cli.Context) error {
	cfg, err := config.GetConfigFile(c)
	if err != nil {
		return err
	}

	chc := common.Cache{
		ConfigFilePath: cfg.Path,
	}

	ctx := common.Context{
		SkipUserInteraction: c.Bool("y"),
		DryRun:              c.Bool("d"),
	}

	migrationFiles, err := migration.GetMigrationFiles(*cfg)
	if err != nil {
		return err
	}

	stateInterface, err := state.GetStateInterface(*cfg, ctx, chc)
	if err != nil {
		return err
	}

	historyInterface, err := history.GetHistoryInterface(*cfg, ctx, chc)
	if err != nil {
		return err
	}

	runner := migration.Runner{
		Context:          &ctx,
		HistoryInterface: historyInterface,
		StateInterface:   stateInterface,
		MigrationFiles:   *migrationFiles,
	}

	command.Runner = &runner

	return nil
}
