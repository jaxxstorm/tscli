## Purpose

Define contract-consistency requirements so `tscli` request/response handling stays aligned with a pinned Tailscale OpenAPI source.

## Requirements

### Requirement: Contract validation against pinned upstream API schema/examples
The project MUST validate CLI request/response model handling against a pinned snapshot of official Tailscale API schema or authoritative API examples.

#### Scenario: Contract snapshot is available
- **WHEN** contract tests run
- **THEN** tests MUST validate key model fields and payload structures against the pinned schema/examples

#### Scenario: Contract mismatch is introduced
- **WHEN** code changes introduce request/response model drift from the pinned contract
- **THEN** contract tests MUST fail with the mismatched field or shape details

### Requirement: Contract source updates are explicit and versioned
The project MUST track contract source version/metadata so updates are intentional and reviewable.

#### Scenario: Contract snapshot is updated
- **WHEN** schema/examples are refreshed from upstream source
- **THEN** repository changes MUST include updated source metadata and passing contract tests

### Requirement: Command-level model decoding failures are covered
Integration tests MUST verify that invalid or unexpected API payload shapes produce clear command errors instead of silent corruption.

#### Scenario: Unexpected payload field shape
- **WHEN** a mock API response returns incompatible field types for a command model
- **THEN** the command MUST return a surfaced error and integration tests MUST assert the failure path

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

#### Scenario: Upstream schema instability is isolated
- **WHEN** upstream OpenAPI changes in incompatible ways
- **THEN** contract tests MUST continue to run against the pinned snapshot until a reviewed snapshot update is merged
