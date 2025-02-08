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
	openaiURL              = "https://api.openai.com/v1/chat/completions"
	openaiApiKeyEnvVarName = "OPENAI_API_KEY"
	openaiTemperature      = 0.7
)

var openaiModels = []string{
	"babbage-002",
	"chatgpt-4o-latest",
	"dall-e-2",
	"dall-e-3",
	"davinci-002",
	"gpt-3.5-turbo",
	"gpt-3.5-turbo-0125",
	"gpt-3.5-turbo-1106",
	"gpt-3.5-turbo-16k",
	"gpt-3.5-turbo-instruct",
	"gpt-3.5-turbo-instruct-0914",
	"gpt-4",
	"gpt-4-0125-preview",
	"gpt-4-0613",
	"gpt-4-1106-preview",
	"gpt-4o",
	"gpt-4o-2024-05-13",
	"gpt-4o-2024-08-06",
	"gpt-4o-2024-11-20",
	"gpt-4o-audio-preview",
	"gpt-4o-audio-preview-2024-10-01",
	"gpt-4o-audio-preview-2024-12-17",
	"gpt-4o-mini",
	"gpt-4o-mini-2024-07-18",
	"gpt-4o-mini-audio-preview",
	"gpt-4o-mini-audio-preview-2024-12-17",
	"gpt-4o-mini-realtime-preview",
	"gpt-4o-mini-realtime-preview-2024-12-17",
	"gpt-4o-realtime-preview",
	"gpt-4o-realtime-preview-2024-10-01",
	"gpt-4o-realtime-preview-2024-12-17",
	"gpt-4-turbo",
	"gpt-4-turbo-2024-04-09",
	"gpt-4-turbo-preview",
	"o1-mini",
	"o1-mini-2024-09-12",
	"o1-preview",
	"o1-preview-2024-09-12",
	"omni-moderation-2024-09-26",
	"omni-moderation-latest",
	"text-embedding-3-large",
	"text-embedding-3-small",
	"text-embedding-ada-002",
	"tts-1",
	"tts-1-1106",
	"tts-1-hd",
	"tts-1-hd-1106",
	"whisper-1",
}

func init() {
	for _, name := range openaiModels {
		Register("OpenAI", name, NewOpenAI)
	}
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
	p.apiKey = os.Getenv(openaiApiKeyEnvVarName)
	if p.apiKey == "" {
		return nil, fmt.Errorf("environment variable %v not set", openaiApiKeyEnvVarName)
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

	request, err := http.NewRequest("POST", openaiURL, bytes.NewBuffer(jsonData))
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
