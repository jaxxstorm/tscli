## MODIFIED Requirements

### Requirement: Configuration structure is documented
The documentation SHALL describe supported configuration keys, including legacy and profile-based tailnet fields, with examples. The documented profile schema SHALL include `name`, optional `tailnet`, `api-key`, `oauth-client-id`, and `oauth-client-secret`, plus `active-tailnet`.

#### Scenario: User reads config docs
- **WHEN** a user opens the configuration documentation page
- **THEN** they SHALL see documented keys for `api-key`, `tailnet`, `tailnets`, `active-tailnet`, `oauth-client-id`, `oauth-client-secret`, and output settings with example YAML

#### Scenario: Precedence is documented
- **WHEN** a user reviews configuration behavior
- **THEN** docs SHALL state precedence order of flags over environment variables over config file values
- **AND** docs SHALL explain that lifecycle commands resolve OAuth credentials from the same precedence layers used for other command auth inputs

### Requirement: Config and auth docs are validated by tests/checks
The project SHALL include checks that required configuration and authentication documentation pages and links exist, including the new OAuth lifecycle guidance.

#### Scenario: Required docs page is missing
- **WHEN** docs validation runs and a required config or auth page is absent
- **THEN** validation SHALL fail with a message identifying missing documentation artifacts

## ADDED Requirements

### Requirement: OAuth client credential authentication workflow is documented
The documentation SHALL explain how to authenticate lifecycle commands using an OAuth client id and secret via CLI flags, environment variables, and config profiles, and SHALL describe how this differs from API-key authentication.

#### Scenario: User reads OAuth auth docs
- **WHEN** a user opens authentication documentation for tailnet lifecycle commands
- **THEN** they SHALL see instructions for `--oauth-client-id`, `--oauth-client-secret`, `TSCLI_OAUTH_CLIENT_ID`, `TSCLI_OAUTH_CLIENT_SECRET`, and profile-based configuration

#### Scenario: User reads create tailnet credential handling guidance
- **WHEN** a user reviews the create tailnet documentation
- **THEN** they SHALL see that the returned OAuth client secret is only shown once
- **AND** the docs SHALL instruct them to store it securely if they plan to use `delete tailnet` or future API interactions for that tailnet

#### Scenario: Security guidance is provided for OAuth credentials
- **WHEN** users follow authentication documentation
- **THEN** docs SHALL include guidance to avoid committing OAuth client secrets and to prefer secure environment or secret management
