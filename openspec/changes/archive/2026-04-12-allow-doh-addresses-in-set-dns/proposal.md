## Why

`tscli set dns nameservers` and `tscli set dns split-dns` currently reject DNS-over-HTTPS endpoint addresses even though the upstream Tailscale DNS configuration accepts them alongside literal IP nameservers. This blocks valid DNS configurations and forces users to fall back to raw API calls instead of using the existing `set` DNS command surface.

## What Changes

- Update `set dns nameservers` so `--nameserver` accepts both literal IP addresses and valid DoH endpoint addresses, while preserving existing `get`, `set`, `delete`, and `list` command semantics for DNS operations.
- Update `set dns split-dns` so `--entry domain=value` accepts DoH endpoint addresses anywhere a nameserver value is currently validated as an IP.
- Add command validation tests and mock-backed success-path tests covering both valid DoH inputs and still-invalid nameserver values.
- Update user-facing help text and examples for the affected `set dns` commands to document the expanded accepted nameserver format.
- Do not add new command groups, flags, config keys, or environment variables; this change only relaxes validation for existing `set dns nameservers` and `set dns split-dns` flags.

## Capabilities

### New Capabilities

- None.

### Modified Capabilities

- `cli-command-test-coverage`: Extend validation and success-path coverage so DNS mutation commands verify acceptance of both IP and DoH nameserver inputs.
- `openapi-command-parity`: Update the DNS command behavior contract so existing `set` DNS commands accept the full upstream-supported nameserver value formats without requiring raw API fallbacks.

## Impact

- Affected code:
  - `cmd/tscli/set/dns/nameservers/cli.go`
  - `cmd/tscli/set/dns/split/cli.go`
  - `test/cli/**` DNS validation and output coverage
  - DNS command help text and examples
- Affected API areas:
  - `POST /tailnet/{tailnet}/dns/nameservers`
  - `PATCH` and `PUT /tailnet/{tailnet}/dns/split-dns`
- Backward compatibility:
  - Existing scripts that pass IP nameservers continue to work unchanged
  - Invalid nameserver values should still fail with actionable validation errors
