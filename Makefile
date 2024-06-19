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
	@kubectl apply -f infra/k8s-local-test/jobs/db-migrate.yaml

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

update_test_image:
	@docker build -f integreted-local-tests/Dockerfile -t test:local .
	@kind load docker-image --name tanken-local-test test:local
	@kubectl delete job test-job
	@kubectl apply -f infra/k8s-local-test/jobs/test.yaml

USERCACHE_POD=$(shell kubectl get pod -l app=usercache -o jsonpath="{.items[0].metadata.name}")
POSTCACHE_POD=$(shell kubectl get pod -l app=postcache -o jsonpath="{.items[0].metadata.name}")
GEOCACHE_POD=$(shell kubectl get pod -l app=geocache -o jsonpath="{.items[0].metadata.name}")

POSTGRES_POD=$(shell kubectl get pod -l app=database -o jsonpath="{.items[0].metadata.name}")

flush_redis:
	@echo "Flushing Redis caches..."
	@kubectl exec $(USERCACHE_POD) -- redis-cli FLUSHALL
	@kubectl exec $(POSTCACHE_POD) -- redis-cli FLUSHALL
	@kubectl exec $(GEOCACHE_POD) -- redis-cli FLUSHALL
	@echo "Redis caches flushed."

clear_postgres:
	@echo "Clearing PostgreSQL database..."
	@kubectl exec -it $(POSTGRES_POD) -- psql -U postgres -d postgres -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "PostgreSQL database cleared."

clear_all_data: flush_redis clear_postgres update_migration_image
	@echo "All caches and databases cleared."