## Context

`tscli` currently has two separate auth paths. Most API commands require `pkg/config.ResolveRuntimeConfig`, which resolves an API key and tailnet, while the tailnet lifecycle commands resolve OAuth client credentials separately and call `pkg/oauth.ExchangeClientCredentials` before using bearer-token requests. Profile-backed config currently stores API keys and OAuth client secrets in plaintext fields, and the authentication docs still describe OAuth support as a lifecycle-only path.

This change is cross-cutting because it affects config schema, config mutation commands, runtime auth resolution, the request path for API commands, and the credentials documentation. It also introduces a new security-sensitive dependency on `age` and adds a second secret-resolution path for decrypting persisted profile values.

## Goals / Non-Goals

**Goals:**
- Let users configure OAuth-backed profiles once and use them across supported tscli API commands without storing a reusable exchanged credential in config.
- Add a minimal, optional AGE-based encryption flow for persisted API keys and OAuth client secrets.
- Keep existing API-key workflows, legacy config, and script-friendly command behavior working as-is.
- Make secret-source precedence and failure modes explicit and testable.

**Non-Goals:**
- Replacing Tailscale API-key support as the default or only auth model.
- Encrypting every config value; only secret fields need encryption.
- Introducing OS-specific keychain integrations in this change.
- Designing a background token cache or long-lived local credential store.

## Decisions

### Reuse the existing OAuth exchange path and generalize command auth selection
Commands that already support bearer-token requests can keep using the existing `pkg/oauth.ExchangeClientCredentials` and `pkg/tscli.DoBearer` path. For the broader API surface, the implementation should add a shared auth-resolution layer that selects one of two runtime auth modes:
- resolved API key plus tailnet for the existing `tsapi.Client` path
- resolved OAuth client credentials plus a just-in-time bearer token for commands that opt into OAuth-backed auth

This keeps exchanged credentials ephemeral and avoids writing any derived key or token back into config. The main alternative was to exchange OAuth credentials once and persist a temporary API credential locally, but that would reintroduce the secret-at-rest problem this change is trying to reduce.

### Add encrypted sibling fields instead of overloading existing secret keys
Persisted profile secrets should use explicit encrypted fields: `api-key-encrypted` and `oauth-client-secret-encrypted`. Plaintext fields remain supported for backward compatibility and for users who choose not to enable encryption. This makes the stored config shape observable, avoids guessing whether a field contains ciphertext or plaintext, and lets validation reject invalid mixed forms such as both `api-key` and `api-key-encrypted` on one profile.

The main alternative was to store AGE ciphertext inline in the existing `api-key` or `oauth-client-secret` fields. That would reduce schema changes but makes config validation, docs, and migration harder because the code would need to infer whether a value is plaintext or encrypted.

### Use a minimal top-level AGE configuration block with explicit precedence
The config should add `encryption.age.public-key` plus one optional configured private-key source: `encryption.age.private-key` or `encryption.age.private-key-command`. At runtime, private-key lookup should use this order:
1. `TSCLI_AGE_PRIVATE_KEY`
2. `encryption.age.private-key-command`
3. `encryption.age.private-key`

This gives users a secure non-persisted option first, supports external secret managers, and still allows a fully local setup for users who prefer convenience. The main alternative was to allow multiple configured fallback sources in config, but that adds complexity and makes failure behavior harder to predict.

### Add a guided setup command rather than requiring manual YAML edits
The `config` command group should grow a focused setup flow, `config encryption setup`, that prompts for the AGE public key and asks how the private key will be provided. The command should only write the chosen settings and should leave encryption disabled until the user explicitly configures it.

The main alternative was to document manual YAML edits only. That would be smaller in code, but it would not meet the requirement for a clear, easy setup path and would make schema discovery harder for users.

### Keep encryption integrated into existing profile mutation commands
`config profiles set` should remain the primary way to create and update API-key-backed and OAuth-backed profiles. When encryption is enabled, it should transparently write encrypted secret fields instead of plaintext secret fields. Listing and active-profile selection do not need new flags; they only need to reflect the resulting auth shape consistently.

The alternative was to add separate encrypted-profile commands or per-command `--encrypt` flags. That would fragment the UX and make the optional encryption model harder to explain.

### Testing should cover both schema transitions and runtime auth behavior
Unit tests should cover profile validation, encryption settings validation, private-key source precedence, encryption and decryption helpers, and auth resolution. Command-level integration tests should cover plaintext and encrypted profile persistence, setup flow behavior, OAuth-backed API execution, and actionable errors when exchange or decryption fails.

## Risks / Trade-offs

- [Command support is uneven] -> Some commands may still be wired only for API-key clients, so implementation should introduce OAuth-backed support incrementally behind shared auth helpers and add coverage where commands opt in.
- [Command-based private-key retrieval adds latency] -> Document the trade-off clearly and keep command execution to a single fetch per CLI invocation.
- [Stored private keys reduce the benefit of encryption-at-rest] -> Prefer `TSCLI_AGE_PRIVATE_KEY` in docs and setup messaging, but still allow config storage because the user explicitly asked for that option.
- [Encrypted field schema increases config complexity] -> Use explicit sibling field names, validation errors, and guided setup docs to keep the config understandable.
- [Bearer-token and API-key request paths can diverge] -> Centralize auth resolution and request setup so command code chooses auth mode without duplicating precedence logic.

## Migration Plan

1. Extend config models and validation to understand encrypted sibling fields and AGE settings without breaking existing plaintext configs.
2. Add encryption helpers plus `config encryption setup` and update `config profiles set` persistence behavior.
3. Introduce shared runtime auth resolution for commands that can use OAuth-backed bearer auth.
4. Update documentation and command docs, then add or expand automated coverage for plaintext, encrypted, and OAuth-backed flows.

Rollback is straightforward because plaintext config remains supported. If the new encryption flow ships with issues, users can disable encryption setup and continue using plaintext API-key or OAuth-backed profiles.

## Open Questions

- Which existing command groups can safely opt into OAuth-backed bearer auth in the first implementation pass without additional API-shape changes?
- Should `config encryption setup` also offer a non-interactive flag-driven mode for automation, or is interactive setup plus manual config editing sufficient for the first release?
