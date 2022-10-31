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

func (app *app) Run() error {
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
		cmd, err := app.parseCommand(strings.TrimSpace(line))
		if err != nil {
			fmt.Printf("[error]: %v\n", err)
			continue
		}
		_ = cmd
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
			return nil, fmt.Errorf("bad input string for command %q: %v", matches[1], matches[2:])
		}
		req := make(CliRequest)
		for i := 1; i < len(matches); i++ {
			req[names[i]] = matches[i]
		}
		if err := cmd.controller.Handle(req); err != nil {
			return nil, fmt.Errorf("parse input string for command %q: %v", matches[1], err)
		}
	}
	if found == nil {
		return nil, fmt.Errorf("can't find command for input string %q", s)
	}
	return found, nil
}
