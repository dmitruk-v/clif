package main

import (
	"fmt"
	"log"

	"github.com/dmitruk-v/clim/v0"
)

func main() {
	commands := clim.Commands{
		clim.NewCommand(`^(?P<command>add) (?P<resource>\S+) (?P<login>\S+) (?P<password>\S+)$`, &CliBlaController{}),
		clim.NewCommand(`^(?P<command>remove) (?P<resource>\S+)$`, nil),
		clim.NewCommand(`^(?P<command>update) (?P<resource>\S+)$`, nil),
	}
	app := clim.NewApp(commands)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

type CliBlaController struct{}

func (ctrl *CliBlaController) Handle(req clim.CliRequest) error {
	fmt.Println(req)
	return nil
}
