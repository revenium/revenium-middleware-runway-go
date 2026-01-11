# AGENTS.md

> Machine-readable instructions for AI agents. Human docs: [README.md](README.md)

## Project Context

**Type**: Go middleware for Runway ML API
**Purpose**: Add Revenium metering to Runway video generation
**Stack**: Go 1.21+
**Module**: `github.com/revenium/revenium-middleware-runway-go`

## Commands

```bash
# Build
go build ./...

# Test
go test ./...

# Run example
cd examples && go run getting_started.go
```

## Architecture

```
revenium/
├── client.go      # Runway client wrapper
├── config.go      # Configuration and validation
├── errors.go      # Error types
├── logger.go      # Logging utilities
├── metering.go    # Revenium metering (fire-and-forget)
├── middleware.go  # Core middleware logic
├── types.go       # Request/response types
└── version.go     # Dynamic version detection
```

## Critical Constraints

1. **Dynamic versioning** - Use `GetMiddlewareSource()` for metering payloads (never hardcode)
2. **Billing fields at TOP LEVEL** - `durationSeconds` NOT in attributes
3. **Fire-and-forget metering** - Use goroutines, never block main request
4. **Auth header** - Use `x-api-key` (not `Authorization: Bearer`)

## Environment Variables

```bash
RUNWAY_API_KEY=...                 # Required: Runway API key
REVENIUM_METERING_API_KEY=hak_...  # Required: Revenium key (starts with hak_)
REVENIUM_DEBUG=true                # Optional: Enable debug logging
```

## Metering Endpoints

| Type | Endpoint | Key Fields |
|------|----------|------------|
| Video | `/meter/v2/ai/video` | `durationSeconds` |

## Supported Operations

| Operation | Method | Notes |
|-----------|--------|-------|
| Image-to-Video | `ImageToVideo()` | gen3a_turbo model |
| Video-to-Video | `VideoToVideo()` | Transformation |
| Upscale | `UpscaleVideo()` | Resolution enhancement |

## Common Errors

| Error | Fix |
|-------|-----|
| `package not found` | `go mod tidy` |
| `metering not tracking` | Check `REVENIUM_METERING_API_KEY` |
| `Runway auth failed` | Check `RUNWAY_API_KEY` |
| `task timeout` | Adjust polling config |

## References

- [Revenium Docs](https://docs.revenium.io)
- [Runway API Docs](https://docs.runwayml.com/)
- [AGENTS.md Spec](https://agents.md/)
