package clim

type Plugin interface {
	Execute(app *App) error
}
