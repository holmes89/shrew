.PHONY: build test wasm clean

build:
	go build -o bin/shrew ./...

test:
	go test ./...

wasm:
	GOOS=js GOARCH=wasm go build -o bin/shrew.wasm ./cmd/wasm/

clean:
	rm -rf bin/
