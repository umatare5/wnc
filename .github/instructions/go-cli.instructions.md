---
description: Go CLI Application Development Instructions
applyTo: "cmd/*.go,internal/**/*.go,pkg/**/*.go"
---

# Go CLI Instructions

**Copilot MUST follow all instructions. General instructions take priority in conflicts.**

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

## 🛠️ Development

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

### Best Practices

- **Flags**: Reuse common flags, support environment variables, validate input
- **Errors**: Return errors from actions, use `log.Fatal()` only for startup issues
- **Output**: Support both table and JSON formats
- **Shell**: Follow Unix conventions, support piping

---

## 🧪 Testing

### Test Structure

```
tests/integration/    → End-to-end CLI tests
internal/.../        → Unit tests (*_test.go)
internal/mock/       → Generated mocks
```

### Key Patterns

```go
// Table-driven tests
func TestFeature(t *testing.T) {
    tests := []struct {
        name string
        test func(t *testing.T)
    }{
        {"scenario", func(t *testing.T) { ... }},
    }
}

// Context safety
type testContextKey string
const testKey testContextKey = "test"

// Skip integration without env
if os.Getenv("WNC_CONTROLLERS") == "" {
    t.Skip("Integration test requires WNC_CONTROLLERS")
}
```

### Coverage

- **Target**: App >80%, Framework >80%, Infrastructure >80%, Packages >80%
- **Commands**: `make test-coverage`, `make test-coverage-filtered`, `make test-coverage-html`
- **Exclusions**: main.go files, mock directories, test utilities

---

## 📦 Package Management

**MANDATORY**: Wrap all external libraries in `pkg/` before use in `internal/`:

```
pkg/
├── logger/         # Wraps logrus
├── tablewriter/    # Wraps tablewriter
└── client/         # Wraps API client
```

```go
// ✅ Correct
import "github.com/org/app/pkg/client"

// ❌ Wrong
import "github.com/external/api-client"
```

---

## 🚀 Development Workflow

1. Create CLI registration in `internal/cli/[category]/`
2. Add flags following existing patterns
3. Implement framework layer for presentation
4. Develop application layer for business logic
5. Create infrastructure layer for data access
6. Add configuration support
7. Write tests (unit + integration)
8. Update documentation
