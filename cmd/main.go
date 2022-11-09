package main

import (
	"log"

	"github.com/dmitruk-v/clim/cmd/controllers"
	"github.com/dmitruk-v/clim/v0"
	"github.com/dmitruk-v/clim/v0/auth"
)

func main() {
	_, err := auth.NewAuthenticator().SignIn()
	if err != nil {
		log.Fatal(err)
	}

	commands := clim.Commands{
		clim.NewCommand(`command:\+ amount:\d+ currency:\w{3}`, controllers.NewCliBlaController()),
		clim.QuitCommand,
	}
	cfg := clim.AppConfig{
		Commands: commands,
	}
	app := clim.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
