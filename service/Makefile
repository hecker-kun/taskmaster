build:
	go build -v -o bin/ server.go

run:
	go run server.go

compile:
	# Linux
	GOOS=linux GOARCH=amd64 go build -v -o bin/ server.go
	# Windows
	GOOS=windows GOARCH=amd64 go build -v -o bin/ server.go

compile-proto:
	protoc --go_out=. proto/taskmaster.proto
	protoc --go-grpc_out=. proto/taskmaster.proto