# Feature: integrations

`source_commit: 32a0438` · `last_verified_at: 2026-10-14` · `verification_status: passed`

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
  serializer errors surface as response errors — not panics; `RoundTrip`
  without a transport errors instead of panicking.
- The request context stays live until the caller closes `Response.Body`
  (cancelling right after `Do` used to abort body reads). SimpleAPI sets `Err`
  on non-2xx and keeps up to 512 bytes of the error body in that error.
  Placeholders in the path are `url.PathEscape`d; placeholders in the query
  (after the first `?`) are `url.QueryEscape`d.
- WorkerPool: submitted jobs execute; pool respects worker count; a `trySpawn`
  panic is recovered inside the spawn loop instead of killing scaling;
  worker-expiry uses a recycled timer; shutdown drains or rejects per its
  contract. Settings fields are private — use `NewDefaultWorkerPoolSettings()`
  or nil; a literal zero struct never spawns workers.

## Failure paths

- Serializer error swallowed or panicking → error-path contract broken.
- Wrong HTTP method/body encoding selected → request-building broken.
- WorkerPool dropping or double-running jobs → scheduling broken.

## Evidence

- `go test -v` output per package.
