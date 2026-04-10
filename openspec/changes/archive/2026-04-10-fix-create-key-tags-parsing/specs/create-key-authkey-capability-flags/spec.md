## MODIFIED Requirements

### Requirement: Auth-key capability flags are exposed on create key
The CLI SHALL provide `--reusable`, `--ephemeral`, `--preauthorized`, and `--tags` flags for `tscli create key` when creating auth keys.

#### Scenario: Create auth key with capability flags
- **WHEN** a user runs `tscli create key --type authkey` with one or more capability flags, including `--tags`
- **THEN** the command SHALL accept those flags and proceed with auth-key creation

#### Scenario: Create auth key without capability flags
- **WHEN** a user runs `tscli create key --type authkey` without capability flags
- **THEN** the command SHALL continue to work with backward-compatible behavior

### Requirement: Capability flags map to auth-key create request payload
For auth-key creation, each capability flag SHALL be mapped to the corresponding Tailscale API auth-key capability field in the outgoing request payload, including `capabilities.devices.create.tags` for `--tags`.

#### Scenario: Reusable value is mapped
- **WHEN** `--reusable` is set for an auth-key create request
- **THEN** the request payload SHALL include the reusable capability value

#### Scenario: Ephemeral value is mapped
- **WHEN** `--ephemeral` is set for an auth-key create request
- **THEN** the request payload SHALL include the ephemeral capability value

#### Scenario: Tags are mapped
- **WHEN** `--tags tag:tsdns` is set for an auth-key create request
- **THEN** the request payload SHALL include `capabilities.devices.create.tags` with `tag:tsdns` rather than `null`

#### Scenario: Preauthorized value is mapped
- **WHEN** `--preauthorized` is set for an auth-key create request
- **THEN** the request payload SHALL include the preauthorized capability value

### Requirement: Capability flag behavior is covered by automated tests
The project SHALL include integration tests for create-key capability flags, including tag serialization regression coverage, and compatibility tests for unchanged behavior.

#### Scenario: Integration test verifies capability payload
- **WHEN** test execution runs against mocked API for `create key`
- **THEN** tests SHALL assert capability booleans are present in auth-key create payload when flags are provided

#### Scenario: Integration test verifies tag payload
- **WHEN** test execution runs `create key --type authkey --tags tag:tsdns` against mocked API
- **THEN** tests SHALL assert the recorded auth-key payload contains `capabilities.devices.create.tags` with the provided tag values

#### Scenario: Compatibility test verifies unchanged default behavior
- **WHEN** tests run auth-key creation without capability flags
- **THEN** tests SHALL assert successful command behavior and no regression in existing key creation flow
