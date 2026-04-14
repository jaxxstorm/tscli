## 1. OAuth auth plumbing

- [x] 1.1 Add config resolution helpers for `oauth-client-id` and `oauth-client-secret` using flag, environment, and active-profile precedence
- [x] 1.2 Extend profile parsing, validation, and persistence to support OAuth-backed profiles and optional explicit `tailnet` values
- [x] 1.3 Add a bearer-authenticated raw request helper that exchanges OAuth client credentials and calls lifecycle endpoints without disturbing existing API-key request paths

## 2. Tailnet lifecycle commands

- [x] 2.1 Add `tscli create tailnet` with `--display-name` validation, lifecycle request/response types, and authoritative output handling
- [x] 2.2 Add `tscli list tailnets` with OAuth credential resolution and authoritative response output
- [x] 2.3 Add `tscli delete tailnet` with tailnet-specific OAuth auth, success output, and actionable error handling
- [x] 2.4 Wire the new commands into the existing `create`, `list`, and `delete` Cobra groups and exempt them from the API-key-only pre-run path

## 3. Profile command updates

- [x] 3.1 Extend `config profiles upsert` to accept OAuth client credentials and validate supported auth shapes
- [x] 3.2 Update `config profiles list` output to show whether each profile is API-key-backed or OAuth-backed while preserving active-profile indication
- [x] 3.3 Verify `config profiles set-active` and `delete` behavior still work with mixed API-key and OAuth-backed profiles

## 4. Tests and documentation

- [x] 4.1 Add unit tests for OAuth profile validation, precedence resolution, and canonical config persistence
- [x] 4.2 Add CLI/integration tests for `create tailnet`, `list tailnets`, and `delete tailnet`, including OAuth exchange failures and API error propagation
- [x] 4.3 Update authentication and configuration docs for OAuth lifecycle commands, profile keys, and one-time secret handling guidance
- [x] 4.4 Update help text and example output coverage for the new commands and profile auth shapes
