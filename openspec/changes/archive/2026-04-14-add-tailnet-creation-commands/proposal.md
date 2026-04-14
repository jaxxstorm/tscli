## Why

Tailscale's API-driven tailnet lifecycle endpoints are now documented, but `tscli` cannot create, list, or delete these tailnets today. Supporting them now lets the CLI cover a newly available organization-level workflow while addressing the fact that these endpoints require OAuth client credentials and bearer tokens instead of the existing API-key auth model.

## What Changes

- Add user-facing CLI commands to create, list, and delete API-driven tailnets against the organization tailnet lifecycle endpoints.
- Add OAuth-backed request support for commands that cannot authenticate with the existing API-key basic auth flow.
- Extend the profile/config model so users can persist credentials for either API-key or OAuth-client-backed workflows without breaking existing `api-key`/`tailnet` behavior.
- Add command validation, output shaping, and tests for the new tailnet lifecycle operations and auth resolution paths.
- Update auth and configuration documentation to explain when to use API keys versus OAuth client credentials for these commands.

## Capabilities

### New Capabilities
- `api-tailnet-lifecycle`: Create, list, and delete API-driven tailnets using the documented organization tailnet lifecycle API endpoints.

### Modified Capabilities
- `multi-tailnet-config-profiles`: Extend profile-backed configuration so runtime auth resolution can support OAuth client credentials in addition to API keys.
- `config-and-auth-documentation`: Document OAuth client credential authentication and profile configuration for API-driven tailnet workflows.

## Impact

- Affected command groups: `create`, `list`, `delete`, and `config profiles`
- Affected flags/config/env: new OAuth client credential flags and matching config/profile keys; existing `--api-key`, `--tailnet`, `TAILSCALE_API_KEY`, and `TAILSCALE_TAILNET` remain supported for current scripts
- Affected code: `cmd/`, `pkg/config/`, `pkg/oauth/`, `pkg/tscli/`, CLI tests, and docs
- External dependency impact: organization tailnet lifecycle endpoints use bearer-token auth and may require raw HTTP support outside the typed SDK surface
- Backward compatibility: existing API-key-based commands and legacy flat config must continue to work unchanged
