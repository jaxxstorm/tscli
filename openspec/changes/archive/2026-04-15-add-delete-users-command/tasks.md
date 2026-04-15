## 1. Command scaffolding and validation

- [x] 1.1 Create `cmd/tscli/delete/users` and register the new plural command in `cmd/tscli/delete/cli.go`
- [x] 1.2 Add `--status`, `--last-seen`, `--devices`, `--admins`, and `--confirm` flags with validation for mutual exclusivity and supported status values
- [x] 1.3 Update command help text and examples to document dry-run behavior, admin exclusion defaults, and filter usage

## 2. Filtering and deletion behavior

- [x] 2.1 Implement user candidate selection from the list-users API using status, last-seen age, device count, and admin-role filtering
- [x] 2.2 Implement dry-run summary output that reports matched, skipped, and failed users in supported output formats
- [x] 2.3 Implement confirmed bulk deletion using per-user delete requests and structured result aggregation

## 3. Parity and test coverage

- [x] 3.1 Add or update command-to-operation mapping/parity coverage for `tscli delete users`
- [x] 3.2 Add unit tests for filter validation, admin exclusion, timestamp handling, and deletion summary generation
- [x] 3.3 Add integration or command-level tests covering dry-run output, confirmed deletion flow, and compatibility with the existing `delete user` command
