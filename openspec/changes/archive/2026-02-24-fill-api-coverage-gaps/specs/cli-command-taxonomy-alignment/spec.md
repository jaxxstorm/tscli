## ADDED Requirements

### Requirement: Canonical command verb taxonomy
The CLI SHALL use canonical verbs based on operation intent: `create` for creation, `set` for updates/mutations, `get` for single-resource retrieval, `list` for multi-resource retrieval, and `delete` for deletion.

#### Scenario: New creation endpoint command
- **WHEN** a command is added for an endpoint that creates a resource
- **THEN** the command SHALL be exposed under the `create` verb tree

#### Scenario: New update endpoint command
- **WHEN** a command is added for an endpoint that mutates existing state
- **THEN** the command SHALL be exposed under the `set` verb tree

#### Scenario: New retrieval endpoint command
- **WHEN** a command is added for a retrieval endpoint
- **THEN** the command SHALL use `get` for single-resource responses and `list` for collection responses

### Requirement: Canonical names are non-abbreviated and consistent
Canonical command paths SHALL use full nouns (for example `nameservers`, `preferences`) rather than abbreviated primary command names.

#### Scenario: Existing abbreviated command path
- **WHEN** an existing command path uses an abbreviation as the primary canonical path
- **THEN** the canonical path SHALL be normalized and the abbreviated path SHALL remain available as a backward-compatible alias

### Requirement: Taxonomy changes preserve script compatibility
When canonical command paths change for taxonomy alignment, aliases and help text SHALL preserve backward compatibility and migration discoverability.

#### Scenario: Legacy command invocation
- **WHEN** a user invokes a legacy command alias
- **THEN** the command SHALL continue to execute successfully and indicate the canonical equivalent in help or command descriptions
