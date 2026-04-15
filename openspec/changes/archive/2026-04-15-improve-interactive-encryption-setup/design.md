## Context

`tscli config setup` already runs through Bubble Tea in interactive terminals, but its model still behaves like a line-oriented prompt loop with a single freeform input buffer and text-only rendering. Encryption setup also assumes key generation on every path-based setup, even though the runtime already knows how to parse AGE identities from files in `pkg/config/encryption.go`.

This change touches two command surfaces: the main `tscli config setup` onboarding flow and the dedicated `tscli config encryption setup` command. Both need to recognize reusable AGE identity files, derive the matching public key from an existing private key file, and keep non-interactive or script-driven behavior stable.

## Goals / Non-Goals

**Goals:**
- Make the interactive `tscli config setup` experience feel like a structured Bubble Tea flow rather than plain prompt text.
- Reuse a single AGE identity-file inspection path across setup commands so existing files can be reused safely.
- Preserve current config persistence behavior by continuing to write `encryption.age.public-key` and `encryption.age.private-key-path`.
- Preserve non-interactive behavior for `tscli config setup` and stable flag behavior for `tscli config encryption setup`.
- Add unit and command-level test coverage for reusable key detection, invalid files, and replacement flows.

**Non-Goals:**
- Changing the persisted encryption schema or introducing new config keys.
- Reworking profile persistence, encryption/decryption primitives, or auth resolution outside the setup flows.
- Turning `config encryption setup` into a full Bubble Tea application.

## Decisions

### 1. Add shared AGE identity inspection helpers in `pkg/config`
The existing runtime decryption path already reads and parses AGE identities from configured file content. This change should add a small shared helper that accepts a candidate file path, expands `~`, reads the file, parses the identity, and returns both the normalized path and derived public recipient.

Why this approach:
- It avoids duplicating AGE parsing rules between `cmd/tscli/config/setup` and `cmd/tscli/config/encryption/setup`.
- It keeps AGE-specific logic near the existing config encryption helpers.
- It lets tests cover parsing and derivation once at the package level.

Alternative considered:
- Keep file parsing inside each command package. Rejected because the two commands would drift in validation and error handling.

### 2. Keep `config setup` as the richer Bubble Tea surface and preserve fallback prompts for non-TTY usage
The current `cmd/tscli/config/setup` implementation already branches between Bubble Tea for TTY use and a prompt loop for non-interactive use. The change should keep that split, but the interactive model should move from a freeform text prompt renderer toward explicit choice-driven views for common steps such as add/delete/quit, encryption yes/no, reuse existing key yes/no, and auth type selection.

Why this approach:
- It improves the interactive experience without breaking redirected input, tests, or script-friendly behavior.
- Most of the requested UI polish applies to interactive terminal use, not to stdin-driven automation.
- It limits scope by improving the existing model instead of replacing the command architecture.

Alternative considered:
- Implement Bubble Tea-only behavior and remove the prompt fallback. Rejected because the project requirement is to keep commands predictable in non-interactive contexts.

### 3. Make path-based encryption setup follow a shared reuse decision flow
When either setup flow targets a private key path, the command should resolve the path, check for an existing file, and inspect it. If the file contains a valid AGE identity, the user should be offered reuse before any overwrite or generation occurs. Accepting reuse should derive and persist the public key from the file; declining reuse should continue with newly generated key material for `config setup`, or with explicit path-based setup behavior for `config encryption setup`.

Why this approach:
- It directly matches the requested behavior and prevents accidental key replacement.
- It treats existing key files as a first-class, observable part of the setup flow.
- It removes the need for users to manually copy the public key when reusing a path-based identity.

Alternative considered:
- Reuse existing files automatically whenever parsing succeeds. Rejected because replacing vs reusing is a user decision and silent reuse could preserve an unintended identity.

### 4. Surface invalid existing files as actionable setup feedback, then continue safely
If a selected path exists but the file cannot be parsed as an AGE identity, the setup flow should tell the user that reuse is unavailable and continue down the generation path rather than persisting invalid encryption settings.

Why this approach:
- It keeps the setup flow moving.
- It avoids partially configured encryption.
- It produces behavior that is easy to cover with integration tests.

Alternative considered:
- Fail immediately on invalid files. Rejected for `config setup` because the user can still complete setup by generating a new key file at the same path.

### 5. Test at both helper and command levels
Unit tests should cover identity-file parsing, public-key derivation, home-path expansion, and invalid-file handling. Command-level tests should cover interactive `config setup` reuse/replace flows and `config encryption setup` path-based reuse behavior alongside existing generated-key coverage.

Why this approach:
- Shared helper tests keep the AGE-specific logic deterministic.
- Command tests verify user-visible prompts, persisted config values, and overwrite vs reuse behavior.

## Risks / Trade-offs

- [Interactive UI changes make snapshot-style expectations brittle] -> Keep rendering changes focused on key decision steps and update command tests to assert important behavior rather than incidental formatting where possible.
- [Shared helper placement could over-couple command UX to config internals] -> Limit shared helpers to path inspection and identity derivation, leaving prompting and command-specific branching in command packages.
- [Reusing an unexpected existing file could preserve the wrong identity] -> Always prompt before reuse and show the selected path in the prompt context.
- [Invalid existing files may be mistaken for reusable AGE identities] -> Parse through the same AGE library used by runtime decryption and treat parse failures as non-reusable.

## Migration Plan

No config migration is required. Existing users keep the same `encryption.age` keys and runtime precedence rules. The rollout is limited to command behavior changes:

1. Add shared identity-file inspection helpers and tests.
2. Update `tscli config setup` to present richer Bubble Tea views and the reuse prompt path.
3. Update `tscli config encryption setup` to support path-based reuse prompts.
4. Extend command-level tests for reuse, replace, and invalid-file cases.

Rollback is low risk because the change is localized to setup flows; reverting the command and helper changes restores the previous always-generate behavior.

## Open Questions

- Whether the interactive Bubble Tea view should display the derived public key before the user confirms reuse, or just the selected path.
- Whether `config encryption setup --private-key-source=path --private-key-path <path>` should prompt on reuse in all cases or only when running interactively.
