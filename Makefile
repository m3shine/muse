#!/usr/bin/make -f

PACKAGES=$(shell go list ./... | grep -v '/simulation')
VERSION := $(shell echo $(shell git describe --always) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')

export GO111MODULE = on
export COSMOS_SDK_TEST_KEYRING = y

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=muse \
	-X github.com/cosmos/cosmos-sdk/version.ServerName=mused \
	-X github.com/cosmos/cosmos-sdk/version.ClientName=musecli \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
    -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"


ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: install lint test

build: go.sum
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/mused.exe ./cmd/mused && \
	go build -mod=readonly $(BUILD_FLAGS) -o build/musecli.exe ./cmd/musecli
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/mused ./cmd/mused && \
	go build -mod=readonly $(BUILD_FLAGS) -o build/musecli ./cmd/musecli
endif

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/mused ./cmd/mused && \
    LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 go build -mod=readonly $(BUILD_FLAGS) -o build/musecli ./cmd/musecli

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/mused && \
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/musecli


########################################
### Tools & dependencies

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

test:
    @go test -mod=readonly $(PACKAGES)

########################################
### Local validator nodes using docker and docker-compose

# Run a 4-node testnet locally
localnet-start: build-linux localnet-stop
	@if ! [ -f build/node0/mused/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/output:Z tygeth/muse:latest mused testnet --v 4 -o /output --starting-ip-address 192.168.10.2 --chain-id musechain ; fi
	docker-compose up

# Stop testnet
localnet-stop:
	docker-compose down