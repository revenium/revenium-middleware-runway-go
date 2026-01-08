# Revenium Middleware for Runway ML (Go)

A lightweight, production-ready middleware that adds **Revenium metering and tracking** to Runway ML API calls.

[![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)](https://golang.org/)
[![Documentation](https://img.shields.io/badge/docs-revenium.io-blue)](https://docs.revenium.io)
[![Website](https://img.shields.io/badge/website-revenium.ai-blue)](https://www.revenium.ai)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- **Seamless Integration** - Drop-in middleware with minimal code changes
- **Automatic Metering** - Tracks all API calls with detailed usage metrics
- **Gen-3 Alpha Support** - Full support for Runway's latest video generation
- **Image-to-Video** - Generate videos from images with text prompts
- **Video-to-Video** - Transform existing videos with AI
- **Video Upscaling** - Enhance video resolution and quality
- **Custom Metadata** - Add custom tracking metadata to any request
- **Production Ready** - Automatic task polling with retry logic

## Getting Started (5 minutes)

### Step 1: Create Your Project

```bash
mkdir my-runway-app
cd my-runway-app
go mod init my-runway-app
```

### Step 2: Install Dependencies

```bash
go get github.com/revenium/revenium-middleware-runway-go
go mod tidy
```

### Step 3: Create Environment File

Create `.env` file in your project root:

```bash
# Required - Get from https://app.runwayml.com/video-tools/gen-3/turbo/settings/api-keys
RUNWAY_API_KEY=your_runway_api_key_here

# Required - Get from Revenium dashboard (https://app.revenium.ai)
REVENIUM_METERING_API_KEY=your_revenium_api_key_here

# Optional - Revenium API base URL (defaults to production)
REVENIUM_METERING_BASE_URL=https://api.revenium.ai
```

**Replace the API keys with your actual keys!**

> **Automatic .env Loading**: The middleware automatically loads `.env` files from your project directory. No need to manually export environment variables!

## Examples

This repository includes runnable examples demonstrating how to use the Revenium middleware with Runway ML:

- **[Examples Guide](./examples/README.md)** - Detailed guide for running examples
- **Go Examples**: `examples/basic/`, `examples/advanced/`

**Run examples after setup:**

```bash
# Clone this repository:
git clone https://github.com/revenium/revenium-middleware-runway-go.git
cd revenium-middleware-runway-go
go mod download
go mod tidy

# Run examples:
go run examples/basic/main.go
go run examples/advanced/main.go
```

See **[Examples Guide](./examples/README.md)** for detailed setup instructions and what each example demonstrates.

## What Gets Tracked

The middleware automatically captures:

- **Video Duration**: Length of generated videos in seconds
- **Operation Type**: Image-to-video, video-to-video, or upscale
- **Request Duration**: Total time for each API call (including polling)
- **Model Information**: Which Runway model was used
- **Credits Consumed**: Runway credits used for the generation
- **Custom Metadata**: Business context you provide
- **Error Tracking**: Failed requests and error details

## Environment Variables

### Required

```bash
RUNWAY_API_KEY=your_runway_api_key_here
REVENIUM_METERING_API_KEY=your_revenium_api_key_here
```

### Optional

```bash
# Runway API base URL (defaults to production)
RUNWAY_BASE_URL=https://api.runwayml.com

# Runway API version
RUNWAY_VERSION=2024-11-06

# Revenium API base URL (defaults to production)
REVENIUM_METERING_BASE_URL=https://api.revenium.ai

# Default metadata for all requests
REVENIUM_ORGANIZATION_ID=my-company
REVENIUM_PRODUCT_ID=my-app

# Debug logging
REVENIUM_LOG_LEVEL=INFO
REVENIUM_VERBOSE_STARTUP=false
```

## Supported Operations

### Image to Video

Generate videos from static images with text prompts.

- **Model**: `gen3a_turbo`
- **Duration**: 5 or 10 seconds
- **Aspect Ratios**: 16:9, 9:16, 1:1

### Video to Video

Transform existing videos with AI-powered effects and styles.

- **Model**: `gen3a_turbo`
- **Duration**: 5 or 10 seconds

### Video Upscale

Enhance video resolution and quality.

- **Model**: `upscale`

## Troubleshooting

### Metering data not appearing in Revenium dashboard

**Problem**: Your app runs successfully but no data appears in Revenium.

**Solution**: The middleware sends metering data asynchronously in the background. If your program exits too quickly, the data won't be sent. Add a delay before exit:

```go
// At the end of your main() function
time.Sleep(2 * time.Second)
```

### "Failed to initialize" error

Check your API keys:

```bash
echo $RUNWAY_API_KEY
echo $REVENIUM_METERING_API_KEY
```

### Task polling timeout

Runway video generation can take several minutes. The middleware polls automatically with exponential backoff. Default timeout is 20 minutes.

### Enable debug logging

```bash
export REVENIUM_LOG_LEVEL=DEBUG
go run main.go
```

## Requirements

- **Go**: 1.21 or higher
- **Runway ML API Key**: Get from [app.runwayml.com](https://app.runwayml.com/video-tools/gen-3/turbo/settings/api-keys)
- **Revenium API Key**: Get from [app.revenium.ai](https://app.revenium.ai)

## Documentation

For more information and advanced usage:

- [Revenium Documentation](https://docs.revenium.io)
- [Runway ML Documentation](https://docs.runwayml.com)

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)

## Code of Conduct

See [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)

## Security

See [SECURITY.md](SECURITY.md)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For issues, feature requests, or contributions:

- **GitHub Repository**: [revenium/revenium-middleware-runway-go](https://github.com/revenium/revenium-middleware-runway-go)
- **Issues**: [Report bugs or request features](https://github.com/revenium/revenium-middleware-runway-go/issues)
- **Documentation**: [docs.revenium.io](https://docs.revenium.io)
- **Contact**: Reach out to the Revenium team for additional support

---

**Built by Revenium**
