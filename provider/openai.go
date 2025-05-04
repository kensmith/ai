package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	openaiAPIURL           = "https://api.openai.com/v1/chat/completions"
	openaiModelsURL        = "https://platform.openai.com/docs/models"
	openaiApiKeyEnvVarName = "OPENAI_API_KEY"
	openaiTemperature      = 0.7
)

func init() {
	Register("openai", NewOpenAI, openaiModelsURL)
}

type OpenAI struct {
	apiKey        string
	selectedModel string
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIChoice struct {
	Index   int           `json:"index"`
	Message OpenAIMessage `json:"message"`
}

type OpenAIResponseBody struct {
	Choices []OpenAIChoice `json:"choices"`
}

type OpenAIRequestBody struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature"`
}

func NewOpenAI(model string) (Provider, error) {
	p := OpenAI{}

	var err error
	p.apiKey, err = getEnvVar(openaiApiKeyEnvVarName)
	if err != nil {
		return nil, err
	}

	p.selectedModel = model

	return &p, nil
}

func (p *OpenAI) Request(question string) (string, error) {
	requestBody := OpenAIRequestBody{
		Model:       p.selectedModel,
		Messages:    []OpenAIMessage{{Role: "user", Content: question}},
		Temperature: openaiTemperature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", openaiAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+p.apiKey)

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

	var responseBody OpenAIResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}

	if len(responseBody.Choices) <= 0 {
		return "", fmt.Errorf("no response")
	}
	return fmt.Sprintf("%v", responseBody.Choices[0].Message.Content), nil
}
