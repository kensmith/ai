package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/kensmith/ai/provider"
)

func handleFlags() (string, string) {
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stdout, "Usage: cat file | ai <model>\nwhere model is one of:")
		for _, providerName := range providers() {
			models := provider.Registry[providerName]
			fmt.Fprintf(os.Stdout, "%s:\n", providerName)
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

func providers() []string {
	p := []string{}
	for k := range provider.Registry {
		p = append(p, k)
	}
	sort.Strings(p)
	return p
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

	fmt.Printf("\n%s's response:\n%v", modelName, response)
}
