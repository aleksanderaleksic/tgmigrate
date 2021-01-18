package command

import "github.com/urfave/cli/v2"

var ApplyCommand = cli.Command{
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
	Action:                 applyCommandAction,
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

func applyCommandAction(c *cli.Context) error  {
	return nil
}