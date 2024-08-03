.PHONY: gen
gen:
	@echo "Generating code..."
	@go generate ./...
	@echo "Code generation complete."
