## Context

`tscli` already exposes `delete user` for one-off deletions and `delete devices` for bulk destructive operations. The new `delete users` command needs to fit the existing `delete` command layout, reuse the current user list and low-level delete API behavior, and remain safe for scripts by defaulting to dry-run output until the caller passes `--confirm`.

The API data needed for filtering is already available from `list users`, including `status`, `lastSeen`, `deviceCount`, and `role`. That makes this change primarily a CLI orchestration problem: list users, filter locally, optionally skip admins, and then issue per-user deletion requests with predictable output.

## Goals / Non-Goals

**Goals:**
- Add `tscli delete users` as a bulk sibling to `delete user` and `delete devices`.
- Support a single primary filter of `--status` or `--last-seen`, with optional `--devices` narrowing.
- Exclude admin users by default and allow explicit inclusion with `--admins=true`.
- Keep deletion behavior safe and script-friendly by following the `delete devices` dry-run plus `--confirm` pattern.
- Add unit and command-level integration coverage for validation, filtering, and deletion results.

**Non-Goals:**
- Replacing or removing the existing `tscli delete user` command.
- Adding configuration-file or environment-variable settings for deletion filters.
- Expanding bulk deletion to additional user attributes beyond status, last-seen age, device count, and admin inclusion.

## Decisions

### Create a new plural command under `cmd/tscli/delete/users`
The implementation should mirror the current singular/plural split already used by `delete device` and `delete devices`. This keeps command discovery intuitive and avoids overloading `delete user` with incompatible bulk-operation flags.

Alternative considered: extend `delete user` with bulk flags.
Why not: it would conflate single-target and bulk semantics, complicate validation, and make help output harder to understand.

### Reuse the `delete devices` safety model
`delete users` should support `--confirm`, default to dry-run output, and return a summary with per-user results. This keeps destructive behavior consistent across bulk delete commands and reduces the chance of accidental user removal.

Alternative considered: delete immediately unless a dry-run flag is provided.
Why not: it would diverge from the existing bulk delete precedent and increase risk for operators.

### Perform filtering client-side from the list-users response
The command should call the existing user list API once, then filter candidates locally using `status`, `lastSeen`, `deviceCount`, and `role`. This matches the available response shape and avoids introducing new API dependencies.

Alternative considered: rely on server-side filtering.
Why not: the required combination of filters is already available from the list response, and the current CLI already consumes that response shape.

### Treat `--status` and `--last-seen` as mutually exclusive, while allowing `--devices` as either a standalone or additive filter
The command should prevent callers from combining `--status` and `--last-seen`, because they represent two different activity-selection models. `--devices` can still be applied on its own or alongside either of those filters because it narrows candidates by inventory rather than introducing a conflicting activity check.

Alternative considered: require either `--status` or `--last-seen` for every invocation.
Why not: the requested CLI examples include `--devices` as a valid standalone cleanup filter, and supporting that case remains safe because the command still defaults to dry-run mode.

### Accept only `inactive` and `suspended` as valid `--status` values
The user request explicitly targets deletion of inactive or suspended users. Limiting accepted status values prevents accidental use of `active` and keeps the destructive workflow aligned with administrative cleanup cases.

Alternative considered: accept any API status value.
Why not: broader acceptance weakens the safety boundary for a destructive command.

### Exclude admin users by default based on the API `role` field
Candidates with `role=admin` should be skipped unless `--admins=true` is set. The output should identify skipped admins so callers can understand why a listed user was not targeted.

Alternative considered: only exclude `owner` or all privileged roles.
Why not: the current request specifically calls out admins; the implementation should start with that exact behavior and expand later only if needed.

### Factor filtering and deletion summary logic for testability
The command package should keep Cobra wiring thin and move candidate selection, skip reason generation, and delete execution into helper functions that can be tested without shelling through Cobra. Command-level tests can then focus on flag validation and output shape.

Alternative considered: keep all logic inline in `RunE`.
Why not: that would make validation and candidate selection harder to test thoroughly.

## Risks / Trade-offs

- [Role semantics may be broader than `admin`] -> Mitigation: implement the requested default exactly, document the behavior in help/specs, and keep the role check isolated for future adjustment.
- [Large user sets could produce many sequential delete calls] -> Mitigation: start with the same simple, predictable request pattern as existing commands; if needed, concurrency can be added later behind the same output contract.
- [`lastSeen` parsing may encounter empty or unexpected values] -> Mitigation: treat unparsable or missing timestamps as skipped users with an explicit reason instead of deleting them.
- [Dry-run output may differ from actual delete results if users change between list and delete] -> Mitigation: report per-user failures individually and preserve aggregate success/failed/skipped counts.

## Migration Plan

1. Add the new `delete users` command and register it in the top-level `delete` command.
2. Implement flag validation, user filtering, and dry-run/deletion summary behavior.
3. Add or update command-operation mapping/parity coverage for the new command.
4. Add unit tests for filter logic and integration tests for command behavior and output.
5. No data migration is required; rollback is removal of the new command and its mapping/tests.

## Open Questions

- Whether `it-admin`, `network-admin`, and other privileged roles should eventually be included in the default-protected set in addition to `admin`.
- Whether bulk user deletion should eventually support concurrent delete requests like `delete devices`, or remain sequential for simpler failure reporting.
