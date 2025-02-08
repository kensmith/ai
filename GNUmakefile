gpt: main.go
	go fmt ./...
	go mod tidy
	go build -o $@ main.go

.PHONY: answer
answer:
	cat question | ./gpt gpt-4o > answer
	echo >> question
	echo "============================================================" >> question
	echo >> question
	cat answer >> question
