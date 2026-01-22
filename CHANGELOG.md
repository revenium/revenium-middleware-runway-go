# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2026-01-22

### Added
- Opt-in prompt capture for analytics via `WithCapturePrompts(true)` option
  - `inputMessages`: JSON array with role/content format for generation prompts
  - `outputResponse`: Generated video URLs from Runway
  - `promptsTruncated`: Flag when prompt exceeds 50K character limit
- Request timeout configuration via `RUNWAY_REQUEST_TIMEOUT` environment variable
- Comprehensive examples demonstrating UsageMetadata fields

## [1.0.0] - 2026-01-09

### Added
- Initial release of Revenium middleware for Runway ML
- Support for Runway ML Gen-3 Alpha API
- Image-to-video generation (`ImageToVideo`)
- Video-to-video transformation (`VideoToVideo`)
- Video upscaling (`UpscaleVideo`)
- Automatic task polling with configurable timeouts and exponential backoff
- Automatic Revenium metering for video generation operations
- Asynchronous metering (fire-and-forget pattern)
- Comprehensive error handling with typed errors
- Context-based metadata support
- Environment variable configuration
- Programmatic configuration options
- Configurable logging (DEBUG, INFO, WARN, ERROR levels)
- Production-ready retry logic for API calls
- Complete API documentation
- Working examples in `examples/` directory

### Configuration
- Environment variables support for API keys and settings
- `.env` file support for local development
- Programmatic configuration via `Config` struct
- Default values for all optional settings

### Metering
- Automatic metering to Revenium API (`/meter/v2/ai/video`)
- Support for custom metadata (organization, product, subscriber)
- Video-specific operation type tracking
- Duration and status tracking
- Error and failure tracking

### Error Handling
- `ConfigError` - Configuration errors
- `MeteringError` - Metering API errors
- `ProviderError` - Runway API errors
- `AuthError` - Authentication errors
- `NetworkError` - Network/HTTP errors
- `TaskError` - Task polling errors
- `ValidationError` - Request validation errors
- `InternalError` - Internal middleware errors

### Documentation
- Comprehensive README with examples
- API reference documentation
- Configuration guide
- Error handling guide
- Metering documentation

