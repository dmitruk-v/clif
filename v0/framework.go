package clim

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type app struct {
	commands []*command
}

func NewApp(commands Commands) *app {
	return &app{
		commands: commands,
	}
}

func (app *app) Run(onstart ...*command) error {
	if len(onstart) > 0 {
		for _, cmd := range onstart {
			if err := app.executeCommand(cmd); err != nil {
				fmt.Printf("[error]: %v\n", err)
			}
		}
	}
	rdr := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := rdr.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		line = strings.TrimSpace(line)
		if line == "quit" || line == "exit" {
			return nil
		}
		cmd, err := app.parseCommand(line)
		if err != nil {
			fmt.Printf("[error]: %v\n", err)
			continue
		}
		if err := app.executeCommand(cmd); err != nil {
			fmt.Printf("[error]: %v\n", err)
			continue
		}
	}
	return nil
}

func (app *app) parseCommand(s string) (*command, error) {
	var found *command
	for _, cmd := range app.commands {
		matches := cmd.rgx.FindStringSubmatch(s)
		if matches == nil {
			continue
		}
		found = cmd
		names := cmd.rgx.SubexpNames()
		if len(matches) != len(names) {
			return nil, fmt.Errorf("parse command: bad input for command %q: %v", matches[1], matches[2:])
		}
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

func (app *app) executeCommand(cmd *command) error {
	req := make(CliRequest)
	for key, val := range cmd.params {
		req[key] = val
	}
	if cmd.params["command"] == "quit" || cmd.params["command"] == "exit" {
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
