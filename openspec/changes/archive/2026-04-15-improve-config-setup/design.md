## Context

`tscli` already supports profile-backed configuration, profile CRUD commands, and `age`-based secret encryption, but those capabilities are split across `config profiles ...` and `config encryption setup` commands with line-oriented prompts. The requested change introduces a cross-cutting onboarding and management flow under `tscli config setup`, touches Cobra command wiring, Bubble Tea UI state, config persistence through Viper, and encrypted secret handling through the existing `pkg/config` package.

The current codebase already has the persistence primitives needed for this change:
- `pkg/config/encryption.go` validates and stores `encryption.age.*` settings and encrypts persisted profile secrets when encryption is enabled.
- `config profiles set`, `set-active`, and `delete` already express the canonical profile mutations.

The new work is therefore primarily orchestration: a guided setup UI that collects user intent, creates any missing `~/.tscli` directories for key material, generates an `age` identity when requested, persists the encryption config, and then creates or deletes profiles using the existing canonical schema.

## Goals / Non-Goals

**Goals:**
- Add a `tscli config setup` Bubble Tea flow that guides users through encryption setup and profile creation without requiring them to know the lower-level config commands.
- Reuse existing profile and encryption persistence logic so saved config shape, precedence, and encrypted field handling remain consistent with the rest of the CLI.
- Support repeated setup runs for profile management by offering add and delete flows when profiles already exist.
- Ensure key file persistence is safe by creating the chosen parent directory before writing the generated `age` identity file.
- Cover the new behavior with unit tests for setup state transitions and command-level integration tests for interactive flows.

**Non-Goals:**
- Replacing existing non-interactive `config profiles` or `config encryption setup` commands.
- Changing runtime credential precedence for operational commands.
- Introducing remote key storage, passphrase-protected keys, or encrypted storage for OAuth client IDs.
- Building a general-purpose TUI framework for unrelated commands.

## Decisions

### Use `tscli config setup` as an orchestration command and keep existing subcommands
The new command should sit under `cmd/tscli/config/setup` and present the interactive experience requested by the user. The existing `config profiles ...` and `config encryption setup` commands should remain available for scriptable and targeted workflows.

Why this approach:
- It adds a friendly first-run path without regressing scriptability.
- It keeps command responsibilities clear: `config setup` orchestrates, while `pkg/config` and existing subcommands continue to own persistence semantics.

Alternative considered:
- Replacing `config encryption setup` and `config profiles set` with Bubble Tea-only behavior. Rejected because it would make narrow automation and existing test helpers harder to preserve.

### Reuse existing `pkg/config` persistence helpers instead of writing setup-specific config files directly
The Bubble Tea flow should collect answers into an in-memory setup model, then call existing helpers such as `config.SetAgeEncryptionConfig`, `config.UpsertTailnetProfile`, and `config.RemoveTailnetProfile`.

Why this approach:
- It preserves the canonical `tailnets` / `active-tailnet` file shape and existing validation rules.
- It keeps encryption behavior centralized in `pkg/config/encryption.go`, including existing encrypted field names and decryption resolution.

Alternative considered:
- Serializing config directly from the TUI model. Rejected because it would duplicate validation, increase divergence risk, and weaken backward compatibility.

### Generate `age` keys in-process and persist the full identity to a user-selected file path
When the user opts into encryption, setup should generate a new X25519 identity with the existing `filippo.io/age` dependency, prompt for a destination path defaulting to `~/.tscli/age.txt`, create the parent directory if needed, write the private key material with restrictive permissions, and persist the public key plus `encryption.age.private-key-path` in config.

Why this approach:
- It matches the requested UX exactly.
- It gives the CLI enough information to encrypt secrets immediately and decrypt them later without requiring extra environment setup.
- It reuses the existing supported `private-key-path` mechanism instead of inventing a new config source.

Alternative considered:
- Storing the private key inline in config or requiring an environment variable. Rejected because the user explicitly asked for path-based persistence and because inline private keys increase config exposure.

### Model setup as a small finite-state Bubble Tea flow
The Bubble Tea program should move through explicit states such as: existing-profiles action, encryption choice, key path entry, auth type choice, profile detail entry, save result, add-another confirmation, and optional delete selection. Existing profile data should be loaded at startup so reruns can branch immediately into add/delete options.

Why this approach:
- The requested behavior is sequential and branch-heavy, which maps cleanly to a state machine.
- It keeps the command implementation testable because unit tests can drive state transitions without terminal snapshot assertions for every branch.

Alternative considered:
- A loose chain of prompt callbacks. Rejected because rerun management, defaults, validation, and graceful exit logic become harder to reason about as branches grow.

### Encrypt only secret-bearing fields and preserve current precedence semantics
If encryption is enabled, setup should continue using the existing persistence behavior that encrypts `api-key` and `oauth-client-secret` while leaving non-secret fields such as profile name, tailnet, and `oauth-client-id` in cleartext. Runtime commands should continue to resolve flags over environment over active profile over legacy config keys.

Why this approach:
- It aligns with the current implementation and minimizes migration risk.
- It avoids introducing hidden changes to runtime authentication behavior while still protecting the most sensitive stored values.

Alternative considered:
- Encrypting the entire profile object. Rejected because the current config model, tests, and troubleshooting workflows depend on visible non-secret metadata.

### Testing strategy combines TUI state tests with command-level integration tests
Unit tests should cover setup state branching, default path resolution, directory creation decisions, and validation errors. Integration tests should exercise first-run plain-text setup, first-run encrypted setup, repeated setup add flow, repeated setup delete flow, and encrypted profile persistence outcomes.

Why this approach:
- Bubble Tea logic is easiest to validate at the state-model layer.
- Interactive command coverage is still needed to prove Cobra wiring, persisted config output, and file-system side effects.

## Risks / Trade-offs

- [Bubble Tea adds UI complexity to a previously simple command] -> Keep the setup state machine small, isolate prompt/view logic from persistence actions, and leave existing low-level commands untouched.
- [Key file write failures could leave config pointing at a missing or unreadable path] -> Create parent directories before writing, use atomic write behavior where practical, and only persist encryption config after the key file write succeeds.
- [Repeated setup runs could accidentally delete or overwrite the wrong profile] -> Show existing profile names in the management step, require explicit selection for deletion, and continue honoring the existing guard that blocks deletion of the active profile unless another profile is selected first.
- [Interactive setup may be less script-friendly] -> Scope the new experience to `config setup` and preserve existing flag-driven commands for automation.
- [Migration confusion between legacy flat config and profile-backed config] -> Keep existing backward-compatible resolution unchanged and ensure setup persists canonical profile-backed config once a profile is created.

## Migration Plan

1. Add the new `config setup` command without removing existing commands.
2. Reuse existing persisted config keys so no config migration is required for current users.
3. On rerun, detect existing profiles and branch into add/delete management while keeping the previous config file valid if the user exits without changes.
4. If rollout issues appear, users can fall back to `config profiles ...` and `config encryption setup` because those commands remain supported.

## Open Questions

- Whether rerun management should also offer `set-active` in the same flow, or leave active-profile changes to the existing `config profiles set-active` command.
- Whether the generated key file should contain only the private identity or both public and private values as the request describes; implementation should prefer the format already accepted by `parseAgeIdentity` and can append the public recipient if needed for operator convenience.
