BIN_DIR=bin
CMD:=$(patsubst cmd/%/main.go,%,$(shell find cmd -type f -name 'main.go'))
LDFLAGS=-ldflags="-s -w"

DICTIONARY_SCRIPT_DIR=scripts/dictionary-loader

vendor:
	if [ ! -d "vendor" ] || [ -z "$(shell ls -A vendor)" ]; then go mod vendor; fi

build-new-way:
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./bin/ --out slack-bot ./cmd/slack-bot

lint:
	golint -set_exit_status ./events/...
	golint -set_exit_status ./cmd/...
	golint -set_exit_status ./internal/...

imports:
	goimports -d -w $(find . -type f -name '*.go' -not -path "./vendor/*")

format:
	go fmt $(go list ./... | grep -v /vendor/)

tests:
	go test ./...

build-dictionary-script:
	go build -o $(DICTIONARY_SCRIPT_DIR)/dictionary-loader $(DICTIONARY_SCRIPT_DIR)/main.go

.PHONY: vendor