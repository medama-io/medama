# Core

The backend API service for Medama written in Golang.

## Development

- Use the Go version specified in the root `mise.toml`
- Ensure you have `gcc` installed as CGO is required
- Install repository tools with `mise install` from the repository root
- Setup gofumpt and golangci-lint to automatically format and lint your code in IDE
- Run `go mod download`

From `core`, start the API server with:

```bash
mise run dev
```

### Tests

- From `core`, run `mise run test`

## License

[Apache License 2.0](LICENSE)
