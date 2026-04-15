## Why

`tscli config setup` currently does not guide users through the profile-based configuration model or the choice between plain-text and encrypted credential storage. Improving setup now will make first-run configuration discoverable, bring encryption into the primary flow, and give users a safe way to manage additional profiles without hand-editing config files.

## What Changes

- Replace the current `tscli config setup` experience with a Bubble Tea driven interactive flow.
- Add setup prompts for enabling credential encryption, generating `age` key material, and persisting that key material to a user-selected path that defaults to `~/.tscli/age.txt`.
- Extend setup so users can create API-key-backed or OAuth-backed profiles, encrypt stored credentials when encryption is enabled, and add multiple profiles in one session.
- Update rerun behavior so `tscli config setup` can manage existing profiles by adding new ones or deleting existing ones instead of acting as a one-time bootstrap only.
- Preserve existing non-interactive config precedence and backward compatibility for users who continue to rely on flags, environment variables, or legacy flat config keys.

## Capabilities

### New Capabilities
- `interactive-config-setup`: Covers the interactive setup workflow, encryption bootstrap, and profile management prompts exposed by `tscli config setup`.

### Modified Capabilities
- `multi-tailnet-config-profiles`: Expand profile requirements to cover setup-driven profile creation, encrypted credential persistence, and profile deletion behavior.

## Impact

- Affected command group: `config`, especially `tscli config setup`.
- Affected config keys and persisted values: `tailnets`, `active-tailnet`, `api-key`, `oauth-client-id`, `oauth-client-secret`, plus encryption-related key path and key material handling.
- Affected code areas will likely include `cmd/config*`, config persistence and credential handling under `pkg/`, and new Bubble Tea UI flow code.
- Adds a dependency on `age`-based key generation and encrypted credential handling during setup.
- Existing scripts remain compatible because interactive changes are scoped to `config setup` and runtime credential precedence remains flags over environment over config.
