package command

import (
	"fmt"
	"github.com/aleksanderaleksic/tgmigrate/common"
	"github.com/aleksanderaleksic/tgmigrate/config"
	"github.com/aleksanderaleksic/tgmigrate/history"
	"github.com/aleksanderaleksic/tgmigrate/migration"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/urfave/cli/v2"
)

type RevertCommand struct {
	Reverter *migration.Reverter
}

func (command *RevertCommand) GetCLICommand() *cli.Command {
	cmd := cli.Command{
		Name:                   "revert",
		Aliases:                nil,
		Usage:                  "Reverts migrations",
		UsageText:              "revert [migration file name]",
		Description:            "If a migration have been applied and you want to revert",
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

func (command *RevertCommand) run(c *cli.Context) error {
	migrationFileName := c.Args().First()
	if migrationFileName == "" {
		return fmt.Errorf("must include the migration file name as argument")
	}
	return command.Reverter.RevertToIncluding(migrationFileName)
}

func (command *RevertCommand) initialize(c *cli.Context) error {
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

	reverter := migration.Reverter{
		Context:          &ctx,
		Config:           cfg,
		HistoryInterface: historyInterface,
		S3StateClient:    s3.New(awsSession),
	}

	command.Reverter = &reverter

	return nil
}

func getAwsSession(region string, assumeRole *string) (*session.Session, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return nil, err
	}

	if assumeRole != nil {
		credentials := stscreds.NewCredentials(sess, *assumeRole)
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String(region),
			Credentials: credentials,
		})

		if err != nil {
			return nil, err
		}
	}

	return sess, nil
}
