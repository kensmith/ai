srcs := $(shell find . -type f -name "*.go")

gpt: $(srcs) $(MAKEFILE_LIST)
	go fmt ./...
	go mod tidy
	go build -o $@ main.go

.PHONY: answer
answer:
	cat question | ./gpt gpt-4o > answer
	echo >> question
	echo "============================================================" >> question
	echo >> question
	cat answer >> question
