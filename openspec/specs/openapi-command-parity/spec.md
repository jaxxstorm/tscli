# openapi-command-parity Specification

## Purpose
TBD - created by archiving change fill-api-coverage-gaps. Update Purpose after archive.
## Requirements
### Requirement: OpenAPI operation parity for CLI commands
The CLI SHALL provide command coverage for every in-scope Tailscale OpenAPI operation from the pinned schema snapshot.

#### Scenario: In-scope operation inventory is generated
- **WHEN** parity analysis runs against the pinned OpenAPI snapshot
- **THEN** each operation SHALL be classified as mapped-to-command, excluded-by-policy, or missing-command

#### Scenario: Missing operation is implemented
- **WHEN** an in-scope operation is classified as missing-command
- **THEN** a CLI command SHALL be added that invokes that operation with validated flags and structured output/error handling

### Requirement: Command-to-operation mapping is explicit
Each implemented CLI command SHALL declare at least one OpenAPI operation mapping in the command-operation map.

#### Scenario: Command mapping is missing
- **WHEN** a CLI command exists without an operation mapping
- **THEN** parity checks SHALL fail and report the unmapped command path

#### Scenario: Operation mapping is stale
- **WHEN** a mapped operation does not exist in the pinned OpenAPI snapshot
- **THEN** parity checks SHALL fail and report the stale mapping entry

### Requirement: Parity coverage includes key API domains
Parity implementation SHALL cover currently uncovered operations across device, invites, DNS, logging, keys, policy, posture integrations, users, contacts, webhooks, services, and settings domains unless explicitly excluded.

#### Scenario: Domain gap report is generated
- **WHEN** parity analysis completes
- **THEN** the report SHALL include uncovered operation counts grouped by API domain

