## Context

`tscli` currently relies primarily on README-based documentation. As command surface area grows, users need a browsable site, command-by-command reference pages, and clearer docs for configuration and API key authentication. The codebase already uses Cobra, which supports Markdown doc generation, making command reference generation a natural fit.

This change introduces a `docs/` site with Docsify and adds automated command reference generation so docs stay synchronized with the CLI tree.

## Goals / Non-Goals

**Goals:**
- Add a Docsify-based documentation site under `docs/`.
- Auto-generate command Markdown pages from the Cobra command tree.
- Provide docs for config file structure, key precedence, multi-tailnet profiles, and API key auth.
- Provide repeatable generation commands suitable for local dev and CI.
- Add tests/checks ensuring command docs generation remains consistent.

**Non-Goals:**
- Changing command behavior or CLI UX beyond documentation generation hooks.
- Replacing OpenSpec docs or other project metadata systems.
- Introducing complex docs build tooling beyond Docsify + simple generation scripts.

## Decisions

### 1. Use Docsify for lightweight docs hosting in-repo

Create `docs/` with:
- `index.html` Docsify shell
- `_sidebar.md` navigation
- landing + topic pages

Rationale: zero compile-step docs hosting is simple for contributors and easy to serve locally.

Alternatives considered:
- Full static-site generator (Docusaurus/MkDocs). Rejected due to higher setup/maintenance overhead.

### 2. Generate command reference from Cobra docs API

Add generation tooling that calls Cobra Markdown generation (for example `doc.GenMarkdownTreeCustom`) against the existing root command builder and writes output to `docs/commands/`.

Rationale: generated docs reduce drift between real flags/usage and written docs.

Alternatives considered:
- Manually writing command pages. Rejected due to maintenance burden and drift risk.

### 3. Keep command docs generation deterministic and CI-friendly

Add scripts/targets to:
- regenerate docs
- optionally verify generated docs are up-to-date

Rationale: deterministic output supports PR review and automated checks.

Alternatives considered:
- Generating docs only at release time. Rejected because drift can persist during development.

### 4. Add explicit config + auth docs as first-class pages

Create focused docs pages for:
- config schema (`tailnet`, `tailnets`, `active-tailnet`, `output`, etc.)
- precedence (flags > env > config)
- API key auth workflows and secure handling practices

Rationale: these are common setup pain points and should be directly discoverable.

### 5. Add documentation behavior tests/checks

Add lightweight tests/checks for:
- generator command succeeds
- expected command docs are produced
- navigation links include command docs and core setup pages

Rationale: prevents silent docs regressions.

## Risks / Trade-offs

- [Generated docs churn may increase PR noise] -> Mitigation: deterministic ordering/naming and clear generation commands.
- [Docsify content can drift from README] -> Mitigation: treat Docsify pages as source-of-truth and reference from README.
- [Command generation may include noisy internal/alias commands] -> Mitigation: configure generator naming/filtering rules and document expected output set.

## Migration Plan

1. Introduce `docs/` scaffold and initial pages.
2. Add command-doc generation tooling and output directory.
3. Add docs generation/verification targets in project workflow.
4. Update README to point to `docs/` and generation commands.
5. Rollback strategy: remove docs generation hooks and `docs/` assets without affecting CLI runtime.

## Open Questions

- Should alias commands be emitted as separate docs pages or redirect/reference canonical command pages?
- Should generated command docs include front matter for Docsify sidebar grouping, or should grouping remain manual?
