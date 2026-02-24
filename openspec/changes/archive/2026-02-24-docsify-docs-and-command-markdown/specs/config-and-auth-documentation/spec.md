## ADDED Requirements

### Requirement: Configuration structure is documented
The documentation SHALL describe supported configuration keys, including legacy and profile-based tailnet fields, with examples.

#### Scenario: User reads config docs
- **WHEN** a user opens the configuration documentation page
- **THEN** they SHALL see documented keys for `api-key`, `tailnet`, `tailnets`, `active-tailnet`, and output settings with example YAML

#### Scenario: Precedence is documented
- **WHEN** a user reviews configuration behavior
- **THEN** docs SHALL state precedence order of flags over environment variables over config file values

### Requirement: API key authentication workflow is documented
The documentation SHALL explain how to authenticate using a Tailscale API key via CLI flags, environment variables, and config.

#### Scenario: User reads authentication docs
- **WHEN** a user opens auth documentation
- **THEN** they SHALL see instructions for `--api-key`, `TAILSCALE_API_KEY`, and config file usage

#### Scenario: Security guidance is provided
- **WHEN** users follow authentication documentation
- **THEN** docs SHALL include guidance to avoid committing secrets and to prefer secure environment/secret management

### Requirement: Config and auth docs are validated by tests/checks
The project SHALL include checks that required configuration/auth documentation pages and links exist.

#### Scenario: Required docs page is missing
- **WHEN** docs validation runs and a required config/auth page is absent
- **THEN** validation SHALL fail with a message identifying missing documentation artifacts
