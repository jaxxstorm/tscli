## Context

The current parity report shows substantial CLI coverage drift versus the pinned Tailscale OpenAPI snapshot (`85` operations, `80` uncovered). Some existing commands are also unmapped or non-canonical in naming, which makes discovery and long-term maintenance difficult.

This change is cross-cutting across most command domains (`device`, `invites`, `dns`, `logging`, `keys`, `policy`, `posture`, `users`, `contacts`, `webhooks`, `services`, `settings`) and must preserve script compatibility while enforcing canonical verb taxonomy:
- create operations -> `create`
- update/mutation operations -> `set`
- single retrieval -> `get`
- collection retrieval -> `list`
- deletions -> `delete`

## Goals / Non-Goals

**Goals:**
- Eliminate in-scope OpenAPI coverage gaps by adding missing CLI commands.
- Normalize canonical command paths to taxonomy-aligned verbs.
- Preserve existing behavior through backward-compatible aliases where command names change.
- Ensure every command has explicit operation mapping and automated mock-backed tests.
- Enforce parity and regression checks in CI.

**Non-Goals:**
- Introducing new API business behavior beyond Tailscale API parity.
- Removing backward-compatible aliases in this change.
- Replacing existing output format architecture (`pretty`/`human`/`json`/`yaml`).

## Decisions

### 1. Operation-first implementation workflow

Decision:
- Use the pinned OpenAPI operation set as the source of truth.
- Classify operations into `mapped`, `excluded`, and `missing`.
- Implement missing operations by API domain batches.

Rationale:
- Prevents ad-hoc command additions and ensures deterministic parity closure.

Alternatives considered:
- Command-first backlog from existing directories only: rejected because it can miss uncovered operations.

### 2. Canonical taxonomy with compatibility aliases

Decision:
- Canonical command paths will follow the required taxonomy and full noun naming.
- Existing abbreviated or legacy paths stay as aliases during migration.

Rationale:
- Preserves user scripts while converging UX to a consistent model.

Alternatives considered:
- Hard rename without aliases: rejected due to breaking-change risk.

### 3. Domain-by-domain command delivery

Decision:
- Implement by domain slices: core resources first (`devices`, `invites`, `keys`, `users`, `webhooks`, `settings`), then supporting domains (`dns`, `logging`, `services`, `contacts`, `policy`, `posture`, `aws-external-id`).

Rationale:
- Keeps pull requests reviewable and testable while steadily reducing uncovered operations.

Alternatives considered:
- One-shot mega implementation: rejected due to review and regression risk.

### 4. Strict mapping and parity checks

Decision:
- Update command-operation map with every command addition.
- Treat unknown mapped commands and uncovered in-scope operations as failing states in CI.

Rationale:
- Makes parity drift visible and enforceable.

Alternatives considered:
- Informational reports only: rejected due to weak enforcement.

### 5. Reuse existing test harness architecture

Decision:
- Use existing shared command harness, API mock server, and output assertion helpers.
- Add table-driven command tests for each newly added command.

Rationale:
- Aligns with established patterns from the prior test-coverage change and keeps implementation velocity high.

Alternatives considered:
- New test framework: rejected as unnecessary complexity.

## Risks / Trade-offs

- [Large command surface expansion] -> Mitigation: domain-batched implementation and required per-domain tests before advancing.
- [Taxonomy normalization breaks scripts] -> Mitigation: retain aliases and document canonical replacements.
- [Ambiguous operation classification] -> Mitigation: explicit classification policy + exclusions file + reviewer validation.
- [Report churn from OpenAPI changes] -> Mitigation: pinned snapshot workflow with explicit refresh process and diff review.
- [Inconsistent command UX across domains] -> Mitigation: shared validation/output helpers and review checklist per command.

## Migration Plan

1. Freeze baseline parity report and classify uncovered operations by domain and command verb.
2. Add canonical command path updates with compatibility aliases.
3. Implement missing operations by domain with command mapping updates.
4. Add/expand unit and mock integration tests for each implemented command.
5. Drive coverage-gap report to zero uncovered in-scope operations and zero unknown mapped commands.
6. Update README command usage and taxonomy guidance.
7. Validate `go test ./...` and coverage-gap CI checks before merge.

Rollback strategy:
- If regressions occur, keep canonical command additions but retain old alias behavior and revert only the affected domain batch.

## Open Questions

- Which operations (if any) should be permanently excluded from parity (for example highly specialized endpoints) and what rationale should be codified?
- Should canonical path normalization for existing abbreviations happen in one release or phased by domain?
- Do we require one command per operation, or allow a single command path to cover multiple closely related operations when behavior is clear?
