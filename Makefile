BIN_DIR=bin
PROJECT_BUILD_DIR=project-build
PROJECT_BUILD_SCRIPTS_DIR=$(PROJECT_BUILD_DIR)/scripts
PROJECT_BUILD_SCRIPTS_INSTALL_DIR=$(PROJECT_BUILD_DIR)/scripts/install
PROJECT_BUILD_SCRIPTS_UPDATE_DIR=$(PROJECT_BUILD_DIR)/scripts/update
CMD:=$(patsubst cmd/%/main.go,%,$(shell find cmd -type f -name 'main.go'))
LDFLAGS=-ldflags="-s -w"

DICTIONARY_SCRIPT_DIR=scripts/dictionary-loader
INSTALL_SCRIPT_DIR=scripts/install
UPDATE_SCRIPT_DIR=scripts/update

vendor:
	if [ ! -d "vendor" ] || [ -z "$(shell ls -A vendor)" ]; then go mod vendor; fi

build-slack-bot-cross-platform:
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(PROJECT_BUILD_DIR)/$(BIN_DIR)/ --out slack-bot ./cmd/slack-bot

build-slack-bot-for-current-system:
	env CGO_ENABLED=1 go build -o ./bin/current-system ./cmd/slack-bot/main.go

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

create-if-not-exists-defined-events:
	if [ ! -f events/defined-events.go ]; then cp events/defined-events.go.dist events/defined-events.go; fi

create-if-not-exists-env:
	if [ ! -f .env ]; then cp .env.example .env; fi

install:
	make create-if-not-exists-defined-events
	make create-if-not-exists-env
	./scripts/install/run

build-dictionary-script:
	go build -o $(DICTIONARY_SCRIPT_DIR)/dictionary-loader $(DICTIONARY_SCRIPT_DIR)/main.go

create-project-build-dirs:
	if [[ ! -d $(PROJECT_BUILD_DIR) ]]; then mkdir $(PROJECT_BUILD_DIR); fi
	if [[ ! -d $(PROJECT_BUILD_SCRIPTS_DIR) ]]; then mkdir $(PROJECT_BUILD_SCRIPTS_DIR); fi
	if [[ ! -d $(PROJECT_BUILD_SCRIPTS_INSTALL_DIR) ]]; then mkdir $(PROJECT_BUILD_SCRIPTS_INSTALL_DIR); fi
	if [[ ! -d $(PROJECT_BUILD_SCRIPTS_UPDATE_DIR) ]]; then mkdir $(PROJECT_BUILD_SCRIPTS_UPDATE_DIR); fi

build-installation-script:
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(PROJECT_BUILD_SCRIPTS_INSTALL_DIR) --out install ./$(INSTALL_SCRIPT_DIR)

build-installation-script-for-current-system:
	go build -o $(INSTALL_SCRIPT_DIR)/run $(INSTALL_SCRIPT_DIR)/main.go

build-update-script:
	env CGO_ENABLED=1 xgo --targets=darwin/*,linux/amd64,linux/386,windows/* --dest ./$(PROJECT_BUILD_SCRIPTS_UPDATE_DIR) --out update ./$(UPDATE_SCRIPT_DIR)

build-update-script-for-current-system:
	go build -o $(UPDATE_SCRIPT_DIR)/run $(UPDATE_SCRIPT_DIR)/main.go

prepare-release:
	make build-project-cross-platform
	cp -R $(INSTALL_SCRIPT_DIR)/database $(PROJECT_BUILD_SCRIPTS_INSTALL_DIR)/database
	cp -R $(UPDATE_SCRIPT_DIR)/migrations $(PROJECT_BUILD_SCRIPTS_UPDATE_DIR)/migrations
	cp .env.example $(PROJECT_BUILD_DIR)/.env.example
	cp devbot.sqlite $(PROJECT_BUILD_DIR)/devbot.sqlite
	tar -czvf $(PROJECT_BUILD_DIR)/devbot.tar.gz $(PROJECT_BUILD_DIR)

build-project-for-current-system:
	make create-if-not-exists-defined-events
	make create-if-not-exists-env
	make build-slack-bot-for-current-system
	make build-installation-script-for-current-system
	make build-update-script-for-current-system

build-project-cross-platform:
	make create-if-not-exists-defined-events
	make create-if-not-exists-env
	make build-slack-bot-cross-platform
	make build-installation-script
	make build-update-script

.PHONY: vendor