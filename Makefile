BINARY_NAME=shell-config

build:
	mkdir -p releases
	GOOS=darwin GOARCH=arm64 go build -o releases/${BINARY_NAME}-darwin-arm64
	GOOS=darwin GOARCH=amd64 go build -o releases/${BINARY_NAME}-darwin-amd64
	GOOS=linux GOARCH=arm go build -o releases/${BINARY_NAME}-linux-arm
	GOOS=linux GOARCH=arm64 go build -o releases/${BINARY_NAME}-linux-arm64
	GOOS=linux GOARCH=amd64 go build -o releases/${BINARY_NAME}-linux-amd64

clean:
	go clean
	rm -r releases
