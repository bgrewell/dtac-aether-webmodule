BINARY_NAME={{module}}
PKG=./cmd/{{module}}

.PHONY: all build clean run

all: build

build:
	go build -o bin/$(BINARY_NAME) $(PKG)

run:
	go run $(PKG)

clean:
	rm -rf bin
