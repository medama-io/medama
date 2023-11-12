package generate

import "embed"

// OpenAPI Codegen
//
//go:generate ./generate.sh

// Embed OpenAPI Specification
//
//go:embed openapi.yaml
var OpenAPIDocument embed.FS
