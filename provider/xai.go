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
	xaiURL              = "https://api.x.ai/v1/chat/completions"
	xaiApiKeyEnvVarName = "XAI_API_KEY"
)

var xaiModels = []string{
	"grok-2-latest",
	"grok-2-vision-latest",
	"grok-beta",
	"grok-vision-beta",
}

func init() {
	for _, name := range xaiModels {
		Register("X.ai", name, NewXai)
	}
}

type Xai struct {
	apiKey        string
	selectedModel string
}

func NewXai(model string) (Provider, error) {
	p := Xai{}
	p.apiKey = os.Getenv(xaiApiKeyEnvVarName)
	if p.apiKey == "" {
		return nil, fmt.Errorf("environment variable %v not set", xaiApiKeyEnvVarName)
	}
	p.selectedModel = model

	return &p, nil
}

type XaiRequestBody struct {
	Model    string              `json:"model"`
	Messages []XaiRequestMessage `json:"messages"`
}

type XaiRequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type XaiResponseMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Refusal string `json:"refusal"`
}

type XaiChoice struct {
	Index        int                `json:"index"`
	Message      XaiResponseMessage `json:"message"`
	FinishReason string             `json:"finish_reason"`
}

type XaiPromptTokensDetails struct {
}

type XaiUsage struct {
	PromptTokens        int                    `json:"prompt_tokens"`
	CompletionTokens    int                    `json:"completion_tokens"`
	TotalTokens         int                    `json:"total_tokens"`
	PromptTokensDetails XaiPromptTokensDetails `json:"prompt_tokens_details"`
}

type XaiResponseBody struct {
	Id                string      `json:"id"`
	Object            string      `json:"object"`
	Created           int         `json:"created"`
	Model             string      `json:"model"`
	Choices           []XaiChoice `json:"choices"`
	Usage             XaiUsage    `json:"usage"`
	SystemFingerprint string      `json:"system_fingerprint"`
}

func (p *Xai) Request(question string) (string, error) {
	requestBody := XaiRequestBody{
		Model:    p.selectedModel,
		Messages: []XaiRequestMessage{{Role: "user", Content: question}},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", xaiURL, bytes.NewBuffer(jsonData))
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

	var responseBody XaiResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}

	if len(responseBody.Choices) <= 0 {
		return "", fmt.Errorf("no response")
	}

	return fmt.Sprintf("%v", responseBody.Choices[0].Message.Content), nil
	return "", nil
}
