## Why

Profile-backed config persistence is currently writing broken mixed state: empty `api-key` values, duplicated `tailnet` and `api-key` flat keys, and profile changes that do not reliably change the effective runtime tailnet. This needs to be fixed now because profile switching is a user-facing workflow and the current file shape makes both the config file and runtime resolution hard to trust.

## What Changes

- Canonicalize persisted profile-based config so profile workflows store only `active-tailnet` plus the `tailnets` array, without copying the active profile into top-level `tailnet` and `api-key`.
- Preserve backward compatibility for legacy flat config files that only use top-level `tailnet` and `api-key`, but treat those keys as compatibility input rather than canonical storage once profile data exists.
- Fix profile mutation and resolution flows so `config profiles upsert`, `config profiles set-active`, and operational command resolution consistently use the selected active profile.
- Add unit and command-level regression coverage for canonical file shape, mixed legacy/profile config handling, and active-profile switching.

## Capabilities

### New Capabilities

### Modified Capabilities
- `multi-tailnet-config-profiles`: tighten the persisted config schema so profile mode uses canonical profile-only storage, legacy flat keys remain backward-compatible input only, and active profile switching deterministically controls runtime credentials.

## Impact

- Affected code: `pkg/config`, Cobra/Viper config initialization and save paths, profile management commands under `cmd/tscli/config/profiles`, and config-related tests.
- Affected config/env keys: `active-tailnet`, `tailnets[].name`, `tailnets[].api-key`, legacy `tailnet`, legacy `api-key`, `TAILSCALE_TAILNET`, and `TAILSCALE_API_KEY`.
- Backward compatibility: legacy flat config files remain supported for reads, but profile-backed configs will be rewritten into the canonical schema so scripts and users no longer depend on duplicated flat keys being present.
