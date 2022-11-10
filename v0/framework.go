package clim

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type app struct {
	config  AppConfig
	canQuit bool
}

func NewApp(cfg AppConfig) *app {
	return &app{
		config:  cfg,
		canQuit: false,
	}
}

func (app *app) Run() error {
	if len(app.config.Commands) == 0 {
		return ErrNoCommands
	}
	if err := app.runInputLoop(); err != nil {
		return app.formatError(err)
	}
	return nil
}

func (app *app) RunCommand(input string) error {
	if len(app.config.Commands) == 0 {
		return ErrNoCommands
	}
	cmd, err := app.matchCommand(input)
	if err != nil {
		return app.formatError(err)
	}
	if err := app.executeCommand(cmd); err != nil {
		return app.formatError(err)
	}
	return nil
}

func (app *app) matchCommand(s string) (*command, error) {
	var found *command
	for _, cmd := range app.config.Commands {
		names := cmd.rgx.SubexpNames()
		matches := cmd.rgx.FindStringSubmatch(s)
		if matches == nil {
			continue
		}
		found = cmd
		found.params = make(map[string]string)
		for i := 1; i < len(matches); i++ {
			found.params[names[i]] = matches[i]
		}
	}
	if found == nil {
		return nil, fmt.Errorf("match command: no match for input %q", s)
	}
	return found, nil
}

func (app *app) executeCommand(cmd *command) error {
	if cmd.ctype == QuitCommand {
		app.canQuit = true
		return nil
	}
	req := make(map[string]string)
	for key, val := range cmd.params {
		req[key] = val
	}
	if cmd.controller == nil {
		return fmt.Errorf("execute command %q: nil controller", cmd.params["command"])
	}
	if err := cmd.controller.Handle(req); err != nil {
		return fmt.Errorf("execute command %q: %v", cmd.params["command"], err)
	}
	return nil
}

func (app *app) runInputLoop() error {
	rdr := bufio.NewReader(os.Stdin)
	for {
		if app.canQuit {
			return nil
		}
		fmt.Print("> ")
		line, err := rdr.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		cmd, err := app.matchCommand(strings.TrimSpace(line))
		if err != nil {
			app.printError(err)
			continue
		}
		if err := app.executeCommand(cmd); err != nil {
			app.printError(err)
			continue
		}
	}
	return nil
}

func (app *app) formatError(err error) error {
	return fmt.Errorf("[ERROR]: %v", err)
}

func (app *app) printError(err error) {
	fmt.Printf("[ERROR]: %v\n", err)
}
