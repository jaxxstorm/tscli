## ADDED Requirements

### Requirement: Tailnet lifecycle commands are available
The CLI SHALL provide `create tailnet`, `list tailnets`, and `delete tailnet` commands for the documented API-driven tailnet lifecycle endpoints.

#### Scenario: Create API-driven tailnet
- **WHEN** a user runs `tscli create tailnet --display-name <display-name>` with valid OAuth client credentials
- **THEN** the CLI SHALL `POST` to `/api/v2/organizations/-/tailnets`
- **AND** the request body SHALL contain the provided `displayName`
- **AND** the command SHALL print the authoritative API response, including the returned `oauthClient` object, in the selected output format

#### Scenario: List organization tailnets
- **WHEN** a user runs `tscli list tailnets` with valid OAuth client credentials
- **THEN** the CLI SHALL `GET` `/api/v2/organizations/-/tailnets`
- **AND** the command SHALL print the authoritative `tailnets` response object in the selected output format

#### Scenario: Delete API-driven tailnet
- **WHEN** a user runs `tscli delete tailnet` with valid tailnet-specific OAuth client credentials
- **THEN** the CLI SHALL `DELETE` `/api/v2/tailnet/-`
- **AND** the command SHALL complete without requiring an API key
- **AND** the command SHALL print a machine-readable success result for structured outputs and a concise success message for human-readable outputs

### Requirement: Tailnet lifecycle commands use OAuth client credential authentication
The `create tailnet`, `list tailnets`, and `delete tailnet` commands SHALL authenticate by exchanging OAuth client credentials for a bearer token and SHALL NOT require `--api-key` or `TAILSCALE_API_KEY`.

#### Scenario: Flags provide OAuth credentials
- **WHEN** a user provides `--oauth-client-id` and `--oauth-client-secret`
- **THEN** the lifecycle command SHALL exchange those credentials for an access token
- **AND** the command SHALL use `Authorization: Bearer <token>` for the lifecycle API request

#### Scenario: Environment provides OAuth credentials
- **WHEN** OAuth flags are absent and `TSCLI_OAUTH_CLIENT_ID` and `TSCLI_OAUTH_CLIENT_SECRET` are set
- **THEN** the lifecycle command SHALL use the environment values for token exchange

#### Scenario: Active profile provides OAuth credentials
- **WHEN** OAuth flags and environment variables are absent and the active profile contains `oauth-client-id` and `oauth-client-secret`
- **THEN** the lifecycle command SHALL use the active profile credentials for token exchange

#### Scenario: OAuth credentials are missing
- **WHEN** the user runs a lifecycle command without a complete OAuth client id and secret from flags, environment, or the active profile
- **THEN** the command SHALL fail before the API request with an actionable missing-credentials error

### Requirement: Tailnet lifecycle command validation and failures are actionable
The lifecycle commands SHALL validate required input before making API requests and SHALL surface OAuth exchange or API failures without replacing them with synthetic success output.

#### Scenario: Create tailnet requires display name
- **WHEN** a user runs `tscli create tailnet` without `--display-name`
- **THEN** the command SHALL fail validation and SHALL NOT issue an API request

#### Scenario: OAuth exchange fails
- **WHEN** token exchange fails for the supplied OAuth client credentials
- **THEN** the command SHALL fail with an error that identifies the OAuth exchange step

#### Scenario: Lifecycle API returns an error
- **WHEN** the lifecycle endpoint returns a non-2xx response
- **THEN** the command SHALL fail with an error that includes the lifecycle operation context and the API response details
