# Research: ProductBoard CLI Read-Only Tool

**Date**: 2026-03-26
**Branch**: `001-productboard-cli-read`

## R1: Language & Framework Choice

**Decision**: Go with Cobra CLI framework

**Rationale**: Go produces single static binaries ideal for Homebrew distribution — no runtime dependencies. Cobra is the de facto standard for Go CLI tools (used by kubectl, gh, docker). It provides built-in help generation, subcommand nesting, flag parsing, and shell completion. The Go ecosystem has excellent HTTP client support and JSON marshaling built into the standard library.

**Alternatives considered**:
- **Rust (clap)**: Excellent performance and binary size, but longer development cycle and smaller contributor pool for a tool like this.
- **Python (click/typer)**: Faster to develop, but requires Python runtime — complicates Homebrew distribution and user setup.
- **Node.js (oclif)**: Good CLI framework, but requires Node runtime — same distribution problem as Python.

## R2: HTTP Client & API Integration

**Decision**: Go standard library `net/http` with a thin wrapper for pagination, auth headers, and error handling.

**Rationale**: The ProductBoard API is a straightforward REST API with Bearer token auth and cursor-based pagination. No need for a generated client — the API surface is read-only with simple GET requests. A thin wrapper keeps dependencies minimal and gives full control over error handling, retry logic, and pagination.

**Alternatives considered**:
- **OpenAPI-generated client (oapi-codegen)**: Would auto-generate types and client from swagger, but adds build complexity and generated code bloat for what is essentially ~50 GET endpoints with simple JSON responses.
- **Third-party HTTP client (resty, req)**: Adds dependency without significant benefit over stdlib for this use case.

## R3: Pagination Strategy

**Decision**: Cursor-based auto-pagination with `--limit` cap.

**Rationale**: The ProductBoard API uses cursor-based pagination (`pageCursor` parameter). The client will follow cursors automatically until all results are fetched or the `--limit` cap is reached. This provides the best UX for scripting (pipe full output) while protecting against accidentally pulling unbounded datasets.

**Implementation pattern**:
- Each `list` operation iterates pages using `pageCursor` from responses
- Accumulate results until: no next cursor, or `--limit` reached
- Default: no limit (fetch all)

## R4: Output Formatting

**Decision**: Table output via `tablewriter` (or similar), JSON via `encoding/json`.

**Rationale**: Table output needs column alignment and truncation for terminal display. JSON output uses Go's built-in marshaling. The `--output` flag (`table` | `json`) controls which formatter is used. Default is `table`.

## R5: Configuration Management

**Decision**: YAML config file at `~/.config/pboard/config.yaml` with `PRODUCTBOARD_API_TOKEN` env var override.

**Rationale**: `~/.config/` follows XDG Base Directory Specification. YAML is human-readable and easy to edit. The `pboard configure` command will write the config file with mode 600. Environment variable takes precedence for CI/CD and scripting use cases.

**Config structure**:
```yaml
api_token: "pb_token_..."
api_url: "https://api.productboard.com"  # optional override
```

## R6: Homebrew Distribution

**Decision**: GoReleaser + GitHub Releases + Homebrew Tap repository.

**Rationale**: GoReleaser automates cross-compilation, checksums, and Homebrew formula generation from a single `goreleaser.yaml` config. It creates GitHub releases with pre-built binaries for macOS (arm64, amd64) and Linux (amd64). A separate `homebrew-tap` repository hosts the formula, updated automatically by GoReleaser on each release tag.

## R7: API Query Parameters (Filtering)

**Decision**: Expose API-native query parameters as CLI flags on relevant `list` commands.

**Rationale**: The swagger defines specific query parameters per endpoint. Mapping these to CLI flags provides filtering without client-side processing.

**Key filter mappings**:

| Resource | Available Filters |
|----------|------------------|
| features | `--status-id`, `--status-name`, `--parent-id`, `--archived`, `--owner-email`, `--note-id` |
| notes | `--date-from`, `--date-to`, `--created-from`, `--created-to`, `--updated-from`, `--updated-to`, `--term`, `--feature-id`, `--company-id`, `--owner-email`, `--source`, `--any-tag`, `--all-tags` |
| companies | `--term`, `--has-notes`, `--feature-id` |
| releases | `--release-group-id` |
| feature-release-assignments | `--feature-id`, `--release-id`, `--release-state`, `--end-date-from`, `--end-date-to` |
| jira-integrations connections | `--issue-key`, `--issue-id` |
| custom-fields | `--type` (required) |
| custom-fields-values | `--type`, `--custom-field-id`, `--hierarchy-entity-id` |

## R8: Error Handling Strategy

**Decision**: Structured error types mapped from HTTP status codes with user-friendly messages.

**Rationale**: The CLI should translate API errors into actionable user guidance.

| HTTP Status | CLI Behavior |
|-------------|-------------|
| 401 | "Authentication failed. Check your API token via `pboard configure` or PRODUCTBOARD_API_TOKEN env var." |
| 403 | "Access denied. Your API token may not have permission for this resource." |
| 404 | "Resource not found: {resource_type} with ID {id}" |
| 429 | "Rate limit exceeded. Wait and retry, or reduce request frequency." |
| 5xx | "ProductBoard API error ({code}). Try again later." |
| Network error | "Network error: unable to reach api.productboard.com. Check your internet connection." |

## R9: Testing Strategy

**Decision**: Unit tests for formatters and config; integration tests with HTTP mocking for API client; no live API tests in CI.

**Rationale**: Live API tests require a real ProductBoard token and workspace, making them unsuitable for CI. HTTP-level mocking (httptest) validates request construction and response parsing. Unit tests cover formatting logic and config file handling.

**Test structure**:
- `*_test.go` files alongside source
- `testdata/` for fixture JSON responses
- `httptest.Server` for API client integration tests
