## Context

`tscli config setup` already walks users through encryption setup and profile creation, then persists profile state immediately through the existing `pkg/config` helpers. Global runtime defaults are split today: `output` is already a persisted config key, while `debug` is accepted from flags, env, or config at runtime but is currently stripped back out when config is rewritten. The requested change is cross-cutting because it touches the Bubble Tea setup state machine, non-interactive prompt flow, persisted config canonicalization, and command/runtime precedence expectations.

## Goals / Non-Goals

**Goals:**
- Extend initial `tscli config setup` so, after profile creation is finished, the flow asks for a default output mode with curated choices `json`, `pretty`, and `human`.
- Prompt in the same setup flow for whether debug HTTP request/response logging should be enabled by default.
- Persist the selected `output` and `debug` values in the config file without disturbing profile, tailnet, or encryption persistence.
- Preserve existing precedence: command flags override environment variables, environment variables override saved config defaults.
- Cover the new flow with unit and command-level integration tests.

**Non-Goals:**
- Changing the global output system or removing support for other formats such as `yaml` outside setup.
- Reworking the existing rerun profile-management flow beyond what is needed to keep initial setup behavior correct.
- Changing debug logging behavior itself; this change only adds setup-time preference capture and persistence.

## Decisions

### Add explicit post-profile setup steps for output and debug
The setup model will gain two new choice-driven steps after initial profile creation is complete: one for output mode and one for debug preference. They will run only for the initial setup path, after the user declines to add another profile, which matches the requested sequencing of "after setting up profiles during initial setup."

Rationale:
- This keeps profile collection focused on authentication data.
- It preserves current rerun flows for add/modify/delete profiles instead of unexpectedly prompting for global preferences during routine profile management.

Alternative considered:
- Prompt for output/debug before profile creation. Rejected because the request explicitly places these prompts after profiles and because setup success is primarily blocked on credentials, not presentation preferences.

### Persist preferences as top-level config values in the existing config file
The setup flow will write `output` and `debug` into the same persisted config used for `active-tailnet`, `tailnets`, and `encryption.age.*`. `output` already fits the existing top-level config model. `debug` will be treated the same way rather than as a setup-only special case.

Rationale:
- Users expect setup-selected defaults to survive later invocations.
- Keeping the values top-level matches current runtime Viper lookup behavior and avoids introducing a new nested settings schema.

Alternative considered:
- Trigger existing `config set` subcommands from setup. Rejected because setup already writes config directly through shared helpers, and bouncing through command execution would add indirection and test complexity.

### Update config canonicalization to preserve persisted debug values
`pkg/config` currently drops `debug` when writing settings. The implementation should preserve boolean `debug` values during canonicalization while still excluding transient keys such as `help` and `base-url`. This is necessary so setup-written debug preferences survive later profile or encryption rewrites.

Rationale:
- Without this change, a setup prompt for debug would appear to succeed but would be silently lost on save.
- The change is minimal and local to existing persistence rules.

Alternative considered:
- Store debug as a string or custom nested config key. Rejected because runtime already consumes `debug` directly as a boolean and no new schema is needed.

### Keep setup choices curated while leaving broader manual configuration intact
The setup prompt will offer `json`, `pretty`, and `human` only, even though manual config or flags may still accept other supported output formats.

Rationale:
- This matches the requested UX and keeps the first-run prompt simple.
- It avoids broadening this change into output-format policy work.

Alternative considered:
- Expose every currently supported output mode in setup. Rejected because the request explicitly names the three desired choices.

### Testing will cover both setup persistence and precedence safety
Add unit coverage for setup-model transitions and config canonicalization, plus command-level integration tests for initial setup persistence and runtime overrides through flags or environment variables.

Rationale:
- The change spans both interactive flow behavior and persisted config semantics.
- Precedence regressions would be user-visible and could break scripts.

## Risks / Trade-offs

- [Persisted `debug` changes config rewrite behavior] -> Mitigation: limit the canonicalization change to preserving explicit boolean debug values and add regression tests for config rewrites that already cover canonical profile persistence.
- [Initial-setup-only gating could be implemented in the wrong branch and prompt during reruns] -> Mitigation: derive the condition from whether setup started without existing profiles and add integration coverage for both first-run and rerun flows.
- [Setup prompt choices diverge from other accepted output formats] -> Mitigation: document that setup exposes a curated subset and leave manual `config set` or direct config editing available for advanced cases.
- [Immediate profile saves happen before final preference prompts] -> Mitigation: write the preference values as a final persistence step after profile creation completes, using the same config file so the finished setup remains atomic from the user's perspective.
