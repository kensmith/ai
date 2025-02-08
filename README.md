# A minimalist's LLM tool

## Installation

```
go install github.com/kensmith/ai@latest
```

`ai` has no dependencies outside of the Go standard library.

## Usage

### List the supported models

```
$ ai
```

### Ask a question

In the tradition of a single purpose unix tool, `ai` just reads your question
from standard input and writes it to standard output. It accepts a single
argument to select the model.

```
$ cat question | ai gpt-4o > answer
```

See the GNUmakefile for this package for a build target that
you can run instead of doing the above by hand. If you copy
this rule into another makefile, you can just delete the
`ai` dependency from the rule.

Create a file called question in the same directory as the
makefile, then:

```
# make answer p=antropic
```

or just

```
$ make answer
```

to use the default provider. The file `answer` will contain the response.

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

You don't need all of these keys. You only need keys for the models you want to
use.

<!--
vim:tw=60:
-->
