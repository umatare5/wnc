---
description: General Instructions
applyTo: "**"
---

# 💡 Special Instructions for Claude Sonnet 4

When you using Claude Sonnet 4:

- **MUST** execute at the start of the prompt:

  - Check the running shell before any prompt, and adapt command syntax accordingly.

- **MUST** comply with the following during prompt processing:

  - Use the command format corresponding to confirmed shell at the start.
  - Use `./tmp` for all temporary files, scripts and directories, including `.test` binaries and coverage reports.

- **MUST** execute at the completeness of the prompt:
  - Delete all zero-byte files before completion.
  - **When modified go code**, run `go build` before completion. Repeat modifications until `go build` completeness.
  - **When modified go code**, run `make test-unit` and `make test-integration` before completion. Repeat modifications until all tests pass.
  - **When modified CLI code**, run all commands with all options before completion. Repeat modifications until all CLI commands work.
  - Summarize results and save to `./.copilot_reports` before completion. The name format should be `<prompt_title>_<timestamp YYYY-MM-DD_HH-mm-ss>.md`.
