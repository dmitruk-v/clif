package clim

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type App struct {
	config  AppConfig
	canQuit bool
	plugins []Plugin
}

func NewApp(cfg AppConfig) *App {
	return &App{
		config:  cfg,
		canQuit: false,
		plugins: cfg.Plugins,
	}
}

func (app *App) Run() error {
	if err := app.runPlugins(); err != nil {
		return app.formatError(err)
	}
	if err := app.runOnStart(); err != nil {
		return app.formatError(err)
	}
	if err := app.runInputLoop(); err != nil {
		return app.formatError(err)
	}
	return nil
}

func (app *App) parseCommand(s string) (*command, error) {
	var found *command
	if QuitCommand.rgx.MatchString(s) {
		return QuitCommand, nil
	}
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
		return nil, fmt.Errorf("parse command: no match for input %q", s)
	}
	return found, nil
}

func (app *App) executeCommand(cmd *command) error {
	req := make(map[string]string)
	for key, val := range cmd.params {
		req[key] = val
	}
	if cmd == QuitCommand {
		app.canQuit = true
		return nil
	}
	if cmd.controller == nil {
		return fmt.Errorf("execute command %q: nil controller", req)
	}
	if err := cmd.controller.Handle(req); err != nil {
		return fmt.Errorf("execute command %q: %v", req, err)
	}
	return nil
}

func (app *App) runPlugins() error {
	for _, plug := range app.plugins {
		if err := plug.Execute(app); err != nil {
			return err
		}
	}
	return nil
}

func (app *App) runOnStart() error {
	if len(app.config.OnStart) > 0 {
		for _, cmd := range app.config.OnStart {
			if err := app.executeCommand(cmd); err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *App) runInputLoop() error {
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
		cmd, err := app.parseCommand(strings.TrimSpace(line))
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

func (app *App) formatError(err error) error {
	return fmt.Errorf("[ERROR]: %v", err)
}

func (app *App) printError(err error) {
	fmt.Printf("[ERROR]: %v\n", err)
}
