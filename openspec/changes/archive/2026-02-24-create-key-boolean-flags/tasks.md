## 1. CLI Flag Surface

- [x] 1.1 Add `--reusable`, `--ephemeral`, and `--preauthorized` boolean flags to `tscli create key` for auth-key creation.
- [x] 1.2 Update create-key command help text/examples to describe new auth-key capability flags and their scope.

## 2. Request Mapping And Validation

- [x] 2.1 Map new boolean flags into the auth-key create request capability fields sent by `CreateAuthKey`.
- [x] 2.2 Preserve existing `--type oauthclient` validation/request behavior and ensure new flags do not regress oauth flow.
- [x] 2.3 Add compatibility checks so auth-key creation without new flags behaves as before.

## 3. Automated Tests

- [x] 3.1 Add/extend command-level integration tests for `create key --type authkey` to assert capability booleans in mocked API request payloads.
- [x] 3.2 Add regression tests covering auth-key creation without capability flags.
- [x] 3.3 Add regression tests confirming oauthclient validation and request behavior remains unchanged.

## 4. Documentation And Verification

- [x] 4.1 Update README usage/examples for `create key` to include new capability flags.
- [x] 4.2 Run relevant CLI and package test suites and ensure command coverage/parity checks remain green.
