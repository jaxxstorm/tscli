## 1. Add shared DNS nameserver validation

- [x] 1.1 Introduce a reusable validator that accepts literal IP addresses and valid HTTPS DoH endpoint addresses for DNS nameserver inputs.
- [x] 1.2 Update `cmd/tscli/set/dns/nameservers/cli.go` to use the shared validator while preserving current request payload and error handling semantics.
- [x] 1.3 Update `cmd/tscli/set/dns/split/cli.go` to use the same validator for `domain=value` entries and keep existing clear/replace behavior intact.

## 2. Cover accepted and rejected input behavior

- [x] 2.1 Add validation-focused CLI tests for `set dns nameservers` covering valid IP inputs, valid DoH endpoint inputs, and malformed nameserver values.
- [x] 2.2 Add validation-focused CLI tests for `set dns split-dns` covering valid DoH endpoint inputs and invalid entry values.
- [x] 2.3 Add mock-backed success-path coverage confirming accepted DoH nameserver inputs are sent to the API unchanged.

## 3. Document and verify the change

- [x] 3.1 Update help text and any user-facing examples for the affected DNS `set` commands to mention DoH endpoint support.
- [x] 3.2 Run the relevant CLI test suites and confirm existing IP-based DNS command behavior still passes.
- [x] 3.3 Review the command behavior against current `get`/`set`/`delete`/`list` DNS semantics to ensure no unrelated command taxonomy or compatibility regressions were introduced.
