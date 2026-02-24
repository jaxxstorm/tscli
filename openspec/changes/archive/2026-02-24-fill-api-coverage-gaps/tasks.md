## 1. Establish parity inventory and policy

- [x] 1.1 Regenerate the OpenAPI coverage-gap report from the pinned schema and snapshot the current uncovered operations by domain.
- [x] 1.2 Define and commit an explicit exclusions policy for any operations intentionally out of scope, including rationale.
- [x] 1.3 Convert uncovered operations into an implementation backlog grouped by command verb (`create`, `set`, `get`, `list`, `delete`) and domain.

## 2. Enforce canonical command taxonomy

- [x] 2.1 Audit existing leaf commands for verb/taxonomy violations and identify canonical path corrections.
- [x] 2.2 Implement canonical command names for non-aligned commands while preserving legacy aliases for compatibility.
- [x] 2.3 Update command help text/examples to reference canonical command paths and alias behavior.

## 3. Implement missing `get` and `list` retrieval commands

- [x] 3.1 Add missing single-resource retrieval commands (`get`) for uncovered single-object operations.
- [x] 3.2 Add missing collection retrieval commands (`list`) for uncovered collection operations.
- [x] 3.3 Ensure each retrieval command supports consistent output modes and clear error messages.

## 4. Implement missing `create` commands

- [x] 4.1 Add missing creation commands for uncovered POST-create operations by domain.
- [x] 4.2 Add request validation for required flags and payload constraints on each new `create` command.
- [x] 4.3 Add response rendering and success output consistency for all new `create` commands.

## 5. Implement missing `set` update/mutation commands

- [x] 5.1 Add missing update/mutation commands for uncovered PATCH/PUT/POST-action operations that change state.
- [x] 5.2 Align update command naming and arguments to the canonical `set` taxonomy and resource nouns.
- [x] 5.3 Ensure mutation commands support idempotent request construction and consistent error wrapping.

## 6. Implement missing `delete` commands

- [x] 6.1 Add missing delete commands for uncovered delete operations.
- [x] 6.2 Ensure delete commands include required target validation and consistent result output payloads.

## 7. Update parity mapping and enforcement tooling

- [x] 7.1 Add command-operation mappings for every implemented canonical command.
- [x] 7.2 Eliminate unknown mapped commands and unknown mapped operations in the coverage report output.
- [x] 7.3 Drive uncovered in-scope operation count to zero (or policy-backed excluded) and regenerate baseline artifacts.

## 8. Expand automated test coverage

- [x] 8.1 Add unit tests for new command flag validation and command wiring.
- [x] 8.2 Add mock-backed integration tests for new command success and API error paths.
- [x] 8.3 Extend command-manifest/coverage checks for all new canonical command paths and aliases.

## 9. Finalize docs and CI checks

- [x] 9.1 Update README command coverage/taxonomy documentation and usage examples for new commands.
- [x] 9.2 Update CI/Make targets as needed to enforce parity and coverage-gap regression checks.
- [x] 9.3 Run full validation (`go test ./...`, coverage-gap generation, regression checks) and resolve remaining gaps before merge.
