---
description: Go CLI Application Development Instructions
applyTo: "cmd/*.go,internal/**/*.go,pkg/**/*.go"
---

# GitHub Copilot Agent Mode – Go CLI Application Development Instructions

Copilot **MUST** comply with all instructions described in this document when editing or creating any Go code in this repository.

However, when there are conflicts between this document and `general.instructions.md`, **ALWAYS** prioritize the instructions in `general.instructions.md`.

---

## 🎯 Primary Goal

Develop and maintain a production-ready CLI application for managing Cisco C9800 WNC. **DO NOT** develop library/SDK code.

---

## 🧭 Architecture & Design Principles

- **Clean Architecture with Clear Separation:**

  ```
  cli/         → CLI layer (commands, flags, UI)
  framework/   → Framework layer (controller logic, data formatting)
  application/ → Application layer (business logic, use cases)
  infrastructure/ → Infrastructure layer (data access, external APIs)
  config/      → Configuration layer (global config management)
  ```

  - Diagram of the architecture:

    ```plaintext
    +----------------------+    +----------------------+    +----------------------+    +---------------------+
    |        cli/          | -> |      framework/      | -> |     application/     | -> |   infrastructure/   |
    |   (CLI & UI Layer)   |    |  (C9800 WNC Access)  |    |   (Business Logic)   |    |    (Data Access)    |
    +--------+-------------+    +--------+-------------+    +--------+-------------+    +--------+------------+
             |                           |                           |                           |
             v                           v                           v                           v
    +---------------------------------------------------------------------------------------------------------+
    |                                                  config/                                                |
    |                                           (Configuration Layer)                                         |
    +---------------------------------------------------------------------------------------------------------+
    ```

- **Dependency Injection Pattern:**
  Always inject dependencies through constructors. Each layer should receive its dependencies explicitly:

  ```go
  // Example pattern used throughout the codebase
  c := config.New()
  r := infrastructure.New(&c)
  u := application.New(&c, &r)
  f := framework.NewShowCli(&c, &r, &u)
  ```

- **Package-per-Command Structure:**
  Organize CLI commands into separate packages (`show/`, `generate/`) with clear responsibilities.

- **Avoid Global State:**
  Pass configuration and state through struct fields and function parameters, never use global variables.

---

## 🛠️ CLI Development Practices

- **Follow Go CLI Best Practices:**
  Conform to [Go CLI best practices](https://go.dev/doc/effective_go) and CLI design principles from the Unix philosophy.

- **Command Structure:**

  ```go
  // Each command should follow this pattern
  func RegisterXxxSubCommand() []*cli.Command {
      return []*cli.Command{
          {
              Name:      "command-name",
              Usage:     "Brief description",
              UsageText: "wnc command-name [options...]",
              Aliases:   []string{"alias"},
              Flags:     registerXxxCmdFlags(),
              Action:    commandAction,
          },
      }
  }
  ```

- **Flag Organization:**
  Group related flags in separate functions. Reuse common flags across commands:

  ```go
  // Reusable flag patterns
  flags = append(flags, registerControllersFlag()...)
  flags = append(flags, registerInsecureFlag()...)
  flags = append(flags, registerPrintFormatFlag()...)
  ```

- **Error Handling for CLI:**
  - Use `log.Fatal()` only for CLI-specific fatal errors (configuration validation, startup issues)
  - Return errors from action functions to be handled by the CLI framework
  - Provide clear, actionable error messages to users

---

## 🏗️ Layer-Specific Guidelines

### CLI Layer (`internal/cli/`)

- **Purpose:** Handle command registration, flag parsing, and user interaction
- **Responsibilities:**

  - Register commands and subcommands
  - Define and validate CLI flags
  - Coordinate dependency injection
  - Handle CLI-specific error cases

- **Key Patterns:**

  ```go
  // Action pattern for commands
  Action: func(ctx context.Context, cmd *cli.Command) error {
      c := config.New()
      r := infrastructure.New(&c)
      u := application.New(&c, &r)
      f := framework.NewShowCli(&c, &r, &u)

      c.SetShowCmdConfig(cmd)
      f.InvokeXxxCli().DoSomething()
      return nil
  }
  ```

### Framework Layer (`internal/framework/`)

- **Purpose:** Handle presentation logic and data formatting
- **Responsibilities:**

  - Format data for display (table, JSON)
  - Sort and filter results
  - Convert raw API responses to user-friendly output
  - Handle output format switching

- **Key Patterns:**
  ```go
  // Display logic with format support
  if isJSONFormat(tc.Config.ShowCmdConfig.PrintFormat) {
      tc.renderAsJSON(data)
  } else {
      tc.renderAsTable(data)
  }
  ```

### Application Layer (`internal/application/`)

- **Purpose:** Implement business logic and use cases
- **Responsibilities:**
  - Coordinate multiple data sources
  - Apply business rules
  - Merge data from multiple controllers
  - Transform data for presentation

### Infrastructure Layer (`internal/infrastructure/`)

- **Purpose:** Handle external API calls and data access
- **Responsibilities:**
  - Create and manage WNC clients
  - Handle API authentication
  - Process API responses
  - Log API-related errors

### Configuration Layer (`internal/config/`)

- **Purpose:** Manage application configuration
- **Responsibilities:**
  - Parse and validate CLI flags
  - Manage controller configurations
  - Handle environment variables
  - Validate user inputs

---

## 🎨 User Experience Guidelines

- **Shell-Friendly Design:**

  - Support piping and shell scripting
  - Provide both table and JSON output formats
  - Use appropriate exit codes
  - Support common shell conventions

- **Multiple Controller Support:**

  - Handle multiple controllers gracefully
  - Continue processing if one controller fails
  - Clearly indicate which controller data comes from

- **Error Messages:**

  - Provide clear, actionable error messages
  - Include troubleshooting hints for common issues
  - Use consistent error formatting

- **Output Consistency:**
  - Maintain consistent column headers across commands
  - Use consistent sorting and filtering patterns
  - Support common sorting options (asc/desc)

---

## 🔧 Configuration and Flag Management

- **Flag Validation:**

  ```go
  // Validate flags before processing
  func (c *Config) validateShowCmdFlags(cli *cli.Command) error {
      if err := c.validateControllersFormat(cli.String(ControllersFlagName)); err != nil {
          log.Fatal(err)
      }
      return nil
  }
  ```

- **Controller Format:**
  Support flexible controller specification: `https://wnc1.example.com:token,wnc2.example.com:token`

- **Environment Variable Support:**
  Use `cli.EnvVars("WNC_CONTROLLERS")` for common flags to support automation

---

## 🎭 Data Presentation Standards

- **Table Output:**

  - Use consistent column alignment
  - Implement sorting by user-specified fields
  - Include visual indicators (✅❌) for status information
  - Keep table headers descriptive but concise

- **JSON Output:**

  - Preserve full data structure
  - Use consistent field names
  - Support programmatic processing

- **Data Conversion:**
  - Convert YANG model enums to human-readable strings
  - Reference official Cisco YANG models in comments
  - Handle missing or null data gracefully

---

## 🧪 CLI Testing Guidelines

- **Integration Testing:**

  - Test with real WNC controllers when possible
  - Use environment variables for test configuration
  - Skip integration tests if controllers are unavailable

- **Flag Testing:**

  - Test flag parsing and validation
  - Verify error handling for invalid inputs
  - Test environment variable integration

- **Output Testing:**
  - Test both table and JSON output formats
  - Verify sorting and filtering functionality
  - Test error message clarity

---

## 📊 Performance Considerations

- **Concurrent Controller Access:**

  - Process multiple controllers concurrently when beneficial
  - Handle timeouts gracefully
  - Implement proper context cancellation

- **Data Processing:**

  - Stream large datasets when possible
  - Minimize memory usage for large responses
  - Implement efficient sorting algorithms

- **Network Optimization:**
  - Reuse HTTP clients when possible
  - Implement appropriate timeouts
  - Support connection pooling

---

## 🔒 Security Guidelines

- **Credential Handling:**

  - Never log authentication tokens
  - Support environment variables for automation
  - Warn about insecure certificate usage

- **Network Security:**
  - Default to secure connections
  - Provide `--insecure` flag only when necessary
  - Validate controller hostnames

---

## 📚 Documentation Standards

- **Command Help:**

  - Provide clear usage examples
  - Document all flags with examples
  - Include troubleshooting information

- **Code Comments:**

  - Reference Cisco YANG models for data conversions
  - Document complex business logic
  - Explain non-obvious CLI patterns

- **Error Documentation:**
  - Document common error scenarios
  - Provide resolution steps
  - Link to relevant Cisco documentation

---

## 🚀 Command Development Workflow

When adding new commands:

1. **Create CLI registration** in `internal/cli/[category]/`
2. **Add flag definitions** following existing patterns
3. **Implement framework layer** for data presentation
4. **Develop application layer** for business logic
5. **Create infrastructure layer** for data access
6. **Add configuration support** in `internal/config/`
7. **Write comprehensive tests** for all layers
8. **Update documentation** with examples

---

## 🔄 Maintenance Guidelines

- **Consistency First:**
  Follow existing patterns and conventions throughout the codebase

- **Error Resilience:**
  Handle partial failures gracefully, especially with multiple controllers

- **User Feedback:**
  Provide progress indicators for long-running operations

- **Backward Compatibility:**
  Maintain CLI interface compatibility across versions

---

## 📦 Package Management Rules

### Third-Party Library Integration

**MANDATORY**: All third-party libraries MUST be wrapped in `pkg/` layer before use in any internal layer.

- **Rule**: External dependencies MUST NOT be imported directly in `internal/` layers
- **Purpose**: Provides abstraction, easier testing, and dependency management
- **Pattern**: `pkg/[library-name]/[library-name].go`

#### Required Package Structure:

```
pkg/
├── logger/              # Wraps github.com/sirupsen/logrus
│   └── logger.go
├── tablewriter/         # Wraps github.com/olekukonko/tablewriter
│   └── tablewriter.go
└── cisco/              # Wraps github.com/umatare5/cisco-ios-xe-wireless-go
    ├── client.go       # Client creation and management
    ├── ap.go           # Access Point operations
    ├── client.go       # Client operations
    ├── wlan.go         # WLAN operations
    └── types.go        # Common types and structures
```

#### Implementation Pattern:

```go
// pkg/cisco/ap.go
package cisco

import (
    "context"
    "github.com/umatare5/cisco-ios-xe-wireless-go/ap"
    wnc "github.com/umatare5/cisco-ios-xe-wireless-go"
)

// GetAccessPointCapwapData wraps the external library call
func GetAccessPointCapwapData(client *wnc.Client, ctx context.Context) (*ap.ApOperCapwapDataResponse, error) {
    return ap.GetApCapwapData(client, ctx)
}
```

#### Usage in Infrastructure Layer:

```go
// internal/infrastructure/ap.go
import (
    "github.com/umatare5/wnc/pkg/cisco"  // ✅ Correct: Use pkg wrapper
    // "github.com/umatare5/cisco-ios-xe-wireless-go/ap"  // ❌ Wrong: Direct import
)

func (r *ApRepository) GetApCapwapData(controller, apikey string, isSecure *bool) *ap.ApOperCapwapDataResponse {
    client, err := cisco.NewClient(controller, apikey, isSecure)
    if err != nil {
        return nil
    }

    return cisco.GetAccessPointCapwapData(client, context.Background())
}
```

### Import Rules

1. **Standard library**: Direct imports allowed everywhere
2. **Internal packages**: Use relative imports within the project
3. **Third-party libraries**: MUST go through `pkg/` wrapper layer
4. **Cisco library**: MUST use `pkg/cisco/` wrapper functions

### Dependency Flow

```
internal/cli/ → internal/framework/ → internal/application/ → internal/infrastructure/
                                                                      ↓
                                                                  pkg/ (wrappers)
                                                                      ↓
                                                              third-party libraries
```

### Benefits of pkg/ Wrapper Pattern

- **Abstraction**: Hide complex external APIs behind simple interfaces
- **Testing**: Easy to mock external dependencies
- **Consistency**: Unified error handling and logging
- **Maintainability**: Changes to external libraries isolated to pkg layer
- **Documentation**: Clear interfaces with application-specific documentation

---
