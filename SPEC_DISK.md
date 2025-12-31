# SimpleDB Disk Storage Specification

## Directory Structure

```
{root}/
├── db.json                     # Database metadata
├── {collection}/               # Collection directory
│   ├── .id                     # Sequential ID state (optional)
│   ├── *.???                   # Tolerate other extensions for extensions.
│   └── {id}.json               # Record files
└── {schema}.schema.json        # Optional JSONSchema docs
```

## Metadata: db.json

```json
{
  "version": "1.0.0",  // database on-disk format version TODO: change to formatVersion
  "schemaVersion": 1, // collections schema version (can be used a migration level, etc.)
  "collections": {
    "users": { "id_algorithm": "uuid464" }
  }
}
```

## ID Storage: Filename IS the ID

**Records do not store IDs in JSON content.** The filename is the single source of truth.

```
people/
└── abc123.json     ← ID is "abc123"
    {               ← JSON has no "id" field
      "name": "Alice",
      "age": 30
    }
```

**Rationale:**
- Impossible to desync filename/content
- Smaller files
- Simpler validation

**Implementation:**
- Write: strip `id` field before writing
- Read: inject `id` from filename into returned object
- Update: validate id parameter matches filename, then strip

**Backward compatibility:** Old files with `id` fields still work; `id` is validated and stripped on write.

This pattern is recommended for similar storage media (cloud object stores, embedded KV) to avoid ID conflicts.

## Sequential ID State: .id

Used by sequential algorithms (seq10, seq36).

- Single-line text file with last generated ID
- Created with "0" if missing
- Updated atomically on `Create()`
- Not present for stateless algorithms

## Atomic Writes

All writes use temp-file-and-rename:
1. Write to `{path}.tmp.{random}`
2. Sync to disk
3. Rename to final path (atomic on POSIX)

## Filesystem Requirements

- POSIX atomic rename (ext4, APFS, NTFS, NFS v4+)
- Avoid FAT32 or older network filesystems
