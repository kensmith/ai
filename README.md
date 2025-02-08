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

### API key

You will need an API key for the provider of the model.
* [Anthropic](https://www.anthropic.com/api)
* [OpenAI](https://openai.com/api/)
* [X.ai](https://x.ai/api)

Then set the API keys in your shell environment.

```
export OPENAI_API_KEY=<your key>
export ANTHROPIC_API_KEY=<your key>
export XAI_API_KEY=<your key>
```
