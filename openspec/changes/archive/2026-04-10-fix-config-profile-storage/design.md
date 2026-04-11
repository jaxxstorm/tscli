## Context

`tscli` now supports both legacy flat config and multi-tailnet profile config, but the persistence path still mixes the two models. Profile mutations can leave top-level `tailnet` and `api-key` keys in the file, copy active-profile state into those keys, and persist empty values that obscure what the runtime resolver should actually use. The implementation must keep legacy files working while making profile-backed config deterministic and script-friendly.

## Goals / Non-Goals

**Goals:**
- Make profile-backed config persist in one canonical shape: `active-tailnet` plus `tailnets`.
- Keep runtime resolution order unchanged for operational commands: flags, environment, active profile, then legacy flat keys.
- Ensure profile mutations and profile switching immediately affect the effective runtime tailnet and API key.
- Add unit and command-level tests that lock the file shape and switching behavior.

**Non-Goals:**
- Removing support for legacy flat config files.
- Changing CLI flag names, environment variable names, or output formats.
- Introducing a new config schema version or external migration tool.

## Decisions

### Persist canonical profile state, not resolved runtime values
Profile-aware writes will persist only `active-tailnet` and `tailnets`. The legacy flat `tailnet` and `api-key` keys represent compatibility input for legacy-only config and will not be rewritten from active profile state.

Alternative considered: keep mirroring the active profile into top-level flat keys.
Why not: it creates duplicate sources of truth, makes empty/partial values ambiguous, and causes profile switching to appear broken when callers inspect the file.

### Normalize mixed config files on profile mutations
When a profile command rewrites config, the saved file will be normalized into canonical profile form. Legacy-only files remain valid until the user enters profile mode, at which point profile data becomes the authoritative persisted shape.

Alternative considered: preserve all pre-existing top-level flat keys indefinitely.
Why not: this keeps the broken mixed-state behavior and makes it impossible to guarantee that the persisted file matches the active profile model.

### Separate read compatibility from save behavior
The resolver can still read legacy flat keys as the final fallback when no flag, env, or active profile provides a value. Save paths must not treat those resolved runtime values as data that should be written back to disk.

Alternative considered: rely on current `viper` in-memory state for both resolution and persistence.
Why not: `viper` merges flags, env, and config state, so using its merged view directly for saves leaks transient runtime values into the config file.

### Test both file shape and runtime effect
Unit tests will cover normalization and precedence behavior in `pkg/config`. Command-level tests will cover `config profiles upsert`, `config profiles set-active`, `config show`, and an operational command path that proves the active profile actually controls resolved runtime values.

Alternative considered: only test the persisted YAML.
Why not: the bug report also includes broken runtime behavior when switching profiles, so file-shape tests alone would miss regressions.

## Risks / Trade-offs

- [Risk] Users with hand-edited mixed config files may see flat legacy keys removed after running a profile command. → Mitigation: keep legacy-only files valid, document canonical profile shape, and limit normalization to save paths.
- [Risk] Resolution logic and persistence logic can drift again if they share mutable `viper` state too directly. → Mitigation: keep normalization and persistence in dedicated config helpers and test them independently from command execution.
- [Risk] Existing tests may have encoded the duplicated flat-key behavior. → Mitigation: update tests to assert canonical profile persistence and add explicit regression cases for mixed files.

## Migration Plan

1. Update config load/save helpers so profile writes normalize to canonical profile storage.
2. Update profile commands to use the new persistence path without writing resolved flat keys.
3. Add regression tests for legacy-only config, mixed config normalization, and active-profile switching.
4. Update config-facing docs or examples to show the canonical profile schema.

Rollback is low-risk because this is local file behavior only: revert the code change and users can still re-add flat legacy keys manually if needed.

## Open Questions

- None; the desired persisted shape and compatibility boundary are clear from the reported breakage.
