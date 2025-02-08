package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	anthropicURL              = "https://api.anthropic.com/v1/messages"
	anthropicApiKeyEnvVarName = "ANTHROPIC_API_KEY"
	anthropicVersion          = "2023-06-01"
	anthropicMaxTokens        = 8192
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

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequestBody struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	Messages  []AnthropicMessage `json:"messages"`
}

type AnthropicResponseBodyContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AnthropicResponseBody struct {
	ID         string                         `json:"id"`
	Type       string                         `json:"type"`
	Role       string                         `json:"role"`
	Model      string                         `json:"model"`
	Content    []AnthropicResponseBodyContent `json:"content"`
	StopReason string                         `json:"stop_reason"`
}

func (p *Anthropic) Request(question string) (string, error) {
	requestBody := AnthropicRequestBody{
		Model:     p.selectedModel,
		MaxTokens: anthropicMaxTokens,
		Messages:  []AnthropicMessage{{Role: "user", Content: question}},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", anthropicURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-api-key", p.apiKey)
	request.Header.Set("anthropic-version", anthropicVersion)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var responseBody AnthropicResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}

	if len(responseBody.Content) <= 0 {
		return "", fmt.Errorf("no response")
	}

	return fmt.Sprintf("%v", responseBody.Content[0].Text), nil
}
