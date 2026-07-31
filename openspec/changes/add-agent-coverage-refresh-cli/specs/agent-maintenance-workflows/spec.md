## ADDED Requirements

### Requirement: Coverage-gap skill runs parity checks before implementation
The repository SHALL provide Codex and OpenCode skills that instruct agents to run coverage-gap checks before implementing OpenAPI command parity, coverage-gap elimination, or schema-driven command work.

#### Scenario: Agent starts parity implementation
- **WHEN** an agent is asked to implement missing OpenAPI-backed CLI command coverage
- **THEN** the coverage-gap skill SHALL instruct the agent to run the supported coverage-gap workflow and inspect the generated gap artifacts before editing commands

#### Scenario: Coverage check fails
- **WHEN** the supported coverage-gap check exits non-zero
- **THEN** the skill SHALL instruct the agent to treat the reported uncovered operations, unmapped commands, unknown mappings, or uncovered properties as the implementation backlog

### Requirement: OpenAPI refresh skill runs refresh before latest-schema work
The repository SHALL provide Codex and OpenCode skills that instruct agents to run the OpenAPI refresh workflow before implementing work that depends on the latest upstream Tailscale API schema.

#### Scenario: Agent starts latest OpenAPI work
- **WHEN** an agent is asked to implement CLI behavior from the latest Tailscale OpenAPI schema
- **THEN** the OpenAPI refresh skill SHALL instruct the agent to run the supported refresh workflow before deriving coverage gaps

#### Scenario: Refresh changes pinned artifacts
- **WHEN** OpenAPI refresh rewrites the pinned schema or snapshot metadata
- **THEN** the skill SHALL instruct the agent to include those changed artifacts in its review and follow-up coverage-gap analysis

### Requirement: Skills identify supported commands and artifacts
The repository skills SHALL name the supported make targets, equivalent `tscli maintenance` commands, and generated artifacts involved in coverage-gap and OpenAPI refresh workflows.

#### Scenario: Agent needs command guidance
- **WHEN** an agent reads the local maintenance workflow skill
- **THEN** the skill SHALL identify the command to run, the artifacts to inspect, and the expected next step after success or failure
