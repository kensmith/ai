SHELL := bash
MAKEFLAGS += -j

# Usage:
# Add the models you want to use in the form <provider>|<model> using the
# spelling for the models you get from the API documentation from the
# provider to $(models). (See `ai -h` for the API documentation URs.)
#
# Write your question in question.md and then run `make`.
#
# The order of $(models) matters. The answers will be concatenated into
# question.md in that order.

models += anthropic|claude-3-7-sonnet-latest
models += openai|gpt-4o
models += x.ai|grok-3-beta
models += ollama|deepseek-r1:32b

.PHONY: ai
ai:

ai := $(strip $(shell which ai 2>/dev/null))

sanitize-model = \
$(strip \
  $(subst /,^, \
    $(subst :,~,$(1)) \
   ) \
 )

get-provider = \
$(strip \
  $(firstword $(subst |, ,$(1))) \
 )

get-model = \
$(strip \
  $(word 2,$(subst |, ,$(1))) \
 )

# Prefer your globally installed `llm-preamble`. If one isn't found use the
# llm-preamble in the current directory.
preamble-name := llm-preamble
preamble := $(strip $(shell which $(preamble-name) 2>/dev/null))
local-preamble = $(strip $(shell ls $(preamble-name) 2>/dev/null))
$(if $(strip $(preamble)), \
  $(comment found a global preamble command), \
  $(if $(strip $(local-preamble)), \
    $(eval preamble := $(CURDIR)/$(local-preamble)), \
    $(eval preamble := echo "") \
   ) \
 )

# to avoid adding the output from llm-preamble:
# make raw=t
$(if $(strip $(raw)), \
  $(eval preamble := echo) \
 )

answers :=
$(foreach directive,$(models), \
  $(eval provider := $(call get-provider,$(directive))) \
  $(eval model := $(call get-model,$(directive))) \
  $(eval sanitized-model := $(call sanitize-model,$(model))) \
  $(eval answer-file := answer-$(provider)-$(sanitized-model).md) \
  $(eval \
    $(answer-file) \
    : question.md $(MAKEFILE_LIST) \
    ; -cat <($(preamble)) question.md | $(ai) -p $(provider) -m $(model) > $(answer-file) \
   ) \
  $(eval answers += $(answer-file)) \
  $(eval \
    ai: $(answer-file) \
   ) \
 )

$(if $(strip $(ai)), \
  $(eval ai: \
  ; cat question.md $(answers) > .question.md \
  ; mv question.md .question.md.orig \
  ; mv .question.md question.md \
  ), \
  $(eval ai:;@echo "please install the ai binary with 'go install github.com/kensmith/ai@latest'") \
 )
