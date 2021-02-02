package command

import (
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/aleksanderaleksic/tgmigrate/state"
	"github.com/urfave/cli/v2"
)

type PlanCommand struct {
	Runner *migration.Runner
}

func (command *PlanCommand) GetCLICommand() *cli.Command {
	cmd := cli.Command{
		Name:                   "plan",
		Aliases:                nil,
		Usage:                  "Plans the migrations",
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

func (command *PlanCommand) run(c *cli.Context) error {
	environment := c.Args().First()
	if environment == "" {
		return command.Runner.Apply(nil)
	}
	return command.Runner.Apply(&environment)
}

func (command *PlanCommand) initialize(c *cli.Context) error {
	cfg, err := config.GetConfigFile(c)
	if err != nil {
		return err
	}

	chc := common.Cache{
		ConfigFilePath: cfg.Path,
	}

	ctx := common.Context{
		SkipUserInteraction: c.Bool("y"),
		DryRun:              true,
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
	}

	command.Runner = &runner

	return nil
}
