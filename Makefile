BIN_DIR=bin
CMD:=$(patsubst cmd/%/main.go,%,$(shell find cmd -type f -name 'main.go'))
LDFLAGS=-ldflags="-s -w"

DICTIONARY_SCRIPT_DIR=scripts/dictionary-loader

build: $(addprefix $(BIN_DIR)/,$(CMD))

vendor:
	if [ ! -d "vendor" ] || [ -z "$(shell ls -A vendor)" ]; then go mod vendor; fi

$(BIN_DIR)/%: cmd/%/main.go vendor
	env GOOS=darwin go build -mod=vendor $(LDFLAGS) -o $@-mac $<
	env GOOS=darwin GOARCH=amd64 go build -mod=vendor $(LDFLAGS) -o $@-darwin-amd64 $<
	env GOOS=linux GOARCH=amd64 go build -mod=vendor $(LDFLAGS) -o $@-linux-amd64 $<

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

.PHONY: vendor build