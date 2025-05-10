.PHONY: generate_openapi
generate_openapi:
	@if ! command -v 'oapi-codegen' &> /dev/null; then \
		echo "Please install oapi-codegen!"; exit 1; \
	fi;

	@mkdir -p analyzer-service/internal/api/openapi/analyzer/v1
	@oapi-codegen -package v1 \
		-generate server,types \
		api/openapi/v1/analyzer-api.yaml > analyzer-service/internal/api/openapi/analyzer/v1/analyzer-api.gen.go
		
	@mkdir -p keeper-service/internal/api/openapi/keeper/v1
	@oapi-codegen -package v1 \
		-generate server,types \
		api/openapi/v1/keeper-api.yaml > keeper-service/internal/api/openapi/keeper/v1/keeper-api.gen.go