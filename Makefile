update_all_proto: update_data_fetcher_proto update_frontend_proto
	@echo "All proto files generated successfully."

update_data_fetcher_proto:
	@echo "Generating data-fetcher proto files..."
	@cd backend/data-fetcher && make update_proto
	@echo "Data-fetcher proto files generated successfully."

update_frontend_proto:
	@echo "Generating frontend proto files..."
	@cd frontend && make update_proto
	@echo "Frontend proto files generated successfully."

update_test_proto:
	@echo "Generating test proto files..."
	@cd backend/test && make update_proto
	@echo "Test proto files generated successfully."

update_local_cluster:
	@echo "Updating local cluster..."
	@cd infra/yaml
	@kubectl apply -f .

update_local_backend_image:
	@echo "Updating local image..."
	@cd backend && docker build -f data-fetcher/Dockerfile -t data-fetcher:local .
	@kind load docker-image data-fetcher:local
	@kubectl delete pod -l app=data-fetcher

update_local_frontend_image:
	@cd frontend && docker build -f Dockerfile -t frontend:local .
	@kind load docker-image frontend:local
	@kubectl delete pod -l app=frontend

all: update_all_proto update_local_cluster update_local_backend_image update_local_frontend_image
	@echo "All updated successfully."

unit_test_all:
	@cd backend/common && make test