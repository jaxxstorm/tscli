## Context

`coverage/coveragegaps` currently inventories OpenAPI paths and methods, compares them to the command-operation map, and reports missing operations or unmapped commands. That catches endpoint parity drift, but it does not inspect request-body or response-body properties for mapped operations, so upstream additions can land without any signal that `tscli` is failing to serialize or surface them.

The gap shows up most clearly on model-heavy commands such as `list devices`, where the CLI reuses `tsapi.Device` and marshals the result directly. If the pinned OpenAPI schema adds a field like `PostureIdentity`, today’s coverage tooling does not tell us whether the field exists in the Go model, whether fixtures exercise it, or whether output coverage has been updated to account for it.

## Goals / Non-Goals

**Goals:**
- Add a deterministic property-level coverage inventory for request and response schemas tied to mapped OpenAPI operations.
- Make uncovered schema properties visible in machine-readable and markdown coverage reports, with CI failure modes for regressions and unresolved gaps.
- Define a reviewable way to mark properties as covered, intentionally excluded, or pending implementation.
- Strengthen contract validation and representative command tests around important response models so field-level drift is easier to catch.
- Keep the existing command surface and runtime behavior stable while improving verification.

**Non-Goals:**
- Replacing the current operation-parity workflow or command-operation map.
- Generating CLI code automatically from OpenAPI schemas.
- Proving every property is semantically rendered in a human-friendly table view; the focus is on request serialization, response decoding, and structured output coverage.
- Auditing every historical endpoint in one undifferentiated pass without a manifest or exclusions model.

## Decisions

### 1. Build property coverage from mapped OpenAPI operations, not from ad hoc model lists

The property inventory should start from the pinned OpenAPI snapshot plus the existing command-operation map. For each mapped operation, the tool should inspect JSON request bodies and success responses, resolve referenced schemas, and flatten property paths into a stable key space such as:
- `get /tailnet/{tailnet}/devices response devices[].postureIdentity`
- `post /tailnet/{tailnet}/keys request capabilities.devices.create.tags`

Rationale: this ties property coverage to the same source of truth already used for endpoint parity, keeping the report deterministic and scoped to the command surface we claim to support.

Alternative considered:
- Maintaining a hand-written list of important models and fields. Rejected because it would drift quickly and would not stay aligned with mapped operations.

### 2. Track property coverage and exclusions explicitly in repository data

Introduce a dedicated property-coverage manifest/exclusions file alongside the existing OpenAPI contract files so the report can distinguish:
- covered properties
- intentionally excluded properties with rationale
- uncovered properties that require code or test work

Coverage entries should point to stable evidence such as command path, model path, or test identifier rather than free-form prose only.

Rationale: property coverage is too granular to infer perfectly from the codebase alone, and reviewers need explicit, diffable declarations when upstream schema changes occur.

Alternative considered:
- Attempting fully automatic detection from Go reflection and grep-based code search. Rejected because it would be brittle, miss indirect serialization paths, and produce noisy false positives.

### 3. Extend the existing coverage report instead of creating a disconnected second workflow

Property-level analysis should plug into `coverage/coveragegaps` or a closely related coverage workflow so existing make/CI entry points continue to produce the canonical reports. The machine-readable output and markdown summary should add property totals, uncovered properties, exclusions, and baseline diffing without removing the current operation-level sections.

Rationale: the repo already has a parity-report workflow and CI expectations around it; property coverage should deepen that report rather than fragment it.

Alternative considered:
- Building a completely separate tool with separate reports and CI jobs. Rejected because it would split reviewer attention and make enforcement easier to ignore.

### 4. Back property coverage with representative contract and integration assertions

The manifest/report should be paired with focused tests for representative commands and models, especially where the CLI passes through SDK models directly. For example, device-list fixtures and contract tests should exercise fields like `postureIdentity` so decoding/output regressions fail before coverage data is manually updated.

Rationale: the report identifies missing coverage, but tests provide executable evidence that covered fields are actually decoded/serialized correctly.

Alternative considered:
- Treating the manifest alone as sufficient proof. Rejected because it would allow stale coverage declarations without runtime validation.

## Risks / Trade-offs

- [Property inventories could explode in size and become noisy] -> Mitigation: scope the report to mapped operations, support explicit exclusions, and present grouped summaries by operation/model.
- [Manifest maintenance may feel burdensome] -> Mitigation: keep the schema-derived keys stable, provide focused diff output, and integrate updates into the existing refresh/review workflow.
- [Directly reused SDK models may still hide semantic output gaps] -> Mitigation: pair manifest entries with representative fixtures/tests for high-value responses.
- [Upstream schema changes could trigger churn in CI] -> Mitigation: continue using the pinned snapshot, baseline diffs, and explicit refresh workflows so property changes are reviewed intentionally.

## Migration Plan

1. Define the property inventory format and coverage/exclusion data files under the existing OpenAPI contract area.
2. Extend coverage analysis to emit request/response property coverage alongside operation parity results.
3. Add baseline diff and fail-on-gap behavior for uncovered properties.
4. Seed the manifest with an initial slice of mapped operations, including device output coverage such as `postureIdentity`.
5. Add or strengthen representative contract/integration tests that assert selected schema properties round-trip correctly.
6. Wire the enhanced report into existing make/CI coverage checks.

Rollback strategy:
- Revert the property-reporting additions and manifest files while preserving the existing operation-level parity workflow.

## Open Questions

- Should property coverage treat query/header/path parameters as part of the same report, or stay focused on JSON body/request-response schemas first?
- What minimum evidence should each covered property entry require: model path only, or a referenced test/assertion as well?
