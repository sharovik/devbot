BIN_DIR=bin
PROJECT_BUILD_DIR=project-build

DICTIONARY_SCRIPT_DIR=scripts/dictionary-loader
INSTALL_SCRIPT_DIR=scripts/install
UPDATE_SCRIPT_DIR=scripts/update

PROJECT_BUILD_SCRIPTS_DIR=$(PROJECT_BUILD_DIR)/scripts
PROJECT_BUILD_SCRIPTS_INSTALL_DIR=$(PROJECT_BUILD_DIR)/$(INSTALL_SCRIPT_DIR)
PROJECT_BUILD_SCRIPTS_UPDATE_DIR=$(PROJECT_BUILD_DIR)/$(UPDATE_SCRIPT_DIR)
PROJECT_BUILD_SCRIPTS_DICTIONARY_DIR=$(PROJECT_BUILD_DIR)/$(DICTIONARY_SCRIPT_DIR)

CMD:=$(patsubst cmd/%/main.go,%,$(shell find cmd -type f -name 'main.go'))
LDFLAGS=-ldflags="-s -w"

install-vendors:
	go mod vendor

build-devbot-cross-platform:
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(PROJECT_BUILD_DIR)/$(BIN_DIR)/ --out devbot ./cmd/devbot

build-devbot-for-current-system:
	env CGO_ENABLED=1 go build -o ./bin/devbot-current-system ./cmd/devbot/main.go

code-check:
	make lint
	env APP_ENV="testing" make tests

code-clean:
	make imports
	make format

lint:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck ./events/...
	staticcheck ./cmd/...
	staticcheck ./internal/...

imports:
	goimports -d -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

format:
	go fmt $(shell go list ./... | grep -v /vendor/)

tests:
	go test ./...

create-if-not-exists-defined-events:
	if [ ! -f events/defined-events.go ]; then cp events/defined-events.go.dist events/defined-events.go; fi

tf-container-definition:
	if [ ! -f terraform/local.container_definition.tf ]; then cp terraform/local.container_definition.tf.dist terraform/local.container_definition.tf; fi

create-if-not-exists-env:
	if [ ! -f .env ]; then cp .env.example .env; fi

install:
	make create-missing-files
	make build-installation-script-for-current-system
	./scripts/install/run

update:
	make build-update-script-for-current-system
	make run-update

run-update:
	./scripts/update/run

create-project-build-dirs:
	if [[ ! -d $(PROJECT_BUILD_DIR) ]]; then mkdir $(PROJECT_BUILD_DIR); fi
	if [[ ! -d $(PROJECT_BUILD_SCRIPTS_DIR) ]]; then mkdir $(PROJECT_BUILD_SCRIPTS_DIR); fi
	if [[ ! -d $(PROJECT_BUILD_SCRIPTS_INSTALL_DIR) ]]; then mkdir $(PROJECT_BUILD_SCRIPTS_INSTALL_DIR); fi
	if [[ ! -d $(PROJECT_BUILD_SCRIPTS_UPDATE_DIR) ]]; then mkdir $(PROJECT_BUILD_SCRIPTS_UPDATE_DIR); fi

build-installation-script:
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(PROJECT_BUILD_SCRIPTS_INSTALL_DIR) --out install ./$(INSTALL_SCRIPT_DIR)

build-installation-script-for-current-system:
	go build -o $(INSTALL_SCRIPT_DIR)/run $(INSTALL_SCRIPT_DIR)/main.go

build-dictionary-script:
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(PROJECT_BUILD_SCRIPTS_INSTALL_DIR) --out install ./$(DICTIONARY_SCRIPT_DIR)

build-dictionary-script-for-current-system:
	go build -o $(DICTIONARY_SCRIPT_DIR)/run $(DICTIONARY_SCRIPT_DIR)/main.go

build-update-script:
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(PROJECT_BUILD_SCRIPTS_UPDATE_DIR) --out update ./$(UPDATE_SCRIPT_DIR)

build-update-script-for-current-system:
	go build -o $(UPDATE_SCRIPT_DIR)/run $(UPDATE_SCRIPT_DIR)/main.go

build-project-archive:
	tar -czvf $(PROJECT_BUILD_DIR)/devbot.tar.gz $(PROJECT_BUILD_DIR)

prepare-release:
	make build-project-cross-platform
	cp -R $(INSTALL_SCRIPT_DIR)/database $(PROJECT_BUILD_SCRIPTS_INSTALL_DIR)/database
	cp -R $(UPDATE_SCRIPT_DIR)/migrations $(PROJECT_BUILD_SCRIPTS_UPDATE_DIR)/migrations
	cp .env.example $(PROJECT_BUILD_DIR)/.env
	cp .env.example $(PROJECT_BUILD_DIR)/.env.example
	make build-project-archive

build:
	make install
	make update
	make build-devbot-for-current-system

refresh-events:
	./scripts/project-tools/update-events.sh

create-missing-files:
	make create-if-not-exists-defined-events
	make create-if-not-exists-env

build-project-cross-platform:
	make create-missing-files
	make build-devbot-cross-platform
	make build-installation-script
	make build-update-script

cleanup:
	rm -rf vendor terraform scripts internal events documentation cmd

check-security:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...