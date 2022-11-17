package clif

import (
	"testing"
)

func TestCommandPatternMatch(t *testing.T) {
	var table = map[*command]string{
		NewCommand(`command:first|fst`, nil, CommandHelp{}):                       `^(?P<command>first|fst)$`,
		NewCommand(`command:second|snd param1:\w+`, nil, CommandHelp{}):           `^(?P<command>second|snd) (?P<param1>\w+)$`,
		NewCommand(`command:third|trd param1:\w+ param2:\d+`, nil, CommandHelp{}): `^(?P<command>third|trd) (?P<param1>\w+) (?P<param2>\d+)$`,
	}
	for cmd, patt := range table {
		if cmd.pattern != patt {
			t.Errorf("got command pattern: %v, want command pattern: %v", cmd.pattern, patt)
		}
	}
}

func TestCommandWrongSyntaxAtCommand(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(`want panic "wrong syntax: first part must be "command:name", got no panics`)
		}
	}()
	NewCommand(`omg:second|snd param1:\w+`, nil, CommandHelp{})
}

func TestCommandWrongSyntaxNoColon(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error(`want panic "wrong syntax: missing colon in pattern", got no panics`)
		}
	}()
	NewCommand(`command:nice|nc param1 \w+ param2:\d+`, nil, CommandHelp{})
}
