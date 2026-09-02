> [!NOTE]
> ## This is a maintained fork
>
> **`phenixrizen/zen-go` is the Go binding for [`phenixrizen/zen`](https://github.com/phenixrizen/zen), a maintained fork of `gorules/zen`, maintained by Phenix Rizen (Nathan Rockhold).**
>
> The bundled static libraries are built from that fork, so this binding exposes features upstream
> does not have:
>
> | addition | what it does |
> | --- | --- |
> | `databaseNode` | Look reference data up from inside a decision, with a host-supplied handler. |
> | SQLite handler for `databaseNode` | Register it once at startup — `zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(load), SqliteConfig: `{"root":"/catalog"}`})`, with a `Loader` *callback* (it cannot be combined with a `FilesystemLoader`) — and every query runs inside the vendored library with no crossing back over the FFI. The driver is `rusqlite` with the bundled SQLite amalgamation, so the archive does contain C; see the fork's `docs/driver-choice.md` for why. |
> | Decision-level `$params` | Static parameters supplied per decision, reachable from switch, expression, and function nodes. |
> | `TZ` is honoured | Date resolution respects `TZ` rather than only `/etc/localtime`. |
> | Exact fractional numbers | Fixes silent truncation of every non-integer value. |
>
> ```
> go get github.com/phenixrizen/zen-go/v2
> ```
>
> Not affiliated with or endorsed by GoRules. MIT licensed, same as upstream, with the original
> copyright retained in [LICENSE](LICENSE).

# Go Rules Engine

**Business logic humans can read and machines can run.** One copy of your rules: the owner reads it, every system runs it.

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/phenixrizen/zen-go/v2.svg)](https://pkg.go.dev/github.com/phenixrizen/zen-go/v2)

ZEN Engine is a cross-platform, open-source Business Rules Engine (BRE) written in **Rust** with native **Go** bindings, alongside Node.js, Python, Java, Kotlin and .NET. Decisions evaluate in microseconds, run identically on every platform, and are stored as portable JSON. Loading the JSON is up to you: file system, database or service call.

## Rules that read like sentences

Conditions are written the way the business says them, in the ZEN Expression Language. The developer view is one toggle away, and the two can never drift apart: there is only one source of truth, and this engine runs it.

## Rules as graphs, or as documents

Model a decision on a visual canvas of decision tables, switches, expressions, functions and reusable sub-decisions. Or write it as a policy document with prose, typed data models and tables. Both compile to the same engine and return the same answers.

A JDM document is either a **graph** (decision tables, switches, expressions, functions and reusable sub-decisions) or a **policy** (prose, typed data models and tables). Both compile to the same engine and return the same answers.

## What's new in 2.0

Version 2.0 is the first stable release of the new engine line:

- **Policy documents**: model decisions as readable documents with typed data models, expressions, decision tables, match blocks and assertions. Policies compile to the same engine as graphs and return the same answers.
- **Workspace analysis**: static type checking across policies and graphs. Type flow, exhaustiveness checking, write-conflict detection and precise diagnostics, all available before anything runs.
- **Per-column collect**: decision table output columns can collect across all matching rows (`tags[]`) while the rest of the table stays first-match.
- **Pre-compiled engine**: decisions are parsed and compiled once at load; evaluation is allocation-light and repeat-safe.
- **Hardened runtime**: out-of-range numbers, arithmetic overflow and malformed inputs return errors or nulls instead of crashing the process.
- **Unified bindings**: configurable loaders, batch evaluation and consistent error envelopes across Node.js, Python, Go and FFI consumers.

## Installation

```bash
go get github.com/phenixrizen/zen-go/v2
```

## Quickstart

```go
package main

import (
	"fmt"
	"os"
	"path"

	zen "github.com/phenixrizen/zen-go/v2"
)

func readTestFile(key string) ([]byte, error) {
	filePath := path.Join("test-data", key)
	return os.ReadFile(filePath)
}

func main() {
	engine := zen.NewEngine(zen.EngineConfig{Loader: zen.Loader(readTestFile)})
	defer engine.Dispose() // Call to avoid leaks

	output, err := engine.Evaluate("rule.json", map[string]any{})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(output)
}
```

### Loader Configurations

`EngineConfig.Loader` accepts either a loader callback wrapped in `zen.Loader` (as above) or a loader configuration of a
known type. With a configuration, decisions are pre-loaded and pre-compiled at engine creation for
faster evaluations.

```go
engine := zen.NewEngine(zen.EngineConfig{Loader: zen.FilesystemLoader{Path: "test-data"}})

engine := zen.NewEngine(zen.EngineConfig{Loader: zen.StaticLoader{
	Content: map[string]json.RawMessage{"rule.json": ruleJson},
}})

engine := zen.NewEngine(zen.EngineConfig{Loader: zen.ZipLoader{Bytes: zipBytes}})
```

If the loader configuration is invalid (e.g. corrupted zip bytes), the error is returned by the
first call to `Evaluate`, `GetDecision` or `CreateDecision`.

The same callback pattern works for loading from a REST API, S3, a database, or anywhere else.

## Other platforms

* **Node.js** — [npm](https://www.npmjs.com/package/@phenixrizen/zen-engine)
* **Python** — [PyPI](https://pypi.org/project/phenixrizen-zen-engine/)
* **Go** — this repository
* **Java / Kotlin** — [source](https://github.com/phenixrizen/zen/tree/master/bindings/uniffi)
* **.NET** — [NuGet](https://www.nuget.org/packages/PhenixRizen.ZenEngine)
* **Rust (core)** — [phenixrizen/zen](https://github.com/phenixrizen/zen) | [crates.io](https://crates.io/crates/phenixrizen-zen-engine)


## Support matrix

| Arch            | Go                 |
|:----------------|:-------------------|
| linux-x64-gnu   | :heavy_check_mark: |
| linux-arm64-gnu | :heavy_check_mark: |
| darwin-x64      | :heavy_check_mark: |
| darwin-arm64    | :heavy_check_mark: |
| win32-x64-msvc  | :heavy_check_mark: |

We do not support linux-musl currently.

## Contribution

**Contributions are welcome here.** This fork exists partly because upstream cannot take them.

Note that the Go code in this repository is a thin binding: most behaviour lives in
[`phenixrizen/zen`](https://github.com/phenixrizen/zen), and the `deps/` static libraries are
built from it. A change to evaluation belongs there; a change to the Go surface belongs here.

```bash
make fmt_check
make test
```

## License

[MIT License](https://opensource.org/licenses/MIT)
