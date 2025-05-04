package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/kensmith/ai/provider"
)

func handleFlags() (string, string) {
	prov := flag.String("p", "", fmt.Sprintf("<provider> in %v", provider.Providers))
	mod := flag.String("m", "", fmt.Sprintf("<model> from %s", provider.ModelURLs))
	flag.Parse()
	if *prov == "" || *mod == "" {
		fmt.Printf("\nexample:\n\necho 'why is the sky blue?' | ai -p anthropic -m claude-3-7-sonnet-latest\n\n")
		flag.Usage()
		os.Exit(1)
	}
	return *mod, *prov
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
	startTime := time.Now()
	modelName, providerName := handleFlags()
	question := slurpStdin()

	providerFunc, ok := provider.Factory[providerName]
	if !ok {
		fmt.Printf("unknown provider %s\n", providerName)
		os.Exit(1)
	}

	provider, err := providerFunc(modelName)
	if err != nil {
		fmt.Printf("failed to create provider for %s: %v\n", providerName, err)
		os.Exit(1)
	}

	response, err := provider.Request(question)
	if err != nil {
		fmt.Printf("failed to complete request: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n\n# response from %s (in %v):\n\n%v\n\n", modelName, elapsed.Round(time.Second), response)
}
