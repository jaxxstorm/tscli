## Context

`tscli` already documents a manual OpenAPI refresh flow in `README.md` and exposes `make coverage-gaps` and `make coverage-gaps-check`, but the workflow is fragmented and still depends on maintainers stitching together curl, hashing, tests, and report generation themselves. At the same time, the repository state suggests the pinned snapshot artifacts referenced by the docs are not consistently present in the worktree, which makes it harder to verify parity after upstream Tailscale API updates.

This change is cross-cutting across repository maintenance assets rather than runtime CLI command groups:
- no end-user `tscli` command group behavior changes
- no new command flags, config keys, or environment variables
- new developer-facing `make` targets for refresh and validation

The design must preserve the current deterministic test model:
- normal `go test` and CI runs continue to use pinned in-repo artifacts only
- network access is limited to an explicit refresh workflow
- existing `make coverage-gaps` and `make coverage-gaps-check` automation remains usable for current scripts

## Goals / Non-Goals

**Goals:**
- Add a standard make-driven workflow to fetch the latest Tailscale OpenAPI schema from the canonical endpoint.
- Regenerate pinned snapshot metadata in a reviewable form whenever the schema is refreshed.
- Provide a make target that refreshes the schema and immediately runs coverage-gap validation against the refreshed snapshot.
- Keep implementation small and composable by reusing existing coverage-gap tooling and current repository conventions.
- Add tests and documentation that make the refresh workflow safe to use and easy to review.

**Non-Goals:**
- Changing runtime `tscli` command behavior, output modes, or flag semantics.
- Replacing the existing coverage-gap analyzer or redefining parity policy.
- Introducing runtime network fetches into standard tests or CI checks.
- Solving unrelated OpenAPI/model drift beyond the refresh and validation workflow itself.

## Decisions

### 1. Add a dedicated refresh target plus a composed validation target

Decision:
- Add a low-level make target that downloads the latest schema from `https://api.tailscale.com/api/v2?outputOpenapiSchema=true` into the pinned OpenAPI snapshot path and rewrites snapshot metadata.
- Add a second make target that depends on the refresh target and then runs the existing coverage-gap check flow so maintainers can perform a single-command latest-schema validation.

Rationale:
- Separating refresh from validation keeps the workflow composable for local debugging while still giving maintainers a one-command path for the common "refresh then test coverage gaps" use case.
- It preserves backward compatibility because current scripts can continue to call existing coverage targets.

Alternatives considered:
- One replacement target only: rejected because it would make ad-hoc snapshot refreshes and reviewable multi-step workflows harder.
- Manual README-only instructions: rejected because the gap is specifically lack of an executable, repeatable workflow.

### 2. Treat snapshot metadata as generated contract state

Decision:
- The refresh target updates both the pinned schema file and snapshot metadata in one step.
- Metadata should include at least source URL, fetch timestamp, and a stable schema identifier such as version and/or digest so reviewers can confirm what changed.

Rationale:
- The existing spec already requires explicit and versioned contract updates, so the refresh workflow should always leave behind reviewable provenance.
- Keeping metadata generation coupled to schema refresh avoids partial updates where the schema changes but provenance does not.

Alternatives considered:
- Hash-only manual review: rejected because it leaves timestamp/source provenance implicit.
- Metadata updates by hand: rejected because it is error-prone and easy to forget.

### 3. Reuse existing coverage-gap tooling instead of adding a parallel checker

Decision:
- The latest-schema validation target should call the existing coverage-gap generator/checker with the refreshed pinned snapshot rather than introducing a second analysis path.
- Any new wiring should stay in `Makefile` and small helper code/scripts near the current contract and coverage tooling.

Rationale:
- Reusing the current analyzer keeps parity rules consistent between baseline checks and latest-schema checks.
- It minimizes disruption to existing commands and scripts and limits the implementation surface.

Alternatives considered:
- Separate "latest" analyzer mode with duplicated logic: rejected because it increases maintenance cost and drift risk.

### 4. Validate the workflow with deterministic tests plus target-level smoke coverage

Decision:
- Keep unit coverage focused on any new metadata-writing or schema-refresh helper logic.
- Extend command-level/developer workflow validation by exercising the refreshed coverage-gap path through existing Go tests or a small smoke-style verification around generated artifacts.
- Update README examples to document both the refresh-only and refresh-plus-coverage commands.

Rationale:
- The behavior users care about is observable artifact generation and failure semantics, not the internals of make itself.
- Documentation is part of the contract here because the workflow is developer-facing rather than exposed as a `tscli` subcommand.

Alternatives considered:
- No automated validation for make targets: rejected because the request is specifically about making the workflow reliable.

## Risks / Trade-offs

- [Upstream schema fetch fails or returns malformed data] -> Mitigation: fail the refresh target immediately, keep existing pinned artifacts unchanged on error, and surface a clear message.
- [Repository currently lacks some documented snapshot artifact paths] -> Mitigation: normalize on one managed snapshot/metadata location during implementation and update docs to match reality.
- [Latest upstream schema introduces many new gaps at once] -> Mitigation: keep the new workflow explicit and review-driven so maintainers can inspect generated reports before follow-up implementation work.
- [Target naming confuses existing automation] -> Mitigation: leave current coverage-gap targets unchanged and add new names rather than repurposing old ones.

## Migration Plan

1. Identify the canonical in-repo OpenAPI snapshot and metadata paths used by coverage-gap tooling.
2. Implement the refresh target to fetch the upstream schema and rewrite snapshot metadata atomically.
3. Implement the composed validation target that refreshes the snapshot and runs coverage-gap checks.
4. Update README/developer docs with the supported workflow and generated artifacts.
5. Add or update tests around metadata generation and coverage-gap validation behavior.
6. Run the refresh workflow and coverage-gap checks to capture the new upstream baseline.

Rollback strategy:
- Revert the refresh workflow wiring and generated snapshot/metadata artifacts together, restoring the previously pinned schema and existing coverage targets.

## Open Questions

- Should the composed latest-schema target run `coverage-gaps` or the stricter `coverage-gaps-check` variant by default?
- Do we want to commit the refreshed coverage-gap baseline as part of this change, or only generate review artifacts and leave baseline updates to a follow-up implementation step?
