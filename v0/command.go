package clim

import (
	"fmt"
	"regexp"
	"strings"
)

type Commands []*command

type command struct {
	pattern    string
	rgx        *regexp.Regexp
	params     map[string]string
	controller CliController
}

func NewCommand(pattern string, controller CliController) *command {
	parts := strings.Split(pattern, " ")
	result := "^"
	for i, pt := range parts {
		namereg := strings.Split(pt, ":")
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
		pattern:    result,
		rgx:        regexp.MustCompile(result),
		controller: controller,
	}
}

var (
	QuitCommand = NewCommand(`command:quit|exit`, nil)
)
