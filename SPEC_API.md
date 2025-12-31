# SimpleDB API Specification

See [api.go](api.go) and [api.ts](api.ts) for interface definitions.

## SimpleDB Interface

| Method | Returns | Description |
|--------|---------|-------------|
| `Get(name)` | ByteCollection | Get collection handle (lazy) |
| `Create(name, algorithm)` | error | Register collection |
| `List()` | []string, error | List collections |
| `Close()` | error | Release resources |
| `Check(fix, args...)` | CheckReport, error | Validate/repair consistency |

## ByteCollection Interface

Composed of ByteStorer, ByteUpdater, ByteLister.

### ByteStorer

| Method | Returns | Description |
|--------|---------|-------------|
| `Store(id, data)` | error | Store with specified ID (errors if exists) |
| `Exists(id)` | bool, error | Check existence |
| `Check(fix, args...)` | CheckReport, error | Validate/repair collection |

### ByteUpdater

| Method | Returns | Description |
|--------|---------|-------------|
| `Create(data)` | id, error | Create with auto-generated ID |
| `Read(id)` | []byte, error | Read bytes |
| `Update(id, data)` | error | Replace record |
| `Delete(id)` | error | Remove record |

### ByteLister

| Method | Returns | Description |
|--------|---------|-------------|
| `Items(start)` | iterator | Iterate (id, data) pairs |
| `IDs(start)` | iterator | Iterate IDs |

## Typed Collection

Composed of Storer[T], Updater[T], Lister[T].

| Method | Returns | Description |
|--------|---------|-------------|
| `Store(id, record)` | error | Store with specified ID |
| `Exists(id)` | bool, error | Check existence |
| `Create(record)` | id, error | Marshal and create |
| `Read(id)` | T, error | Read and unmarshal |
| `Update(id, record)` | error | Marshal and update |
| `Delete(id)` | error | Remove record |
| `Items(start)` | iterator | Iterate typed records |
| `IDs(start)` | iterator | Iterate IDs |

## Errors

| Error | Condition |
|-------|-----------|
| `ErrNotFound` | Record doesn't exist |
| `ErrExists` | Record/collection already exists |
| `ErrClosed` | SimpleDB closed |
| `ErrInvalidID` | Invalid ID format |
| `ErrIO` | Other I/O error |

## Check Args

| Arg | Scope |
|-----|-------|
| (empty) | No checks, just invoke framework |
| `meta` | Metadata and collection consistency |
| `fast` | All checks except reading every record |
| `slow` | All checks including reading every record |
| `all` | Same as `meta` + `slow` |

## Behavior Notes

- `Store` and `Create` never overwrite existing records
- `Update` requires record to exist
- All operations are atomic
- Iteration is lexicographic by ID
- `start` parameter for pagination (skip IDs < start)
