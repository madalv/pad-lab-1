proto:
	protoc -I./proto --go_out=./ --go-grpc_out=./ **/*.proto

run:
	go run main.go

.PHONY: proto run