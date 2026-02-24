## ADDED Requirements

### Requirement: Multi-tailnet profile schema is supported
The CLI configuration SHALL support an `active-tailnet` key and a `tailnets` collection of profile objects with `name` and `api-key` fields.

#### Scenario: Config with profiles is loaded
- **WHEN** the config file contains `active-tailnet` and at least one matching profile in `tailnets`
- **THEN** the CLI SHALL treat that profile as available for runtime credential resolution

#### Scenario: Duplicate profile names are provided
- **WHEN** two profile entries use the same `name`
- **THEN** profile validation SHALL fail with an actionable duplicate-name error

### Requirement: Runtime credential resolution follows deterministic precedence
For operational commands, effective values SHALL be resolved in this order: flags, environment variables, active profile, legacy flat config keys.

#### Scenario: Flags override profile values
- **WHEN** `--api-key` or `--tailnet` is provided and profile data exists
- **THEN** the CLI SHALL use flag values as the effective runtime values

#### Scenario: Environment overrides profile values
- **WHEN** `TAILSCALE_API_KEY` or `TAILSCALE_TAILNET` is set and no corresponding flag is set
- **THEN** the CLI SHALL use environment values as the effective runtime values

#### Scenario: Active profile is used by default
- **WHEN** flags and environment variables are absent and `active-tailnet` maps to an existing profile
- **THEN** the CLI SHALL use that profile's `name` and `api-key` as effective runtime values

#### Scenario: API key is missing from all sources
- **WHEN** flags, environment, active profile, and legacy config do not provide an API key
- **THEN** command execution SHALL fail with the existing required API key error

### Requirement: Legacy single-tailnet configuration remains valid
Existing config files that only define `tailnet` and `api-key` SHALL continue to work without modification.

#### Scenario: Legacy-only config executes commands
- **WHEN** `tailnets` and `active-tailnet` are not set and legacy keys are present
- **THEN** the CLI SHALL resolve runtime credentials from legacy keys exactly as before

#### Scenario: Active profile and legacy values coexist
- **WHEN** profile keys and legacy keys are both present
- **THEN** the CLI SHALL prefer profile-derived values over legacy keys when flags/env are not set

### Requirement: Config commands manage tailnet profiles
The `config` command group SHALL provide operations to list profiles, set active profile, upsert profile credentials, and remove profiles.

#### Scenario: List profiles displays active selection
- **WHEN** multiple profiles exist and one is active
- **THEN** the list output SHALL include all profile names and indicate which profile is active

#### Scenario: Set active profile validates existence
- **WHEN** a user selects an active profile name that does not exist in `tailnets`
- **THEN** the command SHALL fail with a validation error that names the missing profile

#### Scenario: Upsert profile writes persistent config
- **WHEN** a user creates or updates a profile with `name` and `api-key`
- **THEN** the profile entry SHALL be persisted to config and available on the next command execution

#### Scenario: Remove active profile is blocked
- **WHEN** a user attempts to remove the currently active profile without first selecting another profile
- **THEN** the command SHALL fail with guidance to switch active profile first

### Requirement: Profile behavior is covered by automated tests
The project SHALL include unit and integration tests for profile parsing, validation, runtime resolution, and profile command behavior.

#### Scenario: Unit tests cover resolver precedence
- **WHEN** automated unit tests execute for config resolution
- **THEN** tests SHALL verify precedence across flags, environment variables, active profile, and legacy keys

#### Scenario: Integration tests cover profile command flows
- **WHEN** command-level integration tests execute for config profile operations
- **THEN** tests SHALL verify success output, validation failures, and persisted state transitions
