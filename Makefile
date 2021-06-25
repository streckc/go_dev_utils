

GO=go
BUILD_FLAGS=-trimpath -pkgdir ./pkg


all: test_mod


test_mod: cmd/test_mod/main.go
	$(GO) build -o bin/test_mod $(BUILD_FLAGS) cmd/test_mod/main.go


clean:
	rm -rf ./bin
