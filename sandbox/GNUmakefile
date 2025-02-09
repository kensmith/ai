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

model-targets :=
$(foreach model,$(use-models), \
  $(eval model-targets += $(model)-target) \
 )

.PHONY: answers
answers: $(model-targets)

$(foreach model,$(use-models), \
  $(eval .PHONY: $(model)-target) \
  $(eval \
    $(model)-target \
    : $(ai) \
    ; cat question | ai $(model) > answer-$(model) \
  ) \
 )
