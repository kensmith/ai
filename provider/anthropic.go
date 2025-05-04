package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	anthropicAPIURL           = "https://api.anthropic.com/v1/messages"
	anthropicModelsURL        = "https://docs.anthropic.com/en/docs/about-claude/models/all-models"
	anthropicApiKeyEnvVarName = "ANTHROPIC_API_KEY"
	anthropicVersion          = "2023-06-01"
	anthropicMaxTokens        = 8192
)

func init() {
	Register("anthropic", NewAnthropic, anthropicModelsURL)
}

type Anthropic struct {
	apiKey        string
	selectedModel string
}

func NewAnthropic(model string) (Provider, error) {
	p := Anthropic{}

	var err error
	p.apiKey, err = getEnvVar(anthropicApiKeyEnvVarName)
	if err != nil {
		return nil, err
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

	request, err := http.NewRequest("POST", anthropicAPIURL, bytes.NewBuffer(jsonData))
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
