proto:
	python3 -m grpc_tools.protoc -I./proto --python_out=./pb --pyi_out=./pb --grpc_python_out=./pb **/*.proto

run: 
	APP_MODE=dev python3 __init__.py

.PHONY: proto run