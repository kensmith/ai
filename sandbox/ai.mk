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

use-models := \
  claude-3-5-sonnet-latest \
  gpt-4o \
  grok-2-latest

MAKEFLAGS = -j

ai := $(strip $(shell which ai 2>/dev/null))

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
      ; cat question.md | ai $(model) > answer-$(model).md \
    ) \
    $(eval \
      ai: $(model)-target \
     ) \
   ) \
 )
