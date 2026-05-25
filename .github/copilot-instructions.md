# Copilot instructions for Watchtower

## Build, test, and lint

This repository is a Go 1.20 module (`go.mod`) with the main binary at the repo root.

```bash
go build
go test ./... -v
go test -v -coverprofile coverage.out -covermode atomic ./...
staticcheck ./...
./build.sh
```

Use `go build` for a quick local binary build. Use `./build.sh` when you want the binary stamped with `internal/meta.Version` from `git describe --tags`, which matches the release/dev workflow behavior more closely.

To run one test, prefer package-scoped `go test`:

```bash
go test ./internal/flags -run TestGetSecretsFromFilesWithFile -v
go test ./internal/actions -run TestActions -v
```

Many packages use Ginkgo v1 suites behind `go test`, so package-level suite entrypoints such as `TestActions`, `TestRegistry`, or `TestContainer` are common.

## High-level architecture

- `main.go` only initializes logging and calls `cmd.Execute()`.
- `cmd/root.go` is the runtime entrypoint. It wires Cobra commands, registers flags from `internal/flags`, reads env/file-backed secrets in `PreRun`, constructs the Docker client and notifier, then chooses between run-once, scheduled polling, and token-protected HTTP API modes.
- `internal/flags` is the source of truth for CLI/env behavior. It defines the flag surface, maps most options to `WATCHTOWER_*` env vars, enforces `--interval` vs `--schedule`, and implements helper aliases such as `--porcelain`.
- `internal/actions.Update` is the core update loop. It lists filtered containers, checks staleness through the Docker/registry layer, sorts containers by dependency order, marks implicit restarts, stops/restarts containers, runs lifecycle hooks, and returns a session report.
- `pkg/container` wraps the Docker SDK and contains the container recreation logic. `client.go` handles list/inspect/stop/start/pull behavior; `container.go` interprets labels, derives dependencies, and reconstructs create config by diffing container config against image defaults.
- `pkg/filters` builds the container-selection pipeline from explicit names, disabled names, enable labels, scopes, and image filters.
- `pkg/api`, `pkg/api/update`, `pkg/api/metrics`, and `pkg/metrics` provide the optional HTTP surfaces. `/v1/update` and `/v1/metrics` share the same bearer-token gate, and scheduled runs plus HTTP-triggered runs are serialized through a shared lock in `cmd/root.go`.
- `pkg/notifications` is the notification abstraction. It builds Shoutrrr-based notifiers plus legacy adapters, and the stable machine-readable `--porcelain v1` output is implemented as a notification-template alias rather than a separate reporting subsystem.
- `pkg/session` owns the update report model that tracks scanned, skipped, updated, and failed containers for notifications and porcelain output.

## Key conventions

- The project README states the project is no longer maintained. Favor small, behavior-preserving changes over broad refactors or ecosystem upgrades unless the task explicitly requires them.
- Watchtower behavior is usually driven by both global flags and per-container labels. The important repo-specific merge rule is in `pkg/container`: global `--monitor-only` / `--no-pull` values win by default, and `--label-take-precedence` flips that behavior for those booleans.
- Multi-container restart order matters. The real dependency graph is not just Docker links: `pkg/container.Container.Links()` also treats `com.centurylinklabs.watchtower.depends-on` and `network_mode: container:<name>` as restart dependencies, and `pkg/sorter` uses that graph to order stops/starts.
- Multiple watchtower instances are guarded intentionally. `internal/actions.CheckForMultipleWatchtowerInstances` and the scope filter expect scope to be applied consistently to watched containers and the watchtower container itself; the special scope value `none` means unscoped containers.
- Secrets are allowed to come from files, not only literals. `internal/flags.GetSecretsFromFiles` replaces selected flag values with file contents before the rest of startup runs, including `notification-url` and `http-api-token`.
- `--porcelain v1` is an alias that mutates notification flags (`logger://`, report mode, stdout logging, template selection). If you change notification/reporting behavior, check porcelain expectations too.
- Logging is centralized through logrus. `internal/flags.SetupLogging` controls formatter selection, while `notifications.NewNotifier(...).AddLogHook()` attaches notification side effects to log output.
- Test style is mixed: many packages use Ginkgo suite entrypoints under `go test`, while others use plain `testing` functions. Match the existing test style in the package you touch instead of introducing a new one.
