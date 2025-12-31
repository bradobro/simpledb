# SimpleDB ID Specification

## Overview

IDs uniquely identify records within a collection. Implementations and drivers choose their ID generation strategy.

## Requirements

- IDs must be unique within a collection
- Character set: alphanumeric plus `_-` (a-z, A-Z, 0-9, underscore, hyphen)
- ID schemes MAY reserve characters (e.g. `-`) for prefix-based grouping schemes (a common pattern in DynamoDB), e.g. <userid>-<classid>
- ID MAY have other characteristics, such as predictable sorting orders, abbreviations, etc. documented by the backends.

## Recommended Default: tid62

**tid62** = Time ID in Base62

**Structure:** 12 bytes
- 4 bytes: seconds since 2020-01-01 00:00:00 UTC (big-endian)
- 8 bytes: cryptographically random

**Encoding:** Base62 (`0-9A-Za-z`), 17 characters fixed width, zero-padded

**Properties:**
- Time-sortable (lexicographic ≈ chronological)
- 64 bits random per second (~10^19 possible, collision-resistant)
- No external deps (uses stdlib BigInt/big.Int)
- Leaves `-` available for prefixed IDs
- Epoch runway: 2020 to ~2156

**Example:**
```
Timestamp: 2025-12-30 12:00:00 UTC (189388800 seconds from epoch)
Random:    8 bytes from crypto/rand
Result:    0B4kZ7mNpQrStUvWx  (17 chars)
```

See [api.go](api.go) and [api.ts](api.ts) for reference implementations.

## Base62 Alphabet

```
0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
```

Position 0-9 = digits, 10-35 = uppercase, 36-61 = lowercase.

## Alternative Algorithms

Implementations may offer these or others:

| Algorithm | Length | Properties |
|-----------|--------|------------|
| uuid464 | 22| base64url encoded UUID64|
| tid62 | 17 | Time-sorted, recommended default |
| seq10 | 10 | Zero-padded ordinal decimal, requires state |
| seq36 | 7 | Zero-padded ordinal base-36, requires state |

## Custom IDs

`CreateWithID(id, data)` accepts any valid ID string:
- Foreign keys from external systems
- User-chosen identifiers (usernames, slugs)
- Prefixed IDs: `user-0ABC123...`, `order-0XYZ789...`

## Prefixed ID Pattern

Use `-` as separator id portions:

```
# to have human-readable namespaces within a collection
{prefix}-{tid62}
user-0B4kZ7mNpQrStUvWx
order-0B4kZ7mNpQrStUvWy

# OR for DynamoDB-style prefix listing/scoping
{userid}-{workoutid}-{exerciseid}
```

## Driver Flexibility

Drivers may:
- Choose default ID algorithm
- Offer configuration options
- Support prefix generation helpers

Constraint: avoid ID conflicts within a collection.
