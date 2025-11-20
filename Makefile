BINARY_NAME=dtac-aether-webmodule
PKG=./cmd/dtac-aether-webmodule

.PHONY: all build clean run

all: build

build:
	go build -o bin/$(BINARY_NAME) $(PKG)

run:
	go run $(PKG)

clean:
	rm -rf bin
