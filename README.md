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
from standard input and writes it to standard output. It accepts a single
argument to select the model.

```
$ cat question | ai gpt-4o > answer
```

See the GNUmakefile for this package for a build target that
you can run instead of doing the above by hand.

Create a file called question in the same directory as the
makefile, then:

```
# make answers
```

It will query all of the models selected in `$(use-models)`
in parallel.

to use the default provider. The files `answer-$(model)` will contain the responses.

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
Go                           6       504       83         0      421         58
License                      1        21        4         0       17          0
Makefile                     1        44        8         8       28          0
Markdown                     1        72       22         0       50          0
───────────────────────────────────────────────────────────────────────────────
Total                        9       641      117         8      516         58
───────────────────────────────────────────────────────────────────────────────
Processed 14284 bytes, 0.014 megabytes (SI)
───────────────────────────────────────────────────────────────────────────────
```

<!--
vim:tw=60:
-->
