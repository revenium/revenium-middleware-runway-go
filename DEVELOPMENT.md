# Development Guide - Revenium Middleware for Runway (Go)

## Quick Start for Developers & QA

### Prerequisites

1. **Go 1.21+** installed
2. **Git** for version control
3. **API Keys** (for examples and integration tests)

---

## Setup Instructions

### 1. Clone and Setup

```bash
git clone <repository-url>
cd revenium-middleware-runway-go
go mod download
```

### 2. Environment Configuration

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your actual API keys
# Required:
RUNWAY_API_KEY=key_your_runway_key_here
REVENIUM_METERING_API_KEY=hak_your_key_here
```

---

## Testing Commands (Copy & Paste Ready)

### Build & Verification

```bash
# Build project (must pass)
go build ./...

# Clean dependencies
go mod tidy

# Verify no compilation errors
go vet ./...

# Format code
go fmt ./...
```

### Unit Tests (No API Keys Required)

```bash
# Test config
go test -v ./tests/config_test.go

# All unit tests at once
go test -v ./tests/...
```

### Examples Testing (Requires .env with API Keys)

```bash
# Basic example
go run examples/basic/main.go
```

### Integration Tests (Requires API Keys)

```bash
# Full integration test suite
go test -v ./tests/e2e/...
```

---

## Expected Results

### Unit Tests Should Show:

```
=== RUN   TestConfigInitialization
--- PASS: TestConfigInitialization (0.00s)
...
PASS
```

### Examples Should Show:

```
=== Revenium Runway Middleware - Basic Example ===

Task ID: task_xxx
Status: SUCCEEDED
Model: gen3a_turbo
Duration: 5s
Output URLs: [https://...]

Example completed successfully!
```

### Build Should Show:

```
# No output = success
go build ./...
```

---

## Development Workflow

### Daily Development:

```bash
# 1. Pull latest changes
git pull

# 2. Build and verify
go build ./...
go mod tidy

# 3. Run unit tests
go test -v ./tests/...

# 4. Test your changes with examples
go run examples/basic/main.go
```

### Before Committing:

```bash
# 1. Format code
go fmt ./...

# 2. Verify code
go vet ./...

# 3. Run all unit tests
go test -v ./tests/...

# 4. Test examples
go run examples/basic/main.go

# 5. Build final verification
go build ./...
```

---

## QA Testing Checklist

### Environment Setup

- [ ] Go 1.21+ installed
- [ ] Repository cloned
- [ ] Dependencies downloaded (`go mod download`)
- [ ] `.env` file created with valid API keys

### Build Verification

- [ ] `go build ./...` - No errors
- [ ] `go mod tidy` - No changes needed
- [ ] `go vet ./...` - No warnings

### Unit Tests (No API Keys)

- [ ] Config tests pass
- [ ] All tests pass

### Examples (With API Keys)

- [ ] Basic example works

### Expected Behaviors

- [ ] Video generation completes
- [ ] Task polling works with exponential backoff
- [ ] Metadata included in payloads
- [ ] No manual `export` commands needed
- [ ] Fire-and-forget metering working
- [ ] Dynamic version detection working

---

## Troubleshooting

### Common Issues:

#### "API key not found"

```bash
# Check .env file exists and has correct format
cat .env
# Should show: RUNWAY_API_KEY=key_...
```

#### "Tests failing"

```bash
# Check if .env is interfering with unit tests
# Unit tests should work without .env
mv .env .env.backup
go test -v ./tests/...
mv .env.backup .env
```

#### "Examples not working"

```bash
# Verify .env file is loaded
go run examples/basic/main.go
# Should show: "Runway API key loaded"
```

#### "Task timeout"

```bash
# Runway video generation can take several minutes
# Default timeout is 20 minutes - adjust if needed
# Check task status via API if example seems hung
```

---

## Project Structure

```
revenium-middleware-runway-go/
├── revenium/           # Core middleware code
├── examples/           # Working examples
│   └── basic/         # Basic example
├── tests/             # Unit tests
│   ├── config_test.go # Config tests
│   └── e2e/           # End-to-end tests
├── .env.example       # Environment template
├── .env               # Your API keys (create this)
├── go.mod             # Go dependencies
├── README.md          # User documentation
└── DEVELOPMENT.md     # This file
```

---

## Success Criteria

**The project is working correctly when:**

1. All unit tests pass
2. All examples run without errors
3. Video generation completes successfully
4. Task polling works with exponential backoff
5. Metadata is included in metering payloads
6. No manual `export` commands needed
7. Build completes without errors
8. Fire-and-forget metering works asynchronously
9. Dynamic version detection returns correct version

**Ready for production when all items above are complete**
