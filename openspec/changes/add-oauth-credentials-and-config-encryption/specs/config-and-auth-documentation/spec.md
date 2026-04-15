## MODIFIED Requirements

### Requirement: Configuration structure is documented
The documentation SHALL describe supported configuration keys, including legacy and profile-based tailnet fields, with examples. The documented profile schema SHALL include `name`, optional `tailnet`, `api-key`, `api-key-encrypted`, `oauth-client-id`, `oauth-client-secret`, `oauth-client-secret-encrypted`, and `active-tailnet`. The documented encryption schema SHALL include `encryption.age.public-key`, optional `encryption.age.private-key`, optional `encryption.age.private-key-command`, and use of `TSCLI_AGE_PRIVATE_KEY`.

#### Scenario: User reads config docs
- **WHEN** a user opens the configuration documentation page
- **THEN** they SHALL see documented keys for `api-key`, `api-key-encrypted`, `tailnet`, `tailnets`, `active-tailnet`, `oauth-client-id`, `oauth-client-secret`, `oauth-client-secret-encrypted`, and AGE encryption settings with example YAML

#### Scenario: Precedence is documented
- **WHEN** a user reviews configuration behavior
- **THEN** docs SHALL state precedence order of flags over environment variables over config file values
- **AND** docs SHALL explain that supported OAuth-backed commands resolve OAuth credentials from the same precedence layers used for other command auth inputs

### Requirement: API key authentication workflow is documented
The documentation SHALL explain how to authenticate using a Tailscale API key via CLI flags, environment variables, and config.

#### Scenario: User reads authentication docs
- **WHEN** a user opens auth documentation
- **THEN** they SHALL see instructions for `--api-key`, `TAILSCALE_API_KEY`, and config file usage

#### Scenario: Security guidance is provided
- **WHEN** users follow authentication documentation
- **THEN** docs SHALL include guidance to avoid committing secrets and to prefer secure environment or secret management

### Requirement: Config and auth docs are validated by tests/checks
The project SHALL include checks that required configuration and authentication documentation pages and links exist, including OAuth-backed API usage guidance and secret-encryption guidance.

#### Scenario: Required docs page is missing
- **WHEN** docs validation runs and a required config/auth page is absent
- **THEN** validation SHALL fail with a message identifying missing documentation artifacts

### Requirement: OAuth client credential authentication workflow is documented
The documentation SHALL explain how to authenticate supported tscli API commands using an OAuth client id and secret via CLI flags, environment variables, and config profiles, and SHALL describe that tscli exchanges those credentials for a short-lived runtime access token instead of storing a reusable API key in the config file.

#### Scenario: User reads OAuth auth docs
- **WHEN** a user opens authentication documentation for supported OAuth-backed API commands
- **THEN** they SHALL see instructions for `--oauth-client-id`, `--oauth-client-secret`, `TSCLI_OAUTH_CLIENT_ID`, `TSCLI_OAUTH_CLIENT_SECRET`, and profile-based configuration

#### Scenario: User reads create tailnet credential handling guidance
- **WHEN** a user reviews the create tailnet documentation
- **THEN** they SHALL see that the returned OAuth client secret is only shown once
- **AND** the docs SHALL instruct them to store it securely if they plan to use future API interactions for that tailnet

#### Scenario: Security guidance is provided for OAuth credentials
- **WHEN** users follow authentication documentation
- **THEN** docs SHALL include guidance to avoid committing OAuth client secrets and to prefer secure environment or secret management

## ADDED Requirements

### Requirement: Secret encryption workflow is documented
The documentation SHALL explain how to enable optional AGE-based config encryption, how `config encryption setup` works, which secret fields are encrypted, and how to provide the AGE private key through config, `TSCLI_AGE_PRIVATE_KEY`, or a configured command.

#### Scenario: User reads encryption setup guidance
- **WHEN** a user opens the credentials or configuration documentation
- **THEN** they SHALL see step-by-step instructions for enabling AGE encryption with an AGE public key and selecting a private-key source

#### Scenario: User reads command-based private key guidance
- **WHEN** a user wants to retrieve the AGE private key from an external tool such as 1Password
- **THEN** the docs SHALL explain how `encryption.age.private-key-command` is executed and SHALL warn that command-based retrieval can add command startup latency

#### Scenario: User reads optional feature guidance
- **WHEN** a user reviews credentials documentation
- **THEN** they SHALL see that OAuth-backed profiles are optional and that config encryption is optional
