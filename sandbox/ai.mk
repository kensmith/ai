# include this makefile from your project makefile to add a target called 'ai'
# which will send question.md to the LLM providers listed in $(use-models)
#
# usage:
#
# cat << EOF >> GNUmakefile
# include ai.mk
# EOF
# cat << EOF > question.md
# How many "r" are there in the word straberrry?"
# EOF
# make ai

MAKEFLAGS = -j

ai := $(wildcard $(CURDIR)/../ai)
$(if $(strip $(ai)), \
  $(comment binary is built), \
  $(eval ai := $(shell which ai)) \
  $(if $(strip $(ai)), \
    $(comment binary found in environment), \
    $(error couldn't find ai binary, please build or install it) \
   ) \
 )

use-models := \
  claude-3-5-sonnet-latest \
  gpt-4o \
  grok-2-latest

.PHONY: ai
ai:
	cat question.md answer*

$(foreach model,$(use-models), \
  $(eval .PHONY: $(model)-target) \
  $(eval \
    $(model)-target \
    : $(ai) \
    ; cat question.md | ai $(model) > answer-$(model).md \
  ) \
  $(eval \
    ai: $(model)-target \
   ) \
 )
