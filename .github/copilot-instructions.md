# GitHub Copilot Agent Mode – Repo Instructions (umatare5/wnc)

### Scope & Metadata

- **Last Updated**: 2025-08-11

- **Precedence**: Highest in this repository (see §2)

- **Goals**:

  - **Primary Goal**: Contribute **only** to the **`wnc` CLI** in this repo.
  - Keep user-facing behavior **clear, stable, and script-friendly**; maximize **readability/maintainability**.
  - Prefer **minimal diffs** (avoid unnecessary churn).
  - Defaults are secure (no secret logging; TLS verify **on** by default).

- **Non-Goals**:

  - Creating additional standalone apps or services outside this CLI.
  - Introducing ad-hoc lint/style rules not configured in this repo.
  - Emitting or persisting secrets/test credentials anywhere (code, logs, artifacts).

> Context: `wnc` is a shell-friendly CLI for Cisco Catalyst **9800** Wireless Network Controllers (WNC/WLC) over **RESTCONF** on IOS-XE **17.12.x**. The README lists supported commands, install paths, testing, and the release flow.

---

## 0. Normative Keywords

- NK-001 (**MUST**) Interpret **MUST / MUST NOT / SHOULD / SHOULD NOT / MAY** per RFC 2119/8174.

## 1. Repository Purpose & Scope

- GP-001 (**MUST**) Treat this repository as a **Go CLI** that operates Cisco Catalyst **C9800** WNC via **RESTCONF** (IOS-XE **17.12.x**).
- GP-002 (**SHOULD**) Prefer using the companion SDK **`cisco-ios-xe-wireless-go`** for controller access when available.
- GP-003 (**MUST NOT**) Assume non‑C9800 targets or pre/post‑17.12 behavior unless explicitly guarded and documented.

## 2. Precedence & Applicability

- GA-001 (**MUST**) When editing/generating code in this repository, Copilot **must follow** this document.
- GA-002 (**MUST**) In this repository, **this file has the highest precedence** over any other instruction set. **On conflict, prioritize this file**.
- GA-003 (**MUST**) Lint/format rules follow **repo configs** only (see §5).
- GA-004 (**MUST**) AI review comments and change scope are governed by §12.

## 3. Expert Personas (for AI edits/reviews)

- EP-001 (**MUST**) Act as a **Go 1.24 expert**.
- EP-002 (**MUST**) Act as a **domain expert for Cisco Catalyst 9800 / IOS‑XE 17.12** with RESTCONF basics.
- EP-003 (**SHOULD**) Be familiar with **YANG** models and typical IOS‑XE YANG layout, using them to validate paths and fields.

## 4. Security & Privacy

- SP-001 (**MUST NOT**) Log tokens or credentials. Mask any accidental echo with `${TOKEN:0:6}...`.
- SP-002 (**MUST**) TLS verification is **on by default**. **MUST NOT** default to `--insecure`; allow **opt‑in** for dev only (keep strong warnings in help/docs).
- SP-003 (**MUST**) Avoid printing controller hostnames, usernames, or MACs unless essential to the command output.

## 5. Editor‑Driven Tooling (single source of truth)

- ED-001 (**MUST**) Lint/format/type checks follow repository settings (e.g., `.golangci.yml`, `.markdownlint.json`, `.markdownlintignore`, `.taplo.toml`, and any workspace files under `.vscode` if present).
- ED-002 (**MUST NOT**) Add flags/rules or inline disables that are not configured.
- ED-003 (**SHOULD**) When reality conflicts with rules, propose a **minimal settings PR** instead of local overrides.

## 6. Coding Principles (Basics)

- GC-001 (**MUST**) Apply **KISS/DRY** and maintain small, composable functions.
- GC-002 (**MUST**) Avoid magic numbers; **use named constants**.
- GC-003 (**MUST**) Prefer **predicate helpers** (`is*`, `has*`) to clarify conditionals.
- GC-004 (**SHOULD**) Keep packages cohesive: `cmd/` for Cobra/CLI wiring, `pkg/` for domain wrappers, `internal/` for helpers.

## 7. Coding Principles (Conditionals)

- CF-001 (**MUST**) Use predicate helpers in conditions to keep `if` readable.
- CF-002 (**MUST**) Prefer **early returns** to reduce nesting.

## 8. Coding Principles (Loops)

- LP-001 (**MUST**) In loops, prefer **early exits** (`return` / `break` / `continue`) to reduce complexity.

## 9. Working Directory / Temp Files

- WD-001 (**MUST**) Place all temporary artifacts (work files, coverage, test binaries, etc.) **under `./tmp`**.
- WD-002 (**MUST**) Before completion, delete **zero‑byte files** (**exception**: keep `.keep`).

## 10. Model‑Aware Execution Workflow (when shell execution is available)

- WF-001 (**MUST**) Before actions: **always launch and use `bash`** (no shell autodetect).
- WF-002 (**MUST**) After editing Go code: run `make build` and fix until it passes.
- WF-003 (**MUST**) After editing Go code: run `make test-unit` and fix until green. (If integration env is unavailable, skip safely.)
- WF-004 (**MUST**) After editing **shell scripts** under `./scripts/` (if any): execute with documented options; ensure **impacted `make` targets** succeed.
- WF-005 (**MUST**) On completion: write a brief run log to `./.copilot_reports/<prompt_title>_<YYYY-MM-DD_HH-mm-ss>.md`.

## 11. Tests / Quality Gate (for AI reviewers)

- QG-001 (**MUST**) Keep CI green. Do not merge changes that violate lint/tests.
- QG-002 (**SHOULD**) Prefer **unit tests** for business logic; use **integration tests** only when controller access is required (guarded by env).
- QG-003 (**SHOULD**) Table/JSON output snapshots (if any) must be kept stable; add explicit fixtures.

## 12. Change Scope & Tone (for AI reviewers)

- CS-001 (**MUST**) Focus on the **diff**; propose broad refactors only with explicit request/label (e.g., `allow-wide`).
- CS-002 (**SHOULD**) Tag comments with **\[BLOCKER] / \[MAJOR] / \[MINOR (Nit)] / \[QUESTION] / \[PRAISE]**.
- CS-003 (**SHOULD**) Structure comments as “**TL;DR → Evidence (rule/proof) → Minimal-diff proposal**”.
- CS-004 (**MUST NOT**) Introduce new top‑level commands that are not documented in README or `docs/commands/` without a design note.

## 13. CLI Design Rules

- CL-001 (**MUST**) Follow existing command taxonomy from README:

  - `generate token` (auth helper)
  - `show overview | ap | ap-tag | client | wlan` (read-only summaries)

- CL-002 (**MUST NOT**) Implement destructive/exec features here unless explicitly green‑lit; the README points to **`telee`** for exec use‑cases.
- CL-003 (**MUST**) Keep outputs stable and pipe‑friendly (fixed columns for tables; avoid gratuitous reflow). If adding machine‑readable output, hide behind an explicit flag (e.g., `--format json`) and document it.
- CL-004 (**SHOULD**) Keep global flags consistent across subcommands (e.g., controller list, timeouts, TLS flags).
- CL-005 (**MUST**) Document each new command under `docs/commands/NAME.md` with synopsis, examples, and flags.

## 14. RESTCONF & Controller Interactions

- RC-001 (**MUST**) Use official RESTCONF paths for IOS‑XE **17.12.x** and validate with YANG when feasible.
- RC-002 (**SHOULD**) Prefer calling the `cisco-ios-xe-wireless-go` SDK to raw HTTP when available.
- RC-003 (**MUST**) Authentication helpers:

  - Support **Basic token** generation workflow (base64 `user:pass`).

- RC-004 (**MUST**) TLS:

  - Verify certificates by default; `--insecure` is a **dev‑only** escape hatch with a prominent warning in help.

## 15. Build, Release & Versioning

- BR-001 (**MUST**) Build via `make build`. Distribute via GitHub Releases / `goreleaser` as configured.
- BR-002 (**MUST**) Do **not** edit tags manually. The flow updates `VERSION`, merges PR, then automation tags; follow that flow. Afterwards, use the release workflow to publish artifacts.
- BR-003 (**SHOULD**) Pre‑`v1.0.0` may include breaking changes; call out **BREAKING** in PR title and release notes.

## 16. Docs

- DC-001 (**MUST**) Keep **README** install/usage consistent with the CLI behavior (Docker, binaries, supported platforms, examples, security flags).
- DC-002 (**MUST**) For each command, update/create `docs/commands/NAME.md` and cross‑link from README’s **CLI Reference**.
- DC-003 (**SHOULD**) Update **TESTING** and **TROUBLESHOOTING** docs if behavior or flags change.

## 17. Working Repo Realities

- WR-001 (**MUST**) This repo is **work‑in‑progress**; some APIs/commands are **not yet implemented**. Preserve placeholders and notes that signal planned coverage (e.g., `exec` section notes).
- WR-002 (**SHOULD**) Where controller/test resources are missing, structure code to allow **safe no‑op/skip** paths behind flags or env checks.

## 18. Quick Checklist (before completion)

- QC-001 (**MUST**) Behavior matches README and command docs; help text accurate.
- QC-002 (**MUST**) Lint/format clean per repo configs; no ad‑hoc disables.
- QC-003 (**MUST**) Required **Make** targets pass (build ± unit/integration if env set).
- QC-004 (**MUST**) Temp artifacts under `./tmp`, zero‑byte files removed, report written to `./.copilot_reports/`.
- QC-005 (**SHOULD**) Outputs remain stable (no breaking column changes) unless release notes mark **BREAKING**.
- QC-006 (**SHOULD**) If you added controller calls, prefer SDK wrappers where available.

---

### Appendix A: References (for maintainers)

- Project README (features, supported env, command list, security flags, testing & release flow)
- Command docs under `docs/commands/`
- Testing and troubleshooting docs under `docs/`
- Cisco IOS‑XE 17.12 programmability/RESTCONF docs (for path semantics)
- `cisco-ios-xe-wireless-go` library (Go SDK for C9800)
