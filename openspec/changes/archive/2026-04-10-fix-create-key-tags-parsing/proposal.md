## Why

`tscli create key --type authkey --tags ...` currently accepts the `--tags` flag but serializes the auth-key request payload with `"tags": null`, causing the Tailscale API to reject tagged tailnet-owned keys. This is a user-visible regression in an existing flag path, so the change needs to restore correct request mapping and add regression coverage now.

## What Changes

- Fix `tscli create key` auth-key request construction so `--tags` values are preserved in the outgoing `capabilities.devices.create.tags` payload.
- Add automated coverage that reproduces the current `tags: null` regression before the fix and verifies the corrected payload after the fix.
- Preserve existing `create key` behavior for other auth-key capability flags (`--reusable`, `--ephemeral`, `--preauthorized`) and for non-authkey flows.
- Keep the current command surface unchanged: no new flags, no removed flags, and no new config or environment settings.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `create-key-authkey-capability-flags`: Extend the auth-key create requirements so `--tags` is serialized into the auth-key capability payload and covered by regression tests.

## Impact

- Affected code:
  - `cmd/tscli/create/key/cli.go`
  - `test/cli/create_key_flags_integration_test.go`
- Affected command groups/flags/config:
  - Command group: `create key`
  - Existing flag behavior corrected: `--tags`
  - Existing related auth-key flags remain supported: `--reusable`, `--ephemeral`, `--preauthorized`
  - No config keys or environment variables change
- Backward compatibility:
  - No breaking CLI changes are intended.
  - Existing scripts that already pass `--tags` for auth-key creation should begin working as documented instead of failing with a server-side validation error.
