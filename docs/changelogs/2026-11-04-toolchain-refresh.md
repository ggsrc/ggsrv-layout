# chore: latest stable toolchain, pinned versions, clean pre-commit

Tracking: alva-ai/mono-meta#1115 (decision G15: keep this template and fix
its sync). Yumin, 2026-11-04: bring ggsrv-layout and its dependencies to
the latest stable versions, following best practice.

## 1. Current state

- Go 1.27.3, gglib v0.1.0, wpgx and dcache v0.5.0, grpc, protobuf, pgx,
  go-redis, zerolog, testify and testcontainers-go were already the latest
  stable versions. `testcontainers-go/modules/redis` was v0.39.0 while the
  core module was v0.44.0.
- **Generated projects' lint failed.** `lint.yml` called
  `public-ci/golang-lint-v2.yml@main`, which runs golangci-lint v2.5. Its
  release binary is built with Go 1.25 and refuses Go 1.27 modules:
  `the Go language version (go1.25) used to build golangci-lint is lower
  than the targeted Go version (1.27.3)`.
- Tool versions were old or unpinned: buf 1.56.0 in
  `install-buf-deps.sh` but 1.34.0 in the generated CI (through the
  archived `bufbuild/buf-setup-action`), protoc 31.1, protoc-gen-go
  v1.36.6, protoc-gen-go-grpc v1.5.1, golangci-lint v2.5 in
  `install-dev-tools.sh`, mockery / sql-formatter / gotestsum `latest`, the
  template CI's sqlc fork at whatever `HEAD` was.
- Actions ran at old majors by tag (checkout v4, setup-go v5, artifacts v4,
  buildx v3, build-push v5; `gha-cookiecutter@main`), the public-ci
  workflows at `@main`. The template CI had `actions: write` and uploaded a
  non-existent `myproject_empty` directory; upload-artifact v4.4+ also
  dropped the generated project's hidden files.
- A freshly generated project failed its own pre-commit: missing final
  newlines, trailing whitespace and unformatted JSON (also with the old
  hook versions); the stage names `commit` / `push` are deprecated. The
  template root failed `git-check` (mixed indent in `cookiecutter.json`)
  and black.
- The generated `dependabot.yml` covered Go modules only, though the README
  says GitHub Actions and Docker images are updated too; it also declared an
  unused registry without a password and named the organisation as
  assignee and reviewer.
- `scripts/install-dependencies.sh` installs `github.com/stumble/v8runner`,
  which is not on the public module proxy, and nothing calls it.
- Runtime image: `distroless/cc-debian12`.
- Template CI noise: the Makefile ran `go list ./...` while parsing
  (`TEST_DIRS :=`), before `make proto`, so every `make proto` printed a
  `does not contain package` error (shown as a failure annotation); the Go
  cache was never restored (no `go.sum` at the workspace root).

## 2. Behavior after this change

- B1 Generated lint runs `public-ci/golang-lint-v2_8.yml` pinned to the
  `v1` commit with `lint-version: v2.15.0` (built with Go 1.27.0;
  `0 issues` on a generated project). `cruft-update.yml` is pinned the same
  way, as the public-ci README recommends.
- B2 Every action is pinned by commit SHA at its latest release; Dependabot
  (template root and generated projects) proposes updates weekly with a
  7-day cooldown, grouped.
- B3 Tool versions: buf 1.75.0 (one value for the script and CI, CI via
  `bufbuild/buf-action`), protoc 36.2, protoc-gen-go v1.36.12,
  protoc-gen-go-grpc v1.6.2, golangci-lint v2.15.0, mockery v2.53.8,
  sql-formatter 15.9.0, gotestsum v1.13.0, sqlc fork v2.4.0 in the template
  CI. The committed sqlc code is what v2.4.0 generates.
- B4 `testcontainers-go/modules/redis` v0.44.0.
- B5 A generated project passes all of its pre-commit hooks
  (pre-commit-hooks v6.0.0, jumanjihouse 3.0.0); the template root passes
  its own (plus black 26.10.1).
- B6 Generated Dependabot covers Go modules (root and `api/go`), GitHub
  Actions and the Docker images.
- B7 Runtime image `distroless/cc-debian13`.
- B8 Generated CI configures Git with the private-module token only when
  the secret is set, and has read-only `contents` permission.
- B9 `TEST_DIRS` is evaluated only by the targets that use it; the template
  CI restores the Go cache from the generated project's `go.sum`.

## 3. Decisions

- D1 Use the existing `golang-lint-v2_8.yml` (has `lint-version`) instead
  of changing `golang-lint-v2.yml`: the v2.5 default stays as it is for its
  other callers.
- D2 Direct dependencies at their latest versions; indirect versions stay
  what `go mod tidy` selects from them (no `go get -u ./...`).
- D3 Unchanged on purpose: local Postgres 14.5 / Redis 7 images (a major
  upgrade breaks existing local volumes and should follow production),
  mockery v2 (v3 changes its configuration), `cookiecutter.go_version`
  (already the latest).
- D4 Markdown hard line breaks become list items or a trailing `\`,
  because `git-check` rejects trailing whitespace.

## 4. Changes

- `.github/workflows/black.yml`, generated `check.yml`, `lint.yml`,
  `cruft-update.yml`: B1, B2, B3, B8.
- `.github/dependabot.yml` (new), generated `dependabot.yml`: B2, B6.
- `scripts/install-buf-deps.sh`, `scripts/install-dev-tools.sh`,
  `Makefile` (`GOTESTSUM_VERSION`), `api/proto/buf.gen.yaml` comments: B3.
- `Makefile` (`TEST_DIRS =`), `black.yml` (`cache-dependency-path`): B9.
- `go.mod`, `go.sum`, `internal/repos/users/*.go`: B3, B4.
- `.pre-commit-config.yaml` (both), `cookiecutter.json`,
  `hooks/pre_gen_project.py` (black), whitespace / final newline / JSON
  formatting in the files the hooks flagged: B5.
- `docker/Dockerfile`: B7. `docker-compose.yml`: obsolete `version` key
  removed.
- `scripts/install-dependencies.sh`: removed.

## 5. Verification

All against a project generated from this branch with the CI's answers:

- `golang:1.27.3`: `scripts/install-buf-deps.sh` installs buf 1.75.0,
  libprotoc 36.2, protoc-gen-go v1.36.12, protoc-gen-go-grpc 1.6.2; sqlc
  v2.4.0 from source; `make proto`, `make repos` (output identical to the
  committed sqlc code), `go mod tidy` (no change apart from line order),
  no direct dependency with a newer version, `make build`, `go vet`,
  golangci-lint v2.15.0 on the root and `api/go` modules: `0 issues`;
  `make test` with Testcontainers: 19 tests pass.
- `docker build -f docker/Dockerfile .`: builds on `cc-debian13`; the
  server starts and stops on the missing Postgres settings.
- pre-commit 4.7.0: every hook passes on the generated project and on the
  template root. With the old hook versions the generated project failed
  the same three hooks on the same 23 files.
- actionlint: no findings on the template and generated workflows.
- golangci-lint v2.5 (current) vs v2.15.0 on a Go 1.27.3 module: the first
  stops with the error in section 1, the second reports `0 issues`.
- CI of #3 (run 41775063091): `build` and `test1` pass; the Node 20 and
  buf token warnings of `main` are gone.

## 6. Human decisions

- Yumin, 2026-11-04: "其他的你按照最佳实践，推荐来搞就行".

## 7. Outcome

- [ ] CI green.
- [ ] Merged.

## 8. Remaining work

- Generated projects take this through `cruft update` after re-linking
  `.cruft.json` (`2026-11-03-public-repository.md`, section 5); the newer
  golangci-lint may report findings the old one did not.
- The Dockerfile still takes the private-module token as a build `ARG`
  (Docker warns `SecretsUsedInArgOrEnv`). A BuildKit secret mount is the
  better form, but every project's image build passes `GITHUB_ACCESS_TOKEN`
  as a build argument today, so it changes together with those pipelines.
