package command

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/urfave/cli/v2"
	"io/ioutil"
)

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

	var filename = "./example/migrations/20200118_test.hcl"
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	file, err := migration.ParseMigrationFile(filename,source)
	if err != nil {
		return err
	}

	fmt.Println(file)

	return nil
}