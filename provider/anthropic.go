package provider

import (
	"fmt"
	"os"
)

const (
	anthropicURL              = "https://api.anthropic.com/v1/messages"
	anthropicApiKeyEnvVarName = "ANTHROPIC_API_KEY"
)

var anthropicModels = []string{
	"claude-3-5-sonnet-latest",
	"claude-3-5-haiku-latest",
	"claude-3-5-sonnet-latest",
	"claude-3-haiku-latest",
	"claude-3-opus-latest",
}

func init() {
	for _, name := range anthropicModels {
		Register("Anthropic", name, NewAnthropic)
	}
}

type Anthropic struct {
	apiKey        string
	selectedModel string
}

func NewAnthropic(model string) (Provider, error) {
	p := Anthropic{}
	p.apiKey = os.Getenv(anthropicApiKeyEnvVarName)
	if p.apiKey == "" {
		return nil, fmt.Errorf("environment variable %v not set", anthropicApiKeyEnvVarName)
	}
	p.selectedModel = model

	return &p, nil
}

func (p *Anthropic) Request(question string) (string, error) {
	return "", nil
}
