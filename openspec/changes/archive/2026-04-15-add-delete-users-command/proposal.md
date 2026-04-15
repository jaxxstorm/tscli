## Why

`tscli` already supports bulk deletion for devices, but user deletion is limited to a single-user command even though the API exposes enough user metadata to filter accounts for cleanup. Adding a bulk `delete users` command closes a user-management parity gap and gives administrators a safe, scriptable way to remove inactive or suspended accounts.

## What Changes

- Add a new `tscli delete users` command for bulk user deletion alongside the existing singular `tscli delete user` command.
- Add mutually exclusive selection flags for `--status` and `--last-seen`, plus a `--devices` filter to constrain deletion candidates by device count.
- Exclude admin users by default and add an `--admins` boolean flag to opt into deleting admin users when explicitly requested.
- Preserve script-friendly behavior with predictable filtering, validation, and machine-readable output for dry-run and deletion results.
- Extend OpenAPI parity coverage for bulk user-management workflows that currently require repeated single-user deletion calls.

## Capabilities

### New Capabilities
- `bulk-user-deletion`: Bulk delete users by status, last-seen age, and device-count filters with admin exclusion by default.

### Modified Capabilities
- `openapi-command-parity`: Expand parity requirements to cover bulk user deletion behavior and the new `delete users` command mapping.

## Impact

- Affected command groups: `delete users`, existing `delete user`, and shared `delete` command registration.
- Affected flags: new `--status`, `--last-seen`, `--devices`, and `--admins`; validation must enforce mutual exclusivity between `--status` and `--last-seen`.
- Affected code: `cmd/tscli/delete/...`, user filtering/deletion logic, output formatting, and command mapping or parity coverage checks.
- Affected API usage: user listing plus per-user deletion endpoints.
- Backward compatibility: existing `tscli delete user --user <id>` scripts continue to work unchanged; the new command adds behavior without removing or renaming existing flags.
