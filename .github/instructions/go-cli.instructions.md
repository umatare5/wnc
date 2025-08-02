---
description: Go CLI Application Development Instructions
applyTo: "cmd/*.go,internal/**/*.go,pkg/**/*.go"
---

# Go CLI Instructions

Guidelines for developing CLI applications.

**Copilot MUST follow all instructions. General instructions take priority in conflicts.**

---

## 🎯 Architecture

Clean Architecture with dependency injection:

```
cli/         → Commands, flags, UI
framework/   → Data formatting, presentation
application/ → Business logic, use cases
infrastructure/ → Data access, external APIs
config/      → Configuration management
```

**Dependency Pattern**:

```go
c := config.New()
r := infrastructure.New(&c)
u := application.New(&c, &r)
f := framework.NewShowCli(&c, &r, &u)
```

---

## 🛠️ CLI Development

### Command Structure

```go
func RegisterXxxSubCommand() []*cli.Command {
    return []*cli.Command{
        {
            Name:      "command-name",
            Usage:     "Brief description",
            UsageText: "app command-name [options...]",
            Flags:     registerXxxCmdFlags(),
            Action:    commandAction,
        },
    }
}
```

### Flag Management

- **Reuse common flags**: `registerEndpointsFlag()`, `registerPrintFormatFlag()`
- **Environment support**: `cli.EnvVars("APP_ENDPOINTS")`
- **Validation**: Validate flags before processing with clear error messages

### Error Handling

- **CLI fatal errors**: Use `log.Fatal()` for startup/config issues only
- **Action functions**: Return errors to CLI framework
- **User messages**: Clear, actionable error messages with troubleshooting hints

---

## 📋 Layer Responsibilities

### CLI Layer (`internal/cli/`)

- Register commands and flags
- Handle user interaction and dependency injection
- CLI-specific error handling

### Framework Layer (`internal/framework/`)

- Format data for display (table/JSON)
- Sort and filter results
- Handle output format switching

### Application Layer (`internal/application/`)

- Business logic and use cases
- Coordinate multiple data sources
- Apply business rules

### Infrastructure Layer (`internal/infrastructure/`)

- External API calls and data access
- Client management and authentication
- API response processing

### Configuration Layer (`internal/config/`)

- Parse and validate CLI flags
- Manage endpoint configurations
- Handle environment variables

---

## 🎨 User Experience

### Shell-Friendly Design

- Support piping and shell scripting
- Provide table and JSON output formats
- Use appropriate exit codes
- Follow Unix conventions

### Multiple Endpoint Support

- Handle multiple endpoints gracefully
- Continue processing if one endpoint fails
- Clearly indicate data source

### Output Standards

- **Table**: Consistent headers, sorting, visual indicators (✅❌)
- **JSON**: Preserve full structure, consistent field names
- **Data**: Convert API enums to human-readable strings

---

## 🔒 Security & Performance

### Security

- Never log authentication tokens
- Default to secure connections
- Provide `--insecure` flag when necessary
- Support environment variables for automation

### Performance

- Process multiple endpoints concurrently
- Handle timeouts gracefully
- Implement efficient sorting
- Reuse HTTP clients

---

## 📦 Package Management

### Third-Party Library Rule

**MANDATORY**: All external libraries MUST be wrapped in `pkg/` before use in `internal/` layers.

```
pkg/
├── logger/              # Wraps logrus
├── tablewriter/         # Wraps tablewriter
└── client/              # Wraps external API client
    ├── client.go
    ├── resource.go
    └── types.go
```

**Usage Pattern**:

```go
// ✅ Correct: Use pkg wrapper
import "github.com/org/app/pkg/client"

// ❌ Wrong: Direct external import in internal/
import "github.com/external/api-client/resource"
```

**Benefits**: Abstraction, easier testing, consistent error handling, maintainability.

---

## 🚀 Development Workflow

When adding new commands:

1. Create CLI registration in `internal/cli/[category]/`
2. Add flag definitions following existing patterns
3. Implement framework layer for data presentation
4. Develop application layer for business logic
5. Create infrastructure layer for data access
6. Add configuration support
7. Write comprehensive tests
8. Update documentation

---

## � Testing & Documentation

### Testing Guidelines

- **Integration**: Test with real endpoints when possible
- **Flags**: Test parsing, validation, and error handling
- **Output**: Test both table and JSON formats

### Documentation Standards

- Provide clear usage examples in command help
- Reference external API models in code comments
- Document common error scenarios with resolution steps
- Maintain CLI interface compatibility across versions

---
