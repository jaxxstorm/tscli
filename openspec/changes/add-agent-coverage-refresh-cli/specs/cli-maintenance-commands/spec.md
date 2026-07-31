## ADDED Requirements

### Requirement: Maintenance command group is local and additive
The CLI SHALL expose a `maintenance` command group for repository maintenance workflows, and commands in that group MUST NOT require Tailscale API credentials.

#### Scenario: Maintenance help without credentials
- **WHEN** a user runs `tscli maintenance --help` without API key configuration
- **THEN** the CLI SHALL display maintenance command help without performing API authentication preflight

#### Scenario: Existing commands remain unchanged
- **WHEN** existing non-maintenance commands are invoked
- **THEN** their command paths, flags, authentication behavior, and output format behavior SHALL remain unchanged

### Requirement: Coverage-gap command mirrors supported tooling
The CLI SHALL expose `tscli maintenance coverage-gaps` for running the same coverage-gap workflow as the supported coverage-gap tooling.

#### Scenario: Generate coverage artifacts
- **WHEN** a user runs `tscli maintenance coverage-gaps` with default flags from the repository root
- **THEN** the CLI SHALL read the pinned OpenAPI schema, command-operation mapping, leaf command manifest, exclusions, and property coverage files and write the JSON and Markdown coverage-gap artifacts

#### Scenario: Strict coverage check
- **WHEN** a user runs `tscli maintenance coverage-gaps --check`
- **THEN** the CLI SHALL compare against the baseline, write the diff artifact, and exit non-zero when regressions or remaining gaps are present

#### Scenario: Override coverage paths
- **WHEN** a user passes coverage-gap path flags
- **THEN** the CLI SHALL use those paths instead of repository-relative defaults

### Requirement: OpenAPI refresh command mirrors supported tooling
The CLI SHALL expose `tscli maintenance openapi-refresh` for running the same OpenAPI snapshot refresh workflow as the supported refresh tooling.

#### Scenario: Refresh pinned OpenAPI snapshot
- **WHEN** a user runs `tscli maintenance openapi-refresh`
- **THEN** the CLI SHALL fetch the schema from the configured source URL and atomically rewrite the pinned schema and snapshot metadata files

#### Scenario: Override refresh paths and source
- **WHEN** a user passes `--source-url`, `--schema-out`, or `--metadata-out`
- **THEN** the CLI SHALL use those values instead of repository-relative defaults

#### Scenario: Refresh failure
- **WHEN** the source URL cannot be fetched or the response is not a valid OpenAPI schema
- **THEN** the CLI SHALL exit non-zero and leave existing pinned artifacts intact
