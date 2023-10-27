package generate

import "embed"

// OpenAPI Codegen
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --target api --clean openapi.yaml --convenient-errors off

// Embed OpenAPI Specification
//
//go:embed openapi.yaml
var OpenAPIDocument embed.FS
