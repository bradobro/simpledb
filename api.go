// Package simpledb defines the SimpleDB interface specification.
// This file is for reference only - implementations provide the actual code.
package simpledb

import (
	"context"

	"iter"
)

// Standard errors
var (
	ErrClosed    error // SimpleDB has been closed
	ErrExists    error // Record/collection already exists
	ErrInvalidID error // ID contains invalid characters
	ErrIO        error // Other I/O error
	ErrNotFound  error // Record doesn't exist
)

// SimpleDB is the main database interface.
type SimpleDB interface {
	// Get returns a collection handle. Does not check filesystem.
	Get(name string) ByteCollection
	// Create registers a new collection with the specified ID algorithm.
	Create(ctx context.Context, name, idAlgorithm string) error
	// List returns the list of registered collection names.
	List(ctx context.Context) ([]string, error)
	// Close releases resources.
	Close() error
	// Check validates consistency of the store and it's collections If fix is
	// true it attempts to repair the problems and reports how successful it was
	// fixing each problem. Otherwise, it only report problems. Optional args
	// allow ways to limit the check in driver-specific ways. But these guidelines
	// should be followed:
	//   - args names are positive, ADDING a check (they shouldn't begin with
	// 		 "no"), thus...
	//   - args = [] (empty) do no check, just invoke the checking framework
	//   - args = ["meta"], check everything about the central store, including
	// 			 - metadata
	// 			 - collection names are consistent with collection spaces
	// 			   (directories), etc.
	//   - args = ["slow"], run all available checks on each collection
	//   - args = ["fast"], run all available checks on each collection EXCEPT
	//     those that involve reading the data of every record. They MAY examine
	//     file names or the equivalent if a store has them.
	//   - args = ["all"], do all available checks, sa me as ["meta","slow"]
	Check(ctx context.Context, fix bool, args ...string) (CheckReport, error)
}

// ByteCollection is a database "table" of untyped bytes.
type ByteCollection interface {
	ByteStorer
	ByteUpdater
	ByteLister
}

// ByteStorer is the lowest-level storage interface of a byte collection.
// Given an id, it can tell you if it has bytes stored under that id
// Given bytes as well, it can store those bytes under the id
type ByteStorer interface {
	// Store stores a record with the specified ID.
	// If id already exists, it errors.
	Store(ctx context.Context, id string, data []byte) error
	// Exists returns true if a record exists.
	Exists(ctx context.Context, id string) (bool, error)
	// Check tests consistency of the collection. SimpleDB.Check() will usually
	// delegate collection checks to this method. It MUST be safe to run
	// checks on multiple collections in a store simultaneously; this method
	// must not alter or contend for anything outside its purview
	Check(ctx context.Context, fix bool, args ...string) (CheckReport, error)
}

// ByteUpdater holds the basic editing interface of a byte collection.
type ByteUpdater interface {
	// Create stores data with an auto-generated ID. Returns the new ID.
	Create(ctx context.Context, data []byte) (string, error)
	// Read returns the raw bytes for a record.
	Read(ctx context.Context, id string) ([]byte, error)
	// Update replaces an existing record.
	Update(ctx context.Context, id string, data []byte) error
	// Delete removes a record.
	Delete(ctx context.Context, id string) error
}

// ByteLister provides methods for traversing a collection of records of bytes
type ByteLister interface {
	// Items returns an iterator over (id, data) pairs in lexicographic order.
	Items(ctx context.Context, start string) iter.Seq2[string, []byte]
	// IDs returns an iterator over record IDs in lexicographic order.
	IDs(ctx context.Context, start string) iter.Seq[string]
}

// Storer provides typed store/exists operations.
type Storer[T any] interface {
	Store(ctx context.Context, id string, record T) error
	Exists(ctx context.Context, id string) (bool, error)
}

// Updater provides typed CRUD operations.
type Updater[T any] interface {
	Create(ctx context.Context, record T) (string, error)
	Read(ctx context.Context, id string) (T, error)
	Update(ctx context.Context, id string, record T) error
	Delete(ctx context.Context, id string) error
}

// Lister provides typed iteration.
type Lister[T any] interface {
	Items(ctx context.Context, start string) iter.Seq2[string, T]
	IDs(ctx context.Context, start string) iter.Seq[string]
}

// Collection combines all typed operations.
type Collection[T any] interface {
	Storer[T]
	Updater[T]
	Lister[T]
}

// CheckReport contains results from SimpleDB.Check().
type CheckReport struct {
	TotalIssues int
	Issues      []Problem
}

// Problem represents an inconsistency in a SimpleDB
type Problem struct {
	Type        string
	Description string
	Path        string
	Fixed       bool // always false if Check
}
