# Documentation Workflows

This file is the single source of truth for how the documentation site is built, previewed, and updated. The authored pages live alongside generated command reference files.

## Local preview

Serve the site with Docsify during development:

```bash
make docs-serve
```

If `docsify-cli` is missing, install it once:

```bash
npm i -g docsify-cli
```

## Command reference generation

The command docs under `docs/commands/` are generated from the current Cobra tree. After changing CLI commands, regenerate them before committing:

```bash
make docs-generate
```

## Verification and CI

CI runs `make docs-check` to ensure generated files match what `make docs-generate` produces. Run the same check locally after generating docs (or if `make docs-check` fails the CD log shows the `diff -ru` output automatically). The target exits non-zero if any file differs, so rerun `make docs-generate` and commit the generated files before pushing.

If you add new key types, CLI verbs, or request/response fields, keep the OpenAPI coverage data in sync by updating `pkg/contract/openapi/command-operation-map.yaml`, `coverage/property-coverage.yaml`, and `coverage/property-exclusions.yaml`, then rerun `make coverage-gaps-check`. CI uses the generated `coverage/coverage-gaps.*` artifacts to detect endpoint-level and property-level regressions.

## Editing guidance

- Edit authored content in `docs/*.md` (except the `docs/commands/` tree).
- Do **not** hand-edit the generated Markdown in `docs/commands/`; let `make docs-generate` rewrite the entire directory.
- After a change that touches CLI commands or the docs generator, always run `make docs-check` to confirm nothing is stale before submitting a PR.
