## Code generation from OpenAPI specification

```sh
# for the binary install
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
oapi-codegen -package api -generate types,server -o gate/api/api.gen.go gate/api/openapi.yaml
```