## MODIFIED Requirements

### Requirement: Config setup provisions and manages profiles interactively
After the encryption decision, `tscli config setup` SHALL prompt the user to choose API-key or OAuth credentials, collect the corresponding values, persist the resulting profile, ask whether another profile should be added, and on rerun SHALL offer management actions that include adding and deleting profiles. When the command starts without existing profiles and the user finishes profile creation, the setup flow SHALL then prompt for a default `output` mode with choices `json`, `pretty`, and `human`, prompt whether debug HTTP request/response logging should be enabled by default, persist those preferences, and only then exit the setup flow.

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

#### Scenario: Initial setup prompts for output mode after profile creation
- **WHEN** the command starts without existing profiles, the user finishes creating profiles, and declines to add another
- **THEN** the CLI SHALL prompt for a default output mode after profile setup completes
- **AND** the available setup choices SHALL be `json`, `pretty`, and `human`

#### Scenario: Initial setup prompts for debug preference after output selection
- **WHEN** the command starts without existing profiles and the user selects a default output mode in setup
- **THEN** the CLI SHALL prompt whether debug HTTP request/response logging should be enabled by default before setup exits

#### Scenario: Setup exits cleanly after collecting initial preferences
- **WHEN** the command starts without existing profiles, the user completes profile creation, selects an output mode, and answers the debug preference prompt
- **THEN** the CLI SHALL persist the completed changes
- **AND** the CLI SHALL exit the setup flow gracefully

#### Scenario: Rerun offers profile management actions
- **WHEN** the user runs `tscli config setup` and at least one profile already exists
- **THEN** the CLI SHALL present management options that include adding a new profile and deleting an existing profile

#### Scenario: Rerun deletes a selected profile
- **WHEN** the user chooses the delete action and selects a removable profile
- **THEN** the CLI SHALL remove the selected profile from persisted config
- **AND** the CLI SHALL exit or continue according to the user's follow-up choice
