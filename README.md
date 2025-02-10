# A minimalist LLM tool

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

In the tradition of a single purpose unix tool, `ai` just reads your question
from standard input and writes the response to standard
output. It accepts a single argument to select the model.

```
$ cat question | ai gpt-4o > answer
```

See `sandbox/ai.mk` in this package for a build target that
you can run instead of doing the above by hand.

Create a file called `question.md` in the same directory as
that makefile, then:

```
$ make
```

It will query all of the models selected in `$(use-models)`
in parallel. The files `answer-$(model).md` will contain the responses.

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

# Why did I build this?

I had been using
[Fabric](https://github.com/danielmiessler/fabric) which is
pretty great but it does a lot more than I want and I found
its UI to be a little heavy. I'm a minimalist so I made a
minimalist tool. This tool has zero dependencies and is
implemented in less than 500 lines of code on its first
release.

```
───────────────────────────────────────────────────────────────────────────────
Language                 Files     Lines   Blanks  Comments     Code Complexity
───────────────────────────────────────────────────────────────────────────────
Go                           6       503       83         0      420         58
License                      1        21        4         0       17          0
Makefile                     1        35        7         0       28          0
Markdown                     1        85       21         0       64          0
───────────────────────────────────────────────────────────────────────────────
Total                        9       644      115         0      529         58
───────────────────────────────────────────────────────────────────────────────
Processed 15901 bytes, 0.016 megabytes (SI)
───────────────────────────────────────────────────────────────────────────────
```


<!--
vim:tw=60:
-->
