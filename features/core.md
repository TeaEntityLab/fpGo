# Feature: core

`source_commit: 32a0438` · `last_verified_at: 2026-10-14` · `verification_status: passed`

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
  (Just(nil) is present-as-nil: `IsNil()` true, `IsPresent()` false). Typed
  nils of any nilable kind (ptr/slice/map/chan/func) are also absent.
- `Maybe.Just(nil).Or(3)` → 3; `Maybe.Just(1).Or(3)` → 1.
- `m.Let(fn)` runs fn only when present. `Map`/`FlatMap` short-circuit the
  same way — fn never receives a nil. `MaybeMap(m, T->R)` crosses types
  (interface methods cannot declare type params) and preserves absence. The
  `Maybe` namespace var itself is an absent Maybe (`Maybe.Or(x)` → x).
- MonadIO `SubscribeOn(h)` schedules the effect; `ObserveOn(h)` schedules
  OnNext delivery (Rx convention).
- MonadIO `Just(1).FlatMap(f).Subscribe(...)` delivers mapped values via OnNext.
- Publisher: `Publish` fans out to subscribers; unsubscribe during publish is safe.
- fp: `Compose`/`Pipe` ordering, `Distinct`, `Intersection`, `PMap` (concurrent;
  default worker count is bounded by NumCPU×8, override via `PMapOption`).
- `InCaseOfRegex` precompiles the pattern at construction; an invalid pattern
  never matches. `MatchFor` panics only when no pattern matches — design an
  `Otherwise` arm; non-`CompData` structs and nil pointers no longer panic
  inside pattern dispatch.

## Failure paths

- `Just(nil)` treated as "present" → the nil-vs-empty distinction is broken.
- `Or` returning the fallback when a value exists → extraction inverted.
- Conversion guards (`TestNegativeSignedToUnsignedRejected`,
  `TestPositiveOverflowGuardsReturnError`) → overflow/negative rejection broken.
- Publisher publish to nil subscriber or during unsubscribe → deadlock/panic.

## Evidence

- `go test -v` output per regex above.
- For conversions, the consistency test file is the oracle — do not weaken it.
