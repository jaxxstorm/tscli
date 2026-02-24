## Context

`tscli` currently has broad command coverage in code (`cmd/tscli/**`) but very limited automated tests, which increases regression risk for flag handling, output formatting, and API error behavior. The change introduces a consistent test architecture across command groups while preserving existing command UX, existing flags, and script compatibility.

Baseline snapshot (2026-02-24):
- Tailscale OpenAPI (`https://api.tailscale.com/api/v2?outputOpenapiSchema=true`) reports 56 paths and 85 operations.
- `tscli` currently has 88 command directories containing `cli.go`.
- Only 1 command directory currently contains tests, leaving broad command coverage gaps.

The design must work with existing Go/Cobra/Viper patterns:
- Cobra command tree is the command surface to test.
- Viper-backed configuration precedence (flags > env > config file) must be testable.
- Output behavior across `pretty`/`human`/`json`/`yaml` modes must be deterministic.
- Tailscale API behavior should be validated using mocks, not live external API calls.

## Goals / Non-Goals

**Goals:**
- Add reusable test harnesses for command execution and Tailscale API mocking.
- Add unit and integration test coverage for all command groups (`get`, `list`, `create`, `set`, `delete`, `config`, `version`).
- Verify required-flag validation, success paths, error paths, and output/exit behavior.
- Add model contract checks to detect request/response drift early.
- Keep CI execution reliable and deterministic without network dependence.

**Non-Goals:**
- Rewriting command implementations.
- Changing command names, flags, or user-facing output semantics except where current behavior is demonstrably incorrect.
- Depending on live Tailscale API in standard test runs.

## Decisions

### 1. Build a shared command test harness

Decision:
- Add shared test helpers that execute Cobra commands with controlled args, env, stdin, and captured stdout/stderr.
- Standardize assertions on exit error, stderr messages, and output format.

Rationale:
- Command tests are currently ad-hoc; shared helpers reduce repetition and improve consistency.
- Makes it practical to scale testing across dozens of command handlers.

Alternatives considered:
- Per-command bespoke tests only: rejected due to duplication and inconsistent assertions.

### 2. Add mock-backed integration tests using `httptest`

Decision:
- Introduce a mock Tailscale API server abstraction backed by `httptest.Server`.
- Provide fixture-driven responses and error injection per endpoint/verb.
- Ensure command tests route all API calls to mock base URL.

Rationale:
- Integration behavior (request composition + response handling) must be validated without external dependencies.
- `httptest` is idiomatic, fast, and deterministic in Go.

Alternatives considered:
- Mocking HTTP at transport interface only: useful for units, but less representative for endpoint-level behavior.
- Live API tests: rejected for reliability and credentials concerns.

### 3. Add explicit command coverage manifest checks

Decision:
- Add a machine-readable test coverage manifest or discovery check that maps leaf commands to at least one test case.
- Failing check signals newly added commands without tests.

Rationale:
- Prevents future command additions from bypassing testing.

Alternatives considered:
- Relying on reviewer discipline alone: rejected as inconsistent over time.

### 4. Add model contract consistency checks from pinned API schema/examples

Decision:
- Add contract tests that validate `tscli` payload handling against a pinned snapshot of the Tailscale OpenAPI schema.
- Use `https://api.tailscale.com/api/v2?outputOpenapiSchema=true` as the canonical schema source, then store a pinned in-repo snapshot for deterministic tests.
- Pin a snapshot version in-repo to keep CI deterministic; updates are explicit.

Rationale:
- Detects breaking API-model drift early while avoiding flaky network fetches in tests.

Alternatives considered:
- No contract tests: rejected due to drift risk.
- Runtime fetch of schema during test: rejected due to network nondeterminism.

### 5. Separate test lanes for speed and confidence

Decision:
- Unit tests: pure command validation and helper behavior.
- Integration tests: mock API server + command execution.
- Provide stable `go test` targets to run both locally and in CI.

Rationale:
- Keeps default feedback fast while preserving broader behavior checks.

### 6. Generate explicit CLI vs OpenAPI coverage-gap reports

Decision:
- Add a coverage analysis step that maps `tscli` commands/tests to OpenAPI operations and produces a machine-readable + human-readable gap report.
- Report MUST identify covered operations, uncovered operations, and ambiguous mappings requiring manual triage.

Rationale:
- Complements pass/fail tests with visibility into where CLI coverage is missing compared to the upstream API surface.
- Enables prioritization for future command/test additions.

Alternatives considered:
- Coverage tracking via narrative docs only: rejected due to drift and lack of enforceable signal.

## Risks / Trade-offs

- [Large initial test volume] -> Mitigation: implement by command group with shared helpers first, then expand incrementally.
- [Brittle output assertions for styled output] -> Mitigation: assert semantic content and stable fields; avoid overly strict ANSI snapshots unless explicitly needed.
- [Schema/source mismatch for model checks] -> Mitigation: pin schema/examples with explicit update workflow and version metadata.
- [Upstream OpenAPI instability] -> Mitigation: treat `https://api.tailscale.com/api/v2?outputOpenapiSchema=true` as a fetch source only, validate/parsing against a pinned snapshot, and gate updates through review.
- [Maintenance overhead for endpoint fixtures] -> Mitigation: central fixture helpers and reusable response builders.
- [False positives in coverage-gap mapping] -> Mitigation: include deterministic mapping rules plus an explicit "unmapped/needs-review" bucket in the generated report.

## Migration Plan

1. Add shared CLI test harness and mock API server helpers.
2. Add command coverage manifest/discovery check and seed with existing command tree.
3. Add tests by top-level command groups, beginning with high-traffic groups (`get`, `list`, `set`) and then completing all groups.
4. Add model contract fixture/schema validation tests and pin source snapshot.
5. Add OpenAPI-driven CLI/API coverage-gap report generation and publish report artifacts in CI output.
6. Add CI commands/targets for unit and integration tests.
6. Rollback strategy: if failures are excessive, keep harness and scope failing groups behind temporary skip lists while fixing incrementally (without removing test architecture).

## Open Questions

- What is the preferred canonical upstream schema source for model contracts (official OpenAPI document path vs curated API examples)?
- Should integration tests run in default `go test ./...` or via explicit tag/target to keep runtime bounded?
- Should styled (`pretty`) output be validated as plain-text semantic output only, or with optional golden snapshots per command?
