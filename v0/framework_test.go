package clim

import (
	"fmt"
	"testing"
)

func TestMatchCommandSuccess(t *testing.T) {
	cfg := AppConfig{
		Commands: Commands{
			NewCommand(`command:first param1:\d+`, nil),
			NewCommand(`command:second|snd param1:\w+ param2:\d+`, nil),
			NewCommand(`command:third|trd param1:\d+ param2:\w+ param3:1|2|3`, nil),
			QuitCommand,
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

func TestMatchCommandFail(t *testing.T) {
	cfg := AppConfig{
		Commands: Commands{
			NewCommand(`command:first param1:\d+`, nil),
			NewCommand(`command:second|snd param1:\w+ param2:\d+`, nil),
			NewCommand(`command:third|trd param1:\d+ param2:\w+ param3:1|2|3`, nil),
			QuitCommand,
		},
	}
	app := NewApp(cfg)
	table := []string{"frst 123", "second|snd one 123", "trd onetwothree two 3"}
	for _, input := range table {
		_, err := app.matchCommand(input)
		if err == nil {
			t.Errorf(`got: no errors, want: error: no match for input %q`, input)
		}
	}
}

func TestExecuteCommandNilController(t *testing.T) {
	app := NewApp(AppConfig{})
	cmd := NewCommand("command:whatever", nil)
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
	cmd := NewCommand("command:whatever", ctrl)
	cmd.params = testParams
	if err := app.executeCommand(cmd); err == nil {
		t.Errorf("got: no errors, want: controller handler error")
	}
}

func TestExecuteCommandQuit(t *testing.T) {
	app := NewApp(AppConfig{})
	if err := app.executeCommand(QuitCommand); err != nil {
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
	cmd := NewCommand("command:whatever", ctrl)
	cmd.params = testParams
	if err := app.executeCommand(cmd); err != nil {
		t.Errorf("got: error: %v, want: no errors", err)
	}
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
