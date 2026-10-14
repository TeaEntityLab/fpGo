# Feature: queues

`source_commit: 32a0438` · `last_verified_at: 2026-10-14` · `verification_status: passed`

User goal: FIFO/LIFO containers — LinkedListQueue (queue+stack),
ChannelQueue, BufferedChannelQueue (timeouts, pool), ConcurrentQueue.

## Preconditions

- Repo root; `go test -mod=mod` works.

## Entry points

| Area | Files | Test files |
| --- | --- | --- |
| All queue impls | `queue.go` | `queue_test.go`, `queue_regression_test.go`, `queue_concurrent_regression_test.go` |

## Drive

```bash
go test -mod=mod -v -run 'TestLinked|TestDoubly|TestQueue|TestNew' .
go test -mod=mod -v -run 'TestChannel' .
go test -mod=mod -v -run 'TestBuffered' .
go test -mod=mod -v -run 'TestConcurrent' .
```

## Observable outcomes

- LinkedListQueue as Queue: `Offer`→tail, `Poll`→head (FIFO: 1,2,3 in → 1,2,3
  out; empty → `ErrQueueIsEmpty`).
- Same impl as Stack: `Push`/`Pop` at tail (LIFO: 1,2,3 in → 3,2,1 out;
  empty → `ErrStackIsEmpty`).
- ChannelQueue: `Poll`/`Take`/`Put`/`Offer` on a closed channel return
  `ErrQueueIsClosed`. A send must not panic, and a receive must not look like
  a real zero.
- BufferedChannelQueue: `Poll` drains the overflow pool too (not just the
  channel), so a fresh `Poll` after over-capacity `Offer`s does not return a
  spurious `ErrQueueIsEmpty`; `TakeWithTimeout` honors the timeout; `Close`
  stops freeNodePool/loadFromPool goroutines
  (`TestBufferedChannelQueueCloseStopsFreeNodeGoroutine`).
- ConcurrentQueue: no loss under concurrent producers/consumers
  (`TestConcurrentQueueNoLossUnderConcurrency`). Wrapping a self-synchronizing
  queue (ChannelQueue/BufferedChannelQueue) bypasses the wrapper lock, so a
  blocked `Take` cannot deadlock producers (`TestConcurrentQueueChannelBackedProducerConsumer`).

## Failure paths

- FIFO/LIFO order swapped or drained early → core contract broken.
- `TakeWithTimeout` returning immediately on empty → timeout path broken.
- `Close` leaking pool goroutines → lifecycle regression.
- Concurrent loss/duplication → race in the wrapper.

## Evidence

- `go test -v` output; concurrent tests are the oracle for loss claims —
  a single-threaded pass does not prove the concurrent contract.
