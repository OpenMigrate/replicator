# Contributing to Replicator

Thanks for your interest in contributing! This guide covers local setup, branch and commit conventions, development workflow, and where to start as a new contributor.

## Local Development

- Prerequisites
  - Go 1.23+ (see `go.mod`)
  - GitHub CLI (`gh`) for issues/PRs
  - Pre-commit, Commitizen, GolangCI-Lint

- macOS
  - Install tools:
    ```bash
    brew install go gh pre-commit golangci-lint
    pipx install commitizen || pip install commitizen
    ```

- Ubuntu/Debian
  - Install tools:
    ```bash
    sudo apt update && sudo apt install -y golang-go python3-pip
    pipx install pre-commit commitizen || pip install pre-commit commitizen
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ~/.local/bin v1.59.1
    export PATH="$HOME/.local/bin:$PATH"
    ```

- Setup hooks
  ```bash
  pre-commit install --hook-type pre-commit --hook-type commit-msg
  ```

- Run the server
  ```bash
  cp -n config.toml config.local.toml 2>/dev/null || true
  go run ./cmd/replicator --config config.toml
  ```

- Run tests
  ```bash
  go test ./...
  ```

## Branch Naming

- Pattern: `<type>/<issue>-<short-desc>` (we use GitHub issue numbers for project tracking)
  - `<type>` one of: `feat, feature, fix, add, chore, docs, refac, wip, hotfix, security`
  - `<issue>` numeric GitHub issue ID (e.g., `29`)
  - Examples:
    - `feat/29-precommit-setup`
    - `fix/102-buggy-disk-check`

Our pre-commit `no-commit-to-branch` hook blocks commits on non-conforming branches.

## Commit Messages

- We use Commitizen with a customized Conventional Commits subset.
- Format: `<type>(<scope>): <subject>`
  - `<type>` must be one of: `feat, feature, fix, add, chore, docs, refac, wip, hotfix, security`
  - Examples:
    - `feat(api): add app grouping endpoints`
    - `fix(storage): handle sqlite busy timeout`

- Recommended:
  ```bash
  cz commit
  ```
  or commit normally and the `commit-msg` hook will validate.

## Development Workflow

- Create or pick up an issue. For new features, propose design via an issue first.
- Create a branch following the naming convention.
- Keep PRs small and focused. Reference issues with `Closes #<issue>` in the PR description.
- Ensure `go fmt`, lint, and tests pass locally before opening the PR.

## Good First Issues

- Look for issues labeled `good first issue` or `help wanted`:
  - https://github.com/OpenMigrate/replicator/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22
  - https://github.com/OpenMigrate/replicator/issues?q=is%3Aissue+is%3Aopen+label%3A%22help+wanted%22
- Unsure where to start? Comment on an issue and we’ll guide you.

## Code Style & Linting

- `go fmt ./...` and `go vet ./...` run via pre-commit.
- `golangci-lint` runs as part of pre-commit; CI integration may run it on PRs as well.

## Security & Secrets

- `detect-secrets` scans for potential credentials.
- Never commit secrets. Use environment variables or secure stores.

Happy hacking!
