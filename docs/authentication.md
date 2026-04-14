# Authentication

`tscli` supports two authentication models:

- API keys for the existing tailnet-scoped API commands
- OAuth client credentials for the API-driven tailnet lifecycle commands (`create tailnet`, `list tailnets`, `delete tailnet`)

## API key methods

1. CLI flag

```bash
tscli --api-key tskey-xxx list devices
```

2. Environment variable

```bash
export TAILSCALE_API_KEY=tskey-xxx
tscli list devices
```

3. Config file

```yaml
api-key: tskey-xxx
```

Or with profiles:

```yaml
active-tailnet: example.com
tailnets:
  - name: example.com
    api-key: tskey-xxx
```

Profile-backed configs use `active-tailnet` plus the `tailnets` array as the canonical stored shape. Top-level `tailnet` and `api-key` remain legacy compatibility keys for older single-tailnet config files and are not required in profile mode.

## OAuth client credential methods

Use OAuth client credentials for the tailnet lifecycle commands:

1. CLI flags

```bash
tscli list tailnets --oauth-client-id cid --oauth-client-secret secret
```

2. Environment variables

```bash
export TSCLI_OAUTH_CLIENT_ID=cid
export TSCLI_OAUTH_CLIENT_SECRET=secret
tscli create tailnet --display-name sandbox
```

3. Config file or profile

```yaml
oauth-client-id: cid
oauth-client-secret: secret
```

Or with profiles:

```yaml
active-tailnet: org-admin
tailnets:
  - name: org-admin
    oauth-client-id: cid
    oauth-client-secret: secret
```

Lifecycle commands use the same precedence layers as other auth inputs: flags override environment variables, which override active profile values, which override matching top-level config values.

## Tailnet lifecycle notes

- `create tailnet` and `list tailnets` require an organization-approved OAuth client.
- `delete tailnet` requires the tailnet-specific OAuth client returned when the tailnet was created.
- The `oauthClient.secret` returned by `create tailnet` is shown only once by the API. Store it securely if you will need to delete or manage that tailnet later.

## Secret handling guidance

- Never commit API keys to git.
- Never commit OAuth client secrets to git.
- Prefer environment variables in CI via a secret manager.
- Rotate leaked or shared credentials immediately.
- Use least-privileged keys where possible.
