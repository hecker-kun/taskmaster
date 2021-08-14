build:
	go build -v -o bin/ client.go

run:
	go run client.go

compile:
	# Linux
	GOOS=linux GOARCH=amd64 go build -v -o bin/ client.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -v -o bin/ client.go

compile-proto:
	protoc --go_out=. proto/taskmaster.proto
	protoc --go-grpc_out=. proto/taskmaster.proto