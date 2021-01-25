package main

import (
	"github.com/aleksanderaleksic/tgmigrate/command"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/aleksanderaleksic/tgmigrate/state"
	"github.com/urfave/cli/v2"
	"github.com/zclconf/go-cty/cty"
	"log"
	"os"
	"strings"
)

// Version is a version number.
var version = "0.1.0"

func main() {
	var runner migration.Runner
	var applyCommand = &command.ApplyCommand{Runner: &runner}

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
			&cli.BoolFlag{
				Name:    "dryrun",
				Aliases: []string{"d"},
				Usage:   "Dont do any permanent changes, like actually uploading the state changes",
			},
			&cli.StringFlag{
				Name:    "conf-variables",
				Aliases: []string{"cv"},
				Usage:   "ACCOUNT=123456789;NAME=test will be applied to the config file strings using ${ACCOUNT} and ${NAME}",
				EnvVars: []string{"TG-MIGRATE_CONFIG_VARIABLES"},
			},
		},
		Before: func(context *cli.Context) error {
			r, err := Initialize(context)
			if err != nil {
				return err
			}
			runner = *r
			return nil
		},
		Commands: []*cli.Command{
			applyCommand.GetCLICommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func Initialize(c *cli.Context) (*migration.Runner, error) {
	configVariables := getConfigVariables(c)
	cfg, err := config.GetConfigFile(c, configVariables)
	if err != nil {
		return nil, err
	}

	ctx := common.Context{
		SkipUserInteraction: c.Bool("y"),
		DryRun:              c.Bool("d"),
	}

	migrationFiles, err := migration.GetMigrationFiles(*cfg)
	if err != nil {
		return nil, err
	}

	stateInterface, err := state.GetStateInterface(*cfg, ctx)
	if err != nil {
		return nil, err
	}

	historyInterface, err := history.GetHistoryInterface(*cfg, ctx)
	if err != nil {
		return nil, err
	}
	_, err = historyInterface.InitializeHistory(ctx)
	if err != nil {
		return nil, err
	}

	runner := migration.Runner{
		Context:          &ctx,
		HistoryInterface: historyInterface,
		StateInterface:   stateInterface,
		MigrationFiles:   *migrationFiles,
	}

	return &runner, nil
}

func getConfigVariables(c *cli.Context) map[string]cty.Value {
	configVariablesFlag := c.String("cv")

	if configVariablesFlag == "" {
		return map[string]cty.Value{}
	}

	rawKeyValue := strings.Split(configVariablesFlag, ";")

	var keyValue = map[string]cty.Value{}

	for _, raw := range rawKeyValue {
		if raw == "" {
			break
		}
		split := strings.Split(raw, "=")
		keyValue[split[0]] = cty.StringVal(split[1])
	}

	return keyValue
}
