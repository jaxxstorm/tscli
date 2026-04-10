## MODIFIED Requirements

### Requirement: Multi-tailnet profile schema is supported
The CLI configuration SHALL support an `active-tailnet` key and a `tailnets` collection of profile objects with `name` and `api-key` fields. When profile-backed configuration is persisted, that profile schema SHALL be the canonical stored representation and SHALL NOT require duplicate top-level `tailnet` or `api-key` keys.

#### Scenario: Config with profiles is loaded
- **WHEN** the config file contains `active-tailnet` and at least one matching profile in `tailnets`
- **THEN** the CLI SHALL treat that profile as available for runtime credential resolution

#### Scenario: Duplicate profile names are provided
- **WHEN** two profile entries use the same `name`
- **THEN** profile validation SHALL fail with an actionable duplicate-name error

#### Scenario: Profile config is persisted canonically
- **WHEN** a config write occurs for a file that contains profile data
- **THEN** the persisted config SHALL store profile state using `active-tailnet` and `tailnets`
- **AND** the persisted config SHALL NOT depend on duplicated top-level `tailnet` or `api-key` keys to represent the active profile

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

#### Scenario: Active profile switch changes effective runtime tailnet
- **WHEN** a user changes `active-tailnet` to another persisted profile and then runs an operational command without overriding flags or environment variables
- **THEN** the CLI SHALL resolve the effective `tailnet` and `api-key` from the newly active profile

#### Scenario: API key is missing from all sources
- **WHEN** flags, environment, active profile, and legacy config do not provide an API key
- **THEN** command execution SHALL fail with the existing required API key error

### Requirement: Legacy single-tailnet configuration remains valid
Existing config files that only define `tailnet` and `api-key` SHALL continue to work without modification. Once profile data exists, those legacy flat keys SHALL be treated as backward-compatibility input only rather than canonical persisted profile state.

#### Scenario: Legacy-only config executes commands
- **WHEN** `tailnets` and `active-tailnet` are not set and legacy keys are present
- **THEN** the CLI SHALL resolve runtime credentials from legacy keys exactly as before

#### Scenario: Active profile and legacy values coexist
- **WHEN** profile keys and legacy keys are both present
- **THEN** the CLI SHALL prefer profile-derived values over legacy keys when flags and environment variables are not set

#### Scenario: Profile rewrite removes duplicated legacy mirrors
- **WHEN** a profile command rewrites a mixed config file that contains both profile data and flat legacy keys copied from the active profile
- **THEN** the rewritten config SHALL preserve the profile data and active selection
- **AND** the rewritten config SHALL NOT re-persist the duplicated flat `tailnet` and `api-key` values as part of the canonical profile representation

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

#### Scenario: Profile mutations keep canonical file shape
- **WHEN** a user runs `config profiles upsert` or `config profiles set-active`
- **THEN** the saved config SHALL retain the selected `active-tailnet` and `tailnets` entries
- **AND** the saved config SHALL NOT add or refresh duplicated top-level `tailnet` and `api-key` keys for the active profile

#### Scenario: Remove active profile is blocked
- **WHEN** a user attempts to remove the currently active profile without first selecting another profile
- **THEN** the command SHALL fail with guidance to switch active profile first

### Requirement: Profile behavior is covered by automated tests
The project SHALL include unit and integration tests for profile parsing, validation, runtime resolution, and profile command behavior.

#### Scenario: Unit tests cover resolver precedence
- **WHEN** automated unit tests execute for config resolution
- **THEN** tests SHALL verify precedence across flags, environment variables, active profile, and legacy keys

#### Scenario: Unit tests cover canonical persistence
- **WHEN** automated config persistence tests execute for profile-backed config
- **THEN** tests SHALL verify that saved files keep `active-tailnet` and `tailnets`
- **AND** tests SHALL verify that duplicated flat `tailnet` and `api-key` values are not reintroduced for profile-backed config

#### Scenario: Integration tests cover profile command flows
- **WHEN** command-level integration tests execute for config profile operations
- **THEN** tests SHALL verify success output, validation failures, persisted state transitions, and effective runtime behavior after switching the active profile
