## 1. Reusable Maintenance Tooling

- [ ] 1.1 Extract coverage-gap execution from `coverage/coveragegaps` into an importable internal package while preserving existing flags and defaults.
- [ ] 1.2 Extract OpenAPI refresh execution from `coverage/openapirefresh` into an importable internal package while preserving existing flags and defaults.
- [ ] 1.3 Keep existing `go run ./coverage/...` entrypoints working for all make targets.

## 2. CLI Maintenance Commands

- [ ] 2.1 Add a local `maintenance` Cobra command group to the root command and bypass API credential preflight for that group.
- [ ] 2.2 Implement `tscli maintenance coverage-gaps` with repository-default path flags and a `--check` mode that enables baseline diffing and strict gap/regression failures.
- [ ] 2.3 Implement `tscli maintenance openapi-refresh` with source URL, schema output, and metadata output flags.
- [ ] 2.4 Add command-level tests proving maintenance help works without credentials and command flags are wired.

## 3. Agent Skills

- [ ] 3.1 Add matching `.codex/skills/coverage-gaps` and `.opencode/skills/coverage-gaps` instructions for running coverage-gap workflows before parity implementation.
- [ ] 3.2 Add matching `.codex/skills/openapi-refresh` and `.opencode/skills/openapi-refresh` instructions for refreshing the pinned OpenAPI snapshot before latest-schema work.
- [ ] 3.3 Ensure skills name supported make targets, equivalent `tscli maintenance` commands, generated artifacts, and failure handling.

## 4. Repository Metadata and Verification

- [ ] 4.1 Update the leaf command manifest and generated command docs for the new maintenance commands.
- [ ] 4.2 Run focused tests for maintenance commands and reusable tooling.
- [ ] 4.3 Run OpenSpec validation/status checks for `add-agent-coverage-refresh-cli`.
