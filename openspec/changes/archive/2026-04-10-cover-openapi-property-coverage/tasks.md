## 1. Define Property Coverage Inputs

- [x] 1.1 Add repository data files for property coverage declarations and property-level exclusions tied to mapped OpenAPI operations.
- [x] 1.2 Seed the initial property coverage inventory for currently mapped operations, including a concrete device-response example such as `devices[].postureIdentity`.

## 2. Extend Coverage Analysis

- [x] 2.1 Update `coverage/coveragegaps` to derive request/response property paths from the pinned OpenAPI snapshot for mapped operations.
- [x] 2.2 Extend the machine-readable and markdown coverage reports to include covered, excluded, uncovered, and regressed property paths alongside existing operation parity output.
- [x] 2.3 Update baseline diff and fail-on-gap behavior so CI can fail on uncovered property regressions as well as operation regressions.

## 3. Strengthen Contract And Command Evidence

- [x] 3.1 Add or expand contract tests that validate representative request/response property paths against the pinned schema at the property level.
- [x] 3.2 Add or expand mock-backed integration tests for representative commands, including device output coverage that exercises newly tracked properties through decode or structured output.

## 4. Finalize Workflow Integration

- [x] 4.1 Update make/CI wiring and any related coverage documentation so property coverage checks run through the supported workflow.
- [x] 4.2 Run the focused coverage, contract, and integration test suites plus the enhanced coverage report generation, then resolve remaining uncovered-property gaps or policy-backed exclusions.
