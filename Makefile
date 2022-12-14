BINARY_NAME=shell-config

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
