## Purpose

Define contract-consistency requirements so `tscli` request/response handling stays aligned with a pinned Tailscale OpenAPI source.

## Requirements

### Requirement: Contract validation against pinned upstream API schema/examples
The project MUST validate CLI request/response model handling against a pinned snapshot of official Tailscale API schema or authoritative API examples at the property level, not only at the endpoint or whole-payload shape level.

#### Scenario: Contract snapshot is available
- **WHEN** contract tests run
- **THEN** tests MUST validate key request/response property paths and payload structures against the pinned schema/examples

#### Scenario: Contract mismatch is introduced
- **WHEN** code changes introduce request/response model drift from the pinned contract
- **THEN** contract tests MUST fail with the mismatched property or shape details

#### Scenario: Pinned schema adds a response property
- **WHEN** the pinned schema introduces a new response property for a mapped command model such as `devices[].postureIdentity`
- **THEN** contract validation or coverage checks MUST surface that property until the CLI coverage declaration and representative tests are updated

### Requirement: Contract source updates are explicit and versioned
The project MUST track contract source version/metadata so updates are intentional, reviewable, and reproducible through a supported make-driven refresh workflow.

#### Scenario: Contract snapshot is updated
- **WHEN** schema/examples are refreshed from upstream source
- **THEN** repository changes MUST include updated source metadata, the refreshed pinned snapshot, and passing contract tests

#### Scenario: Refresh workflow is invoked
- **WHEN** a maintainer runs the supported OpenAPI refresh make target
- **THEN** the workflow MUST fetch the schema from the canonical Tailscale OpenAPI source and rewrite the in-repo snapshot metadata in the same operation

### Requirement: Command-level model decoding failures are covered
Integration tests MUST verify that invalid or unexpected API payload shapes produce clear command errors instead of silent corruption, and representative command fixtures MUST exercise important covered response properties.

#### Scenario: Unexpected payload field shape
- **WHEN** a mock API response returns incompatible field types for a command model
- **THEN** the command MUST return a surfaced error and integration tests MUST assert the failure path

#### Scenario: Representative covered property is exercised
- **WHEN** a command claims coverage for an important response property in automated coverage data
- **THEN** contract or integration tests MUST include a fixture/assertion that exercises that property through decode or structured output

### Requirement: Contract checks remain deterministic in CI
Contract tests MUST run without runtime network fetches and MUST use in-repo pinned artifacts for deterministic CI behavior.

#### Scenario: CI executes contract tests
- **WHEN** contract tests run in CI
- **THEN** tests MUST complete using repository fixtures/schema snapshots only

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
