# pboard skill — Design Document

## Goal

Add a `pboard skill` command to the CLI itself that installs a Claude Code skill file onto the current machine, making the pboard usage guide available to all Claude Code sessions.

## Usage

```bash
pboard skill install            # Install the skill
pboard skill install --force    # Overwrite without prompting
pboard skill install --dry-run  # Preview without writing
pboard skill uninstall          # Remove the skill
```

## Implementation

### New Go files

| File | Purpose |
|------|---------|
| `internal/cli/skill.go` | `skill` parent command + `install` and `uninstall` subcommands |
| `internal/cli/skill_content.go` | Embedded skill markdown content as a Go string constant |

### Command structure

```
pboard skill
  ├── install [--force] [--dry-run]
  └── uninstall
```

Registered in `root.go` alongside other commands:
```go
rootCmd.AddCommand(newSkillCmd())
```

### `pboard skill install` behavior

```
1. Resolve target path: ~/.claude/commands/pboard.md
2. Create ~/.claude/commands/ directory if it doesn't exist (mode 0755)
3. If file already exists and --force not set:
   → Print "Skill already installed at <path>. Use --force to overwrite."
   → Exit with code 0
4. If --dry-run:
   → Print "Would install skill to <path>" and the content
   → Exit with code 0
5. Write the embedded skill content to the file (mode 0644)
6. Print "Skill installed to <path>"
```

### `pboard skill uninstall` behavior

```
1. Resolve target path: ~/.claude/commands/pboard.md
2. If file doesn't exist:
   → Print "No skill found at <path>. Nothing to remove."
   → Exit with code 0
3. Remove the file
4. Print "Skill removed from <path>"
```

### Flags

| Command | Flag | Type | Description |
|---------|------|------|-------------|
| `install` | `--force` | bool | Overwrite existing skill without prompting |
| `install` | `--dry-run` | bool | Show what would be installed without writing |

## What Gets Installed

A single file: `~/.claude/commands/pboard.md`

This is a Claude Code user-level command. Once installed, any Claude Code session on this machine can invoke `/pboard` or be guided by the skill's description to use the pboard CLI.

### Skill file content

- **Frontmatter:** `description` field for Claude Code discovery
- **Full command reference:** Every command, subcommand, flag, and output column
- **Agent best practices:** Always use `-o json` for programmatic parsing, use `-l N` to limit, use `get` for full details
- **Common workflows:** Recent feedback, release scope exploration, strategic alignment traversal

The content is embedded as a Go constant in `skill_content.go` — no external files, no network requests.

## What Does NOT Happen

- No files are written outside `~/.claude/commands/`
- No system packages installed
- No config files modified
- No network requests
- No `sudo` required
- No API calls to ProductBoard

## Future Extensibility

A `--target` flag could support other agents (e.g., `--target cursor`), each writing to the appropriate location. Out of scope for v1.

## Review Checklist

- [ ] Command is `pboard skill install` / `pboard skill uninstall`
- [ ] Skill content embedded in Go source (no external dependencies)
- [ ] Only writes to `~/.claude/commands/pboard.md`
- [ ] `--force` overwrites without prompting
- [ ] `--dry-run` previews without writing
- [ ] `uninstall` cleanly removes the file
- [ ] No side effects beyond that single file
- [ ] Skill content accurately reflects current CLI capabilities
