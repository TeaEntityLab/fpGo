# Feature: integrations

`source_commit: d7ae751` · `last_verified_at: 2026-09-23` · `verification_status: passed`

User goal: Retrofit-inspired SimpleHTTP client API and ExecutorService-like
WorkerPool — the two subpackages.

## Preconditions

- Repo root; `go test -mod=mod ./network ./worker` works.
- network tests use local httptest servers — no external network needed.

## Entry points

| Area | Files | Test files |
| --- | --- | --- |
| SimpleHTTP | `network/` | `network/*_test.go` |
| WorkerPool | `worker/` | `worker/*_test.go` |

## Drive

```bash
go test -mod=mod -v ./network
go test -mod=mod -v ./worker
```

## Observable outcomes

- SimpleHTTP: method selection (JSON vs multipart), header/body options,
  serializer errors surface as response errors — not panics.
- WorkerPool: submitted jobs execute; pool respects worker count; shutdown
  drains or rejects per its contract.

## Failure paths

- Serializer error swallowed or panicking → error-path contract broken.
- Wrong HTTP method/body encoding selected → request-building broken.
- WorkerPool dropping or double-running jobs → scheduling broken.

## Evidence

- `go test -v` output per package.
