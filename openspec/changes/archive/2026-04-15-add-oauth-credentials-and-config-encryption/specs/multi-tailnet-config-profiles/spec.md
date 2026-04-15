## MODIFIED Requirements

### Requirement: Multi-tailnet profile schema is supported
The CLI configuration SHALL support an `active-tailnet` key and a `tailnets` collection of profile objects. Each profile SHALL include a `name` field and MAY include `tailnet`, `api-key`, `api-key-encrypted`, `oauth-client-id`, `oauth-client-secret`, and `oauth-client-secret-encrypted`. A valid profile SHALL contain either exactly one of `api-key` or `api-key-encrypted`, or `oauth-client-id` plus exactly one of `oauth-client-secret` or `oauth-client-secret-encrypted`. When profile-backed configuration is persisted, that profile schema SHALL be the canonical stored representation and SHALL NOT require duplicate top-level `tailnet` or `api-key` keys.

#### Scenario: Config with API-key profile is loaded
- **WHEN** the config file contains `active-tailnet` and a matching profile with `name` and `api-key`
- **THEN** the CLI SHALL treat that profile as available for API-key runtime credential resolution

#### Scenario: Config with encrypted API-key profile is loaded
- **WHEN** the config file contains `active-tailnet` and a matching profile with `name` and `api-key-encrypted`
- **THEN** the CLI SHALL treat that profile as available for API-key runtime credential resolution after decryption

#### Scenario: Config with OAuth profile is loaded
- **WHEN** the config file contains `active-tailnet` and a matching profile with `name`, `oauth-client-id`, and `oauth-client-secret`
- **THEN** the CLI SHALL treat that profile as available for OAuth-backed command authentication

#### Scenario: Config with encrypted OAuth profile is loaded
- **WHEN** the config file contains `active-tailnet` and a matching profile with `name`, `oauth-client-id`, and `oauth-client-secret-encrypted`
- **THEN** the CLI SHALL treat that profile as available for OAuth-backed command authentication after decryption

#### Scenario: Profile tailnet defaults to profile name
- **WHEN** a profile omits `tailnet`
- **THEN** commands that require a tailnet value SHALL treat the profile `name` as the effective tailnet

#### Scenario: Invalid profile auth shape is provided
- **WHEN** a profile contains neither an API key value nor a complete OAuth client id and secret pair
- **THEN** profile validation SHALL fail with an actionable invalid-auth-shape error

#### Scenario: Duplicate profile names are provided
- **WHEN** two profile entries use the same `name`
- **THEN** profile validation SHALL fail with an actionable duplicate-name error

#### Scenario: Profile config is persisted canonically
- **WHEN** a config write occurs for a file that contains profile data
- **THEN** the persisted config SHALL store profile state using `active-tailnet` and `tailnets`
- **AND** the persisted config SHALL preserve the relevant auth fields for each profile
- **AND** the persisted config SHALL NOT depend on duplicated top-level `tailnet` or `api-key` keys to represent the active profile

### Requirement: Runtime credential resolution follows deterministic precedence
For operational commands, effective values SHALL be resolved in this order: flags, environment variables, active profile, legacy flat config keys where supported by the command. API-key-authenticated commands SHALL continue to resolve `api-key` and `tailnet` using that precedence. Commands that support OAuth-backed authentication SHALL resolve `oauth-client-id` and `oauth-client-secret` using the same precedence model, exchange those credentials for a runtime access token, and SHALL NOT persist the exchanged token to config.

#### Scenario: Flags override OAuth profile values
- **WHEN** `--oauth-client-id` or `--oauth-client-secret` is provided and OAuth profile data exists
- **THEN** the CLI SHALL use the flag values as the effective OAuth credential input

#### Scenario: Environment overrides OAuth profile values
- **WHEN** OAuth flags are absent and `TSCLI_OAUTH_CLIENT_ID` or `TSCLI_OAUTH_CLIENT_SECRET` is set while OAuth profile data exists
- **THEN** the CLI SHALL use the environment values as the effective OAuth credential input

#### Scenario: Active OAuth profile is used by default for supported API commands
- **WHEN** API-key inputs are absent and `active-tailnet` maps to a profile with `oauth-client-id` and an OAuth client secret value
- **THEN** commands that support OAuth-backed authentication SHALL use that profile for runtime token exchange

#### Scenario: API-key command keeps current precedence
- **WHEN** an existing API-key-authenticated command runs and flags, environment variables, profile data, and legacy keys are present
- **THEN** the CLI SHALL continue to resolve `api-key` and `tailnet` using flags over environment over active profile over legacy config

#### Scenario: OAuth-backed general API command succeeds without a stored API key
- **WHEN** a user runs a supported API command with no resolved API key and a complete OAuth client credential pair is resolved from flags, environment variables, or the active profile
- **THEN** the CLI SHALL exchange the OAuth client credentials at runtime and complete the API request without writing a new API key or access token to the config file

#### Scenario: OAuth credential exchange fails
- **WHEN** a supported API command resolves OAuth client credentials but the token exchange request fails
- **THEN** the command SHALL fail with an actionable authentication error and SHALL NOT attempt to persist the failed exchange result

#### Scenario: OAuth credentials are missing from all sources
- **WHEN** an OAuth-authenticated command runs and flags, environment variables, and the active profile do not provide a complete OAuth client id and secret pair
- **THEN** command execution SHALL fail with an actionable required OAuth credentials error

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
The `config` command group SHALL provide operations to list profiles, set active profile, upsert profile credentials, and remove profiles for both API-key-backed and OAuth-backed profile entries. When secret encryption is enabled, profile mutation commands SHALL persist encrypted secret fields instead of plaintext secret fields.

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

#### Scenario: Upsert command writes encrypted secret fields when encryption is enabled
- **WHEN** a user runs `config profiles set` for an API-key-backed or OAuth-backed profile while secret encryption is enabled
- **THEN** the saved config SHALL persist ciphertext in the matching `*-encrypted` secret field
- **AND** the saved config SHALL omit the plaintext secret field for that profile

#### Scenario: Profile mutations keep canonical file shape
- **WHEN** a user runs `config profiles set` or `config profiles set-active`
- **THEN** the saved config SHALL retain the selected `active-tailnet` and `tailnets` entries
- **AND** the saved config SHALL NOT add or refresh duplicated top-level `tailnet` and `api-key` keys for the active profile

#### Scenario: Remove active profile is blocked
- **WHEN** a user attempts to remove the currently active profile without first selecting another profile
- **THEN** the command SHALL fail with guidance to switch active profile first

### Requirement: Profile behavior is covered by automated tests
The project SHALL include unit and integration tests for profile parsing, validation, runtime resolution, encrypted secret handling, and profile command behavior across both API-key-backed and OAuth-backed profiles.

#### Scenario: Unit tests cover OAuth resolver precedence
- **WHEN** automated unit tests execute for config resolution
- **THEN** tests SHALL verify precedence across OAuth flags, environment variables, active profile values, and missing-credential failures

#### Scenario: Unit tests cover canonical persistence for mixed auth profiles
- **WHEN** automated config persistence tests execute for profile-backed config
- **THEN** tests SHALL verify that saved files keep `active-tailnet` and `tailnets`
- **AND** tests SHALL verify that persisted profiles retain the relevant API-key, encrypted API-key, OAuth client id, and OAuth secret fields without reintroducing duplicated flat keys

#### Scenario: Integration tests cover OAuth-backed API command flows
- **WHEN** command-level integration tests execute for config profile operations and OAuth-backed commands
- **THEN** tests SHALL verify validation failures, persisted state transitions, runtime token exchange, encrypted-secret decryption, and successful API requests using active profiles
