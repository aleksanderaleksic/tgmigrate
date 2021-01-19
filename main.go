package main

import (
	"github.com/aleksanderaleksic/tgmigrate/config"
	"io/ioutil"
	"log"
	"os"

	"github.com/aleksanderaleksic/tgmigrate/command"
	"github.com/urfave/cli/v2"
)

// Version is a version number.
var version = "0.0.0"

var defaultConfigFile = ".tgmigrate.hcl"
var applyCommand command.ApplyCommand
var planCommand command.PlanCommand

func main() {
	app := &cli.App{
		Version: version,
		Name:    "tgmigrate",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
		Before: func(context *cli.Context) error {
			confFilePath := getConfigFilePathFromFlags(context)
			source, err := ioutil.ReadFile(confFilePath)
			if err != nil {
				return err
			}

			cfg, err := config.ParseConfigFile(confFilePath,source)
			if err != nil {
				return err
			}

			applyCommand = command.ApplyCommand{Config: *cfg}
			planCommand = command.PlanCommand{Config: *cfg}

			return nil
		},
		Commands: []*cli.Command{
			applyCommand.GetCLICommand(),
			planCommand.GetCLICommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfigFilePathFromFlags(c *cli.Context) string {
	configFlagValue := c.String("config")

	if configFlagValue != "" {
		return configFlagValue
	}
	return defaultConfigFile
}
