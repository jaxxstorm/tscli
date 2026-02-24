## 1. Docsify Site Foundation

- [x] 1.1 Create `docs/` structure with Docsify entrypoint (`index.html`), sidebar navigation, and landing page.
- [x] 1.2 Add core top-level docs pages for getting started, command reference index, configuration, and authentication.
- [x] 1.3 Add local docs serving workflow (documented command or Make target) for contributor preview.

## 2. Command Markdown Generation

- [x] 2.1 Implement command reference generation using Cobra docs APIs to emit Markdown pages per command path under `docs/commands/`.
- [x] 2.2 Add deterministic generation behavior (stable naming/order) and document how to regenerate command docs.
- [x] 2.3 Add docs verification workflow/check that fails when generated command docs are stale or missing.

## 3. Config And Auth Documentation Content

- [x] 3.1 Document configuration schema including legacy and profile-based keys (`api-key`, `tailnet`, `tailnets`, `active-tailnet`, output settings).
- [x] 3.2 Document precedence rules (flags > env > config file) with practical examples.
- [x] 3.3 Document API key authentication methods (`--api-key`, `TAILSCALE_API_KEY`, config file) and secret-handling guidance.

## 4. Validation And Integration

- [x] 4.1 Add tests/checks ensuring required docs pages and command-doc artifacts exist.
- [x] 4.2 Update README and developer docs to point to Docsify docs and generation/verification commands.
- [x] 4.3 Run docs generation, docs verification, and relevant test suites to confirm no regressions.
