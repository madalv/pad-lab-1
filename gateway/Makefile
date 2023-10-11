proto: 
	protoc -I ./lib/course --elixir_out=plugins=grpc:./lib/course course_svc.proto ; \
	protoc -I ./lib/rec --elixir_out=plugins=grpc:./lib/rec rec_svc.proto

run: 
	mix run --no-halt

.PHONY: proto run