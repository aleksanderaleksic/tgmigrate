package command

import "github.com/urfave/cli/v2"

// Version is a version number.
var version = "0.1.5"

func GetApp() *cli.App {
	var applyCommand = ApplyCommand{}
	var planCommand = PlanCommand{}
	var revertCommand = RevertCommand{}

	return &cli.App{
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
			revertCommand.GetCLICommand(),
		},
	}
}
