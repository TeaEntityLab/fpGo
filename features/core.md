# Feature: core

`source_commit: d7ae751` · `last_verified_at: 2026-09-23` · `verification_status: passed`

User goal: Optional/Maybe semantics, Rx-like MonadIO, Publisher pub/sub,
fp function helpers, pattern matching.

## Preconditions

- Repo root; `go test -mod=mod` works (VERIFY.md → Launch, Doctor).
- No server, no network, no secrets — pure in-process library.

## Entry points

| Area | Files | Test files |
| --- | --- | --- |
| Maybe/Optional + conversions | `maybe.go` | `maybe_test.go`, `maybe_conversion_consistency_test.go` |
| MonadIO (Rx-like) | `monadIO.go` | `monadIO_test.go` |
| Publisher | `publisher.go` | `publisher_test.go` |
| fp functions + pattern matching | `fp.go` | `fp_test.go`, `fp_edge_test.go` |

## Drive

```bash
go test -mod=mod -v -run 'TestMaybe|TestJust|TestOr|TestLet|TestIs|TestNone|TestFlat|TestCast|TestClone|TestSome|TestTo' .
go test -mod=mod -v -run 'TestMonadIO|TestAtomBool' .
go test -mod=mod -v -run 'TestPublisher' .
go test -mod=mod -v -run 'TestCompose|TestPipe|TestFPFunctions|TestTrampoline|TestConcat|TestReject|TestMapIndexed|TestReduceIndexed|TestFlattenAndPrepend|TestPartition|TestSplitEvery|TestVariadic|TestCurry|TestComp|TestPatternMatching|TestZipErrors|TestPMap|TestDefProduct|TestNilTypeMatches|TestMergeForInterface|TestMinusForInterface|TestKeysForInterface|TestValuesForInterface|TestDistinctForInterface|TestExistsForInterface|TestIntersection' .
```

## Observable outcomes

- `Maybe.Just(1).IsPresent()` → true; `Maybe.Just(nil).IsPresent()` → **false**
  (Just(nil) is present-as-nil: `IsNil()` true, `IsPresent()` false).
- `Maybe.Just(nil).Or(3)` → 3; `Maybe.Just(1).Or(3)` → 1.
- `m.Let(fn)` runs fn only when present.
- MonadIO `Just(1).FlatMap(f).Subscribe(...)` delivers mapped values via OnNext.
- Publisher: `Publish` fans out to subscribers; unsubscribe during publish is safe.
- fp: `Compose`/`Pipe` ordering, `Distinct`, `Intersection`, `PMap` (concurrent,
  order not guaranteed — `TestPMapNoOrderConcurrency`).

## Failure paths

- `Just(nil)` treated as "present" → the nil-vs-empty distinction is broken.
- `Or` returning the fallback when a value exists → extraction inverted.
- Conversion guards (`TestNegativeSignedToUnsignedRejected`,
  `TestPositiveOverflowGuardsReturnError`) → overflow/negative rejection broken.
- Publisher publish to nil subscriber or during unsubscribe → deadlock/panic.

## Evidence

- `go test -v` output per regex above.
- For conversions, the consistency test file is the oracle — do not weaken it.
