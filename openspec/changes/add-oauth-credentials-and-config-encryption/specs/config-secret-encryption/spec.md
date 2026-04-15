## ADDED Requirements

### Requirement: AGE encryption configuration is supported
The CLI SHALL support optional secret-encryption settings under `encryption.age` in the config file. Supported keys SHALL include `public-key`, optional `private-key`, and optional `private-key-command`. The runtime SHALL also accept `TSCLI_AGE_PRIVATE_KEY` as an external private-key source. Config validation SHALL reject an `encryption.age` block that omits `public-key` when secret encryption is enabled or that defines both `private-key` and `private-key-command`.

#### Scenario: Config file defines AGE encryption settings
- **WHEN** a user saves `encryption.age.public-key` and one private-key source
- **THEN** the CLI SHALL treat the config as encryption-enabled for persisted secrets

#### Scenario: Environment variable supplies the AGE private key
- **WHEN** `TSCLI_AGE_PRIVATE_KEY` is set and encrypted secrets are present in config
- **THEN** the CLI SHALL use the environment value as the decryption key source without requiring `encryption.age.private-key` in the config file

#### Scenario: Conflicting private-key sources are configured in config
- **WHEN** the config defines both `encryption.age.private-key` and `encryption.age.private-key-command`
- **THEN** config validation SHALL fail with an actionable error describing the supported private-key source choices

#### Scenario: Invalid AGE public key is configured
- **WHEN** `encryption.age.public-key` is not a valid AGE recipient key
- **THEN** the CLI SHALL fail setup or config validation with an actionable invalid-public-key error

### Requirement: Config encryption setup provides a minimal guided flow
The `config` command group SHALL provide a setup flow for secret encryption that prompts for an AGE public key and allows the user to choose whether the AGE private key is stored directly in config, supplied via `TSCLI_AGE_PRIVATE_KEY`, or retrieved at runtime by executing `encryption.age.private-key-command`. The setup flow SHALL persist only the selected configuration values and SHALL allow users to leave encryption disabled.

#### Scenario: User enables encryption with a stored private key
- **WHEN** a user runs the encryption setup flow, enters an AGE public key, and chooses to store the AGE private key in config
- **THEN** the CLI SHALL persist `encryption.age.public-key` and `encryption.age.private-key`

#### Scenario: User enables encryption with a command-based private key source
- **WHEN** a user runs the encryption setup flow, enters an AGE public key, and chooses command-based retrieval
- **THEN** the CLI SHALL persist `encryption.age.public-key` and `encryption.age.private-key-command`

#### Scenario: User relies on environment-based private key retrieval
- **WHEN** a user runs the encryption setup flow, enters an AGE public key, and declines to store any private-key source in config
- **THEN** the CLI SHALL persist only `encryption.age.public-key`
- **AND** docs and command output SHALL direct the user to provide `TSCLI_AGE_PRIVATE_KEY` at runtime when decryption is needed

### Requirement: Encrypted profile secrets are persisted and decrypted transparently
When encryption is enabled, config writes that persist profile-backed secrets SHALL encrypt secret values with the configured AGE public key. API keys SHALL be persisted in `api-key-encrypted` instead of `api-key`, and OAuth client secrets SHALL be persisted in `oauth-client-secret-encrypted` instead of `oauth-client-secret`. The CLI SHALL decrypt those values before runtime auth resolution, SHALL leave non-secret fields such as `name`, `tailnet`, and `oauth-client-id` unencrypted, and SHALL NOT persist any exchanged runtime access token.

#### Scenario: API-key profile is persisted with encryption enabled
- **WHEN** a user runs `tscli config profiles upsert sandbox --api-key tskey-xxx` while AGE encryption is enabled
- **THEN** the saved profile SHALL contain `api-key-encrypted`
- **AND** the saved profile SHALL NOT contain the plaintext `api-key`

#### Scenario: OAuth profile is persisted with encryption enabled
- **WHEN** a user runs `tscli config profiles upsert org-admin --oauth-client-id cid --oauth-client-secret secret` while AGE encryption is enabled
- **THEN** the saved profile SHALL contain plaintext `oauth-client-id`
- **AND** the saved profile SHALL contain `oauth-client-secret-encrypted`
- **AND** the saved profile SHALL NOT contain the plaintext `oauth-client-secret`

#### Scenario: Encrypted profile secret is decrypted at runtime
- **WHEN** the active profile contains an encrypted API key or encrypted OAuth client secret and a valid AGE private key source is available
- **THEN** the CLI SHALL decrypt the secret in memory before resolving command authentication

#### Scenario: No AGE private key source is available at runtime
- **WHEN** a command needs to decrypt an encrypted profile secret and `TSCLI_AGE_PRIVATE_KEY`, `encryption.age.private-key-command`, and `encryption.age.private-key` are all unavailable
- **THEN** the command SHALL fail with an actionable error that names the supported AGE private-key sources

#### Scenario: Command-based private key retrieval fails
- **WHEN** `encryption.age.private-key-command` is configured and the command exits non-zero or returns an empty key
- **THEN** the CLI SHALL fail with an actionable decryption-source error and SHALL NOT fall back to plaintext profile fields that are not present

### Requirement: Secret encryption behavior is covered by automated tests
The project SHALL include unit and integration tests for AGE configuration validation, config encryption setup, encrypted profile persistence, runtime decryption, and decryption error handling.

#### Scenario: Unit tests cover AGE configuration precedence and validation
- **WHEN** automated unit tests execute for encryption settings
- **THEN** tests SHALL verify public-key validation, mutually exclusive config private-key sources, and runtime precedence of `TSCLI_AGE_PRIVATE_KEY` over config-supplied private-key sources

#### Scenario: Integration tests cover encrypted profile command flows
- **WHEN** command-level integration tests execute for encrypted API-key-backed and OAuth-backed profiles
- **THEN** tests SHALL verify encrypted config persistence, successful runtime decryption, and actionable failures when the private key source is unavailable
