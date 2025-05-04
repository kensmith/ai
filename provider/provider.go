package provider

import (
	"fmt"
	"os"
)

type FactoryFunc func(string) (Provider, error)

var Factory = map[string]FactoryFunc{}
var ModelURLs = []string{}
var Providers = []string{}

type Provider interface {
	Request(string) (string, error)
}

func Register(provider string, new FactoryFunc, modelURL string) {
	Factory[provider] = new
	Providers = append(Providers, provider)
	ModelURLs = append(ModelURLs, modelURL)
}

func getEnvVar(name string) (string, error) {
	key := os.Getenv(name)
	if key == "" {
		return "", fmt.Errorf("environment variable %v not set", openaiApiKeyEnvVarName)
	}
	return key, nil
}
