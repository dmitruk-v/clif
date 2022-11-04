package main

import (
	"log"

	"github.com/dmitruk-v/clim/cmd/controllers"
	"github.com/dmitruk-v/clim/v0"
	"github.com/dmitruk-v/clim/v0/plugins"
)

func main() {
	commands := clim.Commands{
		clim.NewCommand(`command:bla`, controllers.NewCliBlaController()),
		clim.NewCommand(`command:add resource:\S+ login:\S+ password:\S+`, nil),
		clim.NewCommand(`command:add resource:\w+ password:\d+`, nil),
		clim.NewCommand(`command:add resource:\d+ password:\w+`, nil),
	}
	cfg := clim.AppConfig{
		Commands: commands,
		OnStart:  clim.Commands{},
		Plugins: []clim.Plugin{
			plugins.NewAuthPlugin([]byte{36, 50, 97, 36, 49, 50, 36, 80, 100, 119, 72, 75, 89, 89, 75, 105, 103, 85, 97, 120, 99, 104, 107, 72, 51, 68, 88, 99, 101, 68, 87, 118, 117, 74, 105, 104, 78, 87, 109, 66, 65, 75, 72, 101, 51, 86, 68, 72, 98, 67, 57, 105, 122, 52, 110, 99, 113, 101, 97, 83}),
			plugins.NewBlaPlugin(),
		},
	}
	app := clim.NewApp(cfg)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
