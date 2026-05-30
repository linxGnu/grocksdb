# Building grocksdb on Windows (static RocksDB, from source)

This guide builds **RocksDB v10.10.1** from source as a **static library**
(`librocksdb.a`) and then compiles and tests the **grocksdb** CGO binding against
it, using the MSYS2 **UCRT64** GCC toolchain + **CMake** + **Ninja**, driven from
**PowerShell**.

Everything is statically linked: RocksDB and its compression libraries are linked
into the grocksdb test binaries with no runtime DLL dependency on `librocksdb`,
`libstdc++`, `libgcc` or `libwinpthread`.

## Layout

RocksDB is cloned and built **inside this repo**, under `./deps`:

```
grocksdb/                      # this repo
├── deps/                      # created/managed by build-deps.ps1 (git-ignored)
│   └── rocksdb/               # RocksDB source + build
│       ├── librocksdb.a       # static archive (copied next to include/)
│       ├── include/           # rocksdb/*.h  (rocksdb/c.h is what grocksdb needs)
│       └── build/             # CMake/Ninja build tree
├── build-deps.ps1             # clones + builds RocksDB into ./deps/rocksdb
└── build-grocksdb.ps1         # builds & tests grocksdb against ./deps/rocksdb
```

The static archive and headers end up at:

```
deps/rocksdb/librocksdb.a
deps/rocksdb/include/                 (rocksdb/*.h)
```

> `./deps/` is already listed in this repo's `.gitignore`, so the cloned/built
> dependency is never committed.

---

## Prerequisites — install these BEFORE running any script

### 1. git and Go (from their homepages)

- Git: https://git-scm.com/download/win
- Go:  https://go.dev/doc/install

### 2. MSYS2 + UCRT64 toolchain + CMake/Ninja

Install MSYS2: https://www.msys2.org/

Then open the **MSYS2 UCRT64** shell **once** and install the toolchain, CMake,
Ninja, and the compression libraries RocksDB needs. After this, you never need the
MSYS2 shell again — everything else runs from PowerShell.

```bash
pacman -Syu
pacman -S --needed \
  base-devel \
  mingw-w64-ucrt-x86_64-toolchain \
  mingw-w64-ucrt-x86_64-cmake \
  mingw-w64-ucrt-x86_64-ninja \
  git \
  mingw-w64-ucrt-x86_64-lz4 \
  mingw-w64-ucrt-x86_64-zstd \
  mingw-w64-ucrt-x86_64-zlib \
  mingw-w64-ucrt-x86_64-snappy \
  mingw-w64-ucrt-x86_64-bzip2
```

This gives you:

- `C:\msys64\ucrt64\bin\gcc.exe`, `g++.exe`, `cmake.exe`, `ninja.exe`
- static `.a` archives for lz4 / zstd / zlib / snappy / bzip2 under
  `C:\msys64\ucrt64\lib`, with headers under `C:\msys64\ucrt64\include`

> Visual Studio / MSBuild and vcpkg are **not** required.

---

## Step 1 — build RocksDB: `build-deps.ps1`

`build-deps.ps1` (in the repo root) does the following:

1. Verifies the UCRT64 toolchain (`gcc.exe`, `g++.exe`, `cmake.exe`, `ninja.exe`).
2. **Removes** an existing `./deps` if present, then **recreates** it.
3. **Clones** RocksDB `v10.10.1` into `./deps/rocksdb`.
4. Builds it as a **static** library with the UCRT64 GCC toolchain + CMake + Ninja.
5. Copies `librocksdb.a` next to `include/` for simple CGO link paths.

Run it from the repo root in PowerShell:

```powershell
.\build-deps.ps1
```

If PowerShell blocks the script (execution policy), run it once as:

```powershell
powershell -ExecutionPolicy Bypass -File .\build-deps.ps1
```

Result:

- `deps\rocksdb\librocksdb.a`
- headers in `deps\rocksdb\include` (notably `deps\rocksdb\include\rocksdb\c.h`)

### Notes on the RocksDB CMake flags

- **`-G Ninja` + `-DCMAKE_C_COMPILER`/`-DCMAKE_CXX_COMPILER` → UCRT64 GCC.** This
  forces a GCC/Ninja build instead of MSVC. `CC`/`CXX` are also exported by the
  script.
- **`-DROCKSDB_BUILD_SHARED=OFF`** + `--target rocksdb` builds only the static
  `librocksdb.a` (no shared DLL, no tools/benchmarks/tests).
- **Compression: `-DWITH_LZ4/ZSTD/ZLIB/SNAPPY/BZ2=ON`.** These match the libraries
  installed via pacman above and the link flags used by `build-grocksdb.ps1`.
- **`-DWITH_LIBURING=OFF`** — io_uring is Linux-only.
- **`-DPORTABLE=ON`** avoids emitting CPU instructions the build host supports but a
  target machine may not.
- **`-DCMAKE_POLICY_VERSION_MINIMUM=3.5`** (quoted) — CMake 4.x removed
  compatibility with projects declaring `cmake_minimum_required(VERSION <3.5)`,
  which some of RocksDB's bundled CMake files still do. The flag tells CMake to
  assume a 3.5 policy baseline and configure anyway. It **must be quoted**;
  unquoted, PowerShell splits the token at the `.` and passes a stray `.5`.

> **Why RocksDB v10.10.1?** This grocksdb checkout is written against the RocksDB
> **10.10.1** C API (`rocksdb/c.h`), so we build that exact tag. Newer RocksDB
> (11.x) changes the C API and won't link against this binding. RocksDB 10.10.1
> also builds cleanly under the GCC in current UCRT64.

---

## Step 2 — build & test grocksdb: `build-grocksdb.ps1`

`build-grocksdb.ps1` points CGO at the static archive under `.\deps\rocksdb`,
compiles the binding, and runs the test suite.

Run it from the repo root in PowerShell:

```powershell
.\build-grocksdb.ps1            # compile + run tests
.\build-grocksdb.ps1 -NoTest    # compile only
```

If PowerShell blocks the script (execution policy):

```powershell
powershell -ExecutionPolicy Bypass -File .\build-grocksdb.ps1
```

What it does:

1. Verifies the UCRT64 toolchain and that `deps\rocksdb\librocksdb.a` +
   `deps\rocksdb\include\rocksdb\c.h` exist (i.e. `build-deps.ps1` ran).
2. Sets the CGO toolchain (`CC`/`CXX` → UCRT64 GCC, `CGO_ENABLED=1`).
3. Sets `CGO_CFLAGS`/`CGO_CXXFLAGS` to `-I deps\rocksdb\include`.
4. Sets `CGO_LDFLAGS` to the full static link line (see below).
5. `go build  -tags grocksdb_no_link ./...`
6. `go test  -tags grocksdb_no_link -count=1 ./...`

### Why `-tags grocksdb_no_link`

By default grocksdb injects its own linker flags
(`-lrocksdb -pthread -lstdc++ -ldl -lm -lzstd -llz4 -lz -lsnappy`). The `-ldl`
entry is a **POSIX-only** library that does not exist on Windows/mingw, so the
default flags fail to link. The `grocksdb_no_link` build tag disables grocksdb's
built-in flags and lets us supply the complete link line via `CGO_LDFLAGS`.

### The static link line (`CGO_LDFLAGS`)

```
-L<deps\rocksdb> -L<C:\msys64\ucrt64\lib>
-lrocksdb
-lzstd -llz4 -lz -lsnappy -lbz2
-lstdc++ -lm -lpthread
-lshlwapi -lrpcrt4 -lbcrypt -ldbghelp
-static -static-libgcc -static-libstdc++
```

- **Link order matters with GCC static libs:** the consumer (`-lrocksdb`) comes
  **before** the libraries it depends on (the compression libs), which come before
  the C++/runtime libs, which come before the Windows system libs.
- **Compression libs** (`-lzstd -llz4 -lz -lsnappy -lbz2`) come from
  `C:\msys64\ucrt64\lib` and match the `-DWITH_*=ON` flags used to build RocksDB.
- **`-lshlwapi -lrpcrt4 -lbcrypt -ldbghelp`** satisfy Windows API references in
  recent RocksDB (`Path*`, `UuidCreate`, `BCrypt*`, `SymInitialize`). Drop any that
  the linker reports as unused; add back if you see undefined `BCrypt*` /
  `SymInitialize` / `Path*` symbols.
- **`-static -static-libgcc -static-libstdc++`** produce binaries with no runtime
  dependency on `libstdc++-6.dll` / `libgcc_s_seh-1.dll` / `libwinpthread-1.dll`.

---

## Using grocksdb in your own project on Windows

After `build-deps.ps1` has produced `deps\rocksdb\librocksdb.a`, point your own
`go build` at it the same way `build-grocksdb.ps1` does — set the UCRT64 GCC
toolchain, then:

```powershell
$RocksDb = "C:\path\to\grocksdb\deps\rocksdb"
$env:CGO_CFLAGS  = "-I$RocksDb\include"
$env:CGO_LDFLAGS = "-L$RocksDb -L C:\msys64\ucrt64\lib -lrocksdb " +
                   "-lzstd -llz4 -lz -lsnappy -lbz2 -lstdc++ -lm -lpthread " +
                   "-lshlwapi -lrpcrt4 -lbcrypt -ldbghelp " +
                   "-static -static-libgcc -static-libstdc++"
go build -tags grocksdb_no_link .\...
```

## Verify static linking

With `objdump` / `ldd` from `C:\msys64\ucrt64\bin` on PATH, build a test binary and
inspect it:

```powershell
go test -c -tags grocksdb_no_link -o grocksdb.test.exe .
ldd .\grocksdb.test.exe
```

It should list only Windows system DLLs (KERNEL32, ws2_32, ...) — **no**
`librocksdb`, `libstdc++-6.dll`, `libgcc_s_seh-1.dll` or `libwinpthread-1.dll`. If
those appear, the static CRT flags didn't take — re-check
`-static -static-libgcc -static-libstdc++` in `build-grocksdb.ps1`.
