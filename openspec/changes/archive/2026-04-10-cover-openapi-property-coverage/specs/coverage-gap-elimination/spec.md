## MODIFIED Requirements

### Requirement: Coverage-gap reports enforce parity completion
Coverage-gap reporting SHALL enforce zero missing in-scope OpenAPI operations, zero unknown mapped commands, and zero uncovered in-scope request/response properties before the change is considered complete.

#### Scenario: Uncovered in-scope operation remains
- **WHEN** the coverage-gap report contains one or more uncovered in-scope operations
- **THEN** parity checks SHALL fail and list each missing operation

#### Scenario: Unknown mapped command remains
- **WHEN** the coverage-gap report contains one or more command mappings not present in the leaf command manifest
- **THEN** parity checks SHALL fail and list each unknown mapped command

#### Scenario: Uncovered in-scope property remains
- **WHEN** the coverage-gap report contains one or more uncovered in-scope request or response properties for mapped operations
- **THEN** parity checks SHALL fail and list each missing property path

### Requirement: Exclusions are explicit and reviewable
Operations or properties intentionally not implemented SHALL be declared in an exclusions policy with rationale.

#### Scenario: Operation excluded from parity target
- **WHEN** an operation is excluded from parity completion
- **THEN** the exclusion policy SHALL include operation identifier, reason, and migration/follow-up note

#### Scenario: Property excluded from parity target
- **WHEN** a request or response property is excluded from property coverage completion
- **THEN** the exclusion policy SHALL include the property identifier, reason, and migration/follow-up note

### Requirement: CI regression checks are strict
CI SHALL run parity checks and fail when new uncovered operations, uncovered properties, or unmapped command regressions are introduced.

#### Scenario: Regression introduced in pull request
- **WHEN** a pull request increases uncovered in-scope operations, uncovered in-scope properties, or unmapped commands relative to baseline
- **THEN** CI SHALL fail and publish diff artifacts describing the regression
