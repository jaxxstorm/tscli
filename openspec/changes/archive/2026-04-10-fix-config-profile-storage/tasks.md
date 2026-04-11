## 1. Config Model And Persistence

- [x] 1.1 Separate canonical saved config state from resolved runtime state in `pkg/config` so profile-backed files persist only `active-tailnet` and `tailnets`
- [x] 1.2 Update profile mutation flows (`upsert`, `set-active`, remove/save helpers) to stop mirroring the active profile into top-level legacy `tailnet` and `api-key`
- [x] 1.3 Preserve legacy-only flat config read behavior while normalizing mixed legacy/profile files into canonical profile form on profile writes

## 2. Runtime Resolution And User-Facing Behavior

- [x] 2.1 Fix runtime resolution so switching `active-tailnet` deterministically changes the effective `tailnet` and `api-key` for operational commands
- [x] 2.2 Update config-facing examples, help text, or docs to show the canonical profile-backed file shape and legacy compatibility expectations
- [x] 2.3 Verify `config show` and related config commands reflect the intended canonical profile schema without leaking transient resolved runtime values into persisted config

## 3. Regression Coverage

- [x] 3.1 Add unit tests for legacy-only config, mixed config normalization, canonical profile persistence, and precedence resolution
- [x] 3.2 Add command-level tests for `config profiles upsert`, `config profiles set-active`, and persisted file shape after profile mutations
- [x] 3.3 Add an operational-command regression test proving that changing the active profile changes the resolved runtime tailnet/API key without requiring duplicated flat config keys
