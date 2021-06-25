ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
OUT_DIR:=bin

BINARIES=$(shell ls cmd)
VERSION=0.1.0
BUILD=`git rev-parse HEAD`
PLATFORMS=darwin linux windows
#ARCHITECTURES=386 amd64
ARCHITECTURES=amd64

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

all: clean build_all install

build:
	$(foreach BINARY, $(BINARIES),\
	go build ${LDFLAGS} -o $(OUT_DIR)/$(BINARY) cmd/$(BINARY)/main.go)

build_all:
	$(foreach BINARY, $(BINARIES),\
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES),\
	$(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build ${LDFLAGS} -o $(OUT_DIR)/$(BINARY)-$(GOOS)-$(GOARCH) cmd/$(BINARY)/main.go))))

install:
	$(foreach BINARY, $(BINARIES),\
	go build ${LDFLAGS} -o ~/bin/$(BINARY) cmd/$(BINARY)/main.go)

# Remove only what we've created
clean:
	find ${ROOT_DIR} -name '${BINARY}[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' -delete

.PHONY: check clean install build_all all
