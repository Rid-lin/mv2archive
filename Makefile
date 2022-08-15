DOCKER_ACCOUNT = vsvegner

PWD := $(shell pwd)
PROJECTNAME = $(shell basename $(PWD))
PROGRAM_NAME = $(shell basename $(PWD))

VERSION=$(shell git describe --tags)
COMMIT=$(shell git rev-parse --short HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
TAG=$(shell git describe --tags |cut -d- -f1)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

PLATFORMS=linux windows
# PLATFORMS=darwin linux windows
# ARCHITECTURES=386 amd64 ppc64 arm arm64
ARCHITECTURES=386 amd64 arm arm64

LDFLAGS = -ldflags "-s -w -extldflags '-static' -X=main.Version=${VERSION} -X=main.Build=${COMMIT} -X main.gitTag=${TAG} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH} -X main.buildTime=${BUILD_TIME}"
# LDFLAGS = -ldflags "-s -w -linkmode external -extldflags '-static' -X=main.Version=${VERSION} -X=main.Build=${COMMIT} -X main.gitTag=${TAG} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH} -X main.buildTime=${BUILD_TIME}"
# LDFLAGS = -ldflags "-s -w -X=main.Version=${VERSION} -X=main.Build=${COMMIT} -X main.gitTag=${TAG} -X main.gitCommit=${COMMIT} -X main.gitBranch=${BRANCH} -X main.buildTime=${BUILD_TIME}"

# Check for required command tools to build or stop immediately
EXECUTABLES = git go find pwd basename
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

.PHONY: help clean dep build build_all install uninstall pack release release_move release_pack

.DEFAULT_GOAL := help

help: ## Display this help screen.
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## Clean bin directory.
	rm -f ${PWD}/bin/*

dep: ## Download the dependencies.
	go mod tidy
	go mod download
	go mod vendor

build: ## Build program executable for linux platform.
	mkdir -p ${PWD}/bin
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor ${LDFLAGS} -o bin/${PROGRAM_NAME}_$(VERSION)_linux_$(COMMIT)_amd64 .
	chmod +x bin/${PROGRAM_NAME}_$(VERSION)_linux_$(COMMIT)_amd64

build_all: ## Build program executable for all platform.
	mkdir -p ./bin
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v $(LDFLAGS) -o ./bin/$(PROJECTNAME)_$(VERSION)_$(GOOS)_$(COMMIT)_$(GOARCH))))
	$(shell find ./bin/ -type f -name '*windows*' -exec mv {} {}.exe \;)


pack: ## Packing all executable files in ${PWD}/bin using UPX 
	upx ${PWD}/bin/${PROGRAM_NAME}*

install: ## Install program executable into /usr/bin directory.
	mkdir -p /usr/bin/${PROGRAM_NAME}
	install -pm 755 bin/${PROGRAM_NAME} /usr/bin/${PROGRAM_NAME}/${PROGRAM_NAME}

uninstall: ## Uninstall program executable from /usr/bin directory.
	rm -rf /usr/bin/${PROGRAM_NAME}

release: clean release_move release_pack ## Move current bin from ${PWD}/bin to ${PWD}/release and pack it

release_move:
	mkdir -p ${PWD}/release
	mv ${PWD}/bin/${PROGRAM_NAME}_$(VERSION)_linux_$(COMMIT)_amd64 ${PWD}/release/${PROGRAM_NAME}

release_pack:
	upx ${PWD}/release/${PROGRAM_NAME}