## 1. CLI implementation

- [x] 1.1 Extend `cmd/tscli/create/key` so `--type` accepts `authkey`, `oauthclient`, and `federated`, validating the input and sending `keyType` accordingly when building the POST payload.
- [x] 1.2 Ensure the command reuses the existing key model/response handling so federated keys print the same metadata and honor `--output` format without extra code paths.

## 2. Testing and coverage

- [x] 2.1 Update existing unit/integration tests for the create-key command to cover the federated branch (mock the Tailscale API, validate payload, and assert output/error behavior).
- [x] 2.2 Refresh the OpenAPI snapshot/command coverage map under `pkg/contract/openapi` so the federated create-key operation is marked as implemented, then rerun `coverage/coveragegaps` to confirm no unmapped operations remain.

## 3. Documentation

- [x] 3.1 Regenerate the command docs to describe the new `--type federated` option and confirm `docsify` pages include the updated create-key entry.
- [x] 3.2 Document in `docs/README.md` (or the new docs landing page) how to request federated keys and mention the coverage map update so contributors know to sync the OpenAPI spec when adding future enums.
