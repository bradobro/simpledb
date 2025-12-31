# Example SimpleDB Database

This directory demonstrates the **actual on-disk layout** of a SimpleDB database using the disk backend.

## Purpose

This is **not** just test fixtures - it's a complete, valid SimpleDB database that:
- Shows exactly how data is organized on disk
- Can be examined by developers to understand the format
- Serves as a reference implementation of the storage structure
- Is used by tests as a source of sample data (tests copy portions of it)

## Structure

```
db/
├── db.json                    # Metadata: version + registered collections
├── house.schema.json          # JSONSchema documentation (optional)
├── person.schema.json         # JSONSchema documentation (optional)
├── order.schema.json          # JSONSchema documentation (optional)
├── vehicle.schema.json        # JSONSchema documentation (optional)
├── houses/                    # Collection: cuid2 IDs (~25 chars)
│   ├── clh3am8kw0000g3s0h8d6a9xq.json
│   └── cm5f9k2x10001jh08qrv5z8gt.json
├── people/                    # Collection: seq10 IDs (10-digit decimal)
│   ├── .id                    # Latest ID state (for sequential generators)
│   ├── 0000000001.json
│   ├── 0000000002.json
│   └── 0000000003.json
├── orders/                    # Collection: seq36 IDs (7-digit base-36)
│   ├── .id                    # Latest ID state
│   ├── 0000001.json           # Orders 1-9 use digits
│   ├── ...
│   └── 000000c.json           # Orders 10+ use letters (a=10, b=11, c=12)
└── vehicles/                  # Collection: uuid464 IDs (base64url-encoded UUID4)
    ├── VQ6EAOKbQdSnFkRmVUQAAQ.json
    ├── a6e4EJ2tEdGAtADAT9QwyA.json
    └── 9HrBBYjMQ3KlZw4CssPUeQ.json
```

## Key Files

### db.json
Stores metadata about the database:
- `version`: Backend format version (currently "1.0.0")
- `collections`: Map of collection names to their configuration
  - `id_algorithm`: ID generation algorithm ("cuid2", "uuid464", "seq10", "seq36", etc.)

### .id Files (Sequential Algorithms Only)
For seq10 and seq36 collections, a `.id` file tracks the latest generated ID:
- One-line text file containing the current sequence number
- Atomically updated on each Create() call
- Not used by stateless algorithms (cuid2, uuid464) which generate IDs without state

### Collections
Each collection is a **directory** containing JSON files:
- Directory name = collection name
- Each file = one record
- Filename = record ID + `.json` extension
- File contents = JSON record data (must include `"id"` field)

### Schema Files (Optional)
JSONSchema files at the root document the expected structure of records.
These are **documentation only** - SimpleDB does not enforce schemas at runtime.

## Design Philosophy

This structure allows:
- **Human inspection**: Easy to browse and debug with standard tools
- **Migration**: Simple to export/import data
- **Transparency**: No proprietary binary formats
- **Resilience**: Deleting a collection directory loses data but preserves metadata
- **Flexibility**: Any JSON-serializable data works

## Maintenance

Keep this directory **synchronized with the current API**:
- Update `db.json` if metadata structure changes
- Add collections here when adding new test fixtures
- Keep schema files current with record structures
