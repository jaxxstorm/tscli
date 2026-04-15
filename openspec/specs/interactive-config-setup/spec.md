## Purpose

Define the interactive `tscli config setup` onboarding and profile-management flow, including encryption provisioning and setup-driven profile persistence.

## Requirements

### Requirement: Config setup provides an interactive onboarding flow
The `tscli config setup` command SHALL launch an interactive Bubble Tea experience that guides the user through encryption selection, credential type selection, and profile persistence instead of requiring manual config editing. The interactive flow SHALL present setup steps, available choices, and status messaging in a structured Bubble Tea UI that is easier to scan than plain prompt-only output, while preserving the existing command behavior when running non-interactively.

#### Scenario: First run starts with encryption selection
- **WHEN** a user runs `tscli config setup` and no existing profiles require immediate management choices
- **THEN** the CLI SHALL open an interactive setup flow
- **AND** the first decision presented SHALL ask whether persisted credentials should be encrypted

#### Scenario: Interactive flow presents structured choices
- **WHEN** a user advances through `tscli config setup` in an interactive terminal
- **THEN** the Bubble Tea interface SHALL show the current step and available actions using structured interactive UI elements rather than plain typed prompt text alone
- **AND** the interface SHALL make the current selection state and resulting status messages visible within the same guided flow

#### Scenario: Non-interactive setup remains prompt-driven
- **WHEN** `tscli config setup` runs without an interactive terminal
- **THEN** the command SHALL continue to use line-oriented prompt input and output
- **AND** the command SHALL not require Bubble Tea-only terminal features to complete setup

#### Scenario: Setup stays within the config command group
- **WHEN** a user completes or exits `tscli config setup`
- **THEN** the CLI SHALL return control without changing the behavior of other `config` subcommands

### Requirement: Config setup provisions encryption when requested
When the user enables encryption in `tscli config setup`, the CLI SHALL prompt for a destination path that defaults to `~/.tscli/age.txt`, check whether an AGE identity file already exists at the selected path, and prompt the user to reuse that identity when it contains a valid AGE private key and derivable public key. If the user reuses the existing identity, the CLI SHALL persist the derived `encryption.age.public-key` and selected `encryption.age.private-key-path` without generating new key material. If the user declines reuse or the file is missing or invalid, the CLI SHALL generate `age` key material, create the parent directory when it does not exist, write the generated identity to the selected path, and persist `encryption.age.public-key` plus `encryption.age.private-key-path` so later profile writes use encryption automatically.

#### Scenario: Encryption setup reuses an existing default key path
- **WHEN** the user chooses encrypted storage, accepts the default key path, and a valid AGE identity file already exists at `~/.tscli/age.txt`
- **THEN** the interactive flow SHALL prompt whether to reuse the existing identity
- **AND** if the user accepts, the CLI SHALL derive the public key from the existing private key file
- **AND** the CLI SHALL persist the derived public key and default private key path in config without generating a new key

#### Scenario: Encryption setup reuses an existing custom key path
- **WHEN** the user chooses encrypted storage, enters a custom key path, and a valid AGE identity file already exists at that path
- **THEN** the interactive flow SHALL prompt whether to reuse the existing identity
- **AND** if the user accepts, the CLI SHALL persist the derived public key and custom private key path in config without overwriting the file

#### Scenario: Encryption setup replaces an existing key file when declined
- **WHEN** the user chooses encrypted storage, selects a path containing a valid AGE identity file, and declines to reuse it
- **THEN** the CLI SHALL generate fresh AGE key material
- **AND** the CLI SHALL overwrite the selected path with the generated identity
- **AND** the CLI SHALL persist the generated public key and selected private key path in config

#### Scenario: Encryption setup falls back to generation for invalid existing key data
- **WHEN** the user chooses encrypted storage and the selected path exists but does not contain a valid AGE identity file
- **THEN** the CLI SHALL report that the existing file cannot be reused
- **AND** the CLI SHALL continue with generated key setup instead of persisting invalid encryption settings

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
