## Context

`tscli create key` already binds `--tags` through Cobra as a shared `[]string` flag, but the `authkey` execution path in `cmd/tscli/create/key/cli.go` only copies the boolean capability fields into `tsapi.CreateKeyRequest`. Because `tailscale.com/client/tailscale/v2.KeyCapabilities` defines `capabilities.devices.create.tags` on the auth-key payload, the current code leaves that field nil and the debug payload is serialized as `"tags": null` instead of the provided tag list.

The regression is localized to the existing `create key` command, but it affects both request correctness and user confidence because the help text currently says `--tags` is for OAuth client and federated credentials only even though the flag is parsed for all key types.

## Goals / Non-Goals

**Goals:**
- Preserve `--tags` values for `tscli create key --type authkey` by mapping them into `capabilities.devices.create.tags`.
- Add a command-level regression test that fails on the current `tags: null` behavior and passes once the payload is fixed.
- Keep existing auth-key boolean capability behavior (`--reusable`, `--ephemeral`, `--preauthorized`) unchanged.
- Keep OAuth client and federated tag handling unchanged.
- Align user-facing help text with the supported auth-key tag behavior.

**Non-Goals:**
- Introducing new flags, config keys, or environment variables.
- Changing the request shape for OAuth client or federated key creation.
- Adding new client-side validation for tag contents or for tailnet-owned-key requirements enforced by the Tailscale API.

## Decisions

### 1. Map auth-key tags through the existing Tailscale client request model

For `--type authkey`, the command will continue to build `tsapi.CreateKeyRequest`, but it will also copy the parsed `tags` slice into `req.Capabilities.Devices.Create.Tags` before calling `CreateAuthKey`.

Rationale: `KeyCapabilities` already models the upstream API field location. Reusing that type keeps the auth-key path idiomatic and avoids a custom request body or separate serialization code.

Alternative considered:
- Setting a top-level `tags` field on `CreateKeyRequest`. Rejected because the auth-key Go client type does not define that payload shape and the API expects tags under `capabilities.devices.create`.

### 2. Cover the regression at the command boundary

Extend `test/cli/create_key_flags_integration_test.go` with a case that invokes `create key --type authkey --tags tag:...` against the mock API, decodes the recorded request body, and asserts that `capabilities.devices.create.tags` is a concrete slice containing the provided tag value.

Rationale: the bug occurs in command-to-request translation, so the most reliable regression test is one that inspects the outgoing JSON body rather than isolated flag parsing.

Alternative considered:
- Unit-testing only the Cobra flag binding. Rejected because the flag binding already works; the failure is in request construction.

### 3. Correct the user-facing flag description while fixing the code path

Update the `--tags` flag help text and any nearby `create key` examples/documentation that describe supported key types so they no longer imply auth-key tags are unsupported.

Rationale: once the request path is fixed, leaving the help text unchanged would preserve a misleading contract around an existing flag.

Alternative considered:
- Leaving help text untouched to keep the code diff minimal. Rejected because it would maintain an avoidable documentation mismatch around the same regression.

## Risks / Trade-offs

- [Nil versus empty slice behavior may still be ambiguous if no tags are provided] -> Mitigation: keep current no-tags behavior unchanged and add a regression test only for the explicit tagged case.
- [Help text updates may broaden user expectations beyond current server-side validation] -> Mitigation: document only the supported request mapping and keep Tailscale API validation authoritative for tag requirements.
- [Regression tests could miss non-authkey flows] -> Mitigation: retain existing OAuth client compatibility tests and avoid changing those code paths.

## Migration Plan

1. Add the failing integration test covering auth-key tag serialization.
2. Update the auth-key request mapping to populate `capabilities.devices.create.tags`.
3. Correct `--tags` help text to reflect auth-key support.
4. Run the focused CLI test suite covering `create key`.

Rollback is straightforward: revert the auth-key tag mapping and the new regression test if the upstream API contract proves different than currently modeled.

## Open Questions

- None at the design level; the Go client model already confirms the expected auth-key payload field for tags.
