## ADDED Requirements

### Requirement: OpenAPI refresh workflow is available to agents and CLI users
OpenAPI snapshot refresh SHALL be runnable through repository agent skills and a local CLI maintenance command in addition to existing make targets.

#### Scenario: Agent invokes OpenAPI refresh workflow
- **WHEN** an agent is preparing schema-driven implementation from the latest upstream Tailscale OpenAPI surface
- **THEN** the agent SHALL have local skill instructions that identify the supported refresh workflow and the pinned schema and metadata artifacts to inspect

#### Scenario: CLI user invokes OpenAPI refresh workflow
- **WHEN** a maintainer runs the CLI OpenAPI refresh maintenance command
- **THEN** it SHALL fetch from the configured OpenAPI source URL and rewrite the pinned schema and metadata with the same validation and atomic replacement behavior as the supported refresh tooling
