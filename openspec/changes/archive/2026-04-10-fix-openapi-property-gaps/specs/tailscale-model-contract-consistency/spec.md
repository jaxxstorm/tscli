## ADDED Requirements

### Requirement: API-documented response properties are never silently dropped
For mapped operations, if the pinned API schema defines a response property for structured command output, the CLI MUST preserve and show that property even when the upstream Go client model omits or mis-tags it.

#### Scenario: Upstream SDK omits or mis-tags a response property
- **WHEN** a mapped response property exists in the pinned schema but is absent or mis-tagged in the upstream client model
- **THEN** `tscli` MUST decode through a schema-aligned local model or adapter so the property remains present in structured command output and coverage evidence

#### Scenario: Command previously returned a synthetic summary
- **WHEN** a command currently prints an echoed request or synthetic success object for an operation whose API response body contains documented properties
- **THEN** the command MUST decode and print the authoritative response body or an equivalent schema-complete representation that includes those properties

### Requirement: Audited write operations use schema-aligned request models
The project MUST serialize audited write-command payloads from request models whose property names match the pinned schema and whose key properties are asserted in command-level tests.

#### Scenario: Audited write operation is exercised
- **WHEN** a write command is placed under request property coverage
- **THEN** its outbound payload MUST use the pinned schema property names and mock-backed integration tests MUST assert the expected request fields
