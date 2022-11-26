package clif

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/inancgumus/screen"
)

type app struct {
	config  AppConfig
	canQuit bool
}

func NewApp(cfg AppConfig) *app {
	cfg.Commands = append(cfg.Commands, helpCommand, quitCommand)
	return &app{
		config:  cfg,
		canQuit: false,
	}
}

func (app *app) Run() error {
	if len(app.config.Commands) == 2 {
		return ErrNoUserCommands
	}
	app.clearConsole()
	if err := app.runOnStart(); err != nil {
		return app.formatError(err)
	}
	if err := app.runInputLoop(); err != nil {
		return app.formatError(err)
	}
	return nil
}

func (app *app) runOnStart() error {
	for _, cmd := range app.config.OnStart {
		if err := app.runCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

func (app *app) runInputLoop() error {
	defer app.clearConsole()
	r := bufio.NewReader(os.Stdin)
	for {
		if app.canQuit {
			return nil
		}
		fmt.Print("> ")
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if err := app.runCommand(strings.TrimSpace(line)); err != nil {
			app.printError(err)
			continue
		}
	}
}

func (app *app) runCommand(src string) error {
	cmd, err := app.matchCommand(src)
	if err != nil {
		return err
	}
	if err := app.executeCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (app *app) matchCommand(src string) (*command, error) {
	var found *command
	for _, cmd := range app.config.Commands {
		names := cmd.rgx.SubexpNames()
		matches := cmd.rgx.FindStringSubmatch(src)
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
		return nil, fmt.Errorf("match command: no match for input %q", src)
	}
	return found, nil
}

func (app *app) executeCommand(cmd *command) error {
	if cmd.ctype == QuitCommand {
		app.canQuit = true
		return nil
	}
	if cmd.ctype == HelpCommand {
		return app.showHelp()
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

func (app *app) showHelp() error {
	var help string
	for _, cmd := range app.config.Commands {
		cmdhelp := fmt.Sprintf("\n%s\n", cmd.help.Info)
		if cmd.help.Usage != nil {
			for _, usage := range cmd.help.Usage {
				cmdhelp += fmt.Sprintf("  %v\n", usage)
			}
		}
		if cmd.help.Examples != nil {
			cmdhelp += "Examples:\n"
			for _, ex := range cmd.help.Examples {
				cmdhelp += fmt.Sprintf("  %v\n", ex)
			}
		}
		help += cmdhelp
	}
	fmt.Println(help)
	return nil
}

func (app *app) clearConsole() {
	screen.Clear()
	screen.MoveTopLeft()
}

func (app *app) formatError(err error) error {
	return fmt.Errorf("[ERROR]: %v", err)
}

func (app *app) printError(err error) {
	fmt.Printf("[ERROR]: %v\n", err)
}
