.PHONY: build
build-debug:
	go build -o ./.bin/api/ cmd/restapi/main.go
	mkdir -p ./.bin/api/files
	cp -r ./files ./.bin/api
	cp	.env ./.bin/api


run: build-debug
	./.bin/api/main.exe

production:

.DEFAULT_GOAL := build-debug