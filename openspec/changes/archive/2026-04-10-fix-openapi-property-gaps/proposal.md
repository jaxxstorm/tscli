## Why

`tscli` can now detect request and response property gaps from the pinned OpenAPI schema, and the current report still carries a large default-exclusion bucket plus known SDK-driven device-response mismatches. Closing those gaps now lets the existing `get`, `list`, `set`, and invite/service command groups prove that they parse and emit the documented request and response properties without changing flags, config keys, or env vars.

## What Changes

- Retire the current explicit property exclusions for device response fields that `tscli` already exposes or claims to support, including `advertisedRoutes` and `multipleConnections`.
- Audit the remaining mapped request and response bodies that are still in the default exclusion bucket, starting with device, route, settings, invite, posture, and service-approval operations already shipped in the CLI.
- Introduce schema-aligned request/response decoding for commands whose current SDK usage or synthetic success output would otherwise silently drop documented API fields.
- Add representative mock-backed integration tests and contract assertions for newly audited properties so regressions fail in CI before new exclusions are added.

## Capabilities

### New Capabilities

### Modified Capabilities
- `openapi-property-coverage`: Require currently mapped command operations to move from rollout-era default exclusions to explicit request/response coverage or narrowly justified exclusions.
- `tailscale-model-contract-consistency`: Require schema-aligned decoding and test evidence when existing command requests or responses depend on upstream SDK models or local placeholder output that would otherwise omit documented properties.

## Impact

- Affected command groups: existing `get`, `list`, `set`, invite, posture, and service-related commands whose mapped OpenAPI operations still rely on property exclusions
- Affected flags: existing `--all` and `--device` behavior stays intact; no new flags are introduced
- Config and env keys: no changes to `tailnet`, `api-key`, `base-url`, `output`, or `debug`
- Affected code and data: request/response decoding paths, property coverage manifests and exclusions, mock API fixtures, and contract/integration tests
- Backward compatibility: no new flags or config changes; structured outputs may gain additional API-documented fields or switch from synthetic summaries to authoritative API responses so documented response properties are never silently dropped
