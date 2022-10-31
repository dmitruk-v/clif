package clim

type CliRequest map[string]string

type CliController interface {
	Handle(req CliRequest) error
}
