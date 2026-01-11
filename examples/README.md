# Examples - Revenium Middleware for Runway ML (Go)

## Prerequisites

- Go 1.21+
- Runway API key ([app.runwayml.com](https://app.runwayml.com/))
- Revenium API key ([app.revenium.ai](https://app.revenium.ai))

## Setup

```bash
git clone https://github.com/revenium/revenium-middleware-runway-go.git
cd revenium-middleware-runway-go
go mod download
cp .env.example .env
# Edit .env with your API keys
```

## Examples

| Example | Description | Run |
|---------|-------------|-----|
| `basic/` | Image-to-video generation with Gen-3 Alpha | `go run examples/basic/main.go` |

## Environment Variables

```bash
RUNWAY_API_KEY=your_runway_api_key
REVENIUM_METERING_API_KEY=your_revenium_key
REVENIUM_METERING_BASE_URL=https://api.revenium.ai
```

## Supported Models

- `gen3a_turbo` - Gen-3 Alpha Turbo (fast)
- `gen3a` - Gen-3 Alpha (higher quality)

## Supported Options

- **Durations:** 5 or 10 seconds
- **Ratios:** 1280:768 (landscape), 768:1280 (portrait)

## Support

- [Documentation](https://docs.revenium.io)
- [Issues](https://github.com/revenium/revenium-middleware-runway-go/issues)
