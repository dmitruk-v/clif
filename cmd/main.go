package main

import (
	"log"

	"github.com/dmitruk-v/clim/cmd/controllers"
	"github.com/dmitruk-v/clim/v0"
)

func main() {
	commands := clim.Commands{
		clim.NewCommand(`command:bla`, controllers.NewCliBlaController()),
		clim.NewCommand(`command:add resource:\S+ login:\S+ password:\S+`, nil),
		clim.NewCommand(`command:add resource:\w+ password:\d+`, nil),
		clim.NewCommand(`command:add resource:\d+ password:\w+`, nil),
		clim.QuitCommand,
	}
	cfg := clim.AppConfig{
		Commands: commands,
	}
	app := clim.NewApp(cfg)
	if err := app.ExecuteCommand(commands[0]); err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
