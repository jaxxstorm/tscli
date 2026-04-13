## Context

The existing DNS mutation commands already line up with the project’s canonical `set` semantics, but both `set dns nameservers` and `set dns split-dns` hard-code `net.ParseIP` validation and therefore reject DNS-over-HTTPS endpoint addresses that the upstream Tailscale API accepts. That creates an avoidable mismatch between `tscli` and the supported API surface, particularly for users configuring public resolvers that are represented by DoH URLs instead of raw IP addresses.

This change is intentionally narrow:
- no new command groups, verbs, flags, config keys, or env vars
- no changes to `get`, `list`, or `delete` DNS command semantics beyond preserving compatibility with the relaxed `set` validation rules
- only the accepted value format for existing nameserver inputs changes

The implementation needs to keep the command surface script-friendly and predictable:
- existing IP-based scripts must continue to pass unchanged
- invalid values still need immediate local validation with actionable errors
- DNS commands should continue to exercise the same output and API request paths once validation succeeds

## Goals / Non-Goals

**Goals:**
- Allow valid DoH endpoint addresses anywhere the DNS mutation commands currently accept nameserver IPs.
- Reuse one validation rule across `set dns nameservers` and `set dns split-dns` so behavior stays consistent.
- Preserve the current `set` command taxonomy and request payload structure for DNS operations.
- Add tests covering valid IP inputs, valid DoH inputs, and still-invalid values.
- Update command help/examples so the accepted nameserver formats are discoverable.

**Non-Goals:**
- Adding new DNS commands or changing existing `get`, `list`, `set`, or `delete` command paths.
- Expanding validation to arbitrary DNS URI formats beyond what upstream Tailscale treats as supported DoH addresses.
- Reworking DNS payload shapes or API response rendering.

## Decisions

### 1. Introduce a shared nameserver validator for IP or DoH inputs

Decision:
- Replace the current inline `net.ParseIP` checks with a small reusable validator that accepts either a literal IP address or a valid DoH endpoint address.
- Use the same validator from both `set dns nameservers` and `set dns split-dns`.

Rationale:
- The current mismatch exists in two places, so one shared rule reduces drift and keeps user-facing behavior consistent.
- A focused helper is smaller and safer than duplicating URL parsing rules in each command.

Alternatives considered:
- Keep separate per-command validators: rejected because the accepted format would drift easily.
- Remove local validation entirely and rely on API errors: rejected because the CLI should fail fast with actionable input errors.

### 2. Treat DoH addresses as HTTPS endpoints with strict local parsing

Decision:
- Validate DoH inputs using URL parsing and restrict acceptance to syntactically valid HTTPS endpoints rather than arbitrary strings containing `https://`.
- Continue to reject malformed URLs, non-HTTPS schemes, and obviously invalid nameserver strings.

Rationale:
- This matches the user’s request for DoH addresses without making validation so loose that typos slip through silently.
- Keeping validation local preserves existing CLI ergonomics for scripts and interactive use.

Alternatives considered:
- Allow any hostname or bare domain: rejected because the upstream reference specifically points to DoH endpoint addresses, not generic resolver names.
- Exact parity with all upstream helper edge cases in the first pass: rejected because a pragmatic HTTPS-endpoint rule should cover the intended inputs with much less implementation risk.

### 3. Extend command tests at both validation and success-path layers

Decision:
- Add validation-focused tests for both commands covering accepted DoH inputs and rejected malformed values.
- Add mock-backed success-path coverage to confirm DoH values survive into the outgoing request payload unchanged.

Rationale:
- This change is primarily about input validation, so both local failure behavior and downstream request behavior need explicit coverage.
- Existing test structure already covers DNS command output and mock API execution, so the smallest change is to extend those suites.

Alternatives considered:
- Validation-only tests: rejected because they do not prove the accepted DoH value is sent correctly.
- Integration-only tests: rejected because they make invalid-input behavior harder to assert precisely.

## Risks / Trade-offs

- [Validation is still narrower than upstream edge-case behavior] -> Mitigation: document DoH endpoint support clearly and keep the validator small so it can be refined if upstream semantics prove broader.
- [Shared helper gets reused too broadly] -> Mitigation: keep the helper scoped to DNS nameserver validation rather than turning it into a generic URL validator.
- [Help text drifts from actual accepted inputs] -> Mitigation: update command descriptions/examples alongside tests in the same change.

## Migration Plan

1. Add a shared validator for literal IP or DoH endpoint values.
2. Swap `set dns nameservers` and `set dns split-dns` to use that validator.
3. Update help text/examples for the affected flags.
4. Add validation and mock-backed tests for valid DoH and invalid nameserver inputs.
5. Run the relevant CLI test suites and confirm existing IP-based behavior still passes.

Rollback strategy:
- Revert the shared validator and command wiring together, restoring IP-only validation.

## Open Questions

- Do we want to accept only path-bearing DoH endpoints, or should a bare HTTPS origin also be considered valid if upstream Tailscale accepts it?
