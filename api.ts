/**
 * SimpleDB TypeScript Interface Specification
 * This file is for reference only - implementations provide the actual code.
 */

// -----------------------------------------------------------------------------
// Errors
// -----------------------------------------------------------------------------

export class NotFoundError extends Error {
  constructor(message = "Record not found") {
    super(message);
    this.name = "NotFoundError";
  }
}

export class ExistsError extends Error {
  constructor(message = "Record already exists") {
    super(message);
    this.name = "ExistsError";
  }
}

export class ClosedError extends Error {
  constructor(message = "SimpleDB is closed") {
    super(message);
    this.name = "ClosedError";
  }
}

export class InvalidIDError extends Error {
  constructor(message = "Invalid ID format") {
    super(message);
    this.name = "InvalidIDError";
  }
}

export class IOError extends Error {
  constructor(message = "I/O error") {
    super(message);
    this.name = "IOError";
  }
}

// -----------------------------------------------------------------------------
// Report Types
// -----------------------------------------------------------------------------

export interface Problem {
  type: string;
  description: string;
  path?: string;
  fixed: boolean; // always false if check(fix=false)
}

export interface CheckReport {
  totalIssues: number;
  issues: Problem[];
}

// -----------------------------------------------------------------------------
// SimpleDB Interface
// -----------------------------------------------------------------------------

export interface SimpleDB {
  /** Get returns a collection handle. Does not check filesystem. */
  get(name: string): ByteCollection;

  /** Create registers a new collection with the specified ID algorithm. */
  create(name: string, idAlgorithm: string): Promise<void>;

  /** List returns the list of registered collection names. */
  list(): Promise<string[]>;

  /** Close releases resources. */
  close(): Promise<void>;

  /**
   * Check validates consistency of the store and its collections.
   * If fix is true, attempts to repair problems.
   * Args: "meta", "fast", "slow", "all"
   */
  check(fix: boolean, ...args: string[]): Promise<CheckReport>;
}

// -----------------------------------------------------------------------------
// Byte Collection Interfaces
// -----------------------------------------------------------------------------

/** ByteStorer provides low-level store/exists operations. */
export interface ByteStorer {
  /** Store stores a record with the specified ID. Errors if ID exists. */
  store(id: string, data: Uint8Array): Promise<void>;

  /** Exists returns true if a record exists. */
  exists(id: string): Promise<boolean>;

  /** Check tests consistency of the collection. */
  check(fix: boolean, ...args: string[]): Promise<CheckReport>;
}

/** ByteUpdater provides CRUD operations. */
export interface ByteUpdater {
  /** Create stores data with an auto-generated ID. Returns the new ID. */
  create(data: Uint8Array): Promise<string>;

  /** Read returns the raw bytes for a record. */
  read(id: string): Promise<Uint8Array>;

  /** Update replaces an existing record. */
  update(id: string, data: Uint8Array): Promise<void>;

  /** Delete removes a record. */
  delete(id: string): Promise<void>;
}

/** ByteLister provides iteration methods. */
export interface ByteLister {
  /** Items returns an iterator over (id, data) pairs in lexicographic order. */
  items(start?: string): AsyncIterable<[string, Uint8Array]>;

  /** IDs returns an iterator over record IDs in lexicographic order. */
  ids(start?: string): AsyncIterable<string>;
}

/** ByteCollection combines all byte-level operations. */
export interface ByteCollection extends ByteStorer, ByteUpdater, ByteLister {}

// -----------------------------------------------------------------------------
// Typed Collection Interfaces
// -----------------------------------------------------------------------------

/** Storer provides typed store/exists operations. */
export interface Storer<T> {
  store(id: string, record: T): Promise<void>;
  exists(id: string): Promise<boolean>;
}

/** Updater provides typed CRUD operations. */
export interface Updater<T> {
  create(record: T): Promise<string>;
  read(id: string): Promise<T>;
  update(id: string, record: T): Promise<void>;
  delete(id: string): Promise<void>;
}

/** Lister provides typed iteration. */
export interface Lister<T> {
  items(start?: string): AsyncIterable<[string, T]>;
  ids(start?: string): AsyncIterable<string>;
}

/** Collection combines all typed operations. */
export interface Collection<T> extends Storer<T>, Updater<T>, Lister<T> {}
