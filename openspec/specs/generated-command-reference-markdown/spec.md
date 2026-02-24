## ADDED Requirements

### Requirement: Command reference Markdown is auto-generated from Cobra commands
The project SHALL generate command reference Markdown files from the current Cobra command tree, with one output page per command path.

#### Scenario: Command docs generation runs
- **WHEN** the command-doc generation workflow is executed
- **THEN** Markdown files for all supported `tscli` commands SHALL be written to the documentation output directory

#### Scenario: Generated docs include command usage details
- **WHEN** a command page is generated
- **THEN** it SHALL include usage, synopsis, and flags as defined by the Cobra command metadata

### Requirement: Generated command docs stay synchronized with CLI changes
The project SHALL include a verification mechanism that detects stale command docs after command tree changes.

#### Scenario: Command changes without doc regeneration
- **WHEN** CLI commands or flags change and generated docs are not refreshed
- **THEN** docs verification SHALL fail with actionable output indicating regeneration is required

#### Scenario: Command docs are regenerated after changes
- **WHEN** generation is re-run after command updates
- **THEN** verification SHALL pass with no stale-doc failures
