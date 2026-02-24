## 1. Build shared test infrastructure

- [x] 1.1 Add a reusable Cobra command test harness for args/env/stdin setup and stdout/stderr/error assertions.
- [x] 1.2 Add shared helpers for asserting output modes (`pretty`, `human`, `json`, `yaml`) without brittle formatting checks.
- [x] 1.3 Add common fixtures/utilities for config precedence tests (flags vs env vs config file).

## 2. Add mocked Tailscale API integration harness

- [x] 2.1 Implement an `httptest`-based mock API server helper with route registration for method/path assertions.
- [x] 2.2 Add fixture builders for success and error responses used across command groups.
- [x] 2.3 Ensure command tests can force API base URL to mock endpoints and fail if non-mock hosts are contacted.

## 3. Establish command coverage guardrails

- [x] 3.1 Add a command discovery test that enumerates leaf commands from Cobra root.
- [x] 3.2 Add a coverage mapping/check ensuring every leaf command is represented by at least one test case.
- [x] 3.3 Add clear failure reporting for uncovered command paths.

## 4. Implement command group test suites

- [x] 4.1 Add/expand tests for `get` commands covering required flags, success paths, and API error paths.
- [x] 4.2 Add/expand tests for `list` commands covering filtering/selection flags and output behavior.
- [x] 4.3 Add/expand tests for `create` commands covering payload construction, validation, and error handling.
- [x] 4.4 Add/expand tests for `set` commands covering mutation payloads, validation, and response handling.
- [x] 4.5 Add/expand tests for `delete` commands covering target validation, confirmation/error semantics, and output behavior.
- [x] 4.6 Add tests for `config` and `version` commands covering local behavior and output expectations.

## 5. Add model contract consistency tests

- [x] 5.1 Fetch the Tailscale OpenAPI schema from `https://api.tailscale.com/api/v2?outputOpenapiSchema=true`, pin a deterministic snapshot in-repo, and record snapshot metadata.
- [x] 5.2 Add tests validating request/response model encoding and decoding against pinned contract data.
- [x] 5.3 Add mismatch/failure-path tests for incompatible payload shapes to ensure surfaced command errors.
- [x] 5.4 Document contract fixture update workflow and source metadata versioning.

## 6. Generate CLI/API coverage-gap reporting

- [x] 6.1 Implement mapping logic from Cobra command/test coverage to OpenAPI operations.
- [x] 6.2 Generate machine-readable and human-readable coverage-gap reports with covered/uncovered/unmapped operation lists.
- [x] 6.3 Add CI checks/artifacts for coverage-gap reports and baseline-diff visibility.

## 7. Integrate into developer and CI workflows

- [x] 7.1 Add stable `go test` invocation targets for unit and mock-backed integration suites.
- [x] 7.2 Update README/development docs with test strategy, OpenAPI snapshot refresh workflow, and coverage-gap report usage.
- [x] 7.3 Run full test suites, fix regressions, and verify no user-facing command/flag behavior changed unintentionally.
