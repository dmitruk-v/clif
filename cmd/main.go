package main

import (
	"fmt"
	"log"

	"github.com/dmitruk-v/clim/v0"
)

func main() {
	cfg := clim.AppConfig{
		Commands: clim.Commands{
			clim.NewCommand(`command:\+ amount:\d+ currency:\w{3}`, &depositController{}),
			clim.NewQuitCommand(`command:quit|exit`),
		},
	}
	app := clim.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

type depositController struct{}

func (ctrl *depositController) Handle(req map[string]string) error {
	fmt.Println("deposit controller got request:", req)
	return nil
}
