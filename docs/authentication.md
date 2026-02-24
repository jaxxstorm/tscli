# Authentication

`tscli` authenticates with a Tailscale API key.

## Methods

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

## Secret handling guidance

- Never commit API keys to git.
- Prefer environment variables in CI via a secret manager.
- Rotate leaked or shared keys immediately.
- Use least-privileged keys where possible.
