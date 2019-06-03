package plugins

type PluginInterface interface {
	GetListeners()
	GetProviders()
}

type Plugin struct {
	PluginInterface
}
