# Configuration

`tscli` reads configuration from:

- `./.tscli.yaml` - your home directory
- `.tscli.yaml` - the current directory you're in

## Supported keys

- `api-key`: legacy single-profile API key
- `api-key-encrypted`: encrypted single-profile or profile API key ciphertext
- `tailnet`: legacy single-profile tailnet
- `oauth-client-id`: optional top-level OAuth client id
- `oauth-client-secret`: optional top-level OAuth client secret
- `oauth-client-secret-encrypted`: encrypted OAuth client secret ciphertext
- `tailnets`: profile list with per-tailnet credentials
- `active-tailnet`: selected profile from `tailnets`
- `encryption.age.public-key`: AGE public key used to encrypt persisted secrets
- `encryption.age.private-key`: optional AGE private key stored in config
- `encryption.age.private-key-command`: optional command that returns the AGE private key at runtime
- `output`: output format (`json`, `yaml`, `human`, `pretty`)
- `debug`: request/response debug logging

## Example config

```yaml
active-tailnet: lbrlabs.com
debug: false
output: pretty
tailnets:
  - api-key: redacted
    name: _lbr_sandbox
  - name: org-admin
    oauth-client-id: redacted-client-id
    oauth-client-secret: redacted-client-secret
  - api-key: redacted
    name: lbrlabs.com
    tailnet: lbrlabs.com
```

Example with encryption enabled:

```yaml
encryption:
  age:
    public-key: age1...
active-tailnet: org-admin
tailnets:
  - name: org-admin
    oauth-client-id: redacted-client-id
    oauth-client-secret-encrypted: |
      -----BEGIN AGE ENCRYPTED FILE-----
      ...
      -----END AGE ENCRYPTED FILE-----
  - name: lbrlabs.com
    api-key-encrypted: |
      -----BEGIN AGE ENCRYPTED FILE-----
      ...
      -----END AGE ENCRYPTED FILE-----
```

`tailnets` + `active-tailnet` is the preferred multi-tailnet shape.
Legacy `api-key` and `tailnet` are still supported for backward compatibility. If a profile omits `tailnet`, `tscli` uses the profile `name` as the effective tailnet.
When encryption is enabled, `config profiles upsert` writes `api-key-encrypted` and `oauth-client-secret-encrypted` instead of the plaintext secret fields.

## Profile commands

```bash
tscli config profiles list
tscli config profiles upsert _lbr_sandbox --api-key tskey-abc123
tscli config profiles upsert org-admin --oauth-client-id cid --oauth-client-secret secret
tscli config profiles set-active _lbr_sandbox
tscli config profiles delete _lbr_sandbox
tscli config encryption setup
```

## Precedence

Runtime values resolve in this order:

1. CLI flags
2. Environment variables
3. Active profile values
4. Top-level config values where supported by the command

For profile-backed config, `active-tailnet` selects the profile used at the config layer when no flag or env override is present.

API-key commands resolve `--api-key` / `TAILSCALE_API_KEY` / profile `api-key` / legacy `api-key` using that order. OAuth-backed commands resolve `--oauth-client-id`, `--oauth-client-secret`, `TSCLI_OAUTH_CLIENT_ID`, `TSCLI_OAUTH_CLIENT_SECRET`, and matching profile or config values using the same precedence model.

Encrypted secret decryption uses this order for the AGE private key:

1. `TSCLI_AGE_PRIVATE_KEY`
2. `encryption.age.private-key-command`
3. `encryption.age.private-key`

### Practical examples

Flags override env:

```bash
TAILSCALE_API_KEY=tskey-env tscli --api-key tskey-flag config get api-key
```

Env overrides config:

```bash
TAILSCALE_API_KEY=tskey-env tscli config get api-key
```
