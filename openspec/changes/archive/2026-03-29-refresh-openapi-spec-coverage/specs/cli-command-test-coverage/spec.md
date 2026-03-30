## MODIFIED Requirements

### Requirement: OpenAPI coverage-gap reporting is generated
The project MUST generate and enforce a coverage-gap report that compares CLI/test coverage against the Tailscale OpenAPI operation surface, blocks parity regressions, and supports a make-driven workflow for validating against a freshly refreshed pinned snapshot.

#### Scenario: Coverage-gap report generation
- **WHEN** coverage analysis is executed
- **THEN** the report MUST include totals and lists for covered operations, uncovered operations, unmapped commands, unknown mapped operations, and unknown mapped commands

#### Scenario: Gap report is actionable in CI
- **WHEN** CI runs command coverage checks
- **THEN** CI MUST fail when uncovered in-scope operations or unknown mapped commands are present and MUST publish report artifacts for reviewer triage

#### Scenario: Baseline diff is enforced
- **WHEN** a change introduces new uncovered operations or new unmapped command regressions relative to baseline
- **THEN** coverage checks MUST fail and produce a baseline-diff report with the new regressions

#### Scenario: Latest snapshot coverage validation is available
- **WHEN** a maintainer runs the supported make target for latest OpenAPI coverage validation
- **THEN** the workflow MUST refresh the pinned schema snapshot first and then generate coverage-gap artifacts against that refreshed snapshot

#### Scenario: Existing coverage targets remain stable
- **WHEN** existing automation runs `make coverage-gaps` or `make coverage-gaps-check`
- **THEN** the project MUST continue to generate the same classes of coverage artifacts without requiring runtime network access
