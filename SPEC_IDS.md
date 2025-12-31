# SimpleDB ID Specification

## Overview

IDs uniquely identify records within a collection. Implementations and drivers choose their ID generation strategy.

## Requirements

- IDs must be unique within a collection
- Character set: alphanumeric plus `_-` (a-z, A-Z, 0-9, underscore, hyphen)
- ID schemes MAY reserve characters (e.g. `-`) for prefix-based grouping schemes
- IDs MAY have other characteristics (sorting, abbreviations) documented by backends

## Recommended Default: CUID2

**CUID2** [Collision-resistant Unique IDentifier](https://github.com/paralleldrive/cuid2?tab=readme-ov-file#the-contenders) is the recommended default. If an easier no-lib version is required,


## Alternative Algorithms

Implementations may offer these or others:

- **cuid2** (~25 chars) - Secure, monotonic, stateless. Pro: audited, widely supported. Con: requires library.
- **uuid4** (36 chars) - Standard format. Pro: ubiquitous. Con: longer, less collision-resistant than cuid2.
- **ulid** (26 chars) - Time-sorted, Crockford Base32. Pro: sortable. Con: leaks timestamp.
- **[tid62](ideas/tid62.md)** (17 chars) - Compact time-sorted base62. Pro: short, no deps. Con: unaudited, leaks time.
- **seq10** (10 chars) - Zero-padded decimal. Pro: human-readable. Con: predictable, needs state file.
- **seq36** (7 chars) - Zero-padded base-36. Pro: very compact. Con: predictable, needs state file.

## Custom IDs

`Store(id, data)` accepts any valid ID string:
- Foreign keys from external systems
- User-chosen identifiers (usernames, slugs)
- Prefixed IDs: `user-clh3am8kw...`, `order-cm5f9k2x1...`

## Prefixed ID Pattern

Use `-` as separator for ID portions:

```
# Human-readable namespaces within a collection
{prefix}-{cuid2}
user-clh3am8kw0000g3s0h8d6a9xq
order-cm5f9k2x10001jh08qrv5z8gt

# DynamoDB-style composite keys
{userid}-{workoutid}-{exerciseid}
```

## Driver Flexibility

Drivers may:
- Choose default ID algorithm
- Offer configuration options
- Support prefix generation helpers

Constraint: avoid ID conflicts within a collection.
