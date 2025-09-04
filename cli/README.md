# 🖥️ Entiqon CLI

## Purpose
The **Entiqon CLI** is a lightweight developer & DevOps toolkit designed to **streamline everyday workflows** inside the Entiqon ecosystem.  
It bridges the gap between **library development**, **release management**, and **runtime utilities**, offering consistent, scriptable commands to help contributors and maintainers manage the project efficiently.

Whereas **Entiqon libraries** (e.g., `db/builder`, `token/field`, `core/contracts`) provide compile-time tools for building queryable Go systems, the CLI provides **runtime developer ergonomics**—automation for Git, tests, coverage, tagging, release notes, and package lifecycle.

---

## Philosophy
- **Minimal dependencies** — plain Bash/POSIX tools (works on macOS/Linux, extensible for Windows WSL).
- **Self-documented** — every script provides `-h`/`--help` output.
- **Composable** — individual commands solve one thing well, can be chained.
- **Versioned** — each CLI script is tied to Entiqon’s semver release cycle.
- **Safe** — designed to *never lose work* (e.g., stash before rebases, confirmations for destructive ops).

---

## Tooling Overview

### Git & Release Automation (`bin/`)
- **`gcpr`** – Create GitHub PRs quickly (auto-fills title, branch).
- **`gce`** – Extract commits between tags (`-s`, `-e` for ranges).
- **`gcr`** – Create GitHub releases with changelogs & notes.
- **`gct`** – Automated tagging (`--title`, `--date`, `--notes`, `--sign`).
- **`gsux`** – Git Stash Utility Extended (stash/apply/pop/drop/clear/list).
- **`gcch`** – Cherry-pick helpers for backports.
- **`ddc`** – Deploy Docker containers with standard flags.

### Testing & Coverage
- **`gotestx`** – Extended test runner: coverage, HTML reports, filters, CI mode.
- **`run-tests.sh`** – Runs all packages with coverage (`go test ./... -cover`).
- **`open-coverage.sh`** – Opens `coverage.html` after generation.
- **CI Integration** – Used in GitHub Actions to enforce thresholds and upload to Codecov.

### Documentation
- **Markdown helpers** – regenerate `README.md`, update `CHANGELOG.md`.
- **Release notes** – auto-generate from commits with semantic prefixes (`feat:`, `fix:`, `docs:`, etc.).

---

## Example Workflow

A typical **release cycle** with Entiqon CLI:

\`\`\`bash
# Run all tests with coverage
gotestx --cover --open

# Stage changes and stash WIP if needed
gsux stash -m "WIP: refactor token.Field validation" -u -v

# Generate changelog entries since last release
gce -s v1.13.0

# Tag and sign a new release
gct -t v1.14.0 --title "Token Enhancements" --notes "Adds ResolveExpressionType and ValidateWildcard" --sign

# Push release to GitHub
gcr v1.14.0
\`\`\`

---

## Roadmap
- **Global installer** (`entiqon install cli`) instead of per-repo `bin/`.
- **Unified entrypoint** (`entiqon <command>`) wrapping all scripts.
- **Go-based CLI rewrite** – current scripts are Bash; migration to Go planned for portability, testability, and richer UX.
- **Plugin architecture** – let projects extend CLI with their own subcommands.
- **Improved test harness** – Bats/shunit2 suites for CLI validation.

---

✅ In short: **Entiqon CLI = developer efficiency + project discipline**.  
It codifies the workflows we already practice (TDD, semantic commits, 100% coverage, structured releases) into **repeatable, versioned, safe automation.**
