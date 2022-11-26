package clif

import (
	"fmt"
	"regexp"
	"strings"
)

type CommandHelp struct {
	Info     string
	Usage    []string
	Examples []string
}

type commandType int

const (
	UnknownCommand commandType = iota
	QuitCommand
	HelpCommand
	UserCommand
)

var (
	helpCommand = newTypedCommand(HelpCommand, `command:help`, nil, CommandHelp{
		Info:  "Show this help:",
		Usage: []string{"help"},
	})
	quitCommand = newTypedCommand(QuitCommand, `command:quit`, nil, CommandHelp{
		Info:  "Quit from app:",
		Usage: []string{"quit"},
	})
)

type Commands []*command

type command struct {
	ctype      commandType
	pattern    string
	rgx        *regexp.Regexp
	params     map[string]string
	controller CliController
	help       CommandHelp
}

func newTypedCommand(ctype commandType, pattern string, controller CliController, help CommandHelp) *command {
	parts := strings.Split(pattern, " ")
	result := "^"
	for i, pt := range parts {
		namereg := strings.SplitN(pt, ":", 2)
		if len(namereg) != 2 {
			panic(fmt.Errorf(`wrong syntax in pattern %q, missing colon: %v`, pattern, namereg))
		}
		if i == 0 {
			if namereg[0] != "command" {
				panic(fmt.Errorf(`wrong syntax in pattern %q, first part must be "command:name", got %q`, pattern, namereg))
			}
		}
		if i > 0 {
			result += " "
		}
		result += fmt.Sprintf("(?P<%v>%v)", namereg[0], namereg[1])
	}
	result += "$"
	return &command{
		ctype:      ctype,
		pattern:    result,
		rgx:        regexp.MustCompile(result),
		controller: controller,
		help:       help,
	}
}

func NewCommand(pattern string, controller CliController, help CommandHelp) *command {
	return newTypedCommand(UserCommand, pattern, controller, help)
}
