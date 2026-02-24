## Why

`tscli create key --type authkey` does not expose important Tailscale auth-key capability booleans, forcing users to fall back to raw API calls or other tooling. Adding first-class flags now closes a practical gap in key creation workflows and improves script ergonomics.

## What Changes

- Add boolean flags to `tscli create key` for auth-key creation:
  - `--reusable`
  - `--ephemeral`
  - `--preauthorized`
- Map these flags to the auth-key create request payload capabilities sent to the Tailscale API.
- Validate and document behavior for auth-key mode, including interaction with existing flags (`--type`, `--expiry`, `--description`).
- Keep existing OAuth client creation behavior unchanged.
- Add/expand unit and command integration tests for new flag parsing and request payload behavior.

## Capabilities

### New Capabilities
- `create-key-authkey-capability-flags`: Add explicit CLI support for auth-key capability booleans during key creation.

### Modified Capabilities
- None.

## Impact

- Affected code:
  - `cmd/tscli/create/key/cli.go` for new flags and request mapping
  - `test/cli/**` create-key command integration coverage
  - any key creation helper logic used by tests/mocks
- Affected command groups/flags/config:
  - Command group: `create key`
  - New flags: `--reusable`, `--ephemeral`, `--preauthorized`
  - Existing flags remain supported: `--type`, `--description`, `--expiry`, `--scope`, `--tags`
  - No new config keys or environment variables
- Backward compatibility:
  - No breaking changes expected; existing scripts without new flags continue to work unchanged.
