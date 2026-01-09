# Revenium Middleware for Runway ML (Go)

Official Go middleware for Runway ML API with automatic Revenium metering and usage tracking.

## Features

- Complete Runway ML API support (Gen-3 Alpha)
  - Image-to-video generation
  - Video-to-video transformation
  - Video upscaling
- Automatic task polling with configurable timeouts
- Automatic Revenium metering for AI usage tracking
- Asynchronous metering (fire-and-forget)
- Comprehensive error handling
- Context-based metadata support
- Production-ready with retry logic

## Installation

```bash
go get github.com/revenium/revenium-middleware-runway-go
```

## Quick Start

### 1. Set up environment variables

Create a `.env` file:

```bash
# Runway API Configuration
RUNWAY_API_KEY=your_runway_api_key_here

# Revenium Metering Configuration
REVENIUM_METERING_API_KEY=hak_your_revenium_api_key_here
```

### 2. Initialize and use the middleware

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/revenium/revenium-middleware-runway-go/revenium"
)

func main() {
    // Initialize middleware
    if err := revenium.Initialize(); err != nil {
        log.Fatal(err)
    }

    // Get client
    client, err := revenium.GetClient()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create context and metadata
    ctx := context.Background()
    metadata := &revenium.UsageMetadata{
        OrganizationID: "org-123",
        ProductID:      "product-abc",
        Subscriber: map[string]interface{}{
            "id":    "user-456",
            "email": "user@example.com",
        },
    }

    // Generate video from image
    req := &revenium.ImageToVideoRequest{
        PromptImage: "https://example.com/image.jpg",
        PromptText:  "A cinematic shot of mountains at sunset",
        Model:       "gen3a_turbo",
        Duration:    5,
        Ratio:       "16:9",
    }

    result, err := client.ImageToVideo(ctx, req, metadata)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Video generated: %v\n", result.OutputURLs)
}
```

## API Reference

### Initialize

```go
// Initialize with environment variables
err := revenium.Initialize()

// Or with options
err := revenium.Initialize(
    revenium.WithRunwayAPIKey("your-key"),
    revenium.WithReveniumAPIKey("hak_your-key"),
)
```

### Image to Video

```go
req := &revenium.ImageToVideoRequest{
    PromptImage: "https://example.com/image.jpg", // URL or base64
    PromptText:  "A cinematic shot",
    Model:       "gen3a_turbo",
    Duration:    5,   // 5 or 10 seconds
    Ratio:       "16:9",
    Seed:        nil, // Optional: for reproducibility
}

result, err := client.ImageToVideo(ctx, req, metadata)
```

### Video to Video

```go
req := &revenium.VideoToVideoRequest{
    PromptVideo: "https://example.com/video.mp4",
    PromptText:  "Transform into anime style",
    Model:       "gen3a_turbo",
    Duration:    10,
}

result, err := client.VideoToVideo(ctx, req, metadata)
```

### Video Upscale

```go
req := &revenium.VideoUpscaleRequest{
    PromptVideo: "https://example.com/video.mp4",
    Model:       "upscale",
}

result, err := client.UpscaleVideo(ctx, req, metadata)
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `RUNWAY_API_KEY` | Runway ML API key | (required) |
| `RUNWAY_BASE_URL` | Runway API base URL | `https://api.runwayml.com` |
| `RUNWAY_VERSION` | Runway API version | `2024-11-06` |
| `REVENIUM_METERING_API_KEY` | Revenium API key | (required) |
| `REVENIUM_METERING_BASE_URL` | Revenium API base URL | `https://api.revenium.ai` |
| `REVENIUM_ORGANIZATION_ID` | Organization ID | (optional) |
| `REVENIUM_PRODUCT_ID` | Product ID | (optional) |
| `REVENIUM_LOG_LEVEL` | Log level (DEBUG/INFO/WARN/ERROR) | `INFO` |
| `REVENIUM_VERBOSE_STARTUP` | Verbose startup logging | `false` |

### Programmatic Configuration

```go
client, err := revenium.NewReveniumRunway(&revenium.Config{
    RunwayAPIKey:    "your-runway-key",
    ReveniumAPIKey:  "hak_your-revenium-key",
    ReveniumBaseURL: "https://api.revenium.ai",
})
```

## Metering

The middleware automatically sends metering data to Revenium for every video generation request. Metering is sent **asynchronously** (fire-and-forget) and includes:

- **Operation Type**: `VIDEO`
- **Provider**: `runway`
- **Model Source**: `RUNWAY`
- **Model**: The specific Runway model used
- **Duration**: Total time from request to completion
- **Status**: Success/failure information
- **Custom Metadata**: Organization, product, subscriber info

### Metadata Structure

```go
metadata := &revenium.UsageMetadata{
    OrganizationID:       "org-123",
    ProductID:            "product-abc",
    TaskType:             "video-generation",
    Agent:                "my-ai-agent",
    SubscriptionID:       "sub-456",
    TraceID:              "trace-789",
    Subscriber: map[string]interface{}{
        "id":    "user-123",
        "email": "user@example.com",
        "name":  "John Doe",
    },
    TaskID:               "task-xyz",
    ResponseQualityScore: &qualityScore, // *float64
    Custom: map[string]interface{}{
        "campaign": "holiday-2024",
    },
}
```

## Task Polling

Video generation is asynchronous. The middleware automatically polls the Runway API until the task completes.

### Default Polling Configuration

- **Max Attempts**: 120
- **Initial Interval**: 2 seconds
- **Max Interval**: 10 seconds (exponential backoff)
- **Timeout**: 20 minutes

### Custom Polling Configuration

```go
pollingConfig := &revenium.PollingConfig{
    MaxAttempts:     60,
    InitialInterval: 3 * time.Second,
    MaxInterval:     15 * time.Second,
    Timeout:         10 * time.Minute,
}

// Use with direct client access
runwayClient := revenium.NewRunwayClient(config)
task, _ := runwayClient.CreateImageToVideo(ctx, req)
result, err := runwayClient.WaitForTaskCompletion(ctx, task.ID, pollingConfig)
```

## Error Handling

The middleware provides typed errors for better error handling:

```go
result, err := client.ImageToVideo(ctx, req, metadata)
if err != nil {
    if revenium.IsConfigError(err) {
        // Configuration issue
    } else if revenium.IsAuthError(err) {
        // Authentication failed
    } else if revenium.IsTaskError(err) {
        // Task polling timeout or failure
    } else if revenium.IsProviderError(err) {
        // Runway API error
    } else if revenium.IsNetworkError(err) {
        // Network issue
    }

    // Get error details
    if revErr, ok := err.(*revenium.ReveniumError); ok {
        details := revErr.GetDetails()
        statusCode := revErr.GetStatusCode()
    }
}
```

## Logging

Control logging verbosity:

```bash
export REVENIUM_LOG_LEVEL=DEBUG
export REVENIUM_VERBOSE_STARTUP=true
```

Or programmatically:

```go
revenium.SetLogger(myCustomLogger) // Implement Logger interface
```

## Examples

See the `examples/` directory for more examples:

- `examples/basic/` - Basic image-to-video generation
- `examples/advanced/` - Advanced features with custom polling

## Testing

```bash
go test ./...
```

## Requirements

- Go 1.21 or higher
- Runway ML API key
- Revenium API key

## Support

- **Documentation**: [https://docs.revenium.io](https://docs.revenium.io)
- **Dashboard**: [https://app.revenium.ai](https://app.revenium.ai)
- **Email**: support@revenium.io

## License

MIT License - see LICENSE file for details
