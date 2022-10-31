package main

import (
	"log"

	"github.com/dmitruk-v/clim/v0"
)

func main() {
	commands := clim.Commands{
		clim.NewCommand("cmd-1", nil),
		clim.NewCommand("cmd-2", nil),
		clim.NewCommand("cmd-3", nil),
	}
	app := clim.NewApp(commands)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
