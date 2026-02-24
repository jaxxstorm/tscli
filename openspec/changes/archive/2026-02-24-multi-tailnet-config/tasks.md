## 1. Config Model And Resolution

- [x] 1.1 Add profile-aware config structures/helpers in `pkg/config` for `tailnets` and `active-tailnet` with validation for unique names and required API keys.
- [x] 1.2 Implement centralized resolver logic for effective `api-key` and `tailnet` using precedence: flags > env > active profile > legacy keys.
- [x] 1.3 Integrate resolver into root command pre-run and preserve existing defaults/errors (`tailnet` default `-`, required API key error).

## 2. Backward Compatibility And Persistence

- [x] 2.1 Add normalization/migration behavior so legacy-only config (`tailnet` + `api-key`) works unchanged.
- [x] 2.2 Implement persistence helpers so profile mutations write valid config and mirror active profile to legacy keys.
- [x] 2.3 Add compatibility checks for mixed config states (profiles + legacy keys) and missing active profile references.

## 3. Config Command Surface

- [x] 3.1 Add `config` operations to list tailnet profiles with active marker.
- [x] 3.2 Add `config` operation to set/switch `active-tailnet` with existence validation.
- [x] 3.3 Add `config` operation to upsert a tailnet profile (`name` + `api-key`) and persist changes.
- [x] 3.4 Add `config` operation to remove a tailnet profile with guardrails for active profile removal.

## 4. Automated Tests

- [x] 4.1 Add unit tests in `pkg/config` for schema validation and resolver precedence across flags/env/profile/legacy.
- [x] 4.2 Add unit tests for normalization and compatibility behavior between legacy and profile-based config.
- [x] 4.3 Add command-level integration tests for profile list/set/upsert/remove flows and persisted state transitions.
- [x] 4.4 Add command-level tests that verify runtime commands use active profile values when no flag/env overrides are present.

## 5. Documentation And UX Updates

- [x] 5.1 Update README/config documentation with new multi-tailnet config format and migration examples from legacy format.
- [x] 5.2 Update command help text/examples for new profile management operations and expected error messages.
- [x] 5.3 Add release-note style summary in change docs describing backward compatibility and rollback behavior.
