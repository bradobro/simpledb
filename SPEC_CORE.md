# SimpleDB Core Specification

**Version:** 1.0.0

## Overview

SimpleDB is a document database for small applications. Store JSON records in named collections.

## Core Concepts

### Collections

Named containers for records. Collection names: alphanumeric only (no `/`, `.`, or `-`).

### Records

JSON objects stored in collections. Each record has a unique ID within its collection.

### IDs

Character set: alphanumeric plus `-` (for prefixed IDs like `user-0ABC123`). See [SPEC_IDS.md](SPEC_IDS.md).

## Record Creation

Two ways to create records:

2. **Specified ID:** `CreateWithID(id, data)` - You provide the ID. ErrExists if record exists.
1. **Auto-generated ID:** `Create(data)` - SimpleDB generates a unique ID. If record exists, retry (configurable N and waiting) another generated ID, then ErrExists after N retries.

**Critical rule:** Neither method overwrites existing records.

## Core API

### Store Operations

| Method | Description |
|--------|-------------|
| `RawCollection(name)` | Get collection handle |
| `AddCollection(name, idAlgorithm)` | Register collection with ID algorithm |
| `Collections()` | List registered collections |
| `Close()` | Release resources |
| `Check()` | Validate consistency |
| `Fix()` | Repair inconsistencies |

### Collection Operations

| Method | Description |
|--------|-------------|
| `Create(data)` | Create record, return generated ID |
| `CreateWithID(id, data)` | Create record with specified ID |
| `RawRead(id)` | Read record bytes |
| `Update(id, data)` | Replace existing record |
| `Delete(id)` | Remove record |
| `Exists(id)` | Check if record exists |
| `IDs(start)` | Iterate record IDs |
| `RawItems(start)` | Iterate (ID, data) pairs |

## Error Semantics

| Error | Condition |
|-------|-----------|
| `ErrNotFound` | Record doesn't exist |
| `ErrExists` | Record/collection already exists (or non-colliding ID couldn't be generated after configured retries.) |
| `ErrClosed` | Store has been closed |
| `ErrInvalidID` | ID contains invalid characters |

## Concurrency

TODO: drivers must document concurrency and consistency promises of each ID type.

In general, strive for

- Concurrent reads: always safe
- Concurrent writes to different IDs: safe
- Concurrent updates to same ID: last write wins
- All writes are atomic

## ID Storage

IDs are managed by the storage backend. See [SPEC_DISK.md](SPEC_DISK.md) for file-based storage details on how IDs are stored separately from record content (in file name, NOT in record--a pattern recommended but not required.).

## ID Generation

ID algorithms are pluggable. See [SPEC_IDS.md](SPEC_IDS.md) for available algorithms and recommendations.
