## 1. Config schema and encryption plumbing

- [x] 1.1 Extend the config models and validation logic to support `encryption.age.public-key`, `encryption.age.private-key`, `encryption.age.private-key-command`, `api-key-encrypted`, and `oauth-client-secret-encrypted` while preserving legacy plaintext config support.
- [x] 1.2 Add AGE helper code to validate recipient keys, encrypt secret values for persistence, decrypt secret values at runtime, and resolve AGE private-key input in the order `TSCLI_AGE_PRIVATE_KEY` -> `encryption.age.private-key-command` -> `encryption.age.private-key`.
- [x] 1.3 Add unit tests for config validation, encrypted-field auth-shape validation, AGE key-source precedence, and decryption failure cases.

## 2. Config command workflows

- [x] 2.1 Add `config encryption setup` with the guided prompt flow for public-key entry and private-key source selection, and persist only the chosen AGE settings.
- [x] 2.2 Update `config profiles upsert` so API keys and OAuth client secrets are written to encrypted sibling fields when encryption is enabled and remain plaintext when encryption is disabled.
- [x] 2.3 Update config profile command tests to cover plaintext and encrypted persistence, active-profile behavior, and actionable validation errors for conflicting or incomplete encryption settings.

## 3. Runtime authentication expansion

- [x] 3.1 Refactor runtime auth resolution so commands can resolve either API-key auth or OAuth client credentials from flags, environment variables, active profiles, and legacy config without duplicating precedence logic.
- [x] 3.2 Reuse the existing OAuth exchange and bearer-request path so supported API commands can authenticate with OAuth-backed profiles at runtime without persisting the exchanged token.
- [x] 3.3 Add unit and integration coverage for OAuth-backed command execution, token-exchange failure handling, encrypted-secret decryption during auth resolution, and compatibility with existing API-key-backed flows.

## 4. Documentation and compatibility checks

- [x] 4.1 Update the credentials and configuration docs to cover OAuth-backed profile setup for supported API commands, AGE encryption setup, private-key sourcing options, and the fact that both OAuth profiles and encryption are optional.
- [x] 4.2 Regenerate or update command help/reference docs for any new or changed config commands and flags, including `config encryption setup`.
- [x] 4.3 Update docs or validation checks so missing OAuth/encryption guidance fails automated verification, and confirm existing scripts and legacy config examples remain valid.
