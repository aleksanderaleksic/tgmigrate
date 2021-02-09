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
		Name:                   "apply",
		Aliases:                nil,
		Usage:                  "Applies the migrations",
		UsageText:              "",
		Description:            "",
		ArgsUsage:              "",
		Category:               "",
		BashComplete:           nil,
		Before:                 command.initialize,
		After:                  nil,
		Action:                 command.run,
		OnUsageError:           nil,
		Subcommands:            nil,
		Flags:                  nil,
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

func (command *ApplyCommand) run(c *cli.Context) error {
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
		DryRun:              false,
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
		Config:           cfg,
		HistoryInterface: historyInterface,
		StateInterface:   stateInterface,
	}

	command.Runner = &runner

	return nil
}
