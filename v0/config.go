package clim

type AppConfig struct {
	Commands Commands
	OnStart  Commands
	OnQuit   Commands
}
