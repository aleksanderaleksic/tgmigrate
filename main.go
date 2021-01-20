package main

import (
	"github.com/aleksanderaleksic/tgmigrate/command"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

// Version is a version number.
var version = "0.0.0"

var applyCommand = &command.ApplyCommand{}
var planCommand = &command.PlanCommand{}

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
			&cli.BoolFlag{
				Name:    "yes",
				Aliases: []string{"y"},
				Usage:   "Skip all yes confirm steps",
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
