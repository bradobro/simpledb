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
  constructor(message = "Store is closed") {
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

// -----------------------------------------------------------------------------
// Report Types
// -----------------------------------------------------------------------------

export interface Issue {
  type: string;
  description: string;
  path?: string;
}

export interface CheckReport {
  totalIssues: number;
  issues: Issue[];
}

export interface Fix {
  type: string;
  description: string;
  path?: string;
}

export interface FixReport {
  totalFixed: number;
  fixes: Fix[];
}

// -----------------------------------------------------------------------------
// Store Interface
// -----------------------------------------------------------------------------

export interface Store {
  rawCollection(name: string): RawCollection;
  addCollection(name: string, idAlgorithm: string): Promise<void>;
  collections(): Promise<string[]>;
  close(): Promise<void>;
  check(): Promise<CheckReport>;
  fix(): Promise<FixReport>;
}

// -----------------------------------------------------------------------------
// Collection Interfaces
// -----------------------------------------------------------------------------

export interface RawCollection {
  create(data: Uint8Array): Promise<string>;
  createWithID(id: string, data: Uint8Array): Promise<void>;
  rawRead(id: string): Promise<Uint8Array>;
  update(id: string, data: Uint8Array): Promise<void>;
  delete(id: string): Promise<void>;
  exists(id: string): Promise<boolean>;
  ids(start?: string): AsyncIterable<string>;
  rawItems(start?: string): AsyncIterable<[string, Uint8Array]>;
}

export interface Collection<T> {
  create(record: T): Promise<string>;
  createWithID(id: string, record: T): Promise<void>;
  read(id: string): Promise<T>;
  update(id: string, record: T): Promise<void>;
  delete(id: string): Promise<void>;
  exists(id: string): Promise<boolean>;
  items(start?: string): AsyncIterable<[string, T]>;
}

// -----------------------------------------------------------------------------
// tid62: Time ID in Base62
// -----------------------------------------------------------------------------

const TID62_EPOCH = 1577836800; // 2020-01-01 00:00:00 UTC
const TID62_ALPHABET =
  "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz";
const TID62_LENGTH = 17;

/**
 * Generate a new time-based ID in base62.
 * Structure: 4 bytes timestamp (seconds since epoch) + 8 bytes random.
 * Result: 17 character base62 string.
 */
export function newTID62(): string {
  const buf = new Uint8Array(12);

  // 4 bytes: seconds since epoch (big-endian for sort order)
  const ts = Math.floor(Date.now() / 1000) - TID62_EPOCH;
  new DataView(buf.buffer).setUint32(0, ts, false);

  // 8 bytes: crypto random
  crypto.getRandomValues(buf.subarray(4));

  return encodeBase62(buf);
}

/**
 * Generate a prefixed ID: "{prefix}-{tid62}"
 */
export function newTID62WithPrefix(prefix: string): string {
  return `${prefix}-${newTID62()}`;
}

function encodeBase62(data: Uint8Array): string {
  // Convert bytes to BigInt
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

  // Pad to fixed length
  while (result.length < TID62_LENGTH) {
    result.push("0");
  }

  // Reverse for big-endian order
  return result.reverse().join("");
}
