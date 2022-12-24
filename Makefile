BINARY_NAME=shell-config

# If the first argument is `test` pass the rest of the line to the target
ifeq (test,$(firstword $(MAKECMDGOALS)))
  # Use the rest as arguments for `test`
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  # ...and turn them into do-nothing targets
  $(eval $(RUN_ARGS):;@:)
endif

.PHONY: test

target: build

build:
	mkdir -p releases
	GOOS=darwin GOARCH=arm64 go build -o releases/${BINARY_NAME}-Darwin-arm64
	GOOS=darwin GOARCH=amd64 go build -o releases/${BINARY_NAME}-Darwin-x86_64
	GOOS=linux GOARCH=arm go build -o releases/${BINARY_NAME}-Linux-arm
	GOOS=linux GOARCH=arm64 go build -o releases/${BINARY_NAME}-Linux-aarch64
	GOOS=linux GOARCH=amd64 go build -o releases/${BINARY_NAME}-Linux-x86_64

clean:
	go clean
	rm -r releases

test:
	docker run --rm -it -v $(shell pwd):$(shell pwd) -w $(shell pwd) golang:latest go run main.go $(RUN_ARGS)