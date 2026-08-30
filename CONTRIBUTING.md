# Contributing

Thank you for considering a contribution.

## Commands

The following `make` commands are available for development and testing:

| Command                     | Description                                            |
| :-------------------------- | :----------------------------------------------------- |
| `make help`                 | Display available targets and requirements             |
| `make build`                | Build the binary to `./tmp/wnc`                        |
| `make lint`                 | Verify the lint config, run golangci-lint, tidy go.mod |
| `make test-unit`            | Run unit tests with coverage using gotestsum           |
| `make test-unit-coverage`   | Generate HTML coverage report                          |
| `make snapshot`             | Build a GoReleaser snapshot                            |
| `make clean`                | Remove build artifacts and backup files                |
| `make image`                | Build Docker image                                     |
| `make pre-commit-install`   | Install the pre-commit hooks                           |
| `make pre-commit-test`      | Run every hook across the tree                         |
| `make pre-commit-uninstall` | Remove the pre-commit hooks                            |

`make test-unit` clears the `WNC_*` environment variables before running, because urfave reads them at flag-parse time and a developer's own shell would otherwise decide what the command-tree tests see. Add any new variable the CLI reads to that list.

Markdown style is enforced by the `markdownlint-cli2` hook that `make pre-commit-install` wires in, and again in CI. Links are checked in CI only, because that run reaches third-party hosts. Run `lychee .` to reproduce a link failure locally.

The hook path is the shared git common directory, so `make pre-commit-install` also arms every other worktree and the `main` checkout. It passes `--allow-missing-config` for that reason.

## Testing

See [docs/TESTING.md](docs/TESTING.md) for how the suite is arranged, what the column invariants assert, and how the TLS test harness stands in for a controller.

## Build

The repository includes a ready to use `Dockerfile`. To build a new Docker image:

```bash
make image
```

This cross-compiles a Linux binary into `./tmp/image`, then builds from that directory because the `Dockerfile` expects the binary at the context root. The image is tagged `$USER/wnc`. Released images are pushed to `ghcr.io/umatare5/wnc` by GoReleaser instead.

## Release

To release a new version, follow these steps:

1. Add the `## [vX.Y.Z]` section to `CHANGELOG.md` above the previous release, matching the version in the `VERSION` file, and add that version's release link at the foot of the file.
2. Update the version in the `VERSION` file.
3. Refresh the coverage badge — `make test-unit` then `octocov badge coverage --config .octocov.yml > docs/assets/coverage.svg`. Nothing automates it — the reusable coverage workflow enforces the floor but writes no badge.
4. Submit a pull request with all three files.

Merging that pull request is the whole release. A push to `main` touching `VERSION` runs the [release workflow](https://github.com/umatare5/wnc/actions/workflows/go-release.yml), which tags the commit and publishes the release in the same run. The workflow has no manual trigger, so there is no step to perform by hand.

## Pull requests

1. [Fork](https://github.com/umatare5/wnc/fork) the repository
2. Create a feature branch
3. Commit your changes, following the surrounding style and signing off with `Signed-off-by:`
4. Add tests — CI enforces a coverage floor — and update the documentation alongside the code
5. Run `make lint` and `make test-unit`
6. Rebase your local changes against the `main` branch
7. Create a new Pull Request
