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

#### Scenario: Existing DNS mutation command matches supported nameserver formats
- **WHEN** the upstream DNS nameserver API accepts DNS-over-HTTPS endpoint addresses in addition to literal IP nameservers
- **THEN** the existing `set dns nameservers` and `set dns split-dns` commands SHALL accept those same supported nameserver value formats without requiring a raw API fallback

### Requirement: Command-to-operation mapping is explicit
Each implemented CLI command SHALL declare at least one OpenAPI operation mapping in the command-operation map.

#### Scenario: Command mapping is missing
- **WHEN** a CLI command exists without an operation mapping
- **THEN** parity checks SHALL fail and report the unmapped command path

#### Scenario: Operation mapping is stale
- **WHEN** a mapped operation does not exist in the pinned OpenAPI snapshot
- **THEN** parity checks SHALL fail and report the stale mapping entry

### Requirement: User deletion parity covers bulk deletion workflow
The CLI SHALL provide parity for user deletion workflows by exposing both single-user deletion and filter-driven bulk user deletion commands backed by the in-scope user deletion operations.

#### Scenario: Bulk user deletion is represented in parity coverage
- **WHEN** parity analysis evaluates uncovered user-domain operations and workflows
- **THEN** the `delete users` command SHALL be treated as the supported CLI path for bulk user cleanup based on user-list filtering plus per-user deletion requests

#### Scenario: Bulk user deletion command is mapped explicitly
- **WHEN** the `tscli delete users` command is implemented
- **THEN** command-to-operation mapping and parity checks SHALL include the command and its associated user-domain operations

### Requirement: Parity coverage includes key API domains
Parity implementation SHALL cover currently uncovered operations across device, invites, DNS, logging, keys, policy, posture integrations, users, contacts, webhooks, services, and settings domains unless explicitly excluded.

#### Scenario: Domain gap report is generated
- **WHEN** parity analysis completes
- **THEN** the report SHALL include uncovered operation counts grouped by API domain
