# Getting Started

## Install

Install `tscli` using your preferred method from the project README.

## First call

```bash
TAILSCALE_API_KEY=tskey-xxx \
TAILSCALE_TAILNET=example.com \
tscli list devices
```

## Global flags

- `--api-key`, `-k`: Tailscale API key
- `--tailnet`, `-n`: tailnet name (defaults to `-`)
- `--output`, `-o`: `json`, `yaml`, `human`, or `pretty`
- `--debug`, `-d`: dump raw HTTP traffic

## Docs workflows

```bash
make docs-generate
make docs-check
make docs-serve
```
