# Smylet
## For best results, run these make targets inside the devcontainer


DATABASE_URI=postgres://postgres:postgres@localhost:5432/smylet?sslmode=disable&search_path=public
#
# Project-specific variables
#
# App name.
# include ./resources/env/app.env
# export

# include ./resources/env/app_test.env
# export

APP=symlet-backend
ifeq ($(shell go env GOOS),windows)
  APP:=$(APP).exe
endif
# Version.
VERSION?=$(shell git describe --tags --always --dirty --match='v*' 2> /dev/null | sed 's/^v//')
# Enable Go Modules.
GO111MODULE=on
# Go ldflags.
# Set version to git tag if available, otherwise use commit hash.
# Strip debug symbols and disable DWARF generation.

# Go build tags.
GO_BUILDTAGS=$(shell cat .go-build-tags 2> /dev/null)
# Archive information.
# Use zip on Windows, tar.gz on Linux and macOS.
# Use GNU tar on macOS if available, to avoid issues with BSD tar.
ifeq ($(shell go env GOOS),windows)
  ARCHIVE_EXT=zip
  ARCHIVE_CMD=zip -r
else
  ARCHIVE_EXT=tar.gz
  ARCHIVE_CMD=tar -czf
  ifeq ($(shell which gtar >/dev/null 2>/dev/null; echo $$?),0)
    ARCHIVE_CMD:=g$(ARCHIVE_CMD)
  endif
endif
ARCHIVE_NAME=dist/smylet-backend_$(shell go env GOOS | sed s/darwin/macos/)_$(shell go env GOARCH | sed s/amd64/x86_64/).$(ARCHIVE_EXT)
ARCHIVE_FILES=$(APP) LICENSE README.md
# Docker compose file.
COMPOSE_FILE=tests/integration/docker-compose.yml
# Docker compose project name.
COMPOSE_PROJECT_NAME=$(APP)-integration-tests


#
# Default target (help)
#
.PHONY: help
help: ## display this help
	@echo "Please use \`make <target>' where <target> is one of:"
	@echo
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-24s\033[0m - %s\n", $$1, $$2}'
	@echo

#
# Linter targets.
#
lint: ## run set of linters over the code.
	@golangci-lint run -v --build-tags $(GO_BUILDTAGS)

#
# Go targets.
#
.PHONY: go-get
go-get: ## get go modules.
	@echo '>>> Getting go modules.'
	@go mod download

.PHONY: go-build
go-build: ## build app binary.
	@echo '>>> Building go binary.'
	@CGO_ENABLED=0 go build -ldflags="$(GO_LDFLAGS)" -tags="$(GO_BUILDTAGS)" -o $(APP) ./main.go

.PHONY: go-format
go-format: ## format go code.
	@echo '>>> Frmatting go code.'
	@gofumpt -l .
	@gofumpt -w .
	@goimports -w -local github.com/Symlet/smylet-backend .

.PHONY: go-dist
go-dist: go-build ## archive app binary.
	@echo '>>> Archiving go binary.'
	@dir=$$(dirname $(ARCHIVE_NAME)); if [ ! -d $$dir ]; then mkdir -p $$dir; fi
	@if [ -f $(ARCHIVE_NAME) ]; then rm -f $(ARCHIVE_NAME); fi
	@$(ARCHIVE_CMD) $(ARCHIVE_NAME) $(ARCHIVE_FILES)


# Migration targets.

.PHONY: install-atlas
install-atlas: ## Install atlas CLI tool
	@which atlas > /dev/null 2>&1; \
	if [ $$? -ne 0 ]; then \
		echo ">>> Installing atlas..."; \
		curl -sSf https://atlasgo.sh | sh; \
	else \
		echo ">>> Atlas is already installed. Skipping installation."; \
	fi

# Define a variable for the migration target
MIGRATION_TARGET := setup-user-auth

.PHONY: create-migrate
create-migrate:  ## create migration files.
	@echo ">>> Creating migration files."
	@atlas migrate diff --env gorm $(MIGRATION_TARGET)



# You can override the variable when calling make like this:
# make create-migrate MIGRATION_TARGET=new-target

.PHONY: migrate
migrate:  ## run migrations.
	@echo ">>> Running migrations."
	@atlas migrate apply --dir file://./migrations --url "$(DATABASE_URI)" --env gorm


.PHONY: lint-migrate
lint-migrate:
	@atlas migrate lint --dev-url "$(DATABASE_URI)" --base=20231009112817
# Tests targets.
#

.PHONY: test-go-unit
test-go-unit: ## run go unit tests.
	@echo ">>> Running unit tests."
	go test -v ./...

.PHONY: test-go-integration
test-go-integration:  ## run go integration tests.
	@echo ">>> Running integration tests."
	docker-compose -f ./tests/integration/docker-compose.yml  up --build

#
# Service test targ
#
.PHONY: service-start
service-start: ## start service in container.
	@echo ">>> Starting up service container."
	@COMPOSE_FILE=$(COMPOSE_FILE) COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
		docker-compose up -d service

.PHONY: service-stop
service-stop: ## stop service in container.
	@echo ">>> Stopping service container."
	@COMPOSE_FILE=$(COMPOSE_FILE) COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
		docker-compose stop service

.PHONY: service-restart
service-restart: service-stop service-start ## restart service in container.

# .PHONY: service-test
# service-test: service-restart ## run integration tests in container.
# 	@echo ">>> Running integration tests in container."
# 	@COMPOSE_FILE=$(COMPOSE_FILE) COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
# 	    docker-compose run integration-tests

.PHONY: service-clean
service-clean: ## clean containers.
	@echo ">>> Cleaning containers."
	@COMPOSE_FILE=$(COMPOSE_FILE) COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) \
		docker-compose down -v --remove-orphans

#
# Mockery targets.
#
.PHONY: mocks-clean
mocks-clean: ## cleans old mocks.
	find . -name "mock_*.go" -type f -print0 | xargs -0 /bin/rm -f

.PHONY: mocks-generate
mocks-generate: mocks-clean ## generate mock based on all project interfaces.
	mockery --all --dir "./pkg/api/mlflow" --inpackage --case underscore

#
# Build targets
# 
PHONY: clean
clean: ## clean the go build artifacts
	@echo ">>> Cleaning go build artifacts."
	rm -Rf $(APP)

PHONY: build
build: go-build ## build the go components



PHONY: format
format: go-format python-format ## format the code

PHONY: run
run:  ## run the Symlet app
	@echo ">>> Running the Smylet app."
	./$(APP)

