## MODIFIED Requirements

### Requirement: Contract source updates are explicit and versioned
The project MUST track contract source version/metadata so updates are intentional, reviewable, and reproducible through a supported make-driven refresh workflow.

#### Scenario: Contract snapshot is updated
- **WHEN** schema/examples are refreshed from upstream source
- **THEN** repository changes MUST include updated source metadata, the refreshed pinned snapshot, and passing contract tests

#### Scenario: Refresh workflow is invoked
- **WHEN** a maintainer runs the supported OpenAPI refresh make target
- **THEN** the workflow MUST fetch the schema from the canonical Tailscale OpenAPI source and rewrite the in-repo snapshot metadata in the same operation

### Requirement: OpenAPI source URL and snapshot metadata are tracked
The project MUST source schema updates from `https://api.tailscale.com/api/v2?outputOpenapiSchema=true` and track snapshot metadata in-repo.

#### Scenario: Snapshot metadata is recorded
- **WHEN** a schema snapshot is generated or refreshed
- **THEN** source URL, fetch timestamp, and schema/version identifiers MUST be recorded with the snapshot

#### Scenario: Refresh target publishes reviewable provenance
- **WHEN** the supported OpenAPI refresh make target completes successfully
- **THEN** the repository MUST contain enough metadata to identify the fetched upstream schema revision and when it was retrieved

#### Scenario: Upstream schema instability is isolated
- **WHEN** upstream OpenAPI changes in incompatible ways
- **THEN** contract tests MUST continue to run against the pinned snapshot until a reviewed snapshot update is merged
