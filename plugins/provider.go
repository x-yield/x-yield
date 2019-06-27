package plugins

type ProviderInterface interface {
	Get()
}

type Provider struct {
	ProviderInterface
}
