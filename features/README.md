# Feature Map — fpGo

One file per feature area. Each answers: what exists, how to drive it, what
usually lies. Sweep order for a broad regression pass:
`core.md` → `collections.md` → `queues.md` → `concurrency.md` → `integrations.md`.

- [core.md](core.md) — Maybe/Optional, MonadIO, Publisher, fp functions, pattern matching
- [collections.md](collections.md) — Stream (generics + interface), sortDescriptor
- [queues.md](queues.md) — LinkedListQueue, ChannelQueue, BufferedChannelQueue, ConcurrentQueue
- [concurrency.md](concurrency.md) — Actor model, Coroutine (yield/yieldFrom), Handler
- [integrations.md](integrations.md) — network/SimpleHTTP, worker/WorkerPool

## Maintenance

Each feature file carries `source_commit` + `last_verified_at` +
`verification_status`. After any change to the covered source files, re-run
the affected feature's Drive section and update the metadata. Do not rewrite
the map on a schedule — re-verify on relevant code change only (event-driven,
not calendar-driven).

A `verification_status` of `failed` means the product changed; classify the
failure per VERIFY.md §When a check fails before editing the map.
