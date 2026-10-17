# Hello Azure Go

A tiny Go HTTP service for learning Azure deployment and GitHub delivery workflows.
Uses only the Go standard library; no external dependencies or database.

## Requirements

Go 1.25 or newer. `.tool-versions` selects the existing asdf Go 1.25.7 installation.
The module path is a local identifier; it does not require a website or DNS setup.

## Run

```bash
go run .
```

In another terminal:

```bash
curl -i http://localhost:8080/health
```

The endpoint returns HTTP 200 with `{"status":"ok"}`. It also supports HEAD.
Other paths return 404; unsupported methods on `/health` return 405.
This is a liveness endpoint: it confirms the HTTP service can respond.

Set `PORT` to change the default port of 8080:

```bash
PORT=8081 go run .
```

The server listens on all network interfaces, sets HTTP timeouts, and allows up to
10 seconds for graceful shutdown on SIGINT or SIGTERM.

## Check

```bash
go test -race ./...
go vet ./...
go test -run '^$' -bench . -benchmem ./...
```

Tests exercise the real HTTP server, response format, methods, and route boundaries.
The benchmark measures the handler plus an in-memory response recorder, not network
latency or end-to-end deployment performance.

## Build

```bash
go build -o bin/server .
./bin/server
```

## Docker

```bash
docker build -t hello-azure-go:local .
docker run -d --name hello-azure-go -p 127.0.0.1:8080:8080 hello-azure-go:local
curl -i http://localhost:8080/health
```

The multi-stage build runs tests and compiles a static binary using Go 1.25.7.
The final image contains only that binary, runs as a non-root user, and has no shell.
The published port is bound to localhost for local development.

```bash
docker logs hello-azure-go
docker stop --time 15 hello-azure-go
docker start hello-azure-go
```

## Delivery

GitHub Actions validates formatting, tests, vetting, compilation, and Docker builds.

The CD workflow publishes immutable commit-tagged images to Azure Container Registry and updates the selected Azure Container App through GitHub OIDC.

The current deployment environment is:

| Environment | Azure target |
|-------------|----------------------------------------------|
| `dev`       | `hello-azure-go-dev` in `tg-hello-azure-dev` |

See the [GitHub deployment guide](.github/README.md) for environment setup, variables, OIDC configuration, RBAC, and deployment behavior.

The current development endpoint is: [https://hello-azure-go-dev.lemonwater-bcb51af9.canadaeast.azurecontainerapps.io/health](https://hello-azure-go-dev.lemonwater-bcb51af9.canadaeast.azurecontainerapps.io/health)

Then replace the existing “Next steps” section with:

## Next steps

Implement Terraform plan & apply within the GitHub Action workflows.