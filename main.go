package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// Version is a version number.
var version = "0.0.0"

func main() {
	app := &cli.App{
		Version: version,
		Name: "tgmigrate",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Load configuration from `FILE`",
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}