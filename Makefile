.PHONY: dev-run install buf lint

export

install:

	@go mod tidy
	@go install github.com/bufbuild/buf/cmd/buf@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

buf:
	@env PATH="$$PATH:$$(go env GOPATH)/bin" buf generate --template proto/buf.gen.yaml proto
	@echo "✅ buf done!"

buf-win:
	@set PATH=%PATH%;%GOPATH%\bin
	@buf generate --template proto\buf.gen.yaml proto
	@echo "✅ buf done!""


run:
	go run ./cmd
	
lint:
	gofumpt -l -w .
	golangci-lint run  -v

test-run:
	go test ./...

docker-build:
	docker build -t finman-transaction-service .

docker-run:
	docker run -p 8082:8082 finman-transaction-service

docker-compose-up:
	docker-compose up --build 

docker-compose-down:
	docker-compose down --volumes