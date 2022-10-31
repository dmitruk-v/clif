package clim

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type app struct {
	commands []*command
}

func NewApp(commands Commands) *app {
	return &app{
		commands: make(Commands, 0),
	}
}

func (app *app) Run() error {
	if err := app.loop(); err != nil {
		return err
	}
	return nil
}

func (app *app) loop() error {
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
		cmd, err := app.parseCommand(line)
		if err != nil {
			return err
		}
		_ = cmd
	}
	return nil
}

func (app *app) parseCommand(s string) (*command, error) {
	for _, cmd := range app.commands {
		matches := cmd.rgx.FindStringSubmatch(s)
		if matches == nil {
			continue
		}
		names := cmd.rgx.SubexpNames()
		if len(matches) != len(names) {
			return nil, fmt.Errorf("bad input for command %q", matches[1])
		}
		req := CliRequest{}
		for i := 1; i < len(matches); i++ {
			req[names[i]] = matches[i]
		}
		if err := cmd.controller.Handle(req); err != nil {
			return nil, fmt.Errorf("parse input for command %q: %v", matches[1], err)
		}
	}
	return nil, nil
}
