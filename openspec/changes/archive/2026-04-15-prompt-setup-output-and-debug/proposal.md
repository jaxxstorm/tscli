## Why

`tscli config setup` can already provision encryption and profiles, but it still leaves global CLI defaults for `output` and `debug` to be configured separately. That creates an incomplete first-run experience and makes it easy for users to miss the output mode that best fits either scripts or interactive use.

## What Changes

- Extend `tscli config setup` so the initial setup flow prompts for a default output mode after profile creation, with explicit choices for `json`, `pretty`, and `human`.
- Extend `tscli config setup` so the initial setup flow also prompts whether debug HTTP request/response logging should be enabled by default.
- Persist the selected `output` and `debug` settings in config alongside the profile and encryption settings written by setup.
- Keep the existing `--output`, `TSCLI_OUTPUT`, `--debug`, and `TSCLI_DEBUG` overrides unchanged so scripts and one-off command invocations continue to take precedence over saved defaults.

## Capabilities

### New Capabilities
- `cli-default-output-and-debug`: Define persisted global CLI defaults for output formatting and debug logging, including how config-backed defaults interact with flags and environment variables.

### Modified Capabilities
- `interactive-config-setup`: Update the guided setup flow so initial setup collects output and debug preferences after profile setup completes.
- `multi-tailnet-config-profiles`: Update setup-driven config persistence requirements so setup can save global `output` and `debug` preferences in the same persisted config used for profiles.

## Impact

- Affected command group: `config setup`.
- Affected flags and env keys: `--output`, `TSCLI_OUTPUT`, `--debug`, `TSCLI_DEBUG`.
- Affected persisted config keys: `output`, `debug`, `active-tailnet`, `tailnets`, and existing `encryption.age.*` keys written by setup.
- Backward compatibility: existing scripts keep working because explicit flags and environment variables still override saved defaults; existing config files remain valid when `output` or `debug` are absent.
- Likely affected code: Bubble Tea setup flow under `cmd/tscli/config/setup/`, persisted config handling under `pkg/config/`, and tests covering setup and runtime config precedence.
