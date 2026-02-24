## Purpose

Define required automated coverage expectations for the `tscli` command surface, including validation, output behavior, configuration precedence, and OpenAPI coverage reporting.

## Requirements

### Requirement: Complete command-surface test coverage
The project MUST maintain automated tests for every leaf CLI command registered under `tscli`, including command construction, argument parsing, and execution flow.

#### Scenario: Every leaf command is covered
- **WHEN** the test suite enumerates the Cobra command tree
- **THEN** each leaf command MUST map to at least one automated test case in the coverage manifest/check

#### Scenario: New command without tests fails verification
- **WHEN** a new leaf command is added without corresponding test mapping
- **THEN** the coverage verification test MUST fail with the uncovered command path

### Requirement: Command validation behavior is tested
Each command with required inputs MUST have tests that verify missing or invalid input handling and actionable error messaging.

#### Scenario: Required flag is missing
- **WHEN** a command is executed without a required flag or argument
- **THEN** the command test MUST assert a non-nil execution error and an error message that identifies the missing input

#### Scenario: Invalid flag value is provided
- **WHEN** a command is executed with an invalid value (for example malformed route, tag, or id format where validation exists)
- **THEN** the command test MUST assert failure behavior and the expected validation error

### Requirement: Output and exit semantics are stable
Command tests MUST validate output behavior and error semantics for supported output modes so scripts can depend on stable behavior.

#### Scenario: Successful command execution
- **WHEN** a command succeeds in mock-backed execution
- **THEN** tests MUST assert expected stdout structure/content and empty stderr

#### Scenario: Command execution error
- **WHEN** a command fails due to API or validation error
- **THEN** tests MUST assert non-zero error behavior and error details written to stderr without ambiguous output

### Requirement: Configuration precedence is verified
Automated tests MUST verify command configuration precedence of flags over environment variables over config files.

#### Scenario: Flag overrides environment and config file
- **WHEN** the same config key is set via flag, environment variable, and config file
- **THEN** the command MUST use the flag value

#### Scenario: Environment overrides config file
- **WHEN** the same config key is set via environment variable and config file without a flag
- **THEN** the command MUST use the environment value

### Requirement: OpenAPI coverage-gap reporting is generated
The project MUST generate a coverage-gap report that compares CLI/test coverage against the Tailscale OpenAPI operation surface.

#### Scenario: Coverage-gap report generation
- **WHEN** coverage analysis is executed
- **THEN** the report MUST include totals and lists for covered operations, uncovered operations, and unmapped operations needing manual review

#### Scenario: Gap report is actionable in CI
- **WHEN** CI runs command coverage checks
- **THEN** the generated coverage-gap artifact MUST be available to reviewers and MUST identify newly uncovered operations relative to the pinned baseline
