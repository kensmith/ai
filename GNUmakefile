# usage:
#
# make answer p=<provider>
#
# where provider is one of:
# * anthropic
# * openai
# * xai

anthropic-model := claude-3-5-sonnet-latest
openai-model := gpt-4o
xai-model := grok-2-latest

default-model := $(openai-model)

model := $($(p)-model)

$(if $(strip $(model)), \
  $(comment model is set), \
  $(eval model := $(default-model)) \
 )

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
