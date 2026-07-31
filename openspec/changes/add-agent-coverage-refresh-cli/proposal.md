## Why

Maintainers currently have working coverage-gap and OpenAPI refresh tooling, but using it requires remembering repository-specific commands and paths. Agent workflows should be able to discover and run the same checks automatically, and `tscli` should expose script-friendly maintenance commands that follow the existing CLI contract.

## What Changes

- Add repository skills under `.codex/skills` and `.opencode/skills` that guide agents to run coverage-gap analysis and OpenAPI refresh workflows before implementing command parity work.
- Add `tscli` CLI maintenance commands for coverage gaps and OpenAPI refresh so developers can run the workflows through the binary as well as through existing make targets.
- Keep existing `make coverage-gaps`, coverage reports, OpenAPI fixtures, and generated artifacts as the source of truth.
- Preserve existing end-user API command behavior; the new commands are additive and scoped to repository maintenance.
- Do not introduce new config files, environment variable requirements, or breaking flag changes for existing scripts.

## Capabilities

### New Capabilities
- `agent-maintenance-workflows`: Defines repo-local agent skills that automatically run coverage-gap and OpenAPI refresh workflows before related implementation work.
- `cli-maintenance-commands`: Defines `tscli` maintenance commands for running coverage-gap and OpenAPI refresh tooling in a script-friendly way.

### Modified Capabilities
- `coverage-gap-elimination`: Require coverage-gap workflows to be available through agent skills and CLI commands while preserving existing generated artifact behavior.
- `tailscale-model-contract-consistency`: Require OpenAPI refresh workflows to be available through agent skills and CLI commands while preserving snapshot metadata behavior.

## Impact

- Affected command groups:
  - New `tscli maintenance coverage-gaps`
  - New `tscli maintenance openapi-refresh`
- Affected flags:
  - New command-local flags for selecting check/update behavior where supported by the underlying tooling.
- Affected config/env keys:
  - None.
- Affected code and docs:
  - `.codex/skills/**`
  - `.opencode/skills/**`
  - `cmd/tscli/**`
  - `internal/cli/root.go`
  - `coverage/coveragegaps/**`
  - `coverage/openapirefresh/**`
  - `pkg/contract/openapi/**`
  - CLI integration tests and generated command docs where applicable
- Backward compatibility:
  - Existing make targets and scripts continue to work.
  - Existing `tscli` command groups, flags, output formats, config files, and environment variables are unchanged.
