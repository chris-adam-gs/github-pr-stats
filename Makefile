.DEFAULT_GOAL := build

.PHONY: build
build:
	go build -o build/bin/pr-activity cmd/pr-activity/main.go

.PHONY: run
run:
	go run main.go -username <username> -filename <filename>

.PHONY: clean
clean:
	rm -f pr-activity
