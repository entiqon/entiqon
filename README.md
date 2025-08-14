<h1 align="center">
  <img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_logo.png?raw=true" align="left" height="140" width="140"> Entiqon
</h1>
<p align="center">A structured, intelligent foundation for building queryable, entity-aware Go systems.</p>

<p align="center">
  <a href="https://pkg.go.dev/github.com/entiqon/entiqon"><img src="https://pkg.go.dev/badge/github.com/entiqon/entiqon.svg" alt="Go Reference" /></a>
  <a href="https://goreportcard.com/report/github.com/entiqon/entiqon"><img src="https://goreportcard.com/badge/github.com/entiqon/entiqon" alt="Go Report Card" /></a>
  <a href="https://github.com/entiqon/entiqon/actions/workflows/ci.yml"><img src="https://github.com/entiqon/entiqon/actions/workflows/ci.yml/badge.svg" alt="Build Status" /></a>
  <a href="https://codecov.io/gh/entiqon/entiqon"><img src="https://codecov.io/gh/entiqon/entiqon/graph/badge.svg?token=6t7ENLuwwt"/></a>
  <a href="https://github.com/entiqon/entiqon/releases"><img src="https://img.shields.io/github/v/release/entiqon/entiqon" alt="Latest Release" /></a>
  <a href="https://entiqon.github.io/entiqon/"><img src="https://img.shields.io/badge/docs-online-blue?logo=github" alt="Documentation" /></a>
  <a href="https://github.com/entiqon/entiqon/blob/main/LICENSE"><img src="https://img.shields.io/github/license/entiqon/entiqon" alt="License" /></a>
</p>

## 📦 Packages

- <a href="https://github.com/entiqon/entiqon/blob/main/common"><img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_sharicon.png?raw=true.png" align="left" height="24" width="24">
  Common</a>: Shared utilities and helper functions used across multiple modules. Installation:
  `go get github.com/entiqon/common`
- <a href="https://github.com/entiqon/entiqon/blob/main/db"><img src="https://github.com/entiqon/entiqon/blob/main/assets/entiqon_datacon.png?raw=true.png" align="left" height="24" width="24">
  Database</a>: Modular SQL query builder focused on database operations. Installation: `go get github.com/entiqon/db`

*Future modules such as `core`, `auth`, `http`, and others will be added following the modular architecture.*

---

## 📏 Principles & Best Practices

* 🧼 Clarity over brevity — use explicit method names
* 🚫 Deprecations are tested and clearly marked
* 🔐 Validate every path — no silent failures
* 🧩 Always quote identifiers through the dialect

---

## 🧩 Design Philosophy

* 📐 Chain → Validate → Compile
* 🧠 Tag errors with `StageToken`
* ⚙️ Compose with safe abstractions
* 📂 Group test methods visually

---

## 📦 Releases

- [v1.10.0](./releases/release-notes-v1.10.0.md)
- [v1.9.0](./releases/release-notes-v1.9.0.md)
- [CHANGELOG](./CHANGELOG.md)

---

## 📄 License

💡 Originally created by [Isidro Lopez](https://github.com/ialopezg)  
🏢 Maintained by the [Entiqon Organization](https://github.com/entiqon)

[MIT](./LICENSE) — © Isidro Lopez / Entiqon Project
