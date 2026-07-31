---
name: coverage-gaps
description: Run tscli OpenAPI command and property coverage-gap analysis before implementing command parity or schema-driven CLI coverage work.
license: MIT
---

Use this skill when work involves OpenAPI command parity, coverage-gap elimination, command-operation mappings, property coverage, or missing CLI commands derived from the Tailscale API schema.

## Steps

1. Run the supported coverage-gap workflow before editing commands:

   ```bash
   make coverage-gaps-check
   ```

   Equivalent CLI command:

   ```bash
   go run ./cmd/tscli -- maintenance coverage-gaps --check
   ```

2. Inspect generated artifacts:

   - `coverage/coverage-gaps.json`
   - `coverage/coverage-gaps.md`
   - `coverage/coverage-gaps-diff.md`

3. Treat reported uncovered operations, unmapped commands, unknown mapped operations, unknown mapped commands, and uncovered properties as the implementation backlog.

4. After implementing changes, rerun `make coverage-gaps-check`.

5. If leaf commands changed, update `test/cli/testdata/leaf_commands.txt`, command-operation mappings, property coverage declarations, tests, and generated command docs as needed.

## Failure Handling

- If the check fails because gaps remain, continue implementation against the reported backlog.
- If the check fails because a mapping or exclusion references a missing command or operation, fix the mapping/exclusion before implementing new commands.
- If the workflow cannot run because the pinned OpenAPI snapshot is stale for the requested work, run the OpenAPI refresh workflow first.
