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
	url              = "https://api.openai.com/v1/chat/completions"
	apiKeyEnvVarName = "OPENAI_API_KEY"
)

var models = []string{
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
	for _, name := range models {
		Register("OpenAI", name, New)
	}
}

type OpenAI struct {
	apiKey        string
	selectedModel string
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

type ResponseBody struct {
	Choices []Choice `json:"choices"`
}

type RequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

func New(model string) (Provider, error) {
	o := OpenAI{}
	o.apiKey = os.Getenv(apiKeyEnvVarName)
	if o.apiKey == "" {
		return nil, fmt.Errorf("environment variable %v not set", apiKeyEnvVarName)
	}
	o.selectedModel = model

	return &o, nil
}

func (o *OpenAI) Request(question string) (string, error) {
	requestBody := RequestBody{
		Model:       o.selectedModel,
		Messages:    []Message{{Role: "user", Content: question}},
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+o.apiKey)

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

	var responseBody ResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return "", err
	}

	if len(responseBody.Choices) <= 0 {
		return "", fmt.Errorf("no response")
	}
	return fmt.Sprintf("%v", responseBody.Choices[0].Message.Content), nil
}
