.PHONY: build
build-debug:
	go build -o ./.bin/api ./cmd/restapi/main.go
	mkdir -p ./.bin/api/files


run: build-debug
	./.bin/api

production:

.DEFAULT_GOAL := build-debug