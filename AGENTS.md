# AGENTS.md

## Overview

`devops-axi` is a Go CLI for Azure DevOps work items, wrapping the Azure CLI and its `azure-devops` extension.

## Current command

```sh
go run ./cmd/devops-axi work-item list
```

`work-item list` supports `--state`, `--assigned-to`, `--me`, `--limit`, `--fields`, and raw `--wiql`. Output is agent-friendly TOON. Azure CLI organization and project defaults are used.

## Project notes

- CLI framework: Cobra.
- Root command: `cmd/devops-axi/main.go`.
- Work-item logic: `cmd/devops-axi/work_items.go`.
- Azure CLI bridge: `internal/az/az.go`.
- Test with `/usr/local/go/bin/go test ./...` if `go` is not on `PATH`.
