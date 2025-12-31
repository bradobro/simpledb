# SimpleDB API Specification

See [api.go](api.go) and [api.ts](api.ts) for example interface definitions.

## Store Interface

| Method | Returns | Description |
|--------|---------|-------------|
| `RawCollection(name)` | RawCollection | Get collection handle (lazy) |
| `AddCollection(name, algorithm)` | error | Register collection |
| `Collections()` | []string, error | List collections |
| `Close()` | error | Release resources |
| `Check()` | CheckReport, error | Validate consistency |
| `Fix()` | FixReport, error | Repair issues |
TODO: consider adding driver info (format version, ID schemes and characteristics

## RawCollection Interface

| Method | Returns | Description |
|--------|---------|-------------|
| `Create(data)` | id, error | Create with auto-generated ID |
| `CreateWithID(id, data)` | error | Create with specified ID |
| `RawRead(id)` | []byte, error | Read raw JSON |
| `Update(id, data)` | error | Replace record |
| `Delete(id)` | error | Remove record |
| `Exists(id)` | bool, error | Check existence |
| `IDs(start)` | iterator | Iterate IDs lexicographically |
| `RawItems(start)` | iterator | Iterate (id, data) pairs |


## Typed Collection

Generic wrapper providing JSON marshal/unmarshal:

| Method | Returns | Description |
|--------|---------|-------------|
| `Create(record)` | id, error | Marshal and create |
| `Read(id)` | T, error | Read and unmarshal |
| `Update(id, record)` | error | Marshal and update |
| `Items(start)` | iterator | Iterate typed records |

## Errors

| Error | Condition |
|-------|-----------|
| `ErrNotFound` | Record doesn't exist |
| `ErrExists` | Record/collection already exists |
| `ErrClosed` | Store closed |
| `ErrInvalidID` | Invalid ID format |

## Behavior Notes

- `Create` and `CreateWithID` never overwrite existing records
- `Update` requires record to exist
- All operations are atomic
- Iteration is lexicographic by ID
- `start` parameter for pagination (skip IDs < start)
