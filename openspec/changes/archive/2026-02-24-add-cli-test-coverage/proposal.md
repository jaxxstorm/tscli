## Why

`tscli` has grown to a large command surface without consistent test coverage, which makes command behavior and output contracts easy to regress. We need a unified unit + integration testing strategy now to stabilize CLI behavior across commands, flags, config precedence, and Tailscale API interactions.

## What Changes

- Add a shared CLI test harness that can execute commands with deterministic stdout/stderr/exit-code assertions.
- Add comprehensive command tests for all current command groups under `cmd/tscli` (`get`, `list`, `create`, `set`, `delete`, `config`, `version`), including required flag validation and key success/error paths.
- Add a reusable mocked Tailscale API layer for integration-style command tests so command behavior can be validated without live API calls.
- Add baseline contract checks for request/response model decoding and encoding consistency between `tscli` code paths and Tailscale API payload shapes, sourced from the Tailscale OpenAPI schema endpoint.
- Add automated coverage-gap reporting that compares the current `tscli` command/test surface against Tailscale OpenAPI operations and highlights uncovered areas.
- Add CI-friendly test targets for fast unit tests and integration tests that run against mocks.
- Keep existing command UX and flags stable; this change introduces tests and test infrastructure, not user-facing command behavior changes.

## Capabilities

### New Capabilities

- `cli-command-test-coverage`: Define required automated test coverage for all CLI commands, including validation behavior, output behavior, exit semantics, and API-operation coverage visibility.
- `tailscale-api-mock-integration-tests`: Define a standard mocked API test architecture for command integration tests, including deterministic fixtures for success and failure paths.
- `tailscale-model-contract-consistency`: Define contract checks ensuring model and payload handling remains consistent with the Tailscale OpenAPI schema (`https://api.tailscale.com/api/v2?outputOpenapiSchema=true`) and expected JSON structures.

### Modified Capabilities

- None.

## Impact

- Affected code:
  - `cmd/tscli/**` command constructors and execution flow
  - `pkg/tscli/client.go` HTTP client integration points
  - `pkg/output/**` formatting/output pathways
  - `pkg/config/**` flag/env/config precedence handling as exercised by tests
- Dependencies/systems:
  - Go testing stack (`go test`, `httptest`, table-driven tests)
  - Possible addition of lightweight helpers for fixtures/golden outputs
  - Pinned OpenAPI schema snapshot from `https://api.tailscale.com/api/v2?outputOpenapiSchema=true`
  - Coverage gap report artifact generated from OpenAPI operations vs CLI/test mappings
- Backward compatibility:
  - No intentional CLI behavior changes; tests are used to lock in current behavior and detect regressions for existing scripts and automation.
