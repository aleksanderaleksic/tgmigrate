package main

import (
	"github.com/aleksanderaleksic/tgmigrate/command"
	"log"
	"os"
)

func main() {
	app := command.GetApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
