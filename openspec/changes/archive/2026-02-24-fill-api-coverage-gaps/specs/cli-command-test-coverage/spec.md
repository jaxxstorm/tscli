## MODIFIED Requirements

### Requirement: OpenAPI coverage-gap reporting is generated
The project MUST generate and enforce a coverage-gap report that compares CLI/test coverage against the Tailscale OpenAPI operation surface and blocks parity regressions.

#### Scenario: Coverage-gap report generation
- **WHEN** coverage analysis is executed
- **THEN** the report MUST include totals and lists for covered operations, uncovered operations, unmapped commands, unknown mapped operations, and unknown mapped commands

#### Scenario: Gap report is actionable in CI
- **WHEN** CI runs command coverage checks
- **THEN** CI MUST fail when uncovered in-scope operations or unknown mapped commands are present and MUST publish report artifacts for reviewer triage

#### Scenario: Baseline diff is enforced
- **WHEN** a change introduces new uncovered operations or new unmapped command regressions relative to baseline
- **THEN** coverage checks MUST fail and produce a baseline-diff report with the new regressions
