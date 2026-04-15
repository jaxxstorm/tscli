## MODIFIED Requirements

### Requirement: Multi-tailnet profile schema is supported
The CLI configuration SHALL support an `active-tailnet` key and a `tailnets` collection of profile objects. Each profile SHALL include a `name` field and MAY include `tailnet`, `api-key`, `api-key-encrypted`, `oauth-client-id`, `oauth-client-secret`, and `oauth-client-secret-encrypted`. A valid profile SHALL contain either `api-key` or `api-key-encrypted`, or both `oauth-client-id` and either `oauth-client-secret` or `oauth-client-secret-encrypted`. A profile SHALL NOT persist both plaintext and encrypted variants of the same secret field at the same time. When profile-backed configuration is persisted, that profile schema SHALL be the canonical stored representation and SHALL NOT require duplicate top-level `tailnet` or `api-key` keys.

#### Scenario: Config with API-key profile is loaded
- **WHEN** the config file contains `active-tailnet` and a matching profile with `name` and `api-key`
- **THEN** the CLI SHALL treat that profile as available for API-key runtime credential resolution

#### Scenario: Config with encrypted API-key profile is loaded
- **WHEN** the config file contains `active-tailnet` and a matching profile with `name` and `api-key-encrypted`
- **THEN** the CLI SHALL decrypt that profile secret using the configured `age` identity during API-key runtime credential resolution

#### Scenario: Config with OAuth profile is loaded
- **WHEN** the config file contains `active-tailnet` and a matching profile with `name`, `oauth-client-id`, and `oauth-client-secret`
- **THEN** the CLI SHALL treat that profile as available for OAuth-backed command authentication

#### Scenario: Config with encrypted OAuth profile is loaded
- **WHEN** the config file contains `active-tailnet` and a matching profile with `name`, `oauth-client-id`, and `oauth-client-secret-encrypted`
- **THEN** the CLI SHALL decrypt that profile secret using the configured `age` identity during OAuth-backed command authentication

#### Scenario: Profile tailnet defaults to profile name
- **WHEN** a profile omits `tailnet`
- **THEN** commands that require a tailnet value SHALL treat the profile `name` as the effective tailnet

#### Scenario: Invalid profile auth shape is provided
- **WHEN** a profile contains neither an API-key shape nor a complete OAuth credential shape
- **THEN** profile validation SHALL fail with an actionable invalid-auth-shape error

#### Scenario: Duplicate profile names are provided
- **WHEN** two profile entries use the same `name`
- **THEN** profile validation SHALL fail with an actionable duplicate-name error

#### Scenario: Mixed plaintext and encrypted secret fields are provided
- **WHEN** a profile contains both `api-key` and `api-key-encrypted`, or both `oauth-client-secret` and `oauth-client-secret-encrypted`
- **THEN** profile validation SHALL fail with an actionable mixed-secret-storage error

#### Scenario: Profile config is persisted canonically
- **WHEN** a config write occurs for a file that contains profile data
- **THEN** the persisted config SHALL store profile state using `active-tailnet` and `tailnets`
- **AND** the persisted config SHALL preserve the relevant auth fields for each profile
- **AND** the persisted config SHALL NOT depend on duplicated top-level `tailnet` or `api-key` keys to represent the active profile

### Requirement: Config commands manage tailnet profiles
The `config` command group SHALL provide operations to list profiles, set active profile, upsert profile credentials, remove profiles, and run an interactive `config setup` flow for both API-key-backed and OAuth-backed profile entries.

#### Scenario: List profiles displays active selection and auth shape
- **WHEN** multiple profiles exist and one is active
- **THEN** the list output SHALL include all profile names
- **AND** the output SHALL indicate which profile is active
- **AND** the output SHALL indicate whether each profile is API-key-backed or OAuth-backed

#### Scenario: Set active profile validates existence
- **WHEN** a user selects an active profile name that does not exist in `tailnets`
- **THEN** the command SHALL fail with a validation error that names the missing profile

#### Scenario: Upsert API-key profile writes persistent config
- **WHEN** a user runs `config profiles set <name> --api-key <key>`
- **THEN** the CLI SHALL persist an API-key-backed profile entry
- **AND** the profile SHALL be available on the next command execution

#### Scenario: Upsert OAuth profile writes persistent config
- **WHEN** a user runs `config profiles set <name> --oauth-client-id <id> --oauth-client-secret <secret>`
- **THEN** the CLI SHALL persist an OAuth-backed profile entry
- **AND** the profile SHALL be available on the next command execution

#### Scenario: Interactive setup writes encrypted API-key profile when encryption is configured
- **WHEN** a user runs `config setup`, enables encryption, and saves an API-key-backed profile
- **THEN** the saved config SHALL persist `api-key-encrypted` for that profile instead of plaintext `api-key`

#### Scenario: Interactive setup writes encrypted OAuth profile when encryption is configured
- **WHEN** a user runs `config setup`, enables encryption, and saves an OAuth-backed profile
- **THEN** the saved config SHALL persist `oauth-client-secret-encrypted` for that profile instead of plaintext `oauth-client-secret`

#### Scenario: Profile mutations keep canonical file shape
- **WHEN** a user runs `config profiles set`, `config profiles set-active`, or `config setup`
- **THEN** the saved config SHALL retain the selected `active-tailnet` and `tailnets` entries
- **AND** the saved config SHALL NOT add or refresh duplicated top-level `tailnet` and `api-key` keys for the active profile

#### Scenario: Remove active profile is blocked
- **WHEN** a user attempts to remove the currently active profile without first selecting another profile
- **THEN** the command SHALL fail with guidance to switch active profile first

#### Scenario: Interactive setup deletes a non-active profile
- **WHEN** a user reruns `config setup`, chooses delete, and selects a non-active profile
- **THEN** the CLI SHALL remove that profile from persisted config

### Requirement: Profile behavior is covered by automated tests
The project SHALL include unit and integration tests for profile parsing, validation, runtime resolution, and profile command behavior across both API-key-backed and OAuth-backed profiles, including encrypted persisted secrets and setup-driven mutations.

#### Scenario: Unit tests cover OAuth resolver precedence
- **WHEN** automated unit tests execute for config resolution
- **THEN** tests SHALL verify precedence across OAuth flags, environment variables, active profile values, and missing-credential failures

#### Scenario: Unit tests cover canonical persistence for mixed auth profiles
- **WHEN** automated config persistence tests execute for profile-backed config
- **THEN** tests SHALL verify that saved files keep `active-tailnet` and `tailnets`
- **AND** tests SHALL verify that persisted profiles retain the relevant API-key or OAuth fields without reintroducing duplicated flat keys

#### Scenario: Unit tests cover encrypted secret field validation
- **WHEN** automated unit tests execute for profile validation and encryption helpers
- **THEN** tests SHALL verify encrypted-field acceptance, mixed plaintext/encrypted rejection, and decryption of persisted secrets through configured `age` identities

#### Scenario: Integration tests cover setup-driven profile command flows
- **WHEN** command-level integration tests execute for `config setup` and profile operations
- **THEN** tests SHALL verify first-run creation, repeated-run add and delete flows, encrypted persistence, validation failures, and successful runtime auth resolution using active profiles
