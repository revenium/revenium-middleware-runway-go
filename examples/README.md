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
| `basic/` | Image-to-video generation with Gen-3 Alpha and prompt capture | `go run examples/basic/main.go` |
| `comprehensive/` | Full metering field test with all UsageMetadata fields (Scenario A) | `go run examples/comprehensive/main.go` |
| `comprehensive-b/` | Hard-coding detection test with different values (Scenario B) | `go run examples/comprehensive-b/main.go` |

### Comprehensive Examples

The `comprehensive/` and `comprehensive-b/` examples demonstrate ALL available metering fields with realistic enterprise values. Run BOTH and compare payloads to verify no hard-coding.

**UsageMetadata Fields:**
- `OrganizationID` - Customer organization identifier
- `ProductID` - Product/feature identifier
- `TaskType` - Type of work being performed
- `TaskID` - Unique task identifier
- `Agent` - Worker/service instance identifier
- `SubscriptionID` - Billing subscription reference
- `TraceID` - Distributed tracing correlation ID
- `ParentTransactionID` - Parent request for nested traces
- `TraceType` - Type of trace (e.g., "distributed")
- `TraceName` - Human-readable trace name
- `Environment` - Deployment environment
- `Region` - Geographic region
- `CredentialAlias` - Credential key alias
- `RetryNumber` - Retry attempt counter
- `ResponseQualityScore` - Quality metric (0.0-1.0)
- `VideoJobID` - Video job correlation ID
- `AudioJobID` - Audio job correlation ID
- `Subscriber` - Detailed subscriber information (nested map)
- `Custom` - Business-specific custom fields (merged to top level)

**Prompt Capture Fields (opt-in):**
- `inputMessages` - JSON array with role/content format
- `outputResponse` - Generated video URLs as JSON array
- `promptsTruncated` - true if prompt exceeded 50K character limit

## Environment Variables

```bash
RUNWAY_API_KEY=your_runway_api_key
REVENIUM_METERING_API_KEY=your_revenium_key
REVENIUM_METERING_BASE_URL=https://api.revenium.ai

# Optional: Enable prompt capture for analytics (default: false)
REVENIUM_CAPTURE_PROMPTS=true
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
