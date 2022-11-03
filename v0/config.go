package clim

type AuthConfig struct {
	Hash []byte
}

type AppConfig struct {
	Commands Commands
	OnStart  Commands
	Plugins  []Plugin
}
