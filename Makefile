update_all_proto:
	@cd backend && make update_all_proto
	@cd frontend && make update_proto
	@echo "All proto files generated successfully."

unit_test_all:
	@cd backend/common && make test



initialize_k8s: initialize_local_new_k8s_cluster update_images initialize_k8s_infras

update_images: update_backend_image update_frontend_image update_migration_image

build_and_load_images: build_and_load_backend_image build_and_load_frontend_image build_and_load_migration_image

build_and_load_backend_image:
	@cd backend && docker build -f data-fetcher/Dockerfile -t data-fetcher:local .
	@kind load docker-image --name tanken-local-test data-fetcher:local

build_and_load_migration_image:
	@cd infra && docker build -f k8s-local-test/database/Dockerfile -t database-migrations:local .
	@kind load docker-image --name tanken-local-test database-migrations:local

build_and_load_frontend_image:
	@cd frontend && docker build -f Dockerfile -t frontend:local .
	@kind load docker-image --name tanken-local-test frontend:local

update_backend_image:
	@cd backend && docker build -f data-fetcher/Dockerfile -t data-fetcher:local .
	@kind load docker-image --name tanken-local-test data-fetcher:local
	@kubectl delete pod -l app=data-fetcher

update_frontend_image: 
	@cd frontend && docker build -f Dockerfile -t frontend:local .
	@kind load docker-image --name tanken-local-test frontend:local
	@kubectl delete pod -l app=frontend

update_migration_image:
	@cd infra && docker build -f k8s-local-test/database/Dockerfile -t database-migrations:local .
	@kind load docker-image --name tanken-local-test database-migrations:local
	@kubectl delete job migrate-job
	@kubectl apply -f infra/k8s-local-test/database/migrate-job.yaml

initialize_local_new_k8s_cluster:
	@kind create cluster --name tanken-local-test
	@kubectl config use-context kind-tanken-local-test

initialize_k8s_infras:
	@kubectl apply -f infra/k8s-local-test/cache/geocache.yaml
	@kubectl apply -f infra/k8s-local-test/cache/postcache.yaml
	@kubectl apply -f infra/k8s-local-test/cache/usercache.yaml
	@kubectl apply -f infra/k8s-local-test/database/database.yaml
	@kubectl apply -f infra/k8s-local-test/jobs/db-migrate.yaml
	@kubectl apply -f infra/k8s-local-test/jobs/test.yaml

initialize_k8s_services:
	@kubectl apply -f infra/k8s-local-test/services/data-fetcher.yaml
	@kubectl apply -f infra/k8s-local-test/services/test.yaml

update_test_image:
	@docker build -f test/Dockerfile -t test:local .
	@kind load docker-image --name tanken-local-test test:local
	@kubectl delete pod -l app=test