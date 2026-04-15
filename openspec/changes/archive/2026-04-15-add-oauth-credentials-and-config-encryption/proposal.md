## Why

API keys expire frequently enough that they create ongoing profile maintenance overhead, while OAuth client credentials avoid that churn but are higher-value secrets if they are stored in plaintext. tscli needs a safer, more seamless auth path that lets users configure OAuth-backed profiles for everyday API use and optionally protect stored secrets without making the common setup path complicated.

## What Changes

- Extend profile-backed authentication so OAuth client credentials can be configured once and used across tscli API interactions by exchanging them for a short-lived API credential at runtime instead of storing a reusable API key locally.
- Add config support for optionally encrypting stored profile secrets with `age`, including a minimal setup flow centered on an AGE public key for encryption and an AGE private key supplied either from config, an environment variable, or a user-provided secret-retrieval command.
- Update the `config` command group and related docs so users can add and use OAuth client credentials in profiles with clear setup guidance, while keeping API-key-backed profiles and unencrypted config as supported options.
- Document the credential security model, including that encryption is optional, OAuth profiles are optional, and runtime decryption may depend on user-selected secret retrieval mechanisms that can trade convenience for startup latency.

## Capabilities

### New Capabilities
- `config-secret-encryption`: Optional encryption and decryption of persisted API keys and OAuth client secrets in tscli config files using `age`.

### Modified Capabilities
- `multi-tailnet-config-profiles`: Expand OAuth-backed profile behavior from lifecycle-oriented authentication to general API command use, including runtime exchange for non-persisted API credentials and encrypted secret storage for profile-backed auth fields.
- `config-and-auth-documentation`: Document OAuth profile setup, runtime credential exchange behavior, optional `age` encryption setup, supported private-key sources, and security guidance for storing or retrieving decryption material.

## Impact

- Affected code in `cmd/` and `pkg/` for auth resolution, config profile management, credential exchange, and secret handling.
- Config schema and persistence behavior for profile secrets, encrypted values, and any supporting AGE private-key configuration or command integration.
- Documentation for credentials, configuration, and secure secret management.
- New dependency on `age`-compatible encryption/decryption handling and related automated test coverage for encrypted and unencrypted flows.
