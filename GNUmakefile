model := claude-3-5-sonnet-latest
model := gpt-4o
model := grok-2-latest

srcs := $(shell find . -type f -name "*.go")

ai: $(srcs) $(MAKEFILE_LIST)
	go fmt ./...
	go mod tidy
	go build -o $@ main.go

.PHONY: answer
answer: ai
	cat question | ./ai $(model) > answer
	echo >> question
	echo "============================================================" >> question
	echo >> question
	cat answer >> question
