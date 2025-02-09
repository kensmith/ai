MAKEFLAGS = -j

srcs := $(shell find . -type f -name "*.go")
binary := ai

$(binary): $(srcs) $(MAKEFILE_LIST)
	go fmt ./...
	go mod tidy
	go build -o $@ main.go

.PHONY: loc
loc:
	scc --no-cocomo
