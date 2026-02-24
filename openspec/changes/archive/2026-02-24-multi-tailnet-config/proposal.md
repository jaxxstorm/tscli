## Why

`tscli` currently assumes a single tailnet/API key in config, which forces users with multiple tailnets to manually edit config or export env vars between runs. Supporting multiple tailnet profiles with an explicit active tailnet will remove repetitive setup while keeping command execution predictable.

## What Changes

- Add multi-tailnet profile support in config with a simple `tailnets` list and an `active-tailnet` selector.
- Add profile-aware config resolution so all existing command groups (`get`, `list`, `set`, `create`, `delete`) automatically use the active tailnet profile unless flags/env override.
- Add `config` subcommands to list profiles, set/change active tailnet, and upsert/remove profile credentials.
- Preserve backward compatibility for existing flat config keys (`tailnet`, `api-key`) and current script behavior.
- Document migration behavior from single-tailnet config to multi-tailnet profiles.

## Capabilities

### New Capabilities
- `multi-tailnet-config-profiles`: Define profile storage, active-tailnet switching, and runtime credential/tailnet resolution rules.

### Modified Capabilities
- None.

## Impact

- Affected code:
  - `pkg/config/**` for config schema loading, defaults, and migration/normalization
  - `cmd/tscli/config/**` for profile management commands
  - `internal/cli/root.go` / CLI pre-run resolution of tailnet + API key
  - `test/cli/**` for config precedence and profile switching behavior
  - README/config docs and examples
- Affected user-facing config/env/flags:
  - Config keys: `tailnets`, `active-tailnet`, legacy `tailnet`, legacy `api-key`
  - Existing flags/env remain supported (`--tailnet`, `--api-key`, `TAILSCALE_TAILNET`, `TAILSCALE_API_KEY`)
- Backward compatibility:
  - Existing single-tailnet config continues to work without changes.
  - Legacy keys are interpreted consistently, with profile-based config taking effect only when configured.

## Release Notes

- Added multi-tailnet profile support with `tailnets` and `active-tailnet`.
- Added profile management commands under `tscli config profiles` (`list`, `upsert`, `set-active`, `delete`).
- Runtime auth resolution now follows `flags > env > active profile > legacy config`.
- Backward compatibility: legacy `tailnet` and `api-key` keys continue to work and are mirrored from the active profile when profile commands update config.
- Rollback: remove `tailnets` and `active-tailnet` from config to return to legacy single-tailnet operation.
