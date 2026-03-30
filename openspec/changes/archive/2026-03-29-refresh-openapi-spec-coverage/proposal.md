## Why

`tscli` coverage-gap checks are still tied to an older pinned Tailscale OpenAPI snapshot, so the project cannot quickly verify whether recent upstream API changes introduced new uncovered operations. We need an explicit developer workflow to refresh the pinned schema from the canonical Tailscale endpoint and immediately re-run coverage-gap checks against the updated snapshot.

## What Changes

- Add a documented developer make target that fetches the latest Tailscale OpenAPI schema from `https://api.tailscale.com/api/v2?outputOpenapiSchema=true`, updates the in-repo pinned snapshot, and records refreshed snapshot metadata.
- Add a make target that runs coverage-gap generation/checks after the snapshot refresh so maintainers can validate parity against the latest upstream API surface before implementing new command coverage.
- Update the OpenAPI snapshot and coverage-gap requirements to define the expected refresh workflow, generated artifacts, and failure behavior for developers and CI.
- Do not change any end-user `tscli` command groups, command flags, config keys, or environment variables; this change is limited to repository maintenance and verification workflows.
- Preserve backward compatibility for existing automation by keeping current `make coverage-gaps` and `make coverage-gaps-check` flows intact while adding a clear path for latest-schema validation.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `tailscale-model-contract-consistency`: Require an explicit, versioned refresh workflow for pulling the latest upstream OpenAPI snapshot and recording metadata in-repo.
- `cli-command-test-coverage`: Require a supported make-driven workflow to regenerate coverage-gap artifacts against a refreshed OpenAPI snapshot.

## Impact

- Affected code:
  - `Makefile`
  - `README.md`
  - OpenAPI snapshot and metadata files under the repository's contract fixtures
  - Coverage-gap tooling and generated coverage artifacts
- Affected systems:
  - Tailscale OpenAPI source endpoint `https://api.tailscale.com/api/v2?outputOpenapiSchema=true`
  - Local developer workflows and CI parity verification
- Backward compatibility:
  - No runtime CLI behavior changes
  - Existing scripts that use current make targets remain supported
