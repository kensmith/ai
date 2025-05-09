package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	ollamaURLEnvVarName = "OLLAMA_URL"
	defaultOllamaURL    = "http://localhost:11434/api/generate"
	ollamaModelsURL     = "https://ollama.com/library"
)

var ollamaURL string = defaultOllamaURL

func init() {
	newOllamaURL, err := getEnvVar(ollamaURLEnvVarName)
	if err == nil && len(newOllamaURL) > 0 {
		ollamaURL = newOllamaURL
	}
	Register("ollama", NewOllama, ollamaModelsURL)
}

type Ollama struct {
	selectedModel string
}

func NewOllama(model string) (Provider, error) {
	p := Ollama{}
	p.selectedModel = model

	return &p, nil
}

type OllamaRequestBody struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponseBody struct {
	Model           string    `json:"model"`
	CreatedAt       time.Time `json:"created_at"`
	Response        string    `json:"response"`
	Done            bool      `json:"done"`
	DoneReason      string    `json:"done_reason"`
	Context         []int     `json:"context"`
	TotalDuration   int       `json:"total_duration"`
	LoadDuration    int       `json:"load_duration"`
	PromptEvalCount int       `json:"prompt_eval_count"`
	EvalCount       int       `json:"eval_count"`
	EvalDuration    int       `json:"eval_duration"`
}

func (p *Ollama) Request(question string) (string, error) {
	requestBody := OllamaRequestBody{
		Model:  p.selectedModel,
		Prompt: question,
		Stream: false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", ollamaURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
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

	var responseBody OllamaResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", nil
	}

	if len(responseBody.Response) <= 0 {
		return "", fmt.Errorf("no response")
	}

	return fmt.Sprintf("%v", responseBody.Response), nil
}
