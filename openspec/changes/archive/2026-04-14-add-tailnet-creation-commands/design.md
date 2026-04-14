## Context

`tscli` currently assumes authenticated API operations use a Tailscale API key resolved through flags, environment, an active profile, or legacy flat config. That model works for the existing `pkg/tscli` client, which always sends basic auth, and the only OAuth-related command today is `create token`, which is explicitly unauthenticated and only exchanges client credentials for an access token.

The new API-driven tailnet lifecycle endpoints introduce two constraints that do not fit the current model cleanly:

- `create tailnet` and `list tailnets` authenticate with an organization-approved OAuth client and bearer token, not an API key.
- `delete tailnet` authenticates with the OAuth client returned when the tailnet was created, so the CLI needs a way to resolve OAuth client credentials for that tailnet-specific workflow.

The repo already has profile-backed config (`active-tailnet` plus `tailnets`) and deterministic config precedence, so the design should reuse that machinery where it helps while avoiding broad disruption to all existing API-key commands.

## Goals / Non-Goals

**Goals:**
- Add `create tailnet`, `list tailnets`, and `delete tailnet` commands that target the documented API-driven tailnet lifecycle endpoints.
- Support OAuth client credential auth for those commands using flags, environment variables, and persisted profiles/config.
- Extend profile-backed config so a profile can represent either API-key auth or OAuth client credentials without breaking existing profile files.
- Keep existing API-key-authenticated commands, flags, and scripts working without behavior changes.
- Cover the new command and auth paths with unit tests and CLI-level integration tests.

**Non-Goals:**
- Introducing a general-purpose auth plugin framework for every command.
- Encrypting persisted secrets or integrating with an OS keychain.
- Automatically provisioning, rotating, or recovering OAuth client secrets after creation.
- Refactoring all existing commands away from the current `pkg/tscli` basic-auth client path.

## Decisions

### 1. Add lifecycle commands under the existing `create`, `list`, and `delete` groups

The new surface will follow the existing command taxonomy:

- `tscli create tailnet --display-name <name>`
- `tscli list tailnets`
- `tscli delete tailnet`

This keeps the commands script-friendly and predictable, and avoids introducing a special top-level command group for one feature.

Alternatives considered:
- Add a separate `tailnet` command tree (`tailnet create`, `tailnet list`, `tailnet delete`). Rejected because the repo already organizes actions primarily by verb.

### 2. Keep API-key pre-run behavior intact and resolve OAuth auth within the new commands

The existing root persistent pre-run resolves API-key runtime config for almost every non-local command. Instead of broadening that path to handle multiple auth types for all commands at once, the new tailnet lifecycle commands will opt out of the API-key pre-run and perform their own OAuth credential resolution.

This minimizes disruption to existing commands and keeps the implementation narrow:

- existing commands continue to rely on `config.ResolveRuntimeConfig` and `pkg/tscli.New()`
- `create token` remains unauthenticated
- the new tailnet lifecycle commands resolve OAuth credentials and exchange them for bearer tokens before issuing raw HTTP requests

Alternatives considered:
- Replace the global pre-run auth pipeline with a command-annotation-driven auth framework. Rejected for now because it adds cross-cutting risk unrelated to shipping the new lifecycle commands.

### 3. Extend profiles to support explicit auth shapes instead of assuming every profile is API-key-only

The current profile schema assumes a profile is a tailnet name plus `api-key`, and runtime resolution treats the profile name as the effective tailnet. That is too restrictive for OAuth-backed organization workflows, where a profile may not map directly to a tailnet path and may need `oauth-client-id` and `oauth-client-secret` instead of `api-key`.

The config/profile model will be extended so profile entries can contain:

- `name`: the profile identifier selected by `active-tailnet`
- `tailnet`: optional explicit tailnet value for commands that need one
- `api-key`: optional API-key credential
- `oauth-client-id`: optional OAuth client id
- `oauth-client-secret`: optional OAuth client secret

Validation will make auth inputs mutually meaningful:

- an API-key-backed profile must provide `api-key`
- an OAuth-backed profile must provide both `oauth-client-id` and `oauth-client-secret`
- existing profile files that omit `tailnet` continue to resolve `tailnet = name`

Alternatives considered:
- Avoid profile support and require OAuth flags only. Rejected because the user explicitly wants this considered in the profile mechanism and the delete flow depends on retaining the returned OAuth credentials.
- Store exchanged bearer access tokens in config. Rejected because access tokens are short-lived and should be derived at runtime from durable client credentials.

### 4. Add dedicated OAuth resolution helpers and a bearer-capable raw request path

The lifecycle commands need two reusable pieces:

- a config resolver that applies the existing precedence model for OAuth credentials: flags, environment, active profile, legacy/specialized config keys if present
- a request helper that sends `Authorization: Bearer <token>` for raw endpoints not covered by the typed SDK

The design will reuse `pkg/oauth` for client-credential exchange and add a small helper for authenticated raw requests rather than trying to force bearer tokens through the current basic-auth-only `pkg/tscli.New()` flow.

Alternatives considered:
- Mutate `pkg/tscli.Client` to support both basic and bearer auth globally. Rejected for this change because most commands do not need it, and a narrower helper keeps behavior easier to reason about.

### 5. Persist the OAuth client returned by tailnet creation as profile-ready output, not as an implicit side effect

The create endpoint returns an OAuth client secret exactly once. The CLI should surface that response faithfully in all output modes and make it easy for users to persist it with `config profiles upsert`, but it should not silently rewrite user config during `create tailnet`.

This keeps command behavior explicit and avoids surprising secret persistence.

Alternatives considered:
- Auto-create or auto-update the active profile after `create tailnet`. Rejected because it couples a remote create operation to local config mutation and could overwrite intentionally separate profile state.

## Risks / Trade-offs

- [Profile schema becomes more complex] -> Keep existing fields working, treat `tailnet` as optional with `name` fallback, and limit new validation to clearly invalid auth combinations.
- [OAuth lifecycle commands drift from the default auth path] -> Isolate the divergence to a small resolver/helper pair and cover it with command-level tests.
- [Users may confuse profile name with tailnet value once `tailnet` becomes optional] -> Document the distinction and preserve the existing `name` fallback when `tailnet` is omitted.
- [Delete flow depends on one-time credentials returned at creation] -> Document the requirement clearly and ensure `create tailnet` output includes the returned OAuth client in machine-readable formats.
- [Raw endpoint coverage may require hand-written request/response structs] -> Keep structs close to the commands or a focused package and test exact payload shapes.
