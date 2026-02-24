# Coverage Gaps Baseline

Date: 2026-02-24

## OpenAPI Snapshot

- Source: `https://api.tailscale.com/api/v2?outputOpenapiSchema=true`
- OpenAPI version: `3.1.0`
- API version: `v2`
- Path count: `56`
- Operation count: `85`

## CLI/Test Baseline

- Command directories with `cli.go`: `88`
- Command directories with tests: `1`
- Covered command directories (co-located tests): `1`
- Uncovered command directories (no co-located tests): `87`

## Sample Uncovered Command Directories

- `cmd/tscli/config`
- `cmd/tscli/config/get`
- `cmd/tscli/config/set`
- `cmd/tscli/config/show`
- `cmd/tscli/create`
- `cmd/tscli/create/integration`
- `cmd/tscli/create/invite`
- `cmd/tscli/create/invite/device`
- `cmd/tscli/create/invite/user`
- `cmd/tscli/create/key`
- `cmd/tscli/create/token`
- `cmd/tscli/create/webhook`
- `cmd/tscli/delete`
- `cmd/tscli/delete/device`
- `cmd/tscli/delete/device/invite`
- `cmd/tscli/delete/device/posture`
- `cmd/tscli/delete/devices`
- `cmd/tscli/delete/integration`
- `cmd/tscli/delete/invite/device`
- `cmd/tscli/delete/invite/user`

## Notes

- This is a coarse baseline and not yet a full command-to-operation mapping.
- The implementation tasks in `tasks.md` define the next step: generate deterministic coverage-gap reports that classify OpenAPI operations as covered, uncovered, or unmapped.
