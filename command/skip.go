package command

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/urfave/cli/v2"
)

type SkipCommand struct {
	Skipper *migration.Skipper
}

func (command *SkipCommand) GetCLICommand() *cli.Command {
	cmd := cli.Command{
		Name:                   "skip",
		Aliases:                nil,
		Usage:                  "Skips the provided migration",
		UsageText:              "",
		Description:            "This is useful if a migration is not able to complete successfully",
		ArgsUsage:              "",
		Category:               "",
		BashComplete:           nil,
		Before:                 command.initialize,
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

func (command *SkipCommand) run(c *cli.Context) error {
	migrationFileName := c.Args().First()
	if migrationFileName == "" {
		return fmt.Errorf("must include the migration file name as argument")
	}
	return command.Skipper.Skip(migrationFileName)
}

func (command *SkipCommand) initialize(c *cli.Context) error {
	cfg, err := config.GetConfigFile(c)
	if err != nil {
		return err
	}

	chc := common.Cache{
		ConfigFilePath: cfg.Path,
	}

	ctx := common.Context{
		SkipUserInteraction: c.Bool("y"),
		DryRun:              false,
	}

	historyInterface, err := history.GetHistoryInterface(*cfg, ctx, chc)
	if err != nil {
		return err
	}

	s3Cfg := cfg.State.Config.(*config.S3StateConfig)
	awsSession, err := getAwsSession(s3Cfg.Region, s3Cfg.AssumeRole)
	if err != nil {
		return err
	}

	skipper := migration.Skipper{
		Context:          &ctx,
		Config:           cfg,
		HistoryInterface: historyInterface,
		S3StateClient:    s3.New(awsSession),
	}

	command.Skipper = &skipper

	return nil
}
