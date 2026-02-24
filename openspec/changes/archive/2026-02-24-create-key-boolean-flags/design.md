## Context

`tscli create key` currently supports auth-key creation but only exposes description and expiry controls for `--type authkey`. Tailscale auth keys also support capability booleans (`reusable`, `ephemeral`, `preauthorized`), and users currently need external tooling or raw API requests to set these values.

This change is scoped to the existing `create key` command in `cmd/tscli/create/key/cli.go`, keeping the current command shape and request execution path while expanding the auth-key flag surface.

## Goals / Non-Goals

**Goals:**
- Add `--reusable`, `--ephemeral`, and `--preauthorized` flags for auth-key creation.
- Map these flags to the `CreateAuthKey` request payload in a deterministic way.
- Preserve existing command behavior for users who do not provide the new flags.
- Keep OAuth client creation (`--type oauthclient`) behavior unchanged.
- Add command behavior tests covering new flags and validation expectations.

**Non-Goals:**
- Redesigning `create key` command taxonomy or splitting authkey/oauthclient into separate commands.
- Adding new config file or environment variable keys for key capabilities.
- Changing output format behavior.

## Decisions

### 1. Add explicit boolean flags on `create key`

Introduce optional boolean flags:
- `--reusable`
- `--ephemeral`
- `--preauthorized`

These apply to `--type authkey` requests.

Rationale: direct flags are script-friendly, discoverable in `--help`, and align with existing CLI style.

Alternative considered:
- Accepting a raw JSON body for auth-key options. Rejected because it weakens ergonomics and validation.

### 2. Preserve backward-compatible defaults when flags are omitted

Only set capability booleans explicitly when a flag is provided (or default values match current behavior). Existing calls without the new flags should continue to create keys as before.

Rationale: avoids surprising behavior changes in existing scripts.

Alternative considered:
- Forcing explicit values for all booleans. Rejected as a breaking and verbose UX change.

### 3. Keep oauthclient path isolated from authkey flags

For `--type oauthclient`, existing scope/tag validation remains unchanged. New auth-key capability flags are ignored for oauthclient flow or validated as incompatible if needed by command validation policy.

Rationale: prevents cross-mode coupling and minimizes regression risk.

Alternative considered:
- Reusing flags across both key modes. Rejected due to semantic mismatch with oauth clients.

### 4. Expand command-level integration coverage

Add tests that assert:
- flags are accepted in authkey mode
- request payload sent to mocked API includes capability booleans
- existing behavior still works when flags are omitted
- oauthclient flow remains unaffected

Rationale: this command is a key API integration point where payload regressions are easy to miss without tests.

## Risks / Trade-offs

- [Boolean defaults may not match upstream API expectations] -> Mitigation: keep defaults aligned with current SDK behavior and add explicit payload tests.
- [Mode confusion between authkey and oauthclient] -> Mitigation: keep validations/messaging explicit by key type.
- [Potential script regressions] -> Mitigation: ensure no-flag behavior is unchanged and covered in tests.

## Migration Plan

1. Add new flags and request mapping behind existing `create key` command.
2. Keep existing usage valid without required migration.
3. Update docs/help examples to include optional capability flag usage.
4. Rollback: remove new flags and mapping code; existing pre-change behavior remains.

## Open Questions

- Should authkey capability flags be hard-rejected in oauthclient mode, or silently ignored with documentation?
