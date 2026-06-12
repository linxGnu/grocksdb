# grocksdb Maintenance Guide — Upgrading the bundled RocksDB

This document describes how to bump **grocksdb** to a new **RocksDB** release. It
is the repeatable recipe used for every `Adapt RocksDB X.Y.Z` commit (e.g.
`5fc3d82` → 10.9.1, `632e53a` → 10.10.1, the 10.10.1 → 11.1.1 bump).

grocksdb is a thin **CGO** binding: the Go code calls RocksDB's C API declared in
RocksDB's `include/rocksdb/c.h`. To track a new RocksDB version you re-vendor that
header and reconcile the Go/C glue with any API that was **added**, **removed**, or
whose **signature changed**.

---

## TL;DR checklist

1. **Bump the version** in `build.sh` (`rocksdb_version="X.Y.Z"`).
2. **Diff RocksDB's `include/rocksdb/c.h`** between the old and new tag.
3. **Apply that diff to the repo-root `c.h`** (the vendored copy grocksdb compiles against).
4. **Reconcile the Go/C glue** for any *removed* or *signature-changed* C symbols
   that this binding actually uses (`grocksdb.c`, `grocksdb.h`, `*.go`).
5. **Update `README.md`** "API Support" with newly added but unwrapped C APIs.
6. **Build & test** (`make test`).
7. Commit as `CHORE Adapt RocksDB X.Y.Z`.

---

## Key files

| File | Role |
|------|------|
| `build.sh` | Builds RocksDB + compression deps from source. Holds `rocksdb_version`. |
| `c.h` (repo root) | **Vendored copy** of RocksDB's `include/rocksdb/c.h`. This is the source of truth the Go code binds to. Must match the target tag. |
| `grocksdb.c` / `grocksdb.h` | Hand-written C shims that wrap callback-based C constructors (comparator, merge operator, slice transform, compaction filter). Sensitive to C API signature changes. |
| `*.go` | The Go wrappers calling `C.rocksdb_*`. |
| `non_builtin.go` / `non_builtin_clean_link.go` | CGO `LDFLAGS` for the `grocksdb_no_link` / `grocksdb_clean_link` build tags. |
| `Makefile` | `make libs` runs `build.sh`; `make test` runs the tagged test suite. |
| `README.md` | "API Support" section listing C APIs **not** wrapped in Go. |

> `deps/` is a **gitignored build artifact** (a RocksDB checkout). Never hand-edit
> `deps/rocksdb/include/rocksdb/c.h` — it is regenerated when you rebuild against
> the new tag. Only the **repo-root** `c.h` matters for the binding.

---

## Step 1 — Bump `build.sh`

```diff
-rocksdb_version="10.10.1"
+rocksdb_version="11.1.1"
```

While you're here, sanity-check the compression library versions pinned above it
(`snappy_version`, `zlib_version`, `lz4_version`, `zstd_version`). Bump them only if
the new RocksDB requires it — usually they can stay. Note the deliberate
`-DPORTABLE=1` (kept for CI runner CPU compatibility) and `-DWITH_BZ2=OFF`.

---

## Step 2 — Diff RocksDB's `c.h` between tags

The only RocksDB file that affects the binding is `include/rocksdb/c.h`. Get a
focused diff of just that header between the current and target tag.

### Option A — from a RocksDB clone

```bash
git clone https://github.com/facebook/rocksdb.git
cd rocksdb
git fetch --tags
# Old tag = what build.sh currently pins; new tag = target.
git diff v10.10.1 v11.1.1 -- include/rocksdb/c.h
```

### Option B — from a full release diff you already have

If you have a full source diff (e.g. `10-11.diff`), the header is one slice of it.
Locate its hunk boundaries and read just that range:

```bash
# Find where the c.h diff starts and where the NEXT file diff starts.
grep -n "^diff --git" 10-11.diff | grep -E "include/rocksdb/(c\.h|cache\.h)"
# e.g. c.h begins at line 33691, the following file (cache.h) at 33964
# → the c.h diff is lines 33691..33963.
sed -n '33691,33963p' 10-11.diff
```

### Option C — GitHub compare URL

`https://github.com/facebook/rocksdb/compare/v10.10.1...v11.1.1` then open
`include/rocksdb/c.h` in the Files Changed tab.

---

## Step 3 — Classify every `c.h` change

Read the diff and bucket each change. The bucket decides how much work it is.

| Change type | What to do |
|-------------|------------|
| **Added** typedef / function / enum | Add it verbatim to repo-root `c.h`. No Go work required (it just becomes available; wrap it later if desired). |
| **Removed** function / enum value | Delete it from repo-root `c.h`. **Then** check whether any Go/C code in this repo used it (Step 4) — if so, that code must change too. |
| **Signature changed** (params added/removed/reordered) | Update the declaration in repo-root `c.h`, **and** update every caller in `grocksdb.c` / `*.go`. |
| Cosmetic (whitespace/indent) | Match it to keep future diffs clean, but harmless. |
| Counter/enum bumps (e.g. `rocksdb_total_metric_count = 80 → 85`) | Update the value so the enum stays in sync. |

Apply additions/removals to the **repo-root `c.h`** to mirror the upstream header
exactly. Keeping it byte-faithful to the tag makes the *next* upgrade a clean diff.

> Example from the 11.1.1 bump: additions included the `writebatch_iterate_ld`
> family, `block_align`, `open_files_async`, file-checksum / SST-partitioner /
> table-properties-collector factories, FIFO `max_data_files_size`, and many
> `compaction_service_options_override_*` setters. Removals were
> `skip_checking_sst_file_sizes_on_db_open` (get/set), `readoptions_set_managed`,
> and the **signature change** of `rocksdb_slicetransform_create` (the `in_range`
> callback parameter was dropped).

---

## Step 4 — Reconcile the Go/C glue (the part that breaks the build)

**Additions never break compilation.** Only **removals** and **signature changes**
can, and only if this binding actually references the affected symbol. For each
such symbol, grep the source (excluding the `deps/` artifact):

```bash
# Does any Go/C source use the removed/changed symbol?
grep -rn "rocksdb_options_set_skip_checking_sst_file_sizes_on_db_open" \
  --include=*.go --include=*.c --include=*.h . | grep -v deps/
```

Then fix what you find:

- **Removed C symbol that Go wrapped** → delete the Go wrapper method and any
  tests calling it.
  Example (11.1.1): removed `Options.SetSkipCheckingSSTFileSizesOnDBOpen` /
  `SkipCheckingSSTFileSizesOnDBOpen` from `options.go`, plus their uses in
  `db_test.go` and `options_test.go`.

- **Signature-changed C constructor wrapped by a shim in `grocksdb.c`** → update
  the shim's call, and drop any now-unused `//export`ed Go callback.
  Example (11.1.1): `rocksdb_slicetransform_create` lost its `in_range` callback,
  so `gorocksdb_slicetransform_create` in `grocksdb.c` stopped passing it, and the
  dead `//export gorocksdb_slicetransform_in_range` was removed from
  `slice_transform.go`. The `InRange` method was **kept** on the public
  `SliceTransform` interface (marked `Deprecated:`) to avoid breaking downstream
  implementers — RocksDB simply no longer calls it.

**Public-API stability rule:** prefer keeping exported Go identifiers when only the
*underlying* C API was removed (mark them `Deprecated:`), unless the wrapper
genuinely cannot function. Remove an exported method only when its sole purpose was
to call a C symbol that no longer exists.

After editing, format-check:

```bash
gofmt -l <files you touched>   # empty output = OK
```

---

## Step 5 — Update `README.md` "API Support"

The README lists C APIs that grocksdb does **not** yet wrap in Go. For each
**newly added** C function that you did not wrap, add a `- [ ]` bullet (mirror the
naming already used in that section). This sets expectations and tracks future work.

Quick way to see which new APIs are unwrapped:

```bash
for s in set_block_align open_files_async file_checksum_gen sst_partitioner \
         table_properties_collector iterate_ld max_data_files_size; do
  echo -n "$s: "; grep -rl "$s" --include=*.go . | grep -v deps/ | tr '\n' ' '; echo
done
# Empty result for a symbol = not wrapped → list it in README.
```

---

## Step 6 — Build & test

```bash
make libs          # runs build.sh: builds compression libs + RocksDB into dist/<os>_<arch>
make test          # go test -race -v -count=1 -tags testing,grocksdb_no_link
```

Or point CGO at an existing build:

```bash
CGO_CFLAGS="-I/path/to/rocksdb/include" \
CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lz -lsnappy -llz4 -lzstd" \
  go test -count=1 ./...
```

A **link error** mentioning an `undefined reference to rocksdb_<symbol>` almost
always means the repo-root `c.h` declares a symbol that the rebuilt `librocksdb.a`
no longer exports (a missed **removal**), or vice-versa (a missed **addition**).
Re-run Step 3 against the exact tag.

---

## Step 7 — Commit

Match the existing history style:

```
CHORE Adapt RocksDB 11.1.1
```

Typical touched files: `build.sh`, `c.h`, `README.md`, plus any `*.go` / `grocksdb.c`
reconciled in Step 4.

---

## Reference: how past bumps looked

Inspect prior upgrade commits to calibrate scope:

```bash
git log --oneline --grep="Adapt RocksDB"
git show --stat 632e53a   # 10.10.1: README + build.sh + c.h only (pure additions)
git show --stat 5fc3d82   # 10.9.1 : also touched several *.go (API churn that round)
```

- A **pure-addition** release (like 10.10.1) touches only `build.sh`, `c.h`, `README.md`.
- A release with **removals/signature changes** (like 11.1.1) additionally touches
  `grocksdb.c` and the affected `*.go`/test files — that's Step 4.

---

## Why these mechanics

- **`grocksdb_no_link` / `grocksdb_clean_link` tags:** the default CGO `LDFLAGS`
  (`non_builtin.go`) include `-ldl`, which is POSIX-only. `-tags grocksdb_no_link`
  lets you supply the full link line via `CGO_LDFLAGS` instead. `make test` uses
  `grocksdb_no_link`.
- **Repo-root `c.h` is the contract.** The Go files `#include "rocksdb/c.h"`, but the
  binding is validated against the vendored repo-root copy. Keeping it identical to
  the chosen RocksDB tag is what makes upgrades a mechanical diff.
- **`deps/` is disposable.** It's a gitignored RocksDB working tree; treat it as
  cache, not source.
