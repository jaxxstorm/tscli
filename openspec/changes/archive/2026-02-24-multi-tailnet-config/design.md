## Context

`tscli` currently reads one `tailnet` and one `api-key` from Viper-backed config/environment/flags. This works for a single tailnet, but users with multiple tailnets must repeatedly edit config files or export env vars before every command. The CLI already has a `config` command group (`show`, `get`, `set`) and a root `PersistentPreRunE` that validates API credentials.

The change introduces profile-based tailnet credentials while preserving existing behavior for users who only have legacy flat keys. The implementation must remain script-friendly and keep current precedence expectations for flags and environment variables.

## Goals / Non-Goals

**Goals:**
- Support multiple tailnet profiles in config via `tailnets` plus `active-tailnet`.
- Resolve runtime `tailnet` and `api-key` from the active profile when flags/env are not set.
- Keep existing single-tailnet config files working without user migration.
- Provide `config` operations to list profiles, set active profile, upsert profile credentials, and remove profiles.
- Add unit and command-level integration tests for profile parsing, resolution, and command behavior.

**Non-Goals:**
- Redesigning global command taxonomy outside the `config` profile workflow.
- Changing output formatting conventions beyond what is needed to display profile state.
- Introducing kubeconfig-level complexity (contexts, clusters, users, merges).

## Decisions

### 1. Add a canonical multi-tailnet config model with legacy compatibility fields

Introduce profile-aware config keys:
- `active-tailnet` (string)
- `tailnets` (list of `{name, api-key}`)

Legacy keys remain supported:
- `tailnet`
- `api-key`

Rationale: this keeps existing config files valid and allows gradual migration without breaking scripts.

Alternatives considered:
- Replace legacy keys entirely. Rejected because it breaks existing users and automation.
- Use a map (`tailnets.<name>.api-key`) instead of a list. Rejected to keep YAML structure explicit and stable for deterministic output ordering.

### 2. Centralize runtime credential resolution in root command pre-run

Add a resolver used by `PersistentPreRunE` to compute effective auth settings with this precedence:
1. CLI flags
2. Environment variables
3. Active profile (`active-tailnet` + matching entry in `tailnets`)
4. Legacy flat config (`api-key`, `tailnet`)

Behavior:
- If no effective tailnet is found, default to `-`.
- If no effective API key is found, return the existing required-key error.
- When `active-tailnet` is set but missing from `tailnets`, return a validation error.

Rationale: one resolver avoids duplicated logic across command packages and preserves existing precedence.

Alternatives considered:
- Resolve inside each command package. Rejected due to drift risk and repetitive code.
- Resolve once during config init before env/flag binding. Rejected because it would ignore final flag/env values.

### 3. Extend config commands with profile operations while preserving generic get/set behavior

Add profile-focused operations under `config`:
- list configured tailnet profiles (including active marker)
- set the active tailnet
- upsert a tailnet profile
- remove a tailnet profile

Existing `config get`/`config set` remain available for generic key access.

Rationale: users need explicit profile management commands while existing scripts using generic key operations continue to work.

Alternatives considered:
- Only generic key editing for `tailnets` YAML blobs. Rejected because it is error-prone and not ergonomic.
- A standalone top-level command group for profile switching. Rejected to keep config management in one place.

### 4. Mirror active profile values into legacy keys on profile writes

When profile commands change the active profile or upsert the active profile, persist corresponding legacy keys (`tailnet`, `api-key`) to match active values.

Rationale: this keeps downstream code paths and older assumptions functioning during transition.

Alternatives considered:
- Keep only new keys after profile updates. Rejected because mixed-version workflows become brittle.

### 5. Add focused tests for resolver and CLI command behavior

Testing strategy:
- Unit tests in `pkg/config` for parsing, validation, profile selection, precedence, and migration/normalization.
- Command-level integration tests under `test/cli` for:
  - profile management command success/error behavior
  - profile switching impact on subsequent command execution
  - precedence assertions (flag > env > profile > legacy config)

Rationale: multi-source configuration is regression-prone and needs direct behavior tests.

## Risks / Trade-offs

- [Profile config complexity increases] -> Mitigation: keep schema minimal (`name`, `api-key`, `active-tailnet`) and document examples.
- [Ambiguous source of truth between new and legacy keys] -> Mitigation: define strict precedence and mirror legacy keys from active profile on profile updates.
- [Potential breakage in scripts relying on implicit defaults] -> Mitigation: preserve default tailnet `-` and existing required API key validation text.
- [User confusion when active profile is missing] -> Mitigation: explicit validation error naming the missing profile and remediation.

## Migration Plan

1. Ship support for reading both legacy and new config keys.
2. If users already have only legacy keys, commands continue unchanged.
3. Users can add profiles incrementally via config profile commands.
4. On setting/activating profiles, persist mirrored legacy keys for compatibility.
5. Rollback strategy: remove new profile keys (`tailnets`, `active-tailnet`) and retain legacy `tailnet`/`api-key` values.

## Open Questions

- Should profile names be case-sensitive or normalized to lower-case for lookup and uniqueness checks?
- Should profile listing redact API keys fully or partially in human-readable output?
