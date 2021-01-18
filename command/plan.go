package command

import "github.com/urfave/cli/v2"

var PlanCommand = cli.Command{
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
	Action:                 planCommandAction,
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

func planCommandAction(c *cli.Context) error  {
	return nil
}