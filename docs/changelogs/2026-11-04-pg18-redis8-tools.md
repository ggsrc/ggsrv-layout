# feat: PostgreSQL 18, Redis 8, mockery v3, sqlc from a pinned release

Tracking: alva-ai/mono-meta#1115. Yumin, 2026-11-04: upgrade PostgreSQL and
Redis; move mockery to v3 if it is worth it; developers should get the right
sqlc without building Stumble/sqlc themselves (download, check the SHA).
Follows `2026-11-04-toolchain-refresh.md` (its decision D3 kept PostgreSQL
14.5, Redis 7 and mockery v2; this change replaces that decision).

Depends on Stumble/wpgx v0.5.1 (`testsuite.WithPostgresImage`, Stumble/wpgx#23)
and on Stumble/sqlc v2.4.1, the first release with binaries and `checksums.txt`
(Stumble/sqlc#21–#23). Both were prepared here as patches and merged and
released from Yumin's machine; v2.4.1 also moves sqlc to Go 1.26.10 and fixes
14 reachable vulnerabilities.

## 1. Current state

- PostgreSQL 14.5 everywhere: `docker-compose.yml`, and the tests through
  wpgx's `NewWPgxTestSuiteTcDefault`, which hard-codes `postgres:14.5`.
  Redis `7-alpine` (compose, tests); the template CI used an unpinned `redis`
  service. Latest: PostgreSQL 18.6, Redis 8.10.3.
- Production (infra-gitops): most CNPG clusters run PostgreSQL 14, jagent 17,
  alcron 15, timescale 18; Redis 7.2.1. PostgreSQL 14 reaches end of life on
  2026-11-12.
- **Generated projects' CI tests could not reach Redis.** With
  `GITHUB_ACTIONS=true` the tests connected to `127.0.0.1:6379` instead of
  starting a container, but the generated `check.yml` has no Redis service
  (reproduced: `connection refused`, `make test` exit 2).
- mockery: installed (v2) but unused; the README advertised `make mock`,
  which did not exist. v2 still receives fixes (v2.53.8, 2026-10-16), new work
  happens in v3 (v3.9.0, 2026-10-28), and v2's flag-based mode is deprecated.
- sqlc: no release binaries; the template CI cloned the fork and ran
  `go install`, which reports `v1.31.1` (the Makefile sets the real version
  through `-ldflags`), so the committed code said `sqlc v1.31.1`.

## 2. Behavior after this change

- B1 Local and test PostgreSQL 18.6 (`postgres:18.6`, Debian, like the
  production images) and Redis 8.10.3 (`redis:8.10.3-alpine`). The test
  images are constants in `internal/svc/svc_test.go`; compose uses the same.
- B2 Compose mounts the PostgreSQL volume at `/var/lib/postgresql` (18+
  layout). A volume from an older major makes the container exit with the
  image's explanation; the setup guide says how to recreate it.
- B3 Tests use `REDIS_HOST` when it is set and otherwise start Redis with
  Testcontainers, in CI as locally. The template CI has no Redis service.
- B4 mockery v3.9.0. `api/go/.mockery.yml` generates testify mocks of the
  service's `Client` into `api/go/client/<service>/mocks`; `make mock`
  (after `make proto`) runs it and tidies `api/go`.
- B5 `scripts/install-sqlc.sh` installs the Stumble/sqlc v2.4.1 release
  binary after checking the archive against SHA-256 values pinned in the
  script; `make repos` stops unless `sqlc version` is that release. The
  committed sqlc code is what it generates (`sqlc v2.4.1-wicked-fork`).
- B6 `make tools` installs buf/protoc/plugins, sqlc, golangci-lint, mockery
  and sql-formatter; the README quick start starts with it.
- B7 Template CI: sqlc through `install-sqlc.sh`; `make repos` must
  reproduce the committed code; `make mock` and `go vet` on `api/go`.
- B8 Dependabot also watches `docker-compose.yml`.

## 3. Decisions

- D1 PostgreSQL 18 for new services, as asked. Existing services keep their
  own test image: wpgx's default stays `postgres:14.5`, and the version is an
  explicit option (`WithPostgresImage`). A service created from the template
  should get a PostgreSQL 18 database (CNPG image
  `ghcr.io/cloudnative-pg/postgresql:18`), otherwise tests accept SQL that
  production rejects.
- D2 mockery v3: the template has no mocks to migrate, so v3 costs nothing
  and avoids starting new services on the major that only gets fixes.
- D3 Mocks of the service's own client, published in `api/go`, because
  callers need them; the gRPC interface itself is in a generated file, which
  mockery v3 skips by default.
- D4 sqlc archive checksums pinned in the project, not read from the release,
  so a replaced release asset is refused.

## 4. Changes

- `docker-compose.yml`, `internal/svc/svc_test.go`, `go.mod`/`go.sum`
  (wpgx v0.5.1): B1–B3.
- `.github/workflows/black.yml`: B3, B7.
- `scripts/install-dev-tools.sh` (`ONLY=`, mockery v3), `api/go/.mockery.yml`,
  `Makefile` (`mock`, `tools`, sqlc version check), `SERVICE_TEST_GUIDE.md`:
  B4, B6.
- `scripts/install-sqlc.sh`, `internal/repos/users/*.go`: B5.
- `README.md`, `PROJECT_SETUP_GUIDE.md`: B2, B5, B6.
- `.github/dependabot.yml` (generated): B8.

## 5. Verification

Release assets: `sha256sum -c checksums.txt` OK for the five archives;
`gh attestation verify -R Stumble/sqlc` passes for each (release workflow on
`main`). The release run built, smoke-tested and installed on all platforms
(run 41799904816).

Generated project from this branch, wpgx v0.5.1 and sqlc v2.4.1 from GitHub:

- `go mod tidy`: no change. `install-sqlc.sh` downloads, checks and installs
  `v2.4.1-wicked-fork`. `make repos` output equals the committed code.
- `make mock`, `go vet` on `api/go`, `make build`, golangci-lint v2.15.0
  (`0 issues`, both modules).
- `make test` with `GITHUB_ACTIONS=true` and no Redis service: 19 tests pass.

Earlier, with the local wpgx branch and a locally built archive behind a fake
release URL:

- `make repos` without sqlc: stops with the install hint.
  `install-sqlc.sh`: installs `v2.4.0-wicked-fork`, a second run is a no-op.
  `make repos` output equals the committed code apart from the old
  `v1.31.1` header.
- `make proto`, `make mock` (mockery v3.9.0 writes `MockClient`; `go vet` on
  `api/go`), `make build`, golangci-lint v2.15.0 (`0 issues`, both modules).
- `make test` with `GITHUB_ACTIONS=true` and no Redis service: 19 tests pass
  on `postgres:18.6` and `redis:8.10.3-alpine` (before: connection refused).
- Compose: PostgreSQL `18.6`, data directory `/var/lib/postgresql/18/docker`,
  Redis `8.10.3`, both healthy. A PostgreSQL 14.5 volume: the 18.6 container
  exits with the 18+ layout message.
- mockery: v3 refuses a config without packages, so `.mockery.yml` lists the
  client; v2.53.8 still runs on Go 1.27 but warns that its flags are going.

## 6. Human decisions

- Yumin, 2026-11-04: "pg和redis的版本都升级吧。mockery看一下有没有必要升级到v3。
  如果有必要的话，还是可以升级的。" and the sqlc request above.

## 7. Outcome

- [x] Stumble/wpgx v0.5.1 released; `go.mod` bumped.
- [x] Stumble/sqlc v2.4.1 released with binaries; SHA-256 values pinned in
  `install-sqlc.sh`.
- [ ] CI green; merged.

## 8. Remaining work

- New services: request a PostgreSQL 18 database. Existing PostgreSQL 14
  clusters reach end of life on 2026-11-12 (separate decision).
