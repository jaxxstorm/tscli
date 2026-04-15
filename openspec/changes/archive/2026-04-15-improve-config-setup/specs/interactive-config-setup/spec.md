## ADDED Requirements

### Requirement: Config setup provides an interactive onboarding flow
The `tscli config setup` command SHALL launch an interactive Bubble Tea experience that guides the user through encryption selection, credential type selection, and profile persistence instead of requiring manual config editing.

#### Scenario: First run starts with encryption selection
- **WHEN** a user runs `tscli config setup` and no existing profiles require immediate management choices
- **THEN** the CLI SHALL open an interactive setup flow
- **AND** the first decision presented SHALL ask whether persisted credentials should be encrypted

#### Scenario: Setup stays within the config command group
- **WHEN** a user completes or exits `tscli config setup`
- **THEN** the CLI SHALL return control without changing the behavior of other `config` subcommands

### Requirement: Config setup provisions encryption when requested
When the user enables encryption in `tscli config setup`, the CLI SHALL generate `age` key material, prompt for a destination path defaulting to `~/.tscli/age.txt`, create the parent directory when it does not exist, and persist `encryption.age.public-key` plus `encryption.age.private-key-path` so later profile writes use encryption automatically.

#### Scenario: Encryption setup accepts the default key path
- **WHEN** the user chooses encrypted storage and accepts the default key path
- **THEN** the CLI SHALL create `~/.tscli` if it does not exist
- **AND** the CLI SHALL write the generated key material to `~/.tscli/age.txt`
- **AND** the CLI SHALL persist the generated public key and default private key path in config

#### Scenario: Encryption setup accepts a custom key path
- **WHEN** the user chooses encrypted storage and enters a custom key path
- **THEN** the CLI SHALL create the custom parent directory when it does not exist
- **AND** the CLI SHALL write the generated key material to that path
- **AND** the CLI SHALL persist the generated public key and custom private key path in config

#### Scenario: Plain-text setup skips encryption configuration
- **WHEN** the user declines encrypted storage
- **THEN** the CLI SHALL continue to credential collection without generating `age` keys
- **AND** the CLI SHALL leave `encryption.age.public-key` and `encryption.age.private-key-path` unchanged when they were not already configured

### Requirement: Config setup provisions and manages profiles interactively
After the encryption decision, `tscli config setup` SHALL prompt the user to choose API-key or OAuth credentials, collect the corresponding values, persist the resulting profile, ask whether another profile should be added, and on rerun SHALL offer management actions that include adding and deleting profiles.

#### Scenario: Setup creates an API-key profile
- **WHEN** the user selects API-key authentication and enters a profile name and API key
- **THEN** the CLI SHALL persist an API-key-backed profile
- **AND** if encryption is enabled, the stored API key SHALL be encrypted before persistence

#### Scenario: Setup creates an OAuth profile
- **WHEN** the user selects OAuth authentication and enters a profile name, OAuth client ID, and OAuth client secret
- **THEN** the CLI SHALL persist an OAuth-backed profile
- **AND** if encryption is enabled, the stored OAuth client secret SHALL be encrypted before persistence

#### Scenario: Setup supports adding multiple profiles in one session
- **WHEN** the user finishes creating a profile and chooses to add another
- **THEN** the CLI SHALL restart the credential-type and profile-entry steps within the same setup session

#### Scenario: Setup exits cleanly after profile creation
- **WHEN** the user finishes creating a profile and declines to add another
- **THEN** the CLI SHALL exit the setup flow gracefully after persisting the completed changes

#### Scenario: Rerun offers profile management actions
- **WHEN** the user runs `tscli config setup` and at least one profile already exists
- **THEN** the CLI SHALL present management options that include adding a new profile and deleting an existing profile

#### Scenario: Rerun deletes a selected profile
- **WHEN** the user chooses the delete action and selects a removable profile
- **THEN** the CLI SHALL remove the selected profile from persisted config
- **AND** the CLI SHALL exit or continue according to the user's follow-up choice
