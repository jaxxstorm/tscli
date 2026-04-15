## Purpose

Define AGE-based config secret encryption behavior for persisted profile credentials and related setup, validation, and test coverage.

## Requirements

### Requirement: AGE encryption configuration is supported
The CLI SHALL support optional secret-encryption settings under `encryption.age` in the config file. Supported keys SHALL include `public-key`, optional `private-key-path`, optional `private-key`, and optional `private-key-command`. The runtime SHALL also accept `TSCLI_AGE_PRIVATE_KEY` as an external private-key source. Config validation SHALL reject an `encryption.age` block that omits `public-key` when secret encryption is enabled or that defines more than one configured private-key source among `private-key-path`, `private-key`, and `private-key-command`.

#### Scenario: Config file defines AGE encryption settings
- **WHEN** a user saves `encryption.age.public-key` and one private-key source
- **THEN** the CLI SHALL treat the config as encryption-enabled for persisted secrets

#### Scenario: Config file defines a private key path
- **WHEN** a user saves `encryption.age.public-key` and `encryption.age.private-key-path`
- **THEN** the CLI SHALL read the AGE identity from that path for runtime decryption

#### Scenario: Environment variable supplies the AGE private key
- **WHEN** `TSCLI_AGE_PRIVATE_KEY` is set and encrypted secrets are present in config
- **THEN** the CLI SHALL use the environment value as the decryption key source without requiring another configured private-key source in the config file

#### Scenario: Conflicting private-key sources are configured in config
- **WHEN** the config defines more than one of `encryption.age.private-key-path`, `encryption.age.private-key`, or `encryption.age.private-key-command`
- **THEN** config validation SHALL fail with an actionable error describing the supported private-key source choices

#### Scenario: Invalid AGE public key is configured
- **WHEN** `encryption.age.public-key` is not a valid AGE recipient key
- **THEN** the CLI SHALL fail setup or config validation with an actionable invalid-public-key error

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
- **WHEN** a user runs `tscli config setup`, enables encryption, and enters a custom key path
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

### Requirement: Encrypted profile secrets are persisted and decrypted transparently
When encryption is enabled, config writes that persist profile-backed secrets SHALL encrypt secret values with the configured AGE public key. API keys SHALL be persisted in `api-key-encrypted` instead of `api-key`, and OAuth client secrets SHALL be persisted in `oauth-client-secret-encrypted` instead of `oauth-client-secret`. The CLI SHALL decrypt those values before runtime auth resolution, SHALL leave non-secret fields such as `name`, `tailnet`, and `oauth-client-id` unencrypted, and SHALL NOT persist any exchanged runtime access token.

#### Scenario: API-key profile is persisted with encryption enabled
- **WHEN** a user runs `tscli config profiles set sandbox --api-key tskey-xxx` while AGE encryption is enabled
- **THEN** the saved profile SHALL contain `api-key-encrypted`
- **AND** the saved profile SHALL NOT contain the plaintext `api-key`

#### Scenario: OAuth profile is persisted with encryption enabled
- **WHEN** a user runs `tscli config profiles set org-admin --oauth-client-id cid --oauth-client-secret secret` while AGE encryption is enabled
- **THEN** the saved profile SHALL contain plaintext `oauth-client-id`
- **AND** the saved profile SHALL contain `oauth-client-secret-encrypted`
- **AND** the saved profile SHALL NOT contain the plaintext `oauth-client-secret`

#### Scenario: Encrypted profile secret is decrypted at runtime
- **WHEN** the active profile contains an encrypted API key or encrypted OAuth client secret and a valid AGE private key source is available
- **THEN** the CLI SHALL decrypt the secret in memory before resolving command authentication

#### Scenario: No AGE private key source is available at runtime
- **WHEN** a command needs to decrypt an encrypted profile secret and `TSCLI_AGE_PRIVATE_KEY`, `encryption.age.private-key-path`, `encryption.age.private-key-command`, and `encryption.age.private-key` are all unavailable
- **THEN** the command SHALL fail with an actionable error that names the supported AGE private-key sources

#### Scenario: Command-based private key retrieval fails
- **WHEN** `encryption.age.private-key-command` is configured and the command exits non-zero or returns an empty key
- **THEN** the CLI SHALL fail with an actionable decryption-source error and SHALL NOT fall back to plaintext profile fields that are not present

### Requirement: Secret encryption behavior is covered by automated tests
The project SHALL include unit and integration tests for AGE configuration validation, config encryption setup, encrypted profile persistence, runtime decryption, existing key-file reuse, invalid identity-file handling, and decryption error handling.

#### Scenario: Unit tests cover AGE configuration precedence and validation
- **WHEN** automated unit tests execute for encryption settings
- **THEN** tests SHALL verify public-key validation, mutually exclusive config private-key sources, generated key-path handling, existing identity-file parsing, reuse-path public-key derivation, and runtime precedence of `TSCLI_AGE_PRIVATE_KEY` over config-supplied private-key sources

#### Scenario: Integration tests cover encrypted profile command flows
- **WHEN** command-level integration tests execute for encrypted API-key-backed and OAuth-backed profiles
- **THEN** tests SHALL verify encrypted config persistence, successful runtime decryption, actionable failures when the private key source is unavailable, and the interactive reuse prompt behavior for existing AGE identity files
