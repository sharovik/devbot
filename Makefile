BIN_DIR=bin
CMD:=$(patsubst cmd/%/main.go,%,$(shell find cmd -type f -name 'main.go'))
LDFLAGS=-ldflags="-s -w"

DICTIONARY_SCRIPT_DIR=scripts/dictionary-loader

build-linux-64:
	env GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -mod=vendor $(LDFLAGS) -o $@-linux-amd64 $<
	env GOOS=freebsd CGO_ENABLED=1 GOARCH=amd64 go build -mod=vendor $(LDFLAGS) -o $@-freebsd-amd64 $<

build-windows:
	env GOOS=windows CGO_ENABLED=1 GOARCH=amd64 go build -mod=vendor $(LDFLAGS) -o $@-windows-amd64 $<
	env GOOS=windows CGO_ENABLED=1 GOARCH=386 go build -mod=vendor $(LDFLAGS) -o $@-windows-386 $<

build-linux-86:
	env GOOS=linux CGO_ENABLED=1 GOARCH=386 go build -mod=vendor $(LDFLAGS) -o $@-linux-386 $<
	env GOOS=freebsd CGO_ENABLED=1 GOARCH=386 go build -mod=vendor $(LDFLAGS) -o $@-freebsd-386 $<

build-mac:
	env GOOS=darwin go build -mod=vendor $(LDFLAGS) -o $@-mac $<

vendor:
	if [ ! -d "vendor" ] || [ -z "$(shell ls -A vendor)" ]; then go mod vendor; fi

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