package main

import (
	"fmt"
	"log"

	"github.com/dmitruk-v/clif/v0"
)

func main() {
	commands := clif.Commands{
		clif.NewCommand(`command:\+ amount:\d+ currency:\w{3}`, &depositController{}, clif.CommandHelp{
			Info:  "Deposit some amount of currency:",
			Usage: []string{"+ AMOUNT CURRENCY"},
			Examples: []string{
				"+ 100 usd     Add 100 USD",
			},
		}),
	}
	cfg := clif.AppConfig{
		Commands: commands,
		OnStart:  []string{"help"},
	}
	app := clif.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

type depositController struct{}

func (ctrl *depositController) Handle(req map[string]string) error {
	fmt.Println("deposit controller got request:", req)
	return nil
}
