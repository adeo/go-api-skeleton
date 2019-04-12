# Multi-OS compliant flags
IS_LINUX=$(shell sed --version > /dev/null 2> /dev/null && echo $$?)
ifeq ($(IS_LINUX),0)
    SED_CMD="sed"
else
    SED_CMD="gsed"
endif


NAME := $(shell cat tom.yml | $(SED_CMD) -r '/name:/!d;s/.*: *'"'"'?([^$$])'"'"'?/\1/')
VERSION := $(shell cat tom.yml | $(SED_CMD) -r '/app:/!d;s/.*: *'"'"'?([^$$])'"'"'?/\1/')
GITHASH := $(shell git rev-parse --short HEAD)
BUILDDATE := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

.PHONY: help
help: ## display this help
	@echo "Name: $(NAME)"
	@echo "Version: $(VERSION)"
	@echo "Git hash: $(GITHASH)"
	@echo "Build date: $(BUILDDATE)"
	@echo
	@echo "This is the list of available make targets:"
	@echo " $(shell cat Makefile | $(SED_CMD) -r '/^[a-zA-Z-]+:.*##.*/!d;s/## *//;s/$$/\\n/')"

.PHONY: start
start: ## start the application
	go run main.go --config config/local.json

.PHONY: start-offline
start-offline: ## start the application in offline mode
	go run main.go --log-level debug --log-format text --db-in-memory --db-in-memory-import-file ./testdata/dataset.json

.PHONY: deps
deps: ## get the golang dependencies in the vendor folder
	GO111MODULE=on go mod vendor

.PHONY: build
build: ##  build the executable and set the version
	go build -o turbine-monitoring -ldflags "-X github.com/adeo/turbine-monitoring/handlers.ApplicationVersion=$(VERSION) -X github.com/adeo/turbine-monitoring/handlers.ApplicationName=$(NAME) -X github.com/adeo/turbine-monitoring/handlers.ApplicationGitHash=$(GITHASH) -X github.com/adeo/turbine-monitoring/handlers.ApplicationBuildDate=$(BUILDDATE)" main.go

.PHONY: test
test: ## run go test
	go test -v ./...

.PHONY: bump
bump: ## bump the version in the info.yaml, tom.yml file
	NEW_VERSION=`standard-version --dry-run | $(SED_CMD) -r '/tagging release/!d;s/.*tagging release *v?(.*)/\1/g'`; \
		$(SED_CMD) -r -i 's/^(.*version: *).*$$/\1'$$NEW_VERSION'/' info.yaml; \
		$(SED_CMD) -r -i 's/^(.*app: *).*$$/\1'$$NEW_VERSION'/' tom.yml

.PHONY: release
release: bump ## bump the version in the info.yaml, tom.yml, and make a release (commit, tag and push)
	git add info.yaml tom.yml
	standard-version --message "chore(release): %s [ci skip]" --commit-all
	git push --follow-tags origin HEAD

.PHONY: openapi
openapi: ## install openapi-parser and generate openapi schema
	go get github.com/alexjomin/openapi-parser
	openapi-parser --output openapi.yaml
	openapi-parser merge --output openapi.yaml --main info.yaml --dir .
	go generate
