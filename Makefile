BINARY_NAME=shell-config

build:
	GOOS=darwin GOARCH=arm64 go build -o ${BINARY_NAME}-darwin-arm64
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY_NAME}-darwin-amd64
	GOOS=linux GOARCH=arm go build -o ${BINARY_NAME}-linux-arm
	GOOS=linux GOARCH=arm64 go build -o ${BINARY_NAME}-linux-arm64
	GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME}-linux-amd64
	GOOS=windows GOARCH=arm go build -o ${BINARY_NAME}-windows-arm
	GOOS=windows GOARCH=arm64 go build -o ${BINARY_NAME}-windows-arm64
	GOOS=windows GOARCH=amd64 go build -o ${BINARY_NAME}-windows-amd64

clean:
	go clean
	rm ${BINARY_NAME}-*
