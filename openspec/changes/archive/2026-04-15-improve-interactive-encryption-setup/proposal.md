## Why

`tscli config setup` already uses Bubble Tea, but the current interaction is mostly plain text prompts and does not make key decisions or existing state easy to scan. The encryption flow also always generates new AGE material, which adds friction and can overwrite a user's intended setup when a reusable key file already exists.

## What Changes

- Improve the user-facing `tscli config setup` experience with a more polished Bubble Tea-driven presentation for setup steps, choices, and status messaging while preserving non-interactive behavior.
- Update the `config setup` encryption step to detect an existing AGE identity file, inspect it for the stored private/public key pair, and prompt the user to reuse it instead of always generating a new key.
- Update the dedicated `tscli config encryption setup` flow to offer the same existing-key reuse path when a file-backed AGE identity is already present.
- Keep existing command names, current config keys under `encryption.age.public-key` and `encryption.age.private-key-path`, and script-friendly flag behavior stable.

## Capabilities

### New Capabilities
- None.

### Modified Capabilities
- `interactive-config-setup`: Change the interactive `tscli config setup` requirements so the Bubble Tea flow presents a richer guided UI and offers reuse of an existing AGE identity file during encryption setup.
- `config-secret-encryption`: Change encryption setup requirements so file-backed AGE setup checks for an existing identity file, derives the public key from that file when valid, and prompts the user before reusing or replacing it.

## Impact

- Affected command groups: `tscli config setup`, `tscli config encryption setup`.
- Affected flags and config keys: existing encryption settings remain `encryption.age.public-key` and `encryption.age.private-key-path`; existing `config encryption setup` flags such as `--public-key`, `--private-key-source`, and `--private-key-path` remain part of the flow.
- Affected code: `cmd/tscli/config/setup`, `cmd/tscli/config/encryption/setup`, and AGE config helpers under `pkg/config`.
- Affected tests: command-level setup/encryption tests and any UI- or file-detection-specific coverage for reuse prompts and existing key parsing.
- Backward compatibility: no breaking CLI rename or config schema change is intended; existing scripts and existing AGE key files should continue to work, with interactive users gaining a reuse prompt instead of implicit regeneration.
