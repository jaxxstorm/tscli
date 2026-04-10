## Why

`tscli` currently checks whether OpenAPI operations are represented by commands, but it does not tell us whether request and response properties from those operations are actually parsed, serialized, and surfaced by the CLI. That gap makes upstream schema additions like device `PostureIdentity` easy to miss, so we need property-level coverage checks now to catch silent contract drift before users do.

## What Changes

- Add property-level coverage analysis for pinned OpenAPI request and response schemas so coverage checks report unmapped request/response fields, not just missing operations.
- Define how commands and tests declare property coverage for fields they intentionally parse, emit, or serialize.
- Extend contract-consistency checks so new upstream properties in important models are reviewable and regressions fail deterministically in CI.
- Keep the existing command surface stable: no command renames, no removed flags, and no new user configuration knobs.
- Use concrete examples such as device output fields to validate that CLI JSON/YAML output and request payload handling stay aligned with the pinned schema.

## Capabilities

### New Capabilities
- `openapi-property-coverage`: Report and enforce coverage for request and response properties defined in the pinned OpenAPI schema.

### Modified Capabilities
- `coverage-gap-elimination`: Extend coverage reporting and CI regression checks beyond operation parity to include uncovered properties.
- `tailscale-model-contract-consistency`: Tighten contract validation so request/response model handling is checked at the property level against the pinned OpenAPI snapshot.

## Impact

- Affected code:
  - `coverage/coveragegaps/**`
  - `coverage/openapirefresh/**`
  - `pkg/contract/openapi/**`
  - request/response contract tests under `pkg/**` and `test/**`
- Affected command groups/flags/config:
  - Broadly affects commands backed by OpenAPI request/response models, with immediate focus on device/list/get/settings-style outputs and payloads
  - No new CLI flags
  - No new config keys or environment variables
- Backward compatibility:
  - No user-facing breaking change is intended for command syntax.
  - CI and coverage workflows will become stricter by failing when newly introduced schema properties are not intentionally covered or excluded.
