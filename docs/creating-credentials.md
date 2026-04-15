# Creating Credentials

This page covers credentials that `tscli` can create for use with the Tailscale API. It is separate from [Authentication](authentication.md), which explains how `tscli` authenticates its own API calls.

## Supported key types

Use `tscli create key` to create:

- Auth keys for device enrollment
- OAuth clients for scoped API access
- Federated credentials for identities backed by an external OIDC issuer

The full flag reference lives in the generated [`tscli create key`](commands/tscli_create_key.md) command page.

## Federated credentials

Use federated credentials when you want Tailscale identities to match an OIDC issuer and subject pattern.

Required inputs:

- `--type federated`
- One or more `--scope` values
- `--issuer`
- `--subject`

Optional inputs:

- `--audience`
- `--tags`
- `--claim key=value`

Example:

```bash
tscli create key \
  --type federated \
  --scope users:read \
  --issuer https://example.com \
  --subject example-* \
  --audience my-app \
  --claim env=prod
```
