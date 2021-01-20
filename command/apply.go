package command

import (
	"github.com/urfave/cli/v2"
)

type ApplyCommand struct{}

func (command ApplyCommand) GetCLICommand() *cli.Command {
	cmd := cli.Command{
		Name:                   "apply",
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

func (command ApplyCommand) run(c *cli.Context) error {
	runner, err := Initialize(c)
	if err != nil {
		return err
	}

	err = runner.Apply()
	if err != nil {
		return err
	}

	return nil
}
