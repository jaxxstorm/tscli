## MODIFIED Requirements

### Requirement: Config encryption setup provides an interactive guided flow
The `config` command group SHALL provide setup flows for secret encryption. `tscli config setup` SHALL act as the primary interactive flow by prompting for a destination path that defaults to `~/.tscli/age.txt`, checking for an existing AGE identity file at the chosen path, and prompting the user to reuse that file when it contains a valid AGE private key and derivable public key. When the user accepts reuse, the CLI SHALL persist `encryption.age.public-key` and `encryption.age.private-key-path` from the existing identity without generating new key material. When the user declines reuse, or when no valid existing file is present, `tscli config setup` SHALL generate AGE key material, create the parent directory when needed, and persist `encryption.age.public-key` plus `encryption.age.private-key-path`. The dedicated `config encryption setup` command SHALL remain available for targeted setup, SHALL allow users to configure other supported private-key sources, and SHALL offer the same reuse flow whenever a path-based AGE identity already exists. Both flows SHALL allow users to leave encryption disabled.

#### Scenario: Interactive setup reuses the default existing key path
- **WHEN** a user runs `tscli config setup`, enables encryption, accepts the default key path, and `~/.tscli/age.txt` already contains a valid AGE identity
- **THEN** the CLI SHALL prompt whether to reuse the existing identity
- **AND** if the user accepts, the CLI SHALL derive and persist `encryption.age.public-key` from that file
- **AND** the CLI SHALL persist `encryption.age.private-key-path` as `~/.tscli/age.txt`

#### Scenario: Interactive setup writes the default generated key path after declining reuse
- **WHEN** a user runs `tscli config setup`, enables encryption, accepts the default key path, and declines reuse of an existing valid identity file at `~/.tscli/age.txt`
- **THEN** the CLI SHALL generate AGE key material
- **AND** the CLI SHALL write the generated identity to `~/.tscli/age.txt`
- **AND** the CLI SHALL persist `encryption.age.public-key` and `encryption.age.private-key-path`

#### Scenario: Interactive setup writes a custom generated key path
- **WHEN** a user runs `tscli config setup`, enables encryption, enters a custom key path, and no reusable AGE identity file exists at that path
- **THEN** the CLI SHALL create the custom parent directory when needed
- **AND** the CLI SHALL write the generated identity to that path
- **AND** the CLI SHALL persist `encryption.age.public-key` and `encryption.age.private-key-path`

#### Scenario: Dedicated encryption setup reuses an existing path-based identity
- **WHEN** a user runs `config encryption setup`, chooses `--private-key-source=path` or the interactive path source, and the selected private key path already contains a valid AGE identity
- **THEN** the CLI SHALL prompt whether to reuse the existing identity
- **AND** if the user accepts, the CLI SHALL persist the derived `encryption.age.public-key` and selected `encryption.age.private-key-path` without requiring manual public-key re-entry

#### Scenario: Dedicated encryption setup configures a command-based private key source
- **WHEN** a user runs `config encryption setup`, enters an AGE public key, and chooses command-based retrieval
- **THEN** the CLI SHALL persist `encryption.age.public-key` and `encryption.age.private-key-command`

#### Scenario: Dedicated encryption setup relies on environment-based private key retrieval
- **WHEN** a user runs `config encryption setup`, enters an AGE public key, and declines to store any configured private-key source in config
- **THEN** the CLI SHALL persist only `encryption.age.public-key`
- **AND** docs and command output SHALL direct the user to provide `TSCLI_AGE_PRIVATE_KEY` at runtime when decryption is needed

### Requirement: Secret encryption behavior is covered by automated tests
The project SHALL include unit and integration tests for AGE configuration validation, config encryption setup, encrypted profile persistence, runtime decryption, existing key-file reuse, invalid identity-file handling, and decryption error handling.

#### Scenario: Unit tests cover AGE configuration precedence and validation
- **WHEN** automated unit tests execute for encryption settings
- **THEN** tests SHALL verify public-key validation, mutually exclusive config private-key sources, generated key-path handling, existing identity-file parsing, reuse-path public-key derivation, and runtime precedence of `TSCLI_AGE_PRIVATE_KEY` over config-supplied private-key sources

#### Scenario: Integration tests cover encrypted profile command flows
- **WHEN** command-level integration tests execute for encrypted API-key-backed and OAuth-backed profiles
- **THEN** tests SHALL verify encrypted config persistence, successful runtime decryption, actionable failures when the private key source is unavailable, and the interactive reuse prompt behavior for existing AGE identity files
