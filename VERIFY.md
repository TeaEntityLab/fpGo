# VERIFY — fpGo

How an agent launches, checks, drives, and cleans up verification of this
library. Read this plus the relevant `features/*.md` file before touching the
code.

`source_commit: a8b8521-era` · `last_verified_at: 2026-09-23` · `verification_status: passed`

## Launch

```bash
cd <repo root>
go test -mod=mod ./...
```

**`-mod=mod` is required.** The checked-in `vendor/` tree is stale relative to
`go.mod` (kr/pretty, pmezard/go-difflib, gopkg.in/check.v1 are required but not
marked explicit in `vendor/modules.txt`), so bare `go test ./...` fails with
`inconsistent vendoring` before compiling anything. Do not "fix" this by
editing `vendor/` during verification — record it; `go mod vendor` is a
separate reviewed change.

## Doctor

```bash
go vet -mod=mod ./...
go test -mod=mod -run 'TestMaybeJust|TestLinkedListQueueOfferPoll' .   # smoke: ~1s
```

A healthy tree: `go vet` clean, smoke tests PASS, full suite PASS in ~5s
(1013 test functions across root + `network/` + `worker/`).

## Drive

- Per feature: run the `-run` regex in the feature file's Drive section.
- For behavior not covered by a test, write a scratch driver **outside** the
  repo (e.g. `/tmp/fpgo-driver/main.go`) with its own `go.mod`:

  ```
  module driver
  go 1.18
  require github.com/TeaEntityLab/fpGo/v2 v0.0.0
  replace github.com/TeaEntityLab/fpGo/v2 => <repo root>
  ```

  then `go run -mod=mod .` — never add driver files inside the repo.

## Evidence

- Save `go test -v -run <regex>` output to files under a scratch dir
  (e.g. `/tmp/fpgo-evidence/`).
- For scratch drivers, record the driver source plus its stdout.

## When a check fails

Classify before fixing — four different failures need four different repairs:

| Class | Meaning | Repair |
| --- | --- | --- |
| Product regression | Library behavior changed for the worse | Report; do not edit the map or tests to match |
| Doc drift | The map/VERIFY.md no longer matches the library | Update the doc; the product is fine |
| Spec/oracle error | The acceptance spec or test expectation is wrong | Oracles are read-only; propose the change for review |
| Harness failure | Toolchain, vendoring, environment, or external dep broke | Fix the harness; the product is fine |

Never "fix" a failing check by editing the map or a test to describe the new
behavior — that hides product regressions. If the new behavior is intentional,
the map/test update is a separate reviewed change.

## Cleanup

- Remove scratch driver dirs and evidence dirs you created.
- `git status` must show no modifications inside the repo.
