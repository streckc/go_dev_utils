ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
OUT_DIR:=build
BIN_DIR:=~/bin

TARGETS=$(shell ls cmd)
PLATFORMS=darwin linux windows
#ARCHITECTURES=386 amd64
ARCHITECTURES=amd64

#PACKAGES=$(shell find pkg -type f)

COMMIT=$(shell git rev-parse HEAD)
LDFLAGS=-ldflags "-X main.Build=$(COMMIT) -s -w"

default: build

build: $(foreach TARGET, $(TARGETS), $(OUT_DIR)/$(TARGET))

test:
	go test -covermode=atomic -coverprofile=coverage.out ./...

all: linux windows darwin

linux: $(foreach TARGET, $(TARGETS),\
        $(foreach ARCH, $(ARCHITECTURES),\
          $(OUT_DIR)/linux-$(ARCH)/$(TARGET)))

windows: $(foreach TARGET, $(TARGETS),\
        $(foreach ARCH, $(ARCHITECTURES),\
          $(OUT_DIR)/windows-$(ARCH)/$(TARGET)))

darwin: $(foreach TARGET, $(TARGETS),\
        $(foreach ARCH, $(ARCHITECTURES),\
          $(OUT_DIR)/darwin-$(ARCH)/$(TARGET)))

# Installs into home directory bin folder
install: build
	$(foreach TARGET, $(TARGETS), cp $(OUT_DIR)/$(TARGET) $(BIN_DIR);)

# Remove only what we've created
clean:
	$(foreach TARGET, $(TARGETS),\
		find ${ROOT_DIR} \( -name '$(TARGET)' -o -name '$(TARGET)[-?][a-zA-Z0-9]*[-?][a-zA-Z0-9]*' \) -delete;)
	rm -r $(OUT_DIR)
	rm -f coverage.out audit-utils-coverage*.html

### Dependencies

$(OUT_DIR)/%: cmd/%/main.go cmd/%/* $(PACKAGES)
	$(call build_command,$@,$<)

$(OUT_DIR)/linux-amd64/%: cmd/%/main.go cmd/%/* $(PACKAGES)
	$(call build_command,$@,$<,linux,amd64)

$(OUT_DIR)/windows-amd64/%: cmd/%/main.go cmd/%/* $(PACKAGES)
	$(call build_command,$@,$<,windows,amd64)
	mv $@ $@.exe

$(OUT_DIR)/darwin-amd64/%: cmd/%/main.go cmd/%/* $(PACKAGES)
	$(call build_command,$@,$<,darwin,amd64)

build_command = \
	$(eval DEST=$(1)) $(eval SRC=$(2)) $(eval GOOS=$(3)) $(eval GOARCH=$(4)) \
	$(eval BINARY=$(shell basename $(DEST))) \
	$(eval TAG=$(shell git tag --list | grep "$(BINARY)-v" | tail -1)) \
	$(eval VERSION=$(or $(shell echo $(TAG) | sed -e 's/^.*-v//'), 0.0.0)) \
	$(eval DATE=$(shell git log -1 --format=%aI $(TAG))) \
	GOOS=$(GOOS) GOARCH=$(GOARCH) \
	go build \
		-ldflags "-X main.Build=$(COMMIT) -X main.Version=$(VERSION) -X main.Date=$(DATE) -s -w" \
		-o $(DEST) $(SRC)

.PHONY: build all test windows linux darwin install clean 
