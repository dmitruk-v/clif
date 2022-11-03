package clim

type CliController interface {
	Handle(req map[string]string) error
}
