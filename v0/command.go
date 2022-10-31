package clim

import "regexp"

type Commands []*command

type command struct {
	pattern    string
	rgx        *regexp.Regexp
	controller CliController
}

func NewCommand(pattern string, controller CliController) *command {
	return &command{
		pattern:    pattern,
		rgx:        regexp.MustCompile(pattern),
		controller: controller,
	}
}
