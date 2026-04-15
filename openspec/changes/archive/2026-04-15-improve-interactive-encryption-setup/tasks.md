## 1. Shared encryption helpers

- [x] 1.1 Add a shared helper in `pkg/config` that resolves a candidate AGE private key path, reads an existing identity file, parses it, and returns the derived public key plus normalized path.
- [x] 1.2 Add unit tests for reusable identity parsing, invalid identity files, home-path expansion, and derived public-key behavior.

## 2. Interactive config setup flow

- [x] 2.1 Refactor `cmd/tscli/config/setup` so interactive TTY mode renders more structured Bubble Tea step views for action choice, encryption choice, auth choice, and status messaging while keeping non-interactive prompt mode intact.
- [x] 2.2 Update the encryption step in `cmd/tscli/config/setup` to detect an existing AGE identity file at the selected path and prompt the user to reuse or replace it before generating new key material.
- [x] 2.3 Ensure `tscli config setup` persists `encryption.age.public-key` and `encryption.age.private-key-path` correctly for both reused and newly generated identities, and reports invalid existing files as non-reusable.

## 3. Dedicated encryption setup flow

- [x] 3.1 Update `cmd/tscli/config/encryption/setup` to reuse the shared identity inspection helper for path-based setup.
- [x] 3.2 Add interactive reuse prompting for existing path-based AGE identity files in `config encryption setup`, while preserving existing `--public-key`, `--private-key-source`, `--private-key-path`, and non-path source behavior.
- [x] 3.3 Keep command output and validation actionable for reuse, replace, invalid-file, command-source, and env-source paths.

## 4. Verification and compatibility

- [x] 4.1 Extend command-level tests for `config setup` to cover structured interactive flow behavior, reusable existing key files, declined reuse with overwrite, and invalid existing file fallback.
- [x] 4.2 Extend command-level tests for `config encryption setup` to cover path-based reuse behavior and unchanged command/env source flows.
- [x] 4.3 Review help text and user-facing setup messaging so the new reuse behavior and interactive flow remain clear without changing command names, flags, or config keys.
