## Context

`tscli` is a Go-based Cobra/Viper CLI that wraps the Tailscale HTTP API and documents every supported operation through auto-generated command reference pages. Today the `create key` flows only expose `authkey` and `oauthclient` credentials, yet the upstream API also accepts `keyType: federated`. Because the tooling ignores that variant, both the CLI and `coverage/coveragegaps` report a coverage gap for the operation, and users cannot provision federated identities from the CLI.

## Goals / Non-Goals

**Goals:**
- Allow `tscli create key` to request federated credentials without altering existing flag behavior for `authkey`/`oauthclient`.
- Update the OpenAPI snapshot/command coverage map so the new `federated` operation is marked as implemented and surfaced in docs/tests.
- Keep docs, tests, and coverage tooling aligned so the federated key path stays covered in future releases.

**Non-Goals:**
- Rewriting the entire command hierarchy or flag parsing framework.
- Adding unrelated coverage features beyond the new key type.

## Decisions

- **Add explicit `--type` option with allowed values**: Instead of hard-coding the existing behavior or introducing separate commands, re-use the current `--type` flag on `create key` to accept `authkey`, `oauthclient`, and now `federated`. This keeps the CLI surface stable while ensuring users can request any valid `keyType` and the API payload stays consistent.
- **Treat federated keys as a first-class coverage mapping**: Update the OpenAPI snapshot under `pkg/contract/openapi` (and any generated mapping used by `coverage/coveragegaps`) to include the federated-key operation as soon as the new flag works. This avoids a regression where the coverage tool continues to report unmapped operations after implementation.
- **Expand tests/docs rather than duplicating flows**: Reuse the existing test harness for `create key` by adding federated-specific cases and add generated docs describing the new flag/the API mapping, ensuring parity without a second code path.

## Risks / Trade-offs

- [Risk] `coverage/coveragegaps` may still flag the operation if the OpenAPI mapping isn’t updated in lockstep → Mitigation: coordinate code changes so the federated path is added in the same commit that updates the coverage metadata.
- [Risk] Federated credentials have additional server-side requirements that could change the shape of the API response → Mitigation: rely on the existing JSON deserialization and focus on key-type enumeration; keep the coverage/test harness limited to request validation and not downstream behavior.

## Migration Plan

1. Extend `cmd/tscli/create/key` to accept `federated` in the `--type` flag and ensure the resulting payload sends `keyType: federated`.
2. Refresh `pkg/contract/openapi` snapshots to include the federated-key operation and regenerate any coverage maps that drive `coveragegaps`.
3. Update coverage documentation/tests that rely on the snapshot so the new command is considered covered.

## Open Questions

- Should the new federated key type document any additional default flags (e.g., default reuse/ephemeral settings) or rely entirely on existing flags?
