package provider

import (
	"fmt"
	"os"
)

type Provider interface {
	Request(string) (string, error)
}

func getEnvVar(name string) (string, error) {
	key := os.Getenv(name)
	if key == "" {
		return "", fmt.Errorf("environment variable %v not set", openaiApiKeyEnvVarName)
	}
	return key, nil
}
