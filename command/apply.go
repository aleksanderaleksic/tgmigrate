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
			command.Runner.Cleanup()
			return nil
		},
		Action:                 command.runAll,
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

func (command ApplyCommand) runAll(c *cli.Context) error {
	environment := c.Args().First()
	if environment == "" {
		return command.Runner.Apply(nil)
	}
	return command.Runner.Apply(&environment)
}
