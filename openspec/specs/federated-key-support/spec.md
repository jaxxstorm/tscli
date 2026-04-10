# federated-key-support Specification

## Purpose

Define the expected CLI, coverage, and test behavior for creating federated keys with `tscli create key --type federated`.

## Requirements

### Requirement: tscli create key accepts federated keyType
The CLI SHALL extend the existing `tscli create key` command so the `--type` flag accepts `authkey`, `oauthclient`, and `federated`. When users choose `--type federated`, the command SHALL send `keyType: federated` in the request payload and parse the response using the existing key model.

#### Scenario: creating a federated key
- **WHEN** the user runs `tscli create key --type federated --reuse --name my-federated-key`
- **THEN** the CLI sends a POST to `/tailnet/{tailnet}/keys` with `{"keyType":"federated","reuse":true,"name":"my-federated-key"}` and prints the new key metadata in the configured output format

#### Scenario: federated type validation
- **WHEN** the user passes `--type` with a value outside the allowed set
- **THEN** the command exits with a validation error that references the supported type list and does not invoke the API

### Requirement: coverage tooling maps the federated key operation
The OpenAPI coverage mapping SHALL include the `federated` variant of the `create tailnet keys` operation so that `coverage/coveragegaps` treats the CLI command as implemented.

#### Scenario: coverage gaps acknowledge coverage
- **WHEN** `make coverage-gaps-check` finishes
- **THEN** the command exits with status 0, and the `coverage/coverage-gaps.md` report lists zero unmapped operations for federated key creation

### Requirement: testing ensures federated key behavior stays supported
Unit and integration tests SHALL cover the new `--type federated` path, including flag validation, API payload, and response parsing.

#### Scenario: tests validate federated creation
- **WHEN** the test suite runs the create-key command logic with mocked API responses for `federated`
- **THEN** the tests pass, demonstrating that the command constructs the correct payload, respects the output format flag, and surfaces API errors consistently with other key types
