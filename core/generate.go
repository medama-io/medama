package generate

import (
	"embed"
	"io/fs"
)

// OpenAPI Codegen
//
//go:generate ./generate.sh

// Embed OpenAPI Specification
//
//go:embed openapi.yaml
var OpenAPIDocument embed.FS

// Embed SPA client
//
//go:embed all:client
var spaClient embed.FS

// SPA Client
func SPAClient() (fs.FS, error) {
	return fs.Sub(spaClient, "client")
}
