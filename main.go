package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kensmith/gpt/provider"
)

func handleFlags() (string, string) {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stdout, "Usage: cat file | gpt <model>\nwhere model is one of:")
		for provider, models := range provider.Registry {
			fmt.Fprintf(os.Stdout, "%s:\n", provider)
			for _, model := range models {
				fmt.Fprintf(os.Stdout, "\t%s\n", model)
			}
		}
		os.Exit(1)
	}

	model := flag.Arg(0)
	provider, ok := provider.ReverseRegistry[model]
	if !ok {
		fmt.Fprintln(os.Stdout, "no provider provides this model")
		os.Exit(1)
	}

	return model, provider
}

func slurpStdin() string {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stdout, "Error reading stdin:", err)
		os.Exit(1)
	}

	userMessage := strings.TrimSpace(string(input))

	return userMessage
}

func main() {
	modelName, providerName := handleFlags()
	question := slurpStdin()

	provider, err := provider.Factory[providerName](modelName)
	if err != nil {
		fmt.Printf("failed to create provider for %s\n", providerName)
		os.Exit(1)
	}

	response, err := provider.Request(question)
	if err != nil {
		fmt.Printf("failed to complete request: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%v", response)
}
