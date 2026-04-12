## MODIFIED Requirements

### Requirement: Output and exit semantics are stable
Command tests MUST validate output behavior and error semantics for supported output modes so scripts can depend on stable behavior.

#### Scenario: Successful command execution
- **WHEN** a command succeeds in mock-backed execution
- **THEN** tests MUST assert expected stdout structure/content and empty stderr

#### Scenario: Command execution error
- **WHEN** a command fails due to API or validation error
- **THEN** tests MUST assert non-zero error behavior and error details written to stderr without ambiguous output

#### Scenario: Pretty and human rendering for nested responses is covered
- **WHEN** a command returns nested collection or object payloads whose pretty or human rendering differs from raw JSON shape
- **THEN** automated tests MUST assert readable rendered structure in `pretty` and `human` modes
- **AND** the tests MUST catch record-collapsing, field duplication, and sibling-value leakage across rendered items

#### Scenario: Service command output regressions are detected
- **WHEN** `list services` or `get service` output handling changes
- **THEN** automated tests MUST verify supported output modes for representative service payloads including addresses, ports, tags, comments, and annotations
- **AND** tests MUST fail if service records are flattened into an unreadable pretty or human rendering
