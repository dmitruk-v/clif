package clif

type CliController interface {
	Handle(req map[string]string) error
}

type ControllerFunc func(req map[string]string) error

func (fn ControllerFunc) Handle(req map[string]string) error {
	return fn(req)
}
