package main

import (
	"github.com/aleksanderaleksic/tgmigrate/command"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

// Version is a version number.
var version = "0.1.1"

func main() {
	var applyCommand = command.ApplyCommand{}
	var planCommand = command.PlanCommand{}

	app := &cli.App{
		Version: version,
		Name:    "tgmigrate",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
			&cli.BoolFlag{
				Name:    "yes",
				Aliases: []string{"y"},
				Usage:   "Skip all user interaction",
			},
			&cli.StringFlag{
				Name:    "config-variables",
				Aliases: []string{"cv"},
				Usage:   "ACCOUNT=123456789;NAME=test will be applied to the config file strings using ${ACCOUNT} and ${NAME}",
				EnvVars: []string{"TG-MIGRATE_CONFIG_VARIABLES"},
			},
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
