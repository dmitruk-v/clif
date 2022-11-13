package main

import (
	"fmt"
	"log"

	"github.com/dmitruk-v/clim/v0"
	"github.com/dmitruk-v/clim/v0/auth"
)

func main() {
	au := auth.NewAuthenticator()
	user, err := au.Run()
	if err != nil {
		log.Fatal(err)
	}
	_ = user
	// ----------------------------------

	cfg := clim.AppConfig{
		Commands: clim.Commands{
			clim.NewCommand(`command:\+ amount:\d+ currency:\w{3}`, NewDepositController()),
			clim.NewQuitCommand(`command:quit|exit`),
		},
	}
	app := clim.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

type depositController struct{}

func NewDepositController() *depositController {
	return &depositController{}
}

func (ctrl *depositController) Handle(req map[string]string) error {
	fmt.Println("deposit controller got request:", req)
	return nil
}
