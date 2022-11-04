package plugins

import (
	"fmt"

	"github.com/dmitruk-v/clim/v0"
)

type blaPlugin struct{}

func NewBlaPlugin() *blaPlugin {
	return &blaPlugin{}
}

func (plug *blaPlugin) Execute(app *clim.App) error {
	fmt.Println("--- hello from bla plugin ---")
	return nil
}
