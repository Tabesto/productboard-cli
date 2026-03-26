# Research: Configure Command Token Argument

**Date**: 2026-03-26
**Branch**: `002-configure-token-arg`

## R1: Cobra Optional Positional Arguments

**Decision**: Use `cobra.MaximumNArgs(1)` to accept 0 or 1 positional arguments.

**Rationale**: Cobra natively supports optional positional args via `Args: cobra.MaximumNArgs(1)`. When `len(args) == 1`, use the argument as the token. When `len(args) == 0`, fall back to interactive prompt. This is the standard Cobra pattern for commands with optional positional input.

**Alternatives considered**:
- **Flag-based (`--token`)**: Would work but is less ergonomic for a paste workflow. Positional args are simpler: `pboard configure pb_token_here`.
- **Stdin pipe**: More complex, less intuitive for the primary use case of copy-paste from a browser.
