## Context

The repository already has make targets and Go programs for OpenAPI refresh and coverage-gap reporting:

- `make openapi-refresh` runs `coverage/openapirefresh`
- `make coverage-gaps` runs `coverage/coveragegaps`
- `make coverage-gaps-check` adds baseline regression and gap failure checks
- `make coverage-gaps-latest` refreshes the OpenAPI snapshot before running strict checks

Those workflows are effective, but they are not discoverable from agent skill metadata and are not exposed as `tscli` commands. Existing coverage tools are `package main`, which makes the logic hard to reuse from Cobra commands without shelling out to `go run`.

## Goals / Non-Goals

**Goals:**

- Add repo-local `.codex` and `.opencode` skills that tell agents to run coverage-gap and OpenAPI refresh workflows before implementing command parity or schema-driven work.
- Add script-friendly `tscli maintenance coverage-gaps` and `tscli maintenance openapi-refresh` commands.
- Reuse the existing coverage and OpenAPI refresh defaults so CLI behavior matches make target behavior.
- Keep maintenance commands local so they bypass API-key preflight checks.
- Preserve deterministic CI and existing make targets.

**Non-Goals:**

- Do not refresh the pinned OpenAPI snapshot as part of this change unless a user explicitly runs the refresh command.
- Do not change normal Tailscale API command behavior, authentication, or output defaults.
- Do not replace make targets or generated coverage artifacts.
- Do not add new runtime dependencies.

## Decisions

### Extract reusable packages for maintenance workflows

Move the core logic from `coverage/coveragegaps` and `coverage/openapirefresh` into importable internal packages, while leaving their existing `main.go` entrypoints as thin flag parsing wrappers.

Rationale: Cobra commands can call the same implementation directly without invoking a nested `go run`, which keeps tests fast and avoids depending on a Go toolchain at runtime. The existing make targets remain compatible because the command packages keep the same flags.

Alternative considered: Have `tscli maintenance` shell out to `make`. This would be simple but less portable, harder to test, and would require `make` and source tree assumptions at runtime.

### Add a local `maintenance` command group

Add `cmd/tscli/maintenance` with:

- `coverage-gaps`
- `openapi-refresh`

The root pre-run treats `maintenance` like `config` and `agent`, so these commands do not require API credentials. Command flags mirror the underlying tools and use the same repository-relative defaults.

Rationale: This keeps maintenance work discoverable through the CLI while preserving script-friendly behavior and avoiding config/env churn.

Alternative considered: Add commands under `agent`. That would blur agent installation concerns with repository maintenance and make the command harder to find for human maintainers.

### Add paired Codex and OpenCode skills

Create matching skill directories under `.codex/skills` and `.opencode/skills`:

- `coverage-gaps`
- `openapi-refresh`

Each skill names the exact command sequence, when to run it, expected artifacts, and how to proceed when checks fail. The instructions prefer existing make targets and allow the new `tscli maintenance` commands when validating the CLI surface.

Rationale: Skills are local workflow documentation for agents, so they should point at the repository's supported automation rather than duplicate implementation details.

Alternative considered: One combined skill. Separate skills are easier to trigger precisely; agents can still use both for schema refresh followed by coverage checks.

## Risks / Trade-offs

- CLI commands are repository-maintenance oriented and depend on repository-relative default paths -> Mitigation: expose path flags mirroring the underlying tools so callers can override paths in scripts.
- OpenAPI refresh uses network access and can fail in offline environments -> Mitigation: keep refresh explicit, return clear errors, and keep contract tests pinned to in-repo fixtures.
- Moving code from `package main` can accidentally change existing make target behavior -> Mitigation: keep current flags and add focused tests for both package logic and command help/pre-run behavior.
- New visible commands affect leaf command manifest tests and docs -> Mitigation: update the manifest and generated command docs.
