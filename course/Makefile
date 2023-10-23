proto:
	protoc -I./proto --go_out=./ --go-grpc_out=./ **/*.proto

run:
	APP_MODE=dev go run main.go

.PHONY: proto run