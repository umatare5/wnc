# 🧪 Testing

This CLI application includes comprehensive tests that validate functionality and behavior across different components:

- **Unit tests**: These tests validate serialization and deserialization between JSON and Go structs used in RESTCONF responses.
- **Mock tests**: REST API interactions are simulated using GoMock to ensure reliable and isolated tests without requiring actual controllers.
- **CLI tests**: CLI behavior is verified using `urfave/cli/v3`, including argument parsing, subcommand execution, and flag validation.
- **Integration tests**: These tests interact with multiple API endpoints to verify API communication and overall functionality.

> [!Note]
> Currently, the test coverage is insufficient. All tests will be covered by the future release `v1.0.0`.

## 🎯 Prerequisites

### 🧩 For Unit, Mock and CLI Tests

Unit tests require no special configuration and can be run in any Go development environment.

| Requirement   | Version/Details  | Description                                          |
| ------------- | ---------------- | ---------------------------------------------------- |
| Go            | 1.24 or later    | Required for running tests and building the project. |
| Testing Tools | Standard library | Built-in Go testing framework.                       |

### 🔗 For Integration Tests

#### 🎛️ 1. Cisco Catalyst 9800 Wireless Network Controller

Integration tests require access to real Cisco Catalyst 9800 WNC(s).

For instructions on setting up WNC, please refer to the [References Section](#references).

> [!CAUTION]
> Integration tests interact with real controllers and may affect their state. Use dedicated test controllers when possible.

#### 🔧 2. Environment Variables

Integration tests require the following environment variables:

| Variable          | Description                         | Example                              |
| ----------------- | ----------------------------------- | ------------------------------------ |
| `WNC_CONTROLLERS` | Controller hostname and token pairs | `192.168.1.100:YWRtaW46cGFzc3dvcmQ=` |

<details><summary>Environment Variable Configuration</summary>

```bash
# Single controller
export WNC_CONTROLLERS="192.168.1.100:YWRtaW46cGFzc3dvcmQ="

# Multiple controllers (comma-separated)
export WNC_CONTROLLERS="192.168.1.100:YWRtaW46cGFzc3dvcmQ=,192.168.1.101:YWRtaW46cGFzc3dvcmQ="
```

**Generating Access Tokens:**

Use the `wnc generate token` command to create Base64 encoded access tokens:

```bash
# Generate token for your controller
wnc generate token -u admin -p password
# Output: YWRtaW46cGFzc3dvcmQ=
```

</details>

## 🚀 Running Tests

The project includes convenient Makefile targets for testing:

| Command                 | Description                                                        |
| ----------------------- | ------------------------------------------------------------------ |
| `make test-unit`        | Run unit tests only with enhanced output formatting.               |
| `make test-integration` | Run integration tests with enhanced output. \* Requires WNC access |

<details><summary>Example of gotestsum Enhanced Output</summary>

```text
📦 github.com/umatare5/wnc/cmd (85.7% coverage)
  ✅ TestMainFunction (0.00s)
  ✅ TestVersionCommand (0.01s)

📦 github.com/umatare5/wnc/internal/application (72.3% coverage)
  ✅ TestShowOverview (0.05s)
  ✅ TestShowAP (0.03s)
  ✅ TestShowClient (0.02s)
    application_test.go:156: Show client request successful

📦 github.com/umatare5/wnc/internal/cli (15.2% coverage)
  🚧 TestIntegrationShowOverview (0.00s)
    integration_test.go:45: WNC_CONTROLLERS not set - skipping integration tests
  ✅ TestIntegrationShowAP (5.23s)
    integration_test.go:89: Integration test completed successfully with 3 controllers
```

</details>

## 📊 Test Data Collection

Integration tests automatically collect and save real WNC data to JSON files for validation and debugging purposes.

- **Location**: `./tmp/test_data/` directory
- **Format**: JSON files with descriptive names (e.g., `show_overview_response.json`)
- **Purpose**: Verify API response structure and enable offline debugging

<details><summary>Example of test data tree structure</summary>

```text
./tmp/test_data/
├── show_overview_response.json
├── show_ap_response.json
├── show_client_response.json
├── show_wlan_response.json
└── generate_token_response.json
```

</details>

## 📈 Coverage Analysis

The project supports comprehensive test coverage analysis:

### 📊 Coverage Reports

| Output Type     | Command                   | Description                                  |
| --------------- | ------------------------- | -------------------------------------------- |
| Terminal Output | `make test-coverage`      | Run tests with coverage analysis.            |
| HTML Report     | `make test-coverage-html` | Run tests and generate HTML coverage report. |

<details><summary>Example of Coverage Output</summary>

```text
Coverage report generated at ./tmp/coverage.out
total: (statements) 67.8%

📦 github.com/umatare5/wnc/cmd (85.7% coverage)
📦 github.com/umatare5/wnc/internal/application (72.3% coverage)
📦 github.com/umatare5/wnc/internal/cli (89.1% coverage)
📦 github.com/umatare5/wnc/internal/config (91.2% coverage)
📦 github.com/umatare5/wnc/internal/framework (68.5% coverage)
📦 github.com/umatare5/wnc/internal/infrastructure (45.6% coverage)
📦 github.com/umatare5/wnc/pkg/cisco (82.3% coverage)
📦 github.com/umatare5/wnc/pkg/log (95.0% coverage)
📦 github.com/umatare5/wnc/pkg/tablewriter (78.9% coverage)
```

</details>

## 🔧 Testing Architecture

This CLI follows a layered testing approach that mirrors its clean architecture:

### 📁 Test Organization

| Layer              | Directory                  | Purpose                                   |
| ------------------ | -------------------------- | ----------------------------------------- |
| **Entrypoint**     | `cmd/`                     | Entrypoint of this command-line interface |
| **Application**    | `internal/application/`    | Business logic and use cases              |
| **CLI Framework**  | `internal/cli/`            | CLI framework and command definitions     |
| **Configuration**  | `internal/config/`         | Configuration parsing and validation      |
| **Framework**      | `internal/framework/`      | Framework adapters and interfaces         |
| **Infrastructure** | `internal/infrastructure/` | External API communication                |
| **Packages**       | `pkg/`                     | Reusable utility packages                 |

### 📋 Test Types by Layer

- **Unit Tests**: Focus on individual functions and components
- **Component Tests**: Test layer interactions and business logic
- **Integration Tests**: Validate end-to-end functionality with real controllers (located in `internal/cli/`)

## 📚️ Appendix

### 💡 Testing Tips

For efficient testing workflow, start with unit tests and gradually move to integration tests:

1. **Install Dependencies**: `make deps` - Install gotestsum and other development tools.
2. **Unit Tests First**: `make test-unit` - Ensure basic functionality with enhanced output.
3. **Code Quality Check**: `make lint` - Run linting to catch potential issues.
4. **Environment Setup**: Configure environment variables for integration tests.
5. **Environment Verification**: Test controller connectivity using `wnc show overview`.
6. **Coverage Analysis**: `make test-coverage` - Run tests with coverage analysis.
7. **HTML Coverage Report**: `make test-coverage-html` - Generate detailed HTML coverage report.
8. **Test Data Review**: Examine generated JSON files in `./tmp/test_data/` to understand API responses.
9. **Integration Tests**: `make test-integration` - Test with real controllers.

> [!TIP]
> For comprehensive testing, run both `make test-unit` and `make test-integration` sequentially to validate all functionality.

### 🛠️ Development Dependencies

The project uses several tools to enhance the testing experience:

- **gotestsum**: Provides emoji-enhanced, human-readable test output
- **golangci-lint**: Code linting and static analysis
- **goreleaser**: Release automation and snapshot builds
- **air**: Hot reload for development (optional)

> [!Note]
> Install all dependencies with: `make deps`

### 📖 References

These references provide additional information on Cisco Catalyst 9800 WNC and related technologies:

- 📖 [Cisco Catalyst 9800-CL Wireless Controller for Cloud Deployment Guide](https://www.cisco.com/c/en/us/td/docs/wireless/controller/9800/technical-reference/c9800-cl-dg.html)
  - A comprehensive guide for deploying Cisco Catalyst 9800-CL WNC in cloud environments.
  - This includes setup instructions, configuration examples, and best practices.
- 📖 [Cisco Catalyst 9800 Series Wireless Controller Programmability Guide](https://www.cisco.com/c/en/us/td/docs/wireless/controller/9800/programmability-guide/b_c9800_programmability_cg/cisco-catalyst-9800-series-wireless-controller-programmability-guide.html)
  - A guide for programming and automating Cisco Catalyst 9800 WNC.
  - This includes information on RESTCONF APIs, YANG models, and automation workflows.
- 📖 [YANG Models and Platform Capabilities for Cisco IOS XE 17.12.1](https://github.com/YangModels/yang/tree/main/vendor/cisco/xe/17121#readme)
  - A repository containing YANG models and platform capabilities for Cisco IOS XE 17.12.1.
  - This is useful for understanding the data structures used in the RESTCONF API.
