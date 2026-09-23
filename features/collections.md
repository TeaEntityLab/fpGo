# Feature: collections

`source_commit: d7ae751` · `last_verified_at: 2026-09-23` · `verification_status: passed`

User goal: Java8Stream-like lazy/eager collection ops (generics + interface
variants) and sort descriptors.

## Preconditions

- Repo root; `go test -mod=mod` works.

## Entry points

| Area | Files | Test files |
| --- | --- | --- |
| Stream (generics) | `stream.go` | `stream_test.go`, `stream_regression_test.go` |
| StreamForInterface | `streamForInterface.go` | `streamForInterface_test.go` |
| sortDescriptor | `sortDescriptor.go` | `sortDescriptor_test.go` |

## Drive

```bash
go test -mod=mod -v -run 'TestStream|TestNew|TestFrom|TestFilter|TestMap|TestSet|TestSort' .
go test -mod=mod -v -run 'TestSortDescriptor|TestComparableStringHouseConvention' .
```

## Observable outcomes

- `StreamFromArray([]int{}).Append(1,1).Extend(StreamFromArray([]int{2,3,4}))`
  → `ToArray()` = `[1 1 2 3 4]`; `.Distinct()` → `[1 2 3 4]`.
- `StreamForInterface` tolerates nil elements; `FilterNotNil()` drops them.
- `Stream.Remove` does **not** mutate the receiver and does not alias it
  (regression-tested: `TestStreamRemoveDoesNotMutateReceiver`,
  `TestStreamRemoveResultDoesNotAliasReceiver`).
- sortDescriptor: ascending/descending actually reorder; mixed keys tie-break.

## Failure paths

- `Distinct` returning input unchanged → dedupe broken.
- `Remove` mutating or aliasing the receiver → immutability contract broken.
- `Append`/`Extend` dropping elements or reordering → sequence contract broken.
- sortDescriptor not reordering → comparator wiring broken.

## Evidence

- `go test -v` output; for aliasing checks the regression tests are the oracle.
