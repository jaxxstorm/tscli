# Configuration

`tscli` reads configuration from:

- `./.tscli.yaml` - your home directory
- `.tscli.yaml` - the current directory you're in

## Supported keys

- `api-key`: legacy single-profile API key
- `tailnet`: legacy single-profile tailnet
- `tailnets`: profile list with per-tailnet credentials
- `active-tailnet`: selected profile from `tailnets`
- `output`: output format (`json`, `yaml`, `human`, `pretty`)
- `debug`: request/response debug logging

## Example config

```yaml
active-tailnet: lbrlabs.com
api-key: ""
debug: false
help: false
output: pretty
tailnet: "-"
tailnets:
  - api-key: redacted
    name: _lbr_sandbox
  - api-key: redacted
    name: lbrlabs.com
```

`tailnets` + `active-tailnet` is the preferred multi-tailnet shape.
Legacy `api-key` and `tailnet` are still supported for backward compatibility.

## Profile commands

```bash
tscli config profiles list
tscli config profiles upsert _lbr_sandbox --api-key tskey-abc123
tscli config profiles set-active _lbr_sandbox
tscli config profiles delete _lbr_sandbox
```

## Precedence

Runtime values resolve in this order:

1. CLI flags
2. Environment variables
3. Config file values

For profile-backed config, `active-tailnet` is used at the config layer when no flag or env override is present.

### Practical examples

Flags override env:

```bash
TAILSCALE_API_KEY=tskey-env tscli --api-key tskey-flag config get api-key
```

Env overrides config:

```bash
TAILSCALE_API_KEY=tskey-env tscli config get api-key
```
