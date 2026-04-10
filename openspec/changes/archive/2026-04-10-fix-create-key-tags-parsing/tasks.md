## 1. Regression Coverage

- [x] 1.1 Extend `test/cli/create_key_flags_integration_test.go` with a failing auth-key case that runs `tscli create key --type authkey --tags tag:tsdns` and asserts the recorded payload currently needs `capabilities.devices.create.tags`.
- [x] 1.2 Keep existing compatibility coverage for auth-key boolean flags and non-authkey create-key flows passing while adding the new tag regression case.

## 2. Auth-Key Request Mapping

- [x] 2.1 Update `cmd/tscli/create/key/cli.go` so the auth-key path copies parsed `--tags` values into `tsapi.CreateKeyRequest.Capabilities.Devices.Create.Tags`.
- [x] 2.2 Update the `--tags` flag description and any nearby create-key help/example text to reflect auth-key tag support without changing OAuth client or federated behavior.

## 3. Verification

- [x] 3.1 Run the focused create-key CLI test suite and confirm the new tag regression test passes alongside existing capability-flag coverage.
- [x] 3.2 Confirm the change remains apply-ready with no additional OpenSpec artifact work required.
