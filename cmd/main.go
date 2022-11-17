package main

import (
	"fmt"
	"log"

	"github.com/dmitruk-v/clif/v0"
)

func main() {
	cfg := clif.AppConfig{
		Commands: clif.Commands{
			clif.NewCommand(`command:\+ amount:\d+ currency:\w{3}`, &depositController{}),
			clif.NewQuitCommand(`command:quit|exit`),
			clif.NewHelpCommand(`command:help`),
		},
		HelpFile: "./help.inf",
	}
	app := clif.NewApp(cfg)
	if err := app.RunCommand("help"); err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

type depositController struct{}

func (ctrl *depositController) Handle(req map[string]string) error {
	fmt.Println("deposit controller got request:", req)
	return nil
}
