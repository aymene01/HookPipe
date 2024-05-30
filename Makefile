.PHONY: proto

build:
	@go build -o bin/blocker

run: build
	@./bin/blocker

clean:
	@rm -rf bin

test:
	@go test -v ./...
