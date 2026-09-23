# Feature: concurrency

`source_commit: d7ae751` · `last_verified_at: 2026-09-23` · `verification_status: passed`

User goal: Akka/Erlang-like actors (send/receive/spawn/states), Pythonic
coroutines (yield/yieldFrom), Handler posting.

## Preconditions

- Repo root; `go test -mod=mod` works.

## Entry points

| Area | Files | Test files |
| --- | --- | --- |
| Actor | `actor.go` | `actor_test.go` |
| Coroutine | `cor.go` | `cor_test.go` |
| Handler | `handler.go` | `handler_test.go` |

## Drive

```bash
go test -mod=mod -v -run 'TestActor|TestAsk' .
go test -mod=mod -v -run 'TestCor' .
go test -mod=mod -v -run 'TestHandler|TestGetDefault' .
```

## Observable outcomes

- Actor: `Send` is async; `Spawn` creates children; `Ask` returns a reply or
  times out (`TestActorAskTimeoutThenReplySucceedsForNextAsk` — a timed-out
  Ask must not poison the next one); `Close` shuts down; `Send`/`Spawn` after
  close do not panic.
- Coroutine: `Yield` suspends/resumes; `YieldFrom` delegates; `IsDone`/
  `IsStarted` track lifecycle; `Close` on a done coroutine is safe.
- Handler: `Post` executes in order; `Post` after `Close` does not panic;
  external channel close exits the run loop.

## Failure paths

- Ask timeout blocking or panicking on late reply → ask lifecycle broken.
- Spawn/Send after Close panicking → close-guard broken.
- Coroutine `Yield` after done resuming → lifecycle broken.
- Handler posts executing out of order → ordering contract broken.

## Evidence

- `go test -v` output; race-sensitive paths have dedicated tests
  (`TestActorConcurrentSendCloseDoesNotPanic`,
  `TestCorDoCloseSafeRaceWithClose`) — run them, don't eyeball.
