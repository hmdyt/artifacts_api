.PHONY: gen
gen:
	@echo "Generating code..."
	@docker run --rm -v "${PWD}:/local" openapitools/openapi-generator-cli generate \
         -i /local/openapi.json \
         -g go \
         -o /local/gen
