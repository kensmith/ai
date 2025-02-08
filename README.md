# A very simple LLM tool

## Installation

go install github.com/kensmith/ai@latest

## Usage

### List the supported models

$ ai

### Ask a question

$ cat question | ai gpt-4o > answer

See the GNUmakefile for this package for build target that you can run instead of doing the above by hand.
