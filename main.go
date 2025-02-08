package main

// quick and dirty llm api thing

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Provider struct {
	URL string
}

var openAI = Provider{
	URL: "https://api.openai.com/v1/chat/completions",
}

var Models = map[string]Provider{
	"gpt-4o": openAI,
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

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stdout, "Usage: cat file | gpt <model>\nwhere model is one of:")
		for name, _ := range Models {
			fmt.Fprintf(os.Stdout, "\t%s\n", name)
			os.Exit(1)
		}
	}

	model := flag.Arg(0)
	provider, ok := Models[model]
	if !ok {
		fmt.Fprintln(os.Stdout, "no provider provides this model")
		os.Exit(1)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stdout, "Error: OPENAI_API_KEY environment variable is not set")
		os.Exit(1)
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error reading stdin:", err)
		os.Exit(1)
	}

	userMessage := strings.TrimSpace(string(input))

	if userMessage == "" {
		fmt.Fprintln(os.Stdout, "Error: No input provided")
		os.Exit(1)
	}

	requestBody := RequestBody{
		Model:       model,
		Messages:    []Message{{Role: "user", Content: userMessage}},
		Temperature: 0.7,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error encoding JSON:", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", provider.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error reading response:", err)
		os.Exit(1)
	}

	var response ResponseBody
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Fprintln(os.Stdout, "Error parsing JSON response:", err)
		os.Exit(1)
	}

	if len(response.Choices) > 0 {
		fmt.Println(response.Choices[0].Message.Content)
	} else {
		fmt.Fprintln(os.Stdout, "Error: No response from OpenAI")
		os.Exit(1)
	}
}
