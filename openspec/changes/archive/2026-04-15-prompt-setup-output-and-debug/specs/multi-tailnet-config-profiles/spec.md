## MODIFIED Requirements

### Requirement: Config commands manage tailnet profiles
The `config` command group SHALL provide operations to list profiles, set active profile, upsert profile credentials, remove profiles, and run an interactive `config setup` flow for both API-key-backed and OAuth-backed profile entries. When secret encryption is enabled, profile mutation commands SHALL persist encrypted secret fields instead of plaintext secret fields. During initial `config setup`, after profile creation is complete, the command SHALL also persist the selected top-level `output` mode and `debug` preference in the config file.

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

#### Scenario: Interactive setup writes output and debug preferences during initial setup
- **WHEN** a user runs `config setup` with no existing profiles, completes profile creation, selects `pretty` or `human` or `json`, and chooses whether debug logging is enabled
- **THEN** the saved config SHALL persist the selected top-level `output` value
- **AND** the saved config SHALL persist the selected top-level `debug` boolean value

#### Scenario: Upsert command writes encrypted secret fields when encryption is enabled
- **WHEN** a user runs `config profiles set` for an API-key-backed or OAuth-backed profile while secret encryption is enabled
- **THEN** the saved config SHALL persist ciphertext in the matching `*-encrypted` secret field
- **AND** the saved config SHALL omit the plaintext secret field for that profile

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
