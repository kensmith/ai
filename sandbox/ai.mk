# include this makefile from your project makefile to add a target called 'ai'
# which will send question.md to the LLM providers listed in $(use-models)

# usage:
#
# cat << EOF >> GNUmakefile
# include ai.mk
# EOF
# cat << EOF > question.md
# How many "r" are there in the word strawberry?"
# EOF
# make ai
#
# if you have a command whose output generates a common prompt preamble you
# like to use set preamble-name to the name of that command which should appear
# in your PATH. See llm-preamble in this directory for an example.

SHELL = bash

use-models := \
  gpt-4o \
  grok-2-latest \
  claude-3-5-sonnet-latest

MAKEFLAGS = -j

ai := $(strip $(shell which ai 2>/dev/null))

preamble-name := llm-preamble
preamble := $(strip $(shell which $(preamble-name) 2>/dev/null))
$(if $(strip $(preamble)), \
  $(comment found a preamble command), \
  $(eval preamble := echo "") \
 )

.PHONY: ai
$(if $(strip $(ai)), \
  $(eval ai: \
  ; cat question.md answer* > .question.md \
  ; mv question.md .question.md.orig \
  ; mv .question.md question.md \
  ), \
  $(eval ai:;@echo "please install the ai binary with 'go install github.com/kensmith/ai@latest'") \
 )

$(if $(strip $(ai)), \
  $(foreach model,$(use-models), \
    $(eval .PHONY: $(model)-target) \
    $(eval \
      $(model)-target \
      : $(ai) \
      ; cat <($(preamble)) question.md | ai $(model) > answer-$(model).md \
    ) \
    $(eval \
      ai: $(model)-target \
     ) \
   ) \
 )
