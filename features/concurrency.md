# Feature: concurrency

`source_commit: 32a0438` · `last_verified_at: 2026-10-14` · `verification_status: passed`

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

- Actor: `Send` hands the message to the mailbox channel and returns — with
  the default unbuffered channel it blocks only until the run loop picks the
  message up, so an actor must never `Send` to itself from inside its own
  effect (rendezvous deadlock). `Spawn` creates children with unique IDs even
  when two spawns share a `time.Now()` tick; `Ask` returns a reply or times
  out (`TestActorAskTimeoutThenReplySucceedsForNextAsk` — a timed-out Ask
  must not poison the next one); `Close` shuts down; `Send`/`Spawn` after
  close do not panic. `Close` does not close the mailbox (it may be caller-owned);
  if that send is already possible, `select` can still deliver one message
  after close — Go has no priority select.
- Coroutine: `Yield` suspends/resumes; `YieldFrom` delegates; `IsDone`/
  `IsStarted` track lifecycle; `Close` on a done coroutine is safe. A send
  that races `close` is dropped instead of panicking, and `close` is not
  blocked behind that send.
- Handler: `Post` executes in order; `Post` after `Close` does not panic;
  external channel close exits the run loop. Same select limitation as Actor
  `Send`: a post already racing `Close` may still be queued.

## Failure paths

- Ask timeout blocking or panicking on late reply → ask lifecycle broken.
- Spawn/Send after Close panicking → close-guard broken.
- Coroutine `Yield` after done resuming → lifecycle broken.
- Handler posts executing out of order → ordering contract broken.

## Evidence

- `go test -v` output; race-sensitive paths have dedicated tests
  (`TestActorConcurrentSendCloseDoesNotPanic`,
  `TestCorDoCloseSafeRaceWithClose`) — run them, don't eyeball.
