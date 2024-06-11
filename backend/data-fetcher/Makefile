GOPATH=$(shell go env GOPATH)

PATH := $(GOPATH)/bin:$(PATH)

update_proto: clean_proto_buf generate_proto_buf

install_protoc_plugins:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "gRPC plugins installed successfully."

generate_proto_grpc:
	@protoc --proto_path=../../proto --go_out=./rpc --go-grpc_out=./rpc ../../proto/*.proto
	@echo "Proto files compiled successfully."

generate_proto_buf:
	@buf generate ../../proto
	@echo "Proto files compiled using buf to Go files successfully."

clean_proto_grpc:
	@rm -rf rpc/general rpc/post
	@echo "Cleaned up generated files."
	@mkdir rpc/general rpc/post

clean_proto_buf:
	@rm -rf rpc/connectrpc rpc/pb
	@echo "Cleaned up generated files."

help:
	@echo "Available commands:"
	@echo "  make install_protoc_plugins - installs the required gRPC plugins for Go."
	@echo "  make generate_proto - compiles .proto files in the /rpc folder to Go files."
