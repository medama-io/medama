# Core

The backend API service for Medama written in Golang.

## Development

- Install the Go version specified in `go.mod`
- Ensure you have `gcc` installed as CGO is required
- Ensure you have `Taskfile` installed
- Setup gofumpt and golangci-lint to automatically format and lint your code in IDE
- Run `go mod download`

To start the API server, run:

```bash
task dev -- start # Passes start command as an argument to dev task
```

### Tests

- Run `task test`

## License

[Apache License 2.0](LICENSE)
