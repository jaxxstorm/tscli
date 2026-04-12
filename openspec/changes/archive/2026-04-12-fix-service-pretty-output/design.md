## Context

`list services` and `get service` currently decode into `json.RawMessage`, re-marshal, and hand the payload to the generic output pipeline. That works for JSON and yaml, but the pretty/human renderer is optimized for top-level arrays or object records that already match the CLI's display expectations. Real service responses include nested service structures and collection wrappers that cause the generic formatter to flatten sibling service data into a single unreadable block.

The change touches multiple modules: service commands under `cmd/tscli`, shared output behavior under `pkg/output`, fixtures in `internal/testutil/apimock`, and example/integration tests under `test/cli`. The design needs to preserve existing flags and script-friendly structured output while fixing interactive readability.

## Goals / Non-Goals

**Goals:**
- Make `list services` pretty/human output render one service per record with stable field grouping.
- Make `get service` pretty/human output render a single service record consistently with other `get` commands.
- Preserve JSON and yaml response structure so existing scripts keep working.
- Add regression tests that exercise representative service payloads, including nested fields such as annotations.

**Non-Goals:**
- Redesign the entire pretty printer for every command shape in the CLI.
- Change service command flags, request paths, or output schema for JSON/yaml.
- Introduce config keys or output-mode-specific flags for service rendering.

## Decisions

### Decode service responses into service-shaped values before printing
The service commands will stop treating responses as opaque `json.RawMessage` for all output modes. Instead, they will decode into local service-oriented models or `map[string]any` structures that reflect the actual list and single-service payload shapes, then pass those values through the output layer.

Rationale:
- The bug starts at the boundary where the command loses shape information by treating the payload as undifferentiated raw JSON.
- Command-local shaping is lower risk than making the pretty printer infer special handling for one API family from arbitrary nested maps.
- JSON and yaml output can still be derived from the decoded structure without changing their semantics.

Alternatives considered:
- Teach the generic pretty printer to detect `vipServices`-style payloads. Rejected because it hardcodes service-specific API knowledge into a global renderer and risks unrelated regressions.
- Leave command decoding as-is and only patch wrapping/indent logic in `pkg/output`. Rejected because the malformed rendering is caused by the service payload shape entering the printer without command-aware normalization.

### Preserve the API response contract for structured modes
The implementation will keep JSON/yaml output schema-aligned with the API response while adapting pretty/human rendering to operate on the decoded service records.

Rationale:
- Scripts depend on structured output staying predictable.
- The user-visible bug is in pretty/human readability, not in the machine-readable modes.

Alternatives considered:
- Emit a simplified synthetic summary for all modes. Rejected because it would drop API-documented service fields and create compatibility risk for scripts.

### Add representative service output fixtures and mode-specific assertions
Tests will use richer service fixtures that include addresses, ports, tags, comments, and annotations. Example output coverage will assert not just JSON shape but also readable pretty/human rendering for `list services` and `get service`.

Rationale:
- The current tests pass because they only assert trivial JSON keys on simplistic fixtures.
- Regressions in interactive rendering need direct text assertions, not just structural JSON checks.

Alternatives considered:
- Rely only on unit tests for the pretty printer. Rejected because the defect is command-path-specific and must be covered end to end.

## Risks / Trade-offs

- [Service payload shape may differ between list and get responses] -> Use separate command-local response models or normalization paths and cover both with tests.
- [Command-local shaping could drift from the pinned schema over time] -> Keep JSON/yaml output schema-aligned and extend fixtures/tests to cover important service fields.
- [Changing shared output helpers could affect unrelated commands] -> Prefer minimal shared helper changes and keep service-specific adaptation close to the service commands.

## Migration Plan

No migration is required. The change is an in-place correction to existing command behavior:
1. Update service command decoding/output handling.
2. Expand fixtures and command-level tests for representative service payloads and pretty/human assertions.
3. Verify JSON/yaml output remains stable.
4. Rollback, if needed, by restoring the previous service command output path; no persisted data or config is affected.

## Open Questions

- Whether `list services devices` shares any of the same rendering issues once richer fixtures are added, and whether it should be covered in the same implementation pass if a related issue appears.
