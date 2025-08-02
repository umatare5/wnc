---
description: GitHub Actions Workflow Instructions
applyTo: ".github/workflows/**/*.yml,.github/workflows/**/*.yaml"
---

# GitHub Actions Instructions

Comprehensive guidance for creating secure, efficient GitHub Actions workflows.

**Copilot MUST follow all instructions. General instructions take priority in conflicts.**

---

## 🎯 Required Workflow Order

Execute steps in this order for optimal fail-fast behavior:

1. **Environment Setup** (Go, Node.js, Python)
2. **Dependencies** (go mod download, npm install, pip install)
3. **Static Analysis** (go vet, eslint, flake8)
4. **Formatting** (gofmt, prettier, black)
5. **Linting** (golangci-lint, eslint, pylint)
6. **Security** (gosec, semgrep, bandit)
7. **Unit Tests**
8. **Integration Tests**
9. **Coverage Analysis**

---

## ✅ Approved Tools & Actions

### GitHub Actions (1,000+ stars required)

- `actions/checkout@v4` - Repository checkout
- `actions/setup-go@v5` - Go environment
- `actions/setup-node@v4` - Node.js environment
- `actions/setup-python@v5` - Python environment
- `github/codeql-action@v3` - Security scanning
- `actions/upload-artifact@v4` - Artifact upload

### Command-Line Tools

**Go**: `gotestsum` (required), `go vet`, `golangci-lint`, `gosec`
**Node.js**: `npm test`, `eslint`, `prettier`
**Python**: `pytest`, `flake8`, `black`, `bandit`

---

## 📊 Coverage Standards

### Native Tools Only

- **Go**: `go tool cover`
- **Node.js**: `nyc` or `c8`
- **Python**: `coverage.py`

**Prohibited**: External services (Codecov, etc.) due to security concerns.

### Coverage Extraction Examples

**Go**:

```yaml
- name: Extract coverage
  run: |
    coverage=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
    echo "percentage=${coverage}" >> $GITHUB_OUTPUT
```

**Node.js**:

```yaml
- name: Extract coverage
  run: |
    coverage=$(npm run coverage:report | grep -o '[0-9]\+\.[0-9]\+' | tail -1)
    echo "percentage=${coverage}" >> $GITHUB_OUTPUT
```

**Python**:

```yaml
- name: Extract coverage
  run: |
    coverage=$(coverage report | grep TOTAL | awk '{print $4}' | sed 's/%//')
    echo "percentage=${coverage}" >> $GITHUB_OUTPUT
```

---

## ⚡ Best Practices

### Security Requirements

- **Pin actions to commit SHAs**: `actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683 # v4.2.2`
- **Use GitHub Secrets**: `${{ secrets.GITHUB_TOKEN }}`
- **Avoid external services**: Use native tools only

### Performance Optimization

- **Enable caching**: Language-specific module caching
- **Conditional execution**: Use `if: github.event_name == 'pull_request'`
- **Non-blocking externals**: `fail_ci_if_error: false`

### Configuration Requirements

**Environment Variables** (centralized at workflow level):

```yaml
env:
  GO_VERSION: "1.24.5"
  COVERAGE_FILE: "./tmp/coverage.out"
  BUILD_DIR: "./tmp"
```

**Language Setup Examples**:

```yaml
# Go
- uses: actions/setup-go@v5
  with:
    go-version: ${{ env.GO_VERSION }}
    cache: true

# Node.js
- uses: actions/setup-node@v4
  with:
    node-version: ${{ env.NODE_VERSION }}
    cache: "npm"

# Python
- uses: actions/setup-python@v5
  with:
    python-version: ${{ env.PYTHON_VERSION }}
    cache: "pip"
```

### Quality Standards

**Clear Naming**:

```yaml
- name: Run tests with coverage # ✅ Clear
- name: Extract coverage percentage # ✅ Specific
- name: Test # ❌ Vague
```

**Simplified Logic**:

```yaml
# ✅ Simple and readable
- name: Check formatting
  run: |
    if [ -n "$(gofmt -l .)" ]; then
      echo "❌ Files need formatting"
      exit 1
    fi
```

### Error Prevention Checklist

Before finalizing workflows:

- [ ] Actions pinned to commits
- [ ] Static analysis before tests
- [ ] External services non-blocking
- [ ] Environment variables centralized
- [ ] Language-specific caching enabled
- [ ] Only 1K+ star actions used

---

## 📚 Reference

For working examples, see `.github/workflows/` directory files.
