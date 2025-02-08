package provider

var Registry = map[string][]string{}
var ReverseRegistry = map[string]string{}

type FactoryFunc func(string) (Provider, error)

var Factory = map[string]FactoryFunc{}

func Register(provider, model string, new FactoryFunc) {
	_, found := Registry[provider]
	if !found {
		Registry[provider] = []string{}
	}

	Registry[provider] = append(Registry[provider], model)

	ReverseRegistry[model] = provider

	Factory[provider] = new
}
