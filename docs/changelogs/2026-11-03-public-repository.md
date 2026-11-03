# chore: make the template public

Tracking: alva-ai/mono-meta#1115 (decision G15: keep this template and fix
its sync).

## 1. Current state

- The template lived in a private repository. Projects generated from it
  run `cruft check` and `cruft update` in their own CI with their own
  token, which cannot read a private repository in this organization, so
  the sync never ran.
- The template's CI fetched dependencies with a personal access token
  (`CICD_DOCKER_BUILD_PAT`); every dependency it fetches is public
  (`ggsrc/gglib`, `Stumble/wpgx`, `Stumble/dcache`, `Stumble/sqlc`).
- `remove-artifacts.yaml` called a reusable workflow in the private
  `ggsrc/cicd` repository. A public repository cannot call it.

## 2. Behavior after this change

- B1 `ggsrc/ggsrv-layout` is public and starts from one commit holding the
  contents of the former repository's `main` (63e6928). The earlier history
  is not published; it stays in the private `ggsrc/ggsrv-layout-archive`.
- B2 Generated projects can read the template without any token.
- B3 The template's CI generates, builds, tests and builds the image of a
  project without secrets.

## 3. Decisions

- D1 A new repository instead of rewriting the old one: the old commits
  stay reachable through the pull request refs of the former repository,
  which cannot be deleted.
- D2 No artifact clean-up workflow: the CI artifacts already expire after
  3 days.

## 4. Changes

- `.github/workflows/black.yml`: no `GOPRIVATE`, no token in Git or the
  image build.
- `.github/workflows/remove-artifacts.yaml`: removed.

## 5. Migration for generated projects

The commit recorded in a project's `.cruft.json` belonged to the former
repository and does not exist here, so `cruft update` fails until the
project is linked again: bring the project up to date with the template by
hand (see `2026-11-03-gglib-v0.1.0.md`), then set `"commit"` in its
`.cruft.json` to the template commit it now matches. Later template changes
arrive through `cruft update` again.

## 6. Verification

- The tree of the import commit equals the tree of the former `main`
  (`451bde8`).
- `gitleaks dir .` on the imported tree: no leaks.
- CI of this change: see the pull request.

## 7. Outcome

- [ ] CI green; merged; repository public.
