# Authentication

`tscli` supports two authentication models:

This page covers how `tscli` authenticates its own requests. For creating Tailscale auth keys, OAuth clients, or federated credentials, see [Creating Credentials](creating-credentials.md).

- API keys for tailnet-scoped API commands
- OAuth client credentials for API-driven tailnet lifecycle commands and for profile-backed API access when you want `tscli` to exchange credentials at runtime instead of storing a reusable API key locally

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

Use OAuth client credentials when you want `tscli` to exchange them for a short-lived bearer token at runtime. This is supported for the tailnet lifecycle commands and for profile-backed API access.

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

Supported OAuth-backed commands use the same precedence layers as other auth inputs: flags override environment variables, which override active profile values, which override matching top-level config values.

When `tscli` uses OAuth-backed auth, it exchanges the client id and secret for a short-lived bearer token during command execution. It does not persist the exchanged token or write a reusable API key back into your config file.

## Optional config encryption with age

OAuth client secrets and API keys can be encrypted in the config file with `age`. Encryption is optional, and OAuth-backed profiles are optional.

1. Generate or obtain an AGE keypair.
2. Run the guided setup command:

```bash
tscli config encryption setup
```

The setup flow asks for:

- an AGE public key, stored as `encryption.age.public-key`
- how the AGE private key will be supplied at runtime:
  - `config`: store it as `encryption.age.private-key`
  - `env`: provide it with `TSCLI_AGE_PRIVATE_KEY`
  - `command`: configure `encryption.age.private-key-command`

Example encrypted config:

```yaml
encryption:
  age:
    public-key: age1...
    private-key-command: op read op://vault/tscli/age-private-key
active-tailnet: org-admin
tailnets:
  - name: org-admin
    oauth-client-id: cid
    oauth-client-secret-encrypted: |
      -----BEGIN AGE ENCRYPTED FILE-----
      ...
      -----END AGE ENCRYPTED FILE-----
```

Command-based private-key lookup is useful with secret managers such as 1Password, but it can add command startup latency because the command runs each time `tscli` needs to decrypt a stored secret.

## Tailnet lifecycle notes

- `create tailnet` and `list tailnets` require an organization-approved OAuth client.
- `delete tailnet` requires the tailnet-specific OAuth client returned when the tailnet was created.
- The `oauthClient.secret` returned by `create tailnet` is shown only once by the API. Store it securely if you will need to delete or manage that tailnet later.

## Secret handling guidance

- Never commit API keys to git.
- Never commit OAuth client secrets to git.
- Prefer `TSCLI_AGE_PRIVATE_KEY` or an external command over storing the AGE private key directly in config when possible.
- Prefer environment variables in CI via a secret manager.
- Rotate leaked or shared credentials immediately.
- Use least-privileged keys where possible.
