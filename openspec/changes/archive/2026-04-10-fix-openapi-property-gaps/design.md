## Context

The property-coverage rollout proved that `tscli` can inventory request and response properties from the pinned Tailscale OpenAPI snapshot, but the current report still contains a large default-exclusion bucket for mapped command operations. Two concrete causes showed up immediately:

- some commands rely on upstream SDK models that do not match the pinned schema, most notably `tsapi.Device`, which still mis-tags `advertisedRoutes` and omits `multipleConnections`
- some mutating commands send the right request but print a local summary or echoed request instead of decoding the API response body that the schema defines

This change needs to reduce exclusions without destabilizing the CLI surface. Existing flags, Viper-backed config flow, and output-mode handling should remain intact. The implementation should prefer narrow shared helpers over a wholesale client rewrite.

## Goals / Non-Goals

**Goals:**
- Remove the current explicit device-property exclusions and reduce the default-exclusion bucket for mapped request/response bodies that `tscli` already supports
- Ensure audited commands decode request and response bodies through schema-aligned models, even when the upstream SDK lags the pinned schema, so API-documented properties are never silently dropped from structured output
- Expand property manifests, fixtures, contract tests, and mock-backed command tests together so coverage evidence matches runtime behavior
- Preserve existing flag names, env vars, config keys, and output mode selection

**Non-Goals:**
- Achieve full endpoint parity for unmapped API operations
- Replace the upstream Tailscale SDK across the entire codebase
- Redesign human/pretty output formats beyond the minimum changes needed to keep audited fields visible in structured output
- Introduce new config or auth behavior

## Decisions

### 1. Remediate property gaps by mapped operation family instead of trying to clear every exclusion at once

The implementation should work through the currently mapped operations that still sit behind default or explicit exclusions, starting with device, route, settings, invite, posture, and service-approval commands.

Why:
- the report already groups exclusions by mapped operation, which makes the work trackable
- device-related operations share the highest-risk schema mismatch and unblock multiple exclusions at once
- request-only write operations can often be covered quickly after the shared response model work is done

Alternatives considered:
- Clear every excluded property in one pass. Rejected because the command surface is too broad for a single low-risk implementation batch.
- Keep exclusions in place until upstream SDK fixes arrive. Rejected because that leaves already-detected contract gaps unresolved.

### 2. Introduce local schema-aligned DTOs and helpers for audited operations whose SDK models are incomplete

For audited command paths, `tscli` should decode through small local request/response DTOs under `pkg/` and use shared HTTP helpers on top of `pkg/tscli.Do` when the upstream SDK type cannot preserve the schema fields the CLI needs.

Why:
- `tsapi.Device` cannot currently represent all audited device response fields
- local DTOs let command output, coverage reflection, and contract tests share the same schema-aligned field definitions
- `pkg/tscli.Do` already exists for endpoints not covered correctly by the SDK, so this extends an existing pattern instead of introducing a second transport stack

Alternatives considered:
- Wait for a newer SDK. Rejected because the current pinned dependency already shows the mismatch.
- Fork or patch the upstream SDK in-place. Rejected because it increases maintenance cost and does not help commands that intentionally discard response bodies.

### 3. For audited write operations, decode and print the authoritative API response when the schema defines one

Commands that currently print echoed requests or synthetic success maps should switch to decoding the API response body once that operation side is placed under property coverage. If a property exists in the API response for a mapped operation, `tscli` should preserve and show it instead of silently dropping it because of an SDK limitation or local placeholder output.

Why:
- the property coverage system is about parsed request and response properties, not only outbound payloads
- synthetic summaries hide schema fields and create false confidence in response coverage
- structured output users benefit from additive, authoritative fields rather than local approximations

Alternatives considered:
- Keep synthetic output and mark the response excluded. Rejected because that preserves the gap instead of fixing it.
- Decode the response but continue printing only a summary. Rejected because coverage evidence would no longer match observable CLI behavior.

### 4. Tie coverage data, fixtures, and tests to the same audited models

Each newly audited operation side should land with:

- an explicit manifest entry in `coverage/property-coverage.yaml`
- removed or narrowed exclusions in `coverage/property-exclusions.yaml`
- representative mock fixtures in `internal/testutil/apimock`
- command-level integration assertions in `test/cli`
- contract/property assertions in `pkg/contract/openapi`

Why:
- the current framework only stays trustworthy if manifest coverage reflects what command code actually decodes and prints
- representative fixture coverage catches accidental field loss that reflection alone cannot see

Alternatives considered:
- Update manifests first and add runtime tests later. Rejected because it weakens the evidence behind claimed coverage.

## Risks / Trade-offs

- [Output compatibility for mutating commands] -> Switching from synthetic summaries to API response bodies may affect scripts that parse exact JSON keys. Mitigation: stage changes by command, keep output structured, and prefer returning the full authoritative response so fields are added rather than silently omitted.
- [Local DTO drift from upstream schema] -> Adding local models creates another source of truth. Mitigation: keep DTOs minimal, back them with pinned-schema contract checks, and reuse them in the coverage manifest.
- [Mixed SDK and non-SDK code paths] -> Some commands will still use the upstream client while audited commands use `tscli.Do` helpers. Mitigation: centralize audited helpers in one shared package and keep command code thin.
- [Change size across many mapped operations] -> Reducing the default-exclusion bucket could expand scope quickly. Mitigation: land the work in operation-family batches and keep exclusions only where a documented blocker still exists.

## Migration Plan

1. Add shared schema-aligned DTOs and helper functions for the first audited operation family, starting with device and route responses.
2. Switch commands that currently depend on broken SDK decoding or synthetic response output to those shared helpers.
3. Expand property manifests and remove the matching exclusions as each operation side gains runtime coverage.
4. Add or update contract and integration tests for every new audited request/response side.
5. Run the coverage-gap report and targeted Go test suites until the selected exclusions are removed with no new regressions.

Rollback is straightforward: revert the command/helper changes for the affected operation family and restore the matching exclusions.

## Open Questions

- None. The change direction is to preserve and show API-documented response properties for mapped operations even when the upstream SDK does not expose them directly.
