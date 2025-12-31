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

1. **Specified ID:** `Store(id, data)` - You provide the ID. ErrExists if record exists.
2. **Auto-generated ID:** `Create(data)` - SimpleDB generates a unique ID. Retries on collision.

**Critical rule:** Neither method overwrites existing records.

## Core API

### SimpleDB Operations

| Method | Description |
|--------|-------------|
| `Get(name)` | Get collection handle |
| `Create(name, idAlgorithm)` | Register collection |
| `List()` | List registered collections |
| `Close()` | Release resources |
| `Check(fix, args...)` | Validate/repair consistency |

### Collection Operations

| Method | Description |
|--------|-------------|
| `Store(id, data)` | Store with specified ID |
| `Create(data)` | Create with auto-generated ID |
| `Read(id)` | Read record bytes |
| `Update(id, data)` | Replace existing record |
| `Delete(id)` | Remove record |
| `Exists(id)` | Check if record exists |
| `Items(start)` | Iterate (id, data) pairs |
| `IDs(start)` | Iterate record IDs |

## Error Semantics

| Error | Condition |
|-------|-----------|
| `ErrNotFound` | Record doesn't exist |
| `ErrExists` | Record/collection already exists |
| `ErrClosed` | SimpleDB has been closed |
| `ErrInvalidID` | ID contains invalid characters |
| `ErrIO` | Other I/O error |

## Concurrency

TODO: drivers must document concurrency and consistency promises.

In general, strive for:

- Concurrent reads: always safe
- Concurrent writes to different IDs: safe
- Concurrent updates to same ID: last write wins
- All writes are atomic

## ID Storage

IDs are managed by the storage backend. See [SPEC_DISK.md](SPEC_DISK.md) for file-based storage details.

## ID Generation

ID algorithms are pluggable. See [SPEC_IDS.md](SPEC_IDS.md) for available algorithms and recommendations.
