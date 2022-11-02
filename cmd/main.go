package main

import (
	"fmt"
	"log"

	"github.com/dmitruk-v/clim/v0"
)

func main() {
	commands := clim.Commands{
		clim.NewCommand(`command:add resource:\S+ login:\S+ password:\S+`, &CliBlaController{}),
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
