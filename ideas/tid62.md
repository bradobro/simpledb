# TID62: Time ID in Base62

**Status:** Idea / Custom Design (not a standard)

## Overview

TID62 is a custom ID format designed for SimpleDB during early development. It was later replaced by CUID2 as the recommended default since CUID2 is an established standard with library support.

## Design Goals

- Time-sortable (rough chronological ordering)
- Base62 encoding (leaves `-` available for prefixed IDs)
- No external dependencies (stdlib only)
- Collision-resistant

## Structure

12 bytes total:
- 4 bytes: seconds since 2020-01-01 00:00:00 UTC (big-endian)
- 8 bytes: cryptographically random

Encoded as 17 Base62 characters, zero-padded.

## Properties

- **Length:** 17 characters
- **Sortable:** Lexicographic ≈ chronological (within same second)
- **Collision resistance:** 64 bits random per second (~10^19 possible)
- **Epoch runway:** 2020 to ~2156 (32-bit seconds)
- **Character set:** `0-9A-Za-z` (62 chars)

## Influences

- **ULID:** millisecond timestamp + random, Crockford Base32, 26 chars
- **KSUID:** seconds timestamp + random, base62, 27 chars (Segment)
- **TypeID:** UUIDv7 with type prefix, base32

TID62 differs by using a 2020 epoch (saves bits) and 17-char output.

## Base62 Alphabet

```
0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
```

Position 0-9 = digits, 10-35 = uppercase, 36-61 = lowercase.

---

## Example Implementation: Go

```go
package main

import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"time"
)

const (
	tid62Epoch    = 1577836800 // 2020-01-01 00:00:00 UTC
	tid62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	tid62Length   = 17
)

func NewTID62() string {
	var buf [12]byte

	// 4 bytes: seconds since epoch (big-endian for sort order)
	ts := uint32(time.Now().Unix() - tid62Epoch)
	binary.BigEndian.PutUint32(buf[0:4], ts)

	// 8 bytes: crypto random
	rand.Read(buf[4:12])

	return encodeBase62(buf[:])
}

func NewTID62WithPrefix(prefix string) string {
	return prefix + "-" + NewTID62()
}

func encodeBase62(data []byte) string {
	n := new(big.Int).SetBytes(data)

	base := big.NewInt(62)
	zero := big.NewInt(0)
	mod := new(big.Int)

	result := make([]byte, 0, tid62Length)

	for n.Cmp(zero) > 0 {
		n.DivMod(n, base, mod)
		result = append(result, tid62Alphabet[mod.Int64()])
	}

	// Pad to fixed length
	for len(result) < tid62Length {
		result = append(result, '0')
	}

	// Reverse for big-endian order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}
```

---

## Example Implementation: TypeScript

```typescript
const TID62_EPOCH = 1577836800; // 2020-01-01 00:00:00 UTC
const TID62_ALPHABET =
  "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
const TID62_LENGTH = 17;

export function newTID62(): string {
  const buf = new Uint8Array(12);

  // 4 bytes: seconds since epoch (big-endian for sort order)
  const ts = Math.floor(Date.now() / 1000) - TID62_EPOCH;
  new DataView(buf.buffer).setUint32(0, ts, false);

  // 8 bytes: crypto random
  crypto.getRandomValues(buf.subarray(4));

  return encodeBase62(buf);
}

export function newTID62WithPrefix(prefix: string): string {
  return `${prefix}-${newTID62()}`;
}

function encodeBase62(data: Uint8Array): string {
  let n = 0n;
  for (const byte of data) {
    n = (n << 8n) | BigInt(byte);
  }

  const base = 62n;
  const result: string[] = [];

  while (n > 0n) {
    result.push(TID62_ALPHABET[Number(n % base)]);
    n = n / base;
  }

  while (result.length < TID62_LENGTH) {
    result.push("0");
  }

  return result.reverse().join("");
}
```

---

## Why Not Used

CUID2 was chosen as the recommended default because:
- Established standard with active maintenance
- Libraries available for Go, TypeScript, and many other languages
- Well-documented security properties
- Community adoption and review

TID62 remains documented here as an alternative for projects that want:
- Smaller IDs (17 vs ~25 chars)
- No external dependencies
- Custom epoch for longer runway
