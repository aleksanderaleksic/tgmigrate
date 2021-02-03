package testutil

import (
	"flag"
	"fmt"
	"github.com/urfave/cli/v2"
	"strings"
)

func GetContext(path string, variables *map[string]string) *cli.Context {
	app := cli.NewApp()

	set := flag.NewFlagSet("apply", 0)
	set.String("config", path, "")

	if variables != nil {
		cvString := ""
		for k, v := range *variables {
			cvString += fmt.Sprintf("%s=%s;", k, v)
		}
		set.String("cv", strings.TrimSuffix(cvString, ";"), "")
	}

	return cli.NewContext(app, set, nil)
}
