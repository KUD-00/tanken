update_proto: clean generate_proto

generate_proto:
	@npx buf generate ../proto/
	@echo "Proto files compiled using buf to ts file successfully."

clean:
	@rm -rf rpc
	@echo "Cleaned up generated files."
	@mkdir rpc