// Package simpledb defines the SimpleDB interface specification.
// This file is for reference only - implementations provide the actual code.
package simpledb

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"iter"
	"math/big"
	"time"
)

// Standard errors
var (
	ErrNotFound  error // Record doesn't exist
	ErrExists    error // Record/collection already exists
	ErrClosed    error // Store has been closed
	ErrInvalidID error // ID contains invalid characters
)

// Store is the main database interface.
type Store interface {
	// RawCollection returns a collection handle. Does not check filesystem.
	RawCollection(name string) RawCollection

	// AddCollection registers a new collection with the specified ID algorithm.
	AddCollection(ctx context.Context, name, idAlgorithm string) error

	// Collections returns the list of registered collection names.
	Collections(ctx context.Context) ([]string, error)

	// Close releases resources.
	Close() error

	// Check validates consistency between metadata and filesystem.
	Check(ctx context.Context) (CheckReport, error)

	// Fix repairs inconsistencies detected by Check.
	Fix(ctx context.Context) (FixReport, error)
}

// RawCollection provides byte-level record operations.
type RawCollection interface {
	// Create stores data with an auto-generated ID. Returns the new ID.
	Create(ctx context.Context, data []byte) (string, error)

	// CreateWithID stores data with the specified ID.
	CreateWithID(ctx context.Context, id string, data []byte) error

	// RawRead returns the raw JSON bytes for a record.
	RawRead(ctx context.Context, id string) ([]byte, error)

	// Update replaces an existing record.
	Update(ctx context.Context, id string, data []byte) error

	// Delete removes a record.
	Delete(ctx context.Context, id string) error

	// Exists checks if a record exists.
	Exists(ctx context.Context, id string) (bool, error)

	// IDs returns an iterator over record IDs in lexicographic order.
	IDs(ctx context.Context, start string) iter.Seq[string]

	// RawItems returns an iterator over (id, data) pairs in lexicographic order.
	RawItems(ctx context.Context, start string) iter.Seq2[string, []byte]
}

// Collection provides typed record operations with JSON marshaling.
type Collection[T any] interface {
	Create(ctx context.Context, record T) (string, error)
	CreateWithID(ctx context.Context, id string, record T) error
	Read(ctx context.Context, id string) (T, error)
	Update(ctx context.Context, id string, record T) error
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	Items(ctx context.Context, start string) iter.Seq2[string, T]
}

// CheckReport contains results from Store.Check().
type CheckReport struct {
	TotalIssues int
	Issues      []Issue
}

type Issue struct {
	Type        string
	Description string
	Path        string
}

// FixReport contains results from Store.Fix().
type FixReport struct {
	TotalFixed int
	Fixes      []Fix
}

type Fix struct {
	Type        string
	Description string
	Path        string
}

// -----------------------------------------------------------------------------
// tid62: Time ID in Base62
// -----------------------------------------------------------------------------

const (
	// Epoch is 2020-01-01 00:00:00 UTC
	tid62Epoch    = 1577836800
	tid62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	tid62Length   = 17
)

// NewTID62 generates a new time-based ID in base62.
// Structure: 4 bytes timestamp (seconds since epoch) + 8 bytes random.
// Result: 17 character base62 string.
func NewTID62() string {
	var buf [12]byte

	// 4 bytes: seconds since epoch (big-endian for sort order)
	ts := uint32(time.Now().Unix() - tid62Epoch)
	binary.BigEndian.PutUint32(buf[0:4], ts)

	// 8 bytes: crypto random
	rand.Read(buf[4:12])

	return encodeBase62(buf[:])
}

// NewTID62WithPrefix generates a prefixed ID: "{prefix}-{tid62}"
func NewTID62WithPrefix(prefix string) string {
	return prefix + "-" + NewTID62()
}

func encodeBase62(data []byte) string {
	// Convert bytes to big.Int
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
