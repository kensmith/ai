# A very simple LLM tool

## Installation

```
go install github.com/kensmith/ai@latest
```

## Usage

### List the supported models

```
$ ai
```

### Ask a question

```
$ cat question | ai gpt-4o > answer
```

See the GNUmakefile for this package for a build target that you can run instead of doing the above by hand. The last value for `model` selects the model. If you copy this rule into another makefile, you can just delete the `ai` dependency from the rule.
