package command

import (
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/urfave/cli/v2"
)

type PlanCommand struct {
	Config config.Config
}

func (command PlanCommand) GetCLICommand() *cli.Command {
	cmd := cli.Command{
		Name:                   "plan",
		Aliases:                nil,
		Usage:                  "",
		UsageText:              "",
		Description:            "",
		ArgsUsage:              "",
		Category:               "",
		BashComplete:           nil,
		Before:                 nil,
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

func (command PlanCommand) run(c *cli.Context) error  {
	return nil
}