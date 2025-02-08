# usage:
#
# make answer p=<provider>
#
# where provider is one of:
# * anthropic
# * openai
# * xai

MAKEFLAGS = -j

use-models := \
  claude-3-5-sonnet-latest \
  gpt-4o \
  grok-2-latest

model-targets :=
$(foreach model,$(use-models), \
  $(eval model-targets += $(model)-target) \
 )

srcs := $(shell find . -type f -name "*.go")
binary := ai

$(binary): $(srcs) $(MAKEFILE_LIST)
	go fmt ./...
	go mod tidy
	go build -o $@ main.go

.PHONY: answers
answers: $(model-targets)

$(foreach model,$(use-models), \
  $(eval .PHONY: $(model)-target) \
  $(eval \
    $(model)-target \
    : $(binary) \
    ; cat question | ./ai $(model) > answer-$(model) \
  ) \
 )

.PHONY: loc
loc:
	scc --no-cocomo
