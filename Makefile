update_all_proto:
	@cd backend && make update_all_proto
	@cd frontend && make update_proto
	@echo "All proto files generated successfully."

unit_test_all:
	@cd backend/common && make test