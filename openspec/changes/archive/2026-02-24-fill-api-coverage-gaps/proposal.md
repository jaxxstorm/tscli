## Why

`tscli` currently has significant OpenAPI surface gaps and command mapping drift, so users cannot manage the full Tailscale API from the CLI consistently. We need to close command coverage gaps and enforce one command taxonomy so behavior is predictable and discoverable.

## What Changes

- Examine the OpenAPI coverage-gap report and command map to identify uncovered API operations and unmapped command surfaces.
- Add CLI commands for all missing in-scope Tailscale API operations so the uncovered operation count is driven to zero.
- Enforce command taxonomy for canonical command verbs:
  - creation actions use `create`
  - update/mutation actions use `set`
  - single-resource retrieval uses `get`
  - multi-resource retrieval uses `list`
  - deletion continues to use `delete`
- Normalize existing command naming where needed to align canonical verbs and noun structure while preserving backward compatibility through aliases for existing scripts.
- Update command-to-OpenAPI mapping and coverage-gap tooling so CI can detect regressions in parity and taxonomy alignment.
- Expand tests for newly added commands and updated command paths using the existing mock API integration harness.

## Capabilities

### New Capabilities

- `openapi-command-parity`: Ensure every in-scope Tailscale OpenAPI operation has a corresponding CLI command path and behavior contract.
- `cli-command-taxonomy-alignment`: Define and enforce canonical command verb rules (`create`, `set`, `get`, `list`, `delete`) across the command surface.
- `coverage-gap-elimination`: Define parity verification and regression checks so new OpenAPI gaps or unmapped command surfaces are blocked in CI.

### Modified Capabilities

- `cli-command-test-coverage`: Extend coverage requirements to enforce parity completion and command taxonomy verification across all leaf commands.

## Impact

- Affected code:
  - `cmd/tscli/**` command tree and command aliases
  - `pkg/contract/openapi/command-operation-map.yaml`
  - `coverage/coveragegaps/**` reporting and regression checks
  - `cmd/tscli/*_test.go` and shared test harness coverage
- Affected API areas:
  - devices, invites, DNS, logging, keys, policy, posture integrations, users, contacts, webhooks, services, tailnet settings, and AWS external-id validation flows
- Backward compatibility:
  - Canonical commands will follow strict taxonomy, while existing command aliases should be preserved where needed to minimize script breakage.
