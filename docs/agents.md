# AI Agent Integrations

`tscli agent init` bootstraps `tscli`-aware instructions, skills, prompts, and commands so AI agents can use `tscli` directly instead of ad hoc API calls.

## What it generates

The command writes a versioned bundle into either your user-level config directories or a target repository.

Global install surfaces:

- `.config/tscli/agent/manifest.yaml`: versioned manifest used by `tscli agent update`
- `.config/tscli/agent/commands.md`: generated catalog of the current `tscli` leaf commands
- `.codex/skills/tscli/SKILL.md`: Codex skill surface
- `.claude/commands/tscli-inspect.md`: Claude Code read-only slash command
- `.claude/commands/tscli-operate.md`: Claude Code mutation slash command
- `.config/opencode/commands/tscli-inspect.md`: OpenCode read-only command
- `.config/opencode/commands/tscli-operate.md`: OpenCode mutation command

Repo-local install surfaces:

- `AGENTS.md`: global instructions for agent-aware tools
- `CLAUDE.md`: Claude Code project memory
- `.tscli/agent/manifest.yaml`: versioned manifest used by `tscli agent update`
- `.tscli/agent/commands.md`: generated catalog of the current `tscli` leaf commands
- `.codex/skills/tscli/SKILL.md`: Codex skill surface
- `.github/skills/tscli/SKILL.md`: GitHub Copilot skill surface
- `.github/prompts/tscli-inspect.prompt.md`: GitHub Copilot read-only prompt
- `.github/prompts/tscli-operate.prompt.md`: GitHub Copilot mutation prompt

## Initialize

Install global user-level integrations:

```bash
tscli agent init
```

Install repo-local integrations into another directory:

```bash
tscli agent init --dir /path/to/repo
```

Restrict generation to selected tool surfaces:

```bash
tscli agent init --dir /path/to/repo --tool codex --tool claude --tool opencode
```

Supported `--tool` values are `generic`, `codex`, `claude`, `opencode`, and `copilot`.

Global installs support `codex`, `claude`, and `opencode`.

Repo-local installs support `generic`, `codex`, `claude`, `opencode`, and `copilot`.

## Update

Refresh the global generated bundle after upgrading `tscli` or adding new commands:

```bash
tscli agent update
```

Refresh a repo-local bundle:

```bash
tscli agent update --dir /path/to/repo
```

`update` reads the scope-appropriate manifest to preserve the original tool selection and rewrite the managed files in place.

## Operating model

- Agents should prefer `tscli list ...`, `tscli get ...`, `tscli create ...`, `tscli set ...`, and `tscli delete ...` over raw API calls when a command exists.
- The generated command catalog is derived from the live Cobra command tree, so it stays aligned with the CLI surface.
- Global command surfaces are reusable user-level integrations; repo-local instruction files remain project-scoped.
- If a generated target path already contains an unmanaged file, rerun with `--force` only if you intend to replace it.
