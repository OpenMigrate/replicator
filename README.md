# replicator
Controller component of the OpenMigrate platform — receives replicated disk data from agents, processes block-level changes, applies compression, and securely uploads to target cloud storage (e.g., S3). Enables scalable, centralized orchestration of migrations across environments.

### Pre-commit, Commitizen, and Linting

- Install pre-commit hooks:
  - pipx: `pipx install pre-commit` (or `pip install pre-commit`)
  - Install hooks: `pre-commit install --hook-type pre-commit --hook-type commit-msg`
- Branch naming enforced via pre-commit: `<type>/(OM|om)-<ticket>_<short-desc>`
  - Allowed types: `feat, feature, fix, bugfix, chore, docs, refactor, test, ci, build, perf, wip, break, revert, hotfix, security`
- Commit messages validated by Commitizen (Conventional Commits) on `commit-msg`.
- Go formatting and vet run on commit; static analysis via `golangci-lint`.

Configs:
- `.pre-commit-config.yaml` — hooks configuration
- `.golangci.yml` — linter configuration
- `.cz.toml` — commitizen configuration
