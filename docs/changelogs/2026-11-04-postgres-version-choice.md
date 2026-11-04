# feat: choose the PostgreSQL version when generating a project

Tracking: alva-ai/mono-meta#1115. Yumin, 2026-11-04, on the PostgreSQL 18
upgrade plan: "在代码层面把pg的版本还是变成可选的吧。最近可能暂时没有时间升级的。"
Follows `2026-11-04-pg18-redis8-tools.md`, which fixed new projects on 18.6.

## 1. Current state

- Generated projects test and run locally on `postgres:18.6`, but most
  production databases run PostgreSQL 14 (cnpg-llm, alfs, juicefs, data5,
  scrap), jagent 17, and only timescale/worldline 18; the production upgrade
  is postponed. A service tested on 18 can rely on SQL its 14 database
  rejects.
- PostgreSQL 18 images keep data under `/var/lib/postgresql/<major>/` and want
  the volume at `/var/lib/postgresql`; older images use
  `/var/lib/postgresql/data` and declare that path as a volume.

## 2. Behavior after this change

- B1 `postgres_version` in `cookiecutter.json`: `14.24` (default), `17.11`,
  `18.6`. It sets the compose image, the matching volume path, and the test
  container image (`defaultPostgresImage` in `internal/svc/svc_test.go`), and is
  recorded in `cookiecutter-config-file.yml`.
- B2 `POSTGRES_TEST_IMAGE` overrides the test image for one run.
- B3 The template CI generates and tests a project for `14.24` (checks `build`,
  `test1`, as the ruleset requires) and for `18.6` (`build (postgres 18.6)`,
  `test1 (postgres 18.6)`).
- B4 The setup guide covers a local volume left by another major version.

## 3. Decisions

- D1 Default 14.24: a new service's database most likely lives on cnpg-llm
  (PostgreSQL 14.22). Pick the major version of the database the service will
  use.
- D2 Choice at generation time, not a runtime switch, because the volume path
  depends on the major version.

## 4. Changes

- `cookiecutter.json`, `cookiecutter-config-file.yml`, `docker-compose.yml`,
  `internal/svc/svc_test.go`, `PROJECT_SETUP_GUIDE.md`: B1, B2, B4.
- `.github/workflows/black.yml`: B3.

## 5. Verification

- Rendered with the default and with `postgres_version=18.6`: compose image
  and volume path as in B1; `defaultPostgresImage` follows.
- Default (14.24) project, `golang:1.27.3`: sqlc v2.4.1 installed from the
  release, `make proto`, `make repos` (= committed code), `make mock`, build,
  golangci-lint v2.15.0 (0 issues), `make test` without a Redis service: 19
  tests pass.
- 18.6 project: `make test` passes (19 tests).
- Default project with `POSTGRES_TEST_IMAGE=postgres:17.11`: the suite starts
  `postgres:17.11` and passes.
- actionlint: clean.

## 6. Human decisions

- Yumin, 2026-11-04 (above).

## 7. Outcome

- [x] CI green (run 41908907180: `build`, `test1`, `build (postgres 18.6)`,
  `test1 (postgres 18.6)`); squash-merged (#5).

## 8. Remaining work

- wpgx's container start timeout is fixed at 60 s and includes creating the
  container; on a slow machine a first run can time out (seen once, 54 s pass
  on retry).
