## 1. Shared contract models

- [x] 1.1 Inventory the currently mapped request and response sides that still rely on explicit or default property exclusions, grouped by shared schema family
- [x] 1.2 Add schema-aligned local DTOs and shared helpers for audited request/response bodies that the upstream SDK or current command output cannot represent without dropping documented properties

## 2. Command remediation

- [x] 2.1 Switch device and route read commands such as `list devices`, `get device`, and `list routes` to schema-aligned response decoding so documented API properties remain present in structured output
- [x] 2.2 Update audited write commands that currently echo requests or synthetic success objects to decode and print the authoritative API response body when the schema defines one
- [x] 2.3 Update audited write-command request builders so serialized property names match the pinned schema and covered request fields are not silently dropped

## 3. Coverage data and tests

- [x] 3.1 Expand `coverage/property-coverage.yaml` for each remediated request/response side and remove the matching explicit or default exclusions from `coverage/property-exclusions.yaml`
- [x] 3.2 Add representative mock fixtures and CLI integration tests that assert preserved response properties and expected request property names for each remediated command family
- [x] 3.3 Extend pinned-schema contract tests for shared audited models and update command docs/help text where structured output changes from synthetic summaries to authoritative API responses

## 4. Verification

- [x] 4.1 Run the targeted Go test suites for touched commands plus `coverage/coveragegaps` and `pkg/contract/openapi`, then fix any regressions
- [x] 4.2 Run the coverage-gap report with regression and gap failures enabled and confirm the selected property exclusions are eliminated without introducing new gaps
