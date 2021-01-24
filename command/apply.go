package command

import (
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/urfave/cli/v2"
)

type ApplyCommand struct {
	Runner *migration.Runner
}

func (command ApplyCommand) GetCLICommand() *cli.Command {
	cmd := cli.Command{
		Name:         "apply",
		Aliases:      nil,
		Usage:        "",
		UsageText:    "",
		Description:  "",
		ArgsUsage:    "",
		Category:     "",
		BashComplete: nil,
		Before:       nil,
		After: func(context *cli.Context) error {
			if command.Runner.StateInterface == nil {
				return nil
			}
			return command.Runner.StateInterface.Complete()
		},
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

func (command ApplyCommand) run(c *cli.Context) error {
	return command.Runner.Apply()
}
