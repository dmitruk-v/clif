package clif

import (
	"fmt"
	"testing"
)

func TestNoCommands(t *testing.T) {
	cfg := AppConfig{
		Commands: Commands{},
	}
	app := NewApp(cfg)
	if err := app.RunCommand("whatever command"); err == nil {
		t.Errorf("got: no errors, want error: no commands")
	}
	if err := app.Run(); err == nil {
		t.Errorf("got: no errors, want error: no commands")
	}
}

func TestRunCommandSuccess(t *testing.T) {
	ctrl := &mockCliController{}
	cfg := AppConfig{
		Commands: Commands{
			NewCommand(`command:first param1:\d+`, ctrl, CommandHelp{}),
			NewCommand(`command:\+ amount:\d+ currency:\w{3}`, ctrl, CommandHelp{}),
			NewCommand(`command:second|snd param1:\w+ param2:\d+`, ctrl, CommandHelp{}),
			NewCommand(`command:third|trd param1:\d+ param2:\w+ param3:1|2|3`, ctrl, CommandHelp{}),
			NewCommand(`command:fourth|fth param1:\d+ param2:\w+ param3:1|2|3 param4:\w{3}`, ctrl, CommandHelp{}),
			quitCommand,
		},
	}
	app := NewApp(cfg)
	table := []string{
		"first 123",
		"+ 100 usd",
		"snd one 123",
		"trd 123 two 3",
		"fourth 123 two 2 qqq",
		"quit",
	}
	for _, input := range table {
		if err := app.RunCommand(input); err != nil {
			t.Errorf("got error: %v, want no errors", err)
		}
	}
}

func TestMatchCommandSuccess(t *testing.T) {
	cfg := AppConfig{
		Commands: Commands{
			NewCommand(`command:first param1:\d+`, nil, CommandHelp{}),
			NewCommand(`command:second|snd param1:\w+ param2:\d+`, nil, CommandHelp{}),
			NewCommand(`command:third|trd param1:\d+ param2:\w+ param3:1|2|3`, nil, CommandHelp{}),
			quitCommand,
		},
	}
	app := NewApp(cfg)
	table := []string{"first 123", "snd one 123", "trd 123 two 3", "quit"}
	for idx, input := range table {
		cmd, err := app.matchCommand(input)
		if err != nil {
			t.Errorf("got error: %v, want no errors", err)
		}
		if cmd != cfg.Commands[idx] {
			t.Errorf("command for input %q, does not match wanted command", input)
		}
	}
}

func TestMatchCommandNotMatch(t *testing.T) {
	cfg := AppConfig{
		Commands: Commands{
			NewCommand(`command:first param1:\d+`, nil, CommandHelp{}),
			NewCommand(`command:second|snd param1:\w+ param2:\d+`, nil, CommandHelp{}),
			NewCommand(`command:third|trd param1:\d+ param2:\w+ param3:1|2|3`, nil, CommandHelp{}),
		},
	}
	app := NewApp(cfg)
	table := []string{"frst 123", "second|snd one 123", "trd onetwothree two 3", "qexit"}
	for _, input := range table {
		_, err := app.matchCommand(input)
		if err == nil {
			t.Errorf(`got: no errors, want: error: no match for input %q`, input)
		}
	}
}

func TestExecuteCommandNilController(t *testing.T) {
	app := NewApp(AppConfig{})
	cmd := NewCommand("command:whatever", nil, CommandHelp{})
	cmd.params = map[string]string{
		"command": "whatever",
	}
	if err := app.executeCommand(cmd); err == nil {
		t.Errorf("got: no errors, want: error: nil controller")
	}
}

func TestExecuteCommandControllerError(t *testing.T) {
	testParams := map[string]string{
		"command": "error",
	}
	ctrl := &stubCliController{t: t, testParams: testParams}
	app := NewApp(AppConfig{})
	cmd := NewCommand("command:whatever", ctrl, CommandHelp{})
	cmd.params = testParams
	if err := app.executeCommand(cmd); err == nil {
		t.Errorf("got: no errors, want: controller handler error")
	}
}

func TestExecuteCommandQuit(t *testing.T) {
	app := NewApp(AppConfig{})
	if err := app.executeCommand(quitCommand); err != nil {
		t.Errorf("got: %v, want: no errors", err)
	}
	if app.canQuit != true {
		t.Errorf("got: app.canQuit: %v, want: %v", app.canQuit, true)
	}
}

func TestExecuteCommandSuccess(t *testing.T) {
	testParams := map[string]string{
		"command": "third",
		"param1":  "123",
		"param2":  "wordor456",
		"param3":  "2",
	}
	ctrl := &stubCliController{t: t, testParams: testParams}
	app := NewApp(AppConfig{})
	cmd := NewCommand("command:whatever", ctrl, CommandHelp{})
	cmd.params = testParams
	if err := app.executeCommand(cmd); err != nil {
		t.Errorf("got: error: %v, want: no errors", err)
	}
}

type mockCliController struct{}

func (ctrl *mockCliController) Handle(req map[string]string) error {
	return nil
}

type stubCliController struct {
	t          *testing.T
	testParams map[string]string
}

func (ctrl *stubCliController) Handle(req map[string]string) error {
	if len(req) != len(ctrl.testParams) {
		ctrl.t.Errorf(`got request length: %v, want: %v`, len(req), len(ctrl.testParams))
	}
	for tkey, tval := range ctrl.testParams {
		val, ok := req[tkey]
		if !ok {
			ctrl.t.Errorf(`got: %q key not in request map, want %q in map as a key`, tkey, tkey)
		}
		if val != tval {
			ctrl.t.Errorf(`got value: %q for key %q, want value: %q`, val, tkey, tval)
		}
	}
	if req["command"] == "error" {
		return fmt.Errorf("some expected controller error")
	}
	return nil
}
