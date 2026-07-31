## ADDED Requirements

### Requirement: Coverage-gap workflow is available to agents and CLI users
Coverage-gap reporting SHALL be runnable through repository agent skills and a local CLI maintenance command in addition to existing make targets.

#### Scenario: Agent invokes coverage-gap workflow
- **WHEN** an agent is preparing to close uncovered OpenAPI operations, unmapped commands, unknown mappings, or uncovered properties
- **THEN** the agent SHALL have local skill instructions that identify the supported coverage-gap workflow and generated artifacts

#### Scenario: CLI user invokes coverage-gap workflow
- **WHEN** a maintainer runs the CLI coverage-gap maintenance command
- **THEN** it SHALL generate the same coverage-gap artifacts and enforce the same strict check behavior as the supported coverage-gap tooling
