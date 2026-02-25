## Why

The CLI, docs, and pinned OpenAPI snapshot still treat "keyType" as `authkey` or `oauthclient`, yet the Tailscale API also supports `federated` trust credentials. As a result, create-key commands cannot provision federated identities and the coverage gap tool never reports them, leaving a blind spot for tailnet automation.

## What Changes

- Add federated identity support to the `tscli create key` command, request payloads, and output handling so users can create federated credentials from the CLI.
- Update the OpenAPI snapshot/command mapping to include the `federated` key type and document how CLI commands map to that operation.
- Extend the CLI tests and docs so the new `--type federated` behavior is validated and discoverable by contributors.

## Capabilities

### New Capabilities
- `federated-key-support`: Requirements for exposing federated identity creation via the CLI and keeping coverage tooling synchronized with the OpenAPI spec.

### Modified Capabilities
- `openapi-command-parity`: Capture the new `federated` key-type operation in the command-operation mapping and ensure coverage tooling treats it as mapped.

## Impact

Affected areas: the CLI create key command (`cmd/tscli/create/key/cli.go` and tests), OpenAPI snapshot files + coverage command map under `pkg/contract/openapi`, generated docs for the create-key flow, and automations verifying coverage gaps. New tests/documentation will explain how to request `--type federated` and how `keyType` surfaces in the API.
