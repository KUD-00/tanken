GOPATH=$(shell go env GOPATH)

PATH := $(GOPATH)/bin:$(PATH)

update_proto: clean_proto_buf generate_proto_buf

generate_proto_buf:
	@buf generate ../../proto
	@echo "Proto files compiled using buf to Go files successfully."

clean_proto_buf:
	@rm -rf rpc/connectrpc rpc/pb
	@echo "Cleaned up generated files."

test_local_new_cluster:
	@kind create cluster --name backend-test
	@kubectl config use-context kind-backend-test

test_local_initialize_databases:
	@kubectl apply -f ../../infra/yaml/database.yaml
	@kubectl apply -f ../../infra/yaml/geocache.yaml
	@kubectl apply -f ../../infra/yaml/postcache.yaml
	@kubectl apply -f ../../infra/yaml/job.yaml