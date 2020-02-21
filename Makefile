# Chatbot Makefile.
#
# Arguments:
#	OS	: Platform for binary build, linux Or darwin (OSX)
#

# Bump these on release.
VERSION_MAJOR ?= 0
VERSION_MINOR ?= 0
VERSION_BUILD ?= 1
VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)

# Default Go binary.
ifndef GOROOT
  GOROOT = /usr/local/go
endif

# Determine the OS to build.
ifeq ($(OS),)
  ifeq ($(shell  uname -s), Darwin)
    GOOS = darwin
  else
    GOOS = linux
  endif
else
  GOOS = $(OS)
endif

GOCMD = GOOS=$(GOOS) go
GOBUILD = CGO_ENABLED=0 $(GOCMD) build
GOTEST = $(GOCMD) test -race
RM = rm -rf
PROJECT = chatbot
DIST_DIR = ./dist
BUILD_PACKAGE = ./cmd/chatbot
GO_PKGS?=$$(go list ./...)

VERSION_PACKAGE = github.com/anthonycorbacho/chatbot/internal/version
GO_LDFLAGS :="
GO_LDFLAGS += -X $(VERSION_PACKAGE).version=$(VERSION)
GO_LDFLAGS += -X $(VERSION_PACKAGE).buildDate=$(shell date +'%Y-%m-%dT%H:%M:%SZ')
GO_LDFLAGS += -X $(VERSION_PACKAGE).gitCommit=$(shell git rev-parse HEAD)
GO_LDFLAGS += -X $(VERSION_PACKAGE).gitTreeState=$(if $(shell git status --porcelain),dirty,clean)
GO_LDFLAGS +="

.PHONY: build

build:			## Builds the code
		mkdir -p $(DIST_DIR)
		$(GOBUILD) -ldflags $(GO_LDFLAGS) -i -o $(DIST_DIR)/$(PROJECT)-$(VERSION)-$(GOOS) -v $(BUILD_PACKAGE)

build-docker:		## Builds the code in docker
		docker build \
			-t $(PROJECT):$(VERSION) \
			--build-arg BUILD_DATE=`date +%Y-%m-%dT%H:%M:%SZ` \
			--build-arg VCS_REF=`git rev-parse --short HEAD` \
			.

tools:			## Install developer tools
		curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.19.1

# See golangci.yml in extra folder for linters setup
linter: 		## Run the linter
		./bin/golangci-lint run -c golangci.yml ./...

test: 			## Test the code
		$(GOTEST) -v $(GO_PKGS)

integration-test: 	## Run Integration test the code
		$(GOTEST) -count=1 -v -tags integration $(GO_PKGS)

bench: 			## Benchmarck the code
		$(GOCMD) test -bench=. ./... -benchmem

clean: 			## Clean project by removing temp file and binary
		find . -type f -name '*~' -exec rm {} +
		find . -type f -name '\#*\#' -exec rm {} +
		find . -type f -name '*.coverprofile' -exec rm {} +
		$(RM) checkstyle.xml
		$(RM) $(DIST_DIR)/*

fclean: clean
		$(RM) $(DIST_DIR)

version: 		## Show de current version
		@echo $(VERSION)

help:           	## Show this help.
		@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'