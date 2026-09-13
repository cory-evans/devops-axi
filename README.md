# devops-axi

A focused Azure DevOps CLI for work items, built as a scoped wrapper around the Azure CLI.

## Status

The initial scaffold provides root help only. The private Azure CLI runner is ready for future work-item commands but is not user-facing yet.

## Requirements

- Go 1.27+
- Azure CLI (`az`) for future commands

## Run

```sh
go run ./cmd/devops-axi --help
```

## Test

```sh
go test ./...
```
