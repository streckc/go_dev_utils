

GO=go
BUILD_FLAGS=-trimpath -pkgdir ./pkg


all: mod_run


mod_run: cmd/mod_run/main.go
	$(GO) build -o bin/mod_run $(BUILD_FLAGS) cmd/mod_run/main.go


clean:
	rm -rf ./bin
