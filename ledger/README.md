# Ledger

Ledger is a small Kafka-inspired segmented commit log written in Go.

## What it was meant to be

Ledger was built as a toy system to understand the mechanics behind a commit log, not just talk about them. The goal was to build a real file-backed system with:

- an append-only store
- a separate index file
- segmented on-disk layout
- rollover between segments
- reopening persisted data across runs
- a small reader/writer API on top of those pieces

## What it was not meant to be

Ledger was not meant to be:

- a production-ready broker
- Kafka-compatible at the protocol level
- a replicated system
- a system with retention, compaction, consumer groups, or crash-recovery machinery
- a CLI-first app; the `cmd/*` binaries are just local harnesses over the package APIs

In short: this project was meant to explore the core storage ideas honestly, not to simulate every production concern.

## Current architecture

The codebase is split into three layers:

- `api`
  - `Writer` appends messages and rolls segments when they reach a configured message limit.
  - `Reader` reads sequentially from an absolute message position and moves across segments as needed.
- `models`
  - `Segment` owns one `index.data` file and one `store.data` file.
  - `Store` manages append/read behavior for the raw message file.
  - `Index` manages append/read behavior for the fixed-width offset file.
  - `Position` resolves an absolute read position to a segment start plus a segment-local cursor.
  - `Entry` is the value returned by reads and writes.
- `concerns`
  - small binary and file helpers used by the storage layer

## Terminology

Ledger currently uses three different coordinates, and they mean different things:

- `Segment` — the absolute starting message position of a segment; also used as the segment directory name
- `Position` / `Cursor` — the message index within a segment
- `Offset` — the byte offset of a record inside `store.data`

For example, if the segment directory is `0000000000000010`, then that segment starts at absolute message position `10`. The first message inside it has:

- `Segment = 10`
- `Position = 0`
- `Offset = 0`

## On-disk layout

A ledger root contains zero-padded segment directories. Each segment contains two files:

```text
<ledger-root>/
  0000000000000000/
    index.data
    store.data
  0000000000000010/
    index.data
    store.data
  0000000000000020/
    index.data
    store.data
```

Segment directory names are zero-padded 16-digit decimal numbers. In the current design they represent the starting absolute message position of that segment. If the writer limit is `10`, segment starts are `0`, `10`, `20`, and so on.

### `store.data`

`store.data` is append-only binary data. Each record is stored as:

```text
[8-byte big-endian message length][message bytes]
```

The byte where a record begins in `store.data` is the record's store offset.

### `index.data`

`index.data` is append-only fixed-width binary data. Each entry is:

```text
[8-byte big-endian store offset]
```

The entry's slot in `index.data` is the segment-local message position.

## Write path

Current write flow:

1. `api.NewWriter(directory, limit)` ensures the ledger root exists and attaches to the latest segment, or segment `0` if none exist yet.
2. `Writer.Append(message)` writes to the current segment.
3. `Segment.Append` appends the message to `store.data` with an 8-byte length prefix.
4. The returned byte offset in `store.data` is appended to `index.data`.
5. When `segment.GetSize() >= limit`, the writer rolls to the next segment.
6. The next segment start is computed as `currentSegment.Start + currentSegment.GetSize()`.

## Read path

Current read flow:

1. `api.NewReader(directory, position)` takes an absolute message position.
2. `PositionForRead` resolves that absolute position to:
   - `Start`: the segment directory to open
   - `Cursor`: the segment-local message position within that segment
3. `Reader.ReadNext()` reads the current entry and then increments the cursor.
4. `Segment.ReadAt(position)` reads the store offset from `index.data`, then reads the record header and payload from `store.data`.
5. When the reader reaches the end of the current segment and another segment exists, it rolls over to the next segment and continues.

## Current assumptions and limits

These are intentional parts of the current design:

- The reader is allowed to wait for the ledger root directory to appear during construction.
- The reader is not allowed to open before the first segment exists.
- Opening a reader against an existing ledger with no initialized segments is treated as invalid.
- The project uses panic-driven error handling in many places instead of a rich error surface.
- The CLI binaries are just thin local harnesses and are not the core abstraction of the project.

## Local harnesses

There are two small commands for manual experimentation:

- `go run ./cmd/writer` — append messages from stdin
- `go run ./cmd/reader` — poll and print messages from the ledger

These are useful for playing with the system, but the main design lives in `api`, `models`, and `concerns`.

## Testing

The project now has package-level tests for all non-main packages:

- `api`
- `models`
- `concerns`

Run them with:

```bash
go test ./...
```

Those tests cover the storage primitives, position resolution, segment reopen/rollover behavior, and the reader/writer APIs.
