lint:
	golangci-lint run -v

test:
	go test -v ./...

fmt:
	go fmt ./... && gofumpt -w .

build:
	go build -o bin/chatgpt-cli .

add_to_path:
	sudo mv bin/chatgpt-cli /usr/local/bin/chatgpt-cli


