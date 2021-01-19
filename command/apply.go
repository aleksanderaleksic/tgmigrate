package command

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/urfave/cli/v2"
	"io/ioutil"
)

type ApplyCommand struct {
	Config config.Config
}

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

func (command ApplyCommand) run(c *cli.Context) error  {
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