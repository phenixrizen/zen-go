# Contributing

Contributions are welcome. This is the Go binding for
[`phenixrizen/zen`](https://github.com/phenixrizen/zen), a maintained fork of `gorules/zen`.

Maintained by Phenix Rizen (Nathan Rockhold).

## Where does a change belong?

This repository is a thin cgo binding. Almost all behaviour lives in the engine.

| change | repository |
| --- | --- |
| Evaluation, expressions, node kinds, the database node | [`phenixrizen/zen`](https://github.com/phenixrizen/zen) |
| The Go API surface, marshalling, handle lifetimes | here |
| The C FFI surface (`zen_engine.h`) | `phenixrizen/zen` under `bindings/c`, then regenerated here |

`deps/*/libzen_ffi.a` and `zen_engine.h` are **build artifacts**, not source. They are produced by
the `Go` workflow in `phenixrizen/zen` and pushed here as a "chore: update deps" pull request.
Do not hand-edit them; change the engine and let the pipeline rebuild.

## Development setup

Requires Go and a C toolchain for cgo. The static libraries are prebuilt and committed, so there
is no Rust toolchain needed to work on the Go side.

```bash
git clone https://github.com/phenixrizen/zen-go.git
cd zen-go
make test
```

## The test gate

```bash
make fmt_check
make test
make memory_test   # runs under the memory_test build tag
```

## Building the libraries yourself

Only needed when you are changing the engine and want to try it before the pipeline runs. From a
checkout of `phenixrizen/zen`:

```bash
cargo build -p zen-ffi --release --all-features --target x86_64-unknown-linux-gnu
cp target/x86_64-unknown-linux-gnu/release/libzen_ffi.a <zen-go>/deps/linux_amd64/
cp bindings/c/zen_engine.h <zen-go>/zen_engine.h
```

`--all-features` matters: it links the SQLite database handler in. Without it,
`zen_engine_new_golang_with_sqlite` is not compiled and setting `EngineConfig.SqliteConfig`
fails to link.

The other platforms are cross-built in CI; a local build only replaces your own.

## Commit messages

[Conventional Commits](https://www.conventionalcommits.org/).

## Pull requests

- Branch off `master`.
- A bug fix needs a regression test, and you should confirm the test **fails without the fix**.
- If your change needs an engine change, link the `phenixrizen/zen` pull request.
