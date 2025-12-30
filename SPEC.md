# SimpleDB Specification

**Version:** 1.0.0
**Last Updated:** 2025-12-27

## Overview

SimpleDB is a simple document database for small applications that want to keep their options open without making premature architectural decisions. This specification describes the storage semantics and API in terms of the **on-disk (file-based) backend**.

**Important:** Future storage backends (in-memory, cloud, embedded database) may significantly alter **how** they store information, but should maintain the **behavioral semantics** described here.

## Reference Implementation

The **Go implementation** is the reference. Other language implementations (TypeScript, Python, etc.) must maintain behavioral compatibility with Go.

---

## Design Philosophy

### Trust Disk Data

SimpleDB assumes **data on disk is correct**. We optimize for performance over paranoia:

- **No validation on every read:** Collection names, IDs, and file formats are trusted to be valid
- **Check() validates when asked:** Explicit consistency checks detect corruption or manual edits
- **Fix() repairs issues:** Automated repair of detected inconsistencies
- **Fast path by default:** Normal operations don't pay validation overhead

This approach provides:
- **Better performance:** No redundant checks on hot paths
- **Developer-friendly:** Easy to inspect and manually edit database files
- **Fail-fast when corrupted:** Errors surface immediately when data is actually corrupt

### Single Source of Truth for IDs

**The filename IS the ID.** IDs are never stored in JSON content, only in filenames.

**Rationale:**
- **Impossible to desync:** Can't have filename/content ID mismatches
- **Smaller files:** No redundant ID storage in JSON
- **Clearer semantics:** Filename = ID (obvious and inspectable)
- **Simpler validation:** Only need to check filename format

**Implementation:**
- **Write operations:** Strip `id` field from JSON before writing to disk
- **Read operations:** Inject `id` field from filename when reading
- **Update operations:** Validate id parameter matches filename, then strip before writing

**Example on disk:**
```
people/
├── 0000000001.json    # Filename is the ID
│   {                  # JSON content (no "id" field)
│     "name": "Alice",
│     "age": 30
│   }
```

**Example in memory (after read):**
```json
{
  "id": "0000000001",  // Injected from filename
  "name": "Alice",
  "age": 30
}
```

**Backward Compatibility:**
- Old databases with `id` fields in JSON still work
- On write, `id` field is stripped (auto-migrates)
- On read with existing `id` field, it's validated against filename then replaced

---

## On-Disk Format

### Directory Structure

```
{root}/
├── db.json                     # Database metadata
├── {collection}/               # Collection directory
│   ├── .id                     # Sequential ID state (seq10/seq36 only)
│   ├── {id}.json               # Record file
│   └── ...
├── {schema}.schema.json        # Optional JSONSchema docs
└── ...
```

### Metadata File: db.json

Location: `{root}/db.json`

Format:
```json
{
  "version": "1.0.0",
  "collections": {
    "{collection-name}": {
      "id_algorithm": "{algorithm}"
    }
  }
}
```

**Fields:**
- `version` (string): Backend format version (currently "1.0.0")
- `collections` (object): Map of collection names to configuration
  - `id_algorithm` (string): ID generation algorithm
    - Valid values: `"cuid2"`, `"seq10"`, `"seq36"`

**Creation:**
- Created automatically on first `AddCollection()` if missing
- Updated atomically via temp file + rename

### Collection Directories

Location: `{root}/{collection-name}/`

**Properties:**
- Directory name = collection name
- Contains one `.json` file per record
- May contain `.id` file for sequential ID algorithms
- Created automatically when collection is added to metadata

### Record Files

Location: `{root}/{collection-name}/{id}.json`

Format:
```json
{
  "id": "{id}",
  ...additional fields...
}
```

**Properties:**
- Filename = record ID + `.json` extension
- Content = JSON object with required `id` field
- The `id` field in JSON must match the filename (minus extension)
- Written atomically via temp file + rename

**Naming Rules:**
- IDs must be alphanumeric (a-z, A-Z, 0-9)
- No length restrictions beyond filesystem limits
- Case-sensitive (though lowercase recommended for consistency)

### Sequential ID State File: .id

Location: `{root}/{collection-name}/.id`

Format: Single-line text file containing current sequence number

**Properties:**
- Used only by `seq10` and `seq36` algorithms
- Contains the last generated ID (not next ID)
- Updated atomically on each `Create()` operation
- Not present for `cuid2` collections
- Created with "0" if missing (next Create generates "1")

**Examples:**
- `seq10`: Contains "0000000042" after 42 records created
- `seq36`: Contains "000000z" after 35 records created (35 in base36 = z)

### Schema Files (Optional)

Location: `{root}/{name}.schema.json`

Format: JSONSchema specification

**Properties:**
- Documentation only (not enforced at runtime)
- Naming convention: `{singular-name}.schema.json`
- Example: `person.schema.json` for `people` collection

---

## ID Generation Algorithms

### CUID2 (Collision-Resistant)

**Algorithm:** `"cuid2"`

**Properties:**
- Generates collision-resistant IDs using cuid2 library
- ID length: ~24 characters
- Character set: lowercase alphanumeric
- No state file required
- Thread-safe by design (stateless)

**Example IDs:**
```
clxyz1000000house0001
clkv6s3gg0000qz08wz9z3z3z
```

**Use Cases:**
- Default choice for most collections
- Distributed systems (no coordination needed)
- Collections with frequent concurrent writes

### SEQ10 (Sequential Decimal)

**Algorithm:** `"seq10"`

**Properties:**
- Generates sequential 10-digit zero-padded decimal IDs
- ID format: `0000000000` to `4294967295` (uint32 range)
- Requires `.id` state file in collection directory
- Thread-safe via mutex protection

**Example IDs:**
```
0000000001
0000000042
0000001337
4294967295
```

**Zero-Padding Rules:**
- Always 10 digits
- Enables lexicographic sorting (text sort = numeric sort)
- Leading zeros dropped for display/reference

**Use Cases:**
- Human-readable IDs
- Testing/debugging (predictable sequence)
- Import order tracking

### SEQ36 (Sequential Base-36)

**Algorithm:** `"seq36"`

**Properties:**
- Generates sequential base-36 encoded IDs
- Character set: `0-9a-z` (36 characters)
- ID format: 7 characters zero-padded
- Range: `0000000` to `1z141z3` (uint32 max in base-36)
- Requires `.id` state file in collection directory
- Thread-safe via mutex protection

**Example IDs:**
```
0000001  # 1
000000z  # 35
0000010  # 36
0000a2p  # 13284
1z141z3  # 4294967295 (max uint32)
```

**Base-36 Encoding:**
```
0-9 → 0-9
a-z → 10-35
```

**Zero-Padding Rules:**
- Always 7 characters (max uint32 in base36)
- Enables lexicographic sorting
- Visual shortening: display as `"a2p"` but store as `"0000a2p"`

**Use Cases:**
- Compact human-readable IDs
- URL-friendly identifiers
- Balance between compactness and readability

### Custom IDs

**Method:** `CreateWithID(id, data)`

**Properties:**
- Accepts any alphanumeric string as ID
- Does NOT update `.id` file (manual ID assignment)
- Enables foreign keys, user-chosen IDs, external system integration
- Returns error if ID already exists

---

## API Semantics

### Store Interface

#### `RawCollection(name) → RawCollection`

Returns a collection handle for byte-level operations.

**Behavior:**
- Returns handle immediately (does not check filesystem)
- Collection uses `id_algorithm` from metadata if registered
- Defaults to `"cuid2"` if collection not in metadata
- Directory created lazily on first write operation

#### `AddCollection(name, idAlgorithm) → error`

Registers a new collection with specified ID algorithm.

**Behavior:**
- Updates `db.json` metadata atomically
- Creates collection directory if missing
- For `seq10`/`seq36`: creates `.id` file with "0"
- Returns `ErrExists` if collection already registered
- Valid algorithms: `"cuid2"`, `"seq10"`, `"seq36"`

**Atomic Guarantee:**
- Metadata and filesystem state updated together
- Failure leaves no partial state

#### `Collections() → []string, error`

Returns list of registered collection names.

**Behavior:**
- Reads from `db.json` metadata
- Returns names in unspecified order
- Empty list if no collections registered

#### `Close() → error`

Releases resources and marks store as closed.

**Behavior:**
- After close, all operations return `ErrClosed`
- Idempotent (safe to call multiple times)
- No automatic flush needed (writes are atomic)

#### `Check() → CheckReport, error`

Validates consistency between metadata and filesystem.

**Checks Performed:**
1. **Metadata checks:**
   - `db.json` exists and is valid JSON
   - All `id_algorithm` values are valid
2. **Collection checks:**
   - All registered collections have directories
   - All collection directories are in metadata
   - Sequential collections have `.id` files
3. **File checks:**
   - All `.json` files parse as valid JSON
   - No unexpected files in collection directories
    - Normal check warns on non-JSON files.

**Return Value:**
```go
type CheckReport struct {
    TotalIssues int
    Issues      []Issue
}

type Issue struct {
    Type        string  // "missing_directory", "invalid_id_algorithm", etc.
    Description string  // Human-readable description
    Path        string  // Filesystem path (if applicable)
}
```

#### `Fix() → FixReport, error`

Repairs inconsistencies detected by `Check()`.

**Repairs Applied:**
1. Creates `db.json` if missing (version "1.0.0", empty collections)
2. Creates missing collection directories
3. Adds unlisted collection directories to metadata (defaults to `"cuid2"`)
4. Fixes invalid `id_algorithm` values to `"cuid2"`
5. Creates missing `.id` files with `"0"`
6. Removes unparseable files (after logging)

**Atomic Guarantee:**
- Each fix is atomic (temp file + rename)
- Failures in one fix don't affect others
- Report lists all successful fixes

**Return Value:**
```go
type FixReport struct {
    TotalFixed int
    Fixes      []Fix
}

type Fix struct {
    Type        string  // "created_directory", "created_id_file", etc.
    Description string  // Human-readable description
    Path        string  // Filesystem path
}
```

### RawCollection Interface

#### `Create(data) → id, error`

Creates a new record with auto-generated ID.

**Behavior:**
1. Generate ID using collection's `id_algorithm`
2. For sequential algorithms: atomically increment `.id` file
3. Write `{id}.json` atomically (temp file + rename)
4. Return generated ID

**Content Parameter:**
- Passed to ID generator function
- `cuid2`: ignores content
- `seq10`/`seq36`: currently ignore, reserved for future hash-based IDs
- Must be valid JSON if used with typed `Collection[T]`

**Concurrency:**
- Thread-safe (mutex-protected for sequential IDs)
- Concurrent Creates may generate non-consecutive IDs (acceptable)

#### `CreateWithID(id, data) → error`

Creates a record with explicitly specified ID.

**Behavior:**
1. Validate ID is alphanumeric
2. Check if `{id}.json` already exists → `ErrExists`
3. Write `{id}.json` atomically
4. Does NOT update `.id` file (manual ID, not part of sequence)

**Use Cases:**
- Foreign keys (referencing external IDs)
- User-chosen IDs (usernames, slugs)
- Migration/import (preserving original IDs)

#### `RawRead(id) → []byte, error`

Reads raw JSON bytes for a record.

**Behavior:**
- Read `{id}.json` file
- Return contents as-is (no validation)
- Return `ErrNotFound` if file doesn't exist

#### `Update(id, data) → error`

Replaces an existing record.

**Behavior:**
- Check `{id}.json` exists → `ErrNotFound` if missing
- Write new data atomically (temp file + rename)
- Old data completely replaced (no merging)

**Atomic Guarantee:**
- Readers see old OR new data, never partial writes
- Failure leaves old data intact

#### `Delete(id) → error`

Removes a record.

**Behavior:**
- Delete `{id}.json` file
- Return `ErrNotFound` if file doesn't exist
- Does NOT update `.id` file (sequence continues forward)

#### `Exists(id) → bool, error`

Checks if a record exists.

**Behavior:**
- Check if `{id}.json` file exists
- Return `true` if exists, `false` otherwise
- Return error only for filesystem errors (not for non-existence)

#### `IDs(start) → iter.Seq[string]`

Returns lazy iterator over record IDs in lexicographic order.

**Behavior:**
- Read collection directory entries
- Sort IDs lexicographically (text sort)
- If `start` empty: iterate all IDs from beginning
- If `start` provided: skip IDs < start, begin at first ID ≥ start (inclusive)
- Lazy evaluation: IDs read as iteration proceeds

**Lexicographic Ordering:**
- Text-based sort (e.g., "0000000002" < "0000000010")
- Zero-padding ensures numeric order = text order for seq10/seq36
- CUID2 IDs sorted alphabetically

**Iterator Properties:**
- Yields IDs one at a time
- No transaction isolation (concurrent writes may/may not appear)
- Empty iteration if no records exist

#### `RawItems(start) → iter.Seq2[string, []byte]`

Returns lazy iterator over (ID, data) pairs in lexicographic order.

**Behavior:**
- Same ordering/pagination as `IDs()`
- Yields both ID and raw JSON bytes
- Reads file contents lazily during iteration

**Use Cases:**
- Scanning entire collection
- Export/backup operations
- Filtered queries (client-side filtering)

### Typed Collection Interface

```go
type Collection[T any] struct {
    raw RawCollection
}
```

Wrapper around `RawCollection` providing type-safe JSON marshaling.

**Methods:**
- `Create(record T) → id, error`: Marshal to JSON, call `raw.Create()`
- `CreateWithID(id, record T) → error`: Marshal, call `raw.CreateWithID()`
- `Read(id) → T, error`: Call `raw.RawRead()`, unmarshal
- `Update(id, record T) → error`: Marshal, call `raw.Update()`
- `Delete(id) → error`: Delegate to `raw.Delete()`
- `Exists(id) → bool, error`: Delegate to `raw.Exists()`
- `Items(start) → iter.Seq2[string, T]`: Iterate, unmarshal each

**Type Constraints:**
- `T` must be JSON-serializable
- `T` should have `id` field (though not enforced by type system)

---

## Error Handling

### Standard Errors

**`ErrNotFound`**
- Record/collection doesn't exist
- Returned by: `RawRead()`, `Update()`, `Delete()`

**`ErrExists`**
- Record/collection already exists
- Returned by: `CreateWithID()`, `AddCollection()`

**`ErrClosed`**
- Store has been closed
- Returned by: all operations after `Close()`

**`ErrInvalidID`**
- ID contains non-alphanumeric characters
- Returned by: `CreateWithID()` with invalid ID

### Error Propagation

- Filesystem errors (permission denied, disk full) propagate as-is
- JSON unmarshal errors propagate with context
- Validation errors return descriptive messages

---

## Concurrency Model

### Thread Safety

**Store Level:**
- Multiple goroutines can safely call all Store methods
- Metadata updates are mutex-protected

**Collection Level:**
- Concurrent reads are always safe
- Concurrent writes to different IDs are safe
- Concurrent writes to same ID: last write wins (no locking)

**ID Generation:**
- `cuid2`: Stateless, inherently thread-safe
- `seq10`/`seq36`: Mutex-protected `.id` file access

### Atomic Writes

All writes use atomic temp-file-and-rename pattern:

1. Write data to `{path}.tmp.{random}`
2. Sync to disk
3. Rename to final path (atomic on POSIX)

**Guarantees:**
- Readers never see partial writes
- Process crash leaves old data OR new data, never corruption
- Rename atomicity relies on POSIX semantics (same filesystem)

---

## Consistency Guarantees

### Single-Record Operations

**Atomicity:**
- Each operation (Create, Update, Delete) is atomic
- No partial writes visible to readers

**Isolation:**
- No transaction isolation across operations
- Concurrent readers may see intermediate states during multi-operation workflows

**Durability:**
- Writes sync to disk before returning (configurable per backend)
- Process crash after successful return preserves data

### Multi-Record Operations

**No Transactions:**
- No multi-record ACID transactions
- Use application-level coordination if needed

**Eventual Consistency:**
- `Check()` may detect transient inconsistencies during writes
- `Fix()` repairs only persistent inconsistencies

---

## Pagination and Iteration

### Start-Based Pagination

Both `IDs(start)` and `RawItems(start)` support start-based pagination:

```go
// First page
for id := range coll.IDs("") {
    if count >= pageSize {
        lastID = id  // Save for next page
        break
    }
    // Process...
}

// Next page (start after last ID)
for id := range coll.IDs(nextID(lastID)) {
    // Process...
}

func nextID(id string) string {
    // Implementation-specific: increment ID lexicographically
}
```

**Properties:**
- Stable if no writes during iteration
- New records may/may not appear if written during iteration
- Deleted records skipped automatically

### Lexicographic Ordering

IDs sorted as strings, not numerically:

```
# seq10 (correct due to zero-padding):
0000000001, 0000000002, 0000000010, 0000000100

# Without padding (WRONG):
1, 10, 100, 2  # Lexicographic ≠ numeric

# cuid2 (alphabetic):
clxyz1..., clxyz2..., clxyz3...

# seq36 (correct due to zero-padding):
0000001, 0000002, 000000a, 000000z, 0000010
```

---

## Filesystem Requirements

### POSIX Compliance

Required for atomicity guarantees:
- Atomic rename within same filesystem
- Directory fsync support (optional but recommended)

### Permissions

- Store must have read/write access to `{root}` directory
- Newly created files inherit directory permissions
- No special permission handling (use umask)

### Path Constraints

- Collection names: alphanumeric, no `/` or `.`
- IDs: alphanumeric (a-z, A-Z, 0-9)
- Max path length subject to OS limits (typically 255 bytes per component)

### Filesystem Types

**Supported:**
- Local filesystems (ext4, APFS, NTFS)
- Network filesystems with POSIX semantics (NFS v4+)

**Not Recommended:**
- FAT32 (no atomic rename)
- Network filesystems without atomic rename (SMB, older NFS)

---

## Migration and Versioning

### Format Version

Current version: `"1.0.0"` in `db.json`

**Future Changes:**
- Minor version: backward-compatible additions
- Major version: breaking changes (require migration)

### Backward Compatibility

**Guaranteed:**
- Can always read older format versions
- `Fix()` can upgrade to newer formats

**Migration Path:**
1. `Check()` detects old version
2. `Fix()` upgrades in-place
3. Version updated in `db.json`

---

## Future Backend Compatibility

Backends other than disk (in-memory, cloud, embedded DB) should maintain these semantics:

**Must Preserve:**
- ID generation behavior (same IDs for same sequence)
- Atomicity of single-record operations
- Lexicographic ordering in iteration
- Error return values for same conditions

**May Alter:**
- Physical storage format (binary, columnar, etc.)
- Metadata location (embedded, separate service)
- Concurrency mechanisms (transactions, MVCC)
- Durability guarantees (configurable sync)

**Example:** An in-memory backend might:
- Store records in hash tables (not files)
- Use atomic reference swaps (not rename)
- Still generate seq10 IDs as "0000000001", "0000000002", etc.
- Still iterate in lexicographic order
- Still return `ErrNotFound` for missing records

---

## Appendix: Example Workflows

### Creating a Collection

```go
store, _ := disk.New("/path/to/db")
defer store.Close()

// Create collection with sequential IDs
store.AddCollection("users", "seq10")

// Add records
users := store.RawCollection("users")
id1, _ := users.Create([]byte(`{"name":"Alice"}`))  // "0000000001"
id2, _ := users.Create([]byte(`{"name":"Bob"}`))    // "0000000002"
```

### Using Custom IDs

```go
products := store.RawCollection("products")

// Use SKU as ID
products.CreateWithID("SKU-12345", []byte(`{"name":"Widget","price":29.99}`))

// Foreign key reference
orders := store.RawCollection("orders")
orders.Create([]byte(`{"product_id":"SKU-12345","quantity":5}`))
```

### Consistency Checking

```go
// Check for issues
report, _ := store.Check()
if report.TotalIssues > 0 {
    for _, issue := range report.Issues {
        log.Printf("%s: %s (%s)", issue.Type, issue.Description, issue.Path)
    }

    // Repair automatically
    fixReport, _ := store.Fix()
    log.Printf("Fixed %d issues", fixReport.TotalFixed)
}
```

### Pagination

```go
coll := store.RawCollection("logs")

pageSize := 100
start := ""

for {
    count := 0
    var lastID string

    for id, data := range coll.RawItems(start) {
        // Process record...
        count++
        lastID = id

        if count >= pageSize {
            break
        }
    }

    if count < pageSize {
        break  // Last page
    }

    // Next page starts after last ID
    start = incrementID(lastID)
}
```

---

**End of Specification**
