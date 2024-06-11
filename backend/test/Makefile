GOPATH=$(shell go env GOPATH)

PATH := $(GOPATH)/bin:$(PATH)

update_proto: clean_proto_buf generate_proto_buf

generate_proto_buf:
	@buf generate ../../proto
	@echo "Proto files compiled using buf to Go files successfully."

clean_proto_buf:
	@rm -rf rpc/connectrpc rpc/pb
	@echo "Cleaned up generated files."