#requires -Version 5.1
<#
  build-deps.ps1 — clone & statically build RocksDB into ./deps

  Builds RocksDB v11.1.1 as a static library (librocksdb.a) using the MSYS2
  UCRT64 GCC toolchain + CMake + Ninja, so that grocksdb can link against it.

  Run from the grocksdb repo root:
      .\build-deps.ps1

  If PowerShell blocks the script (execution policy):
      powershell -ExecutionPolicy Bypass -File .\build-deps.ps1
#>

$ErrorActionPreference = 'Stop'

# --- Version ---
# grocksdb (this repo) targets the RocksDB 11.1.1 C API. Build the matching tag so
# the CGO wrappers in c.h / grocksdb.c link cleanly. Using a different RocksDB minor
# version may add/remove C API symbols and fail to link against this binding.
$RocksDbTag = 'v11.1.1'

# --- Toolchain locations (MSYS2 UCRT64) ---
$Ucrt64 = 'C:\msys64\ucrt64'
$Gcc    = Join-Path $Ucrt64 'bin\gcc.exe'
$Gpp    = Join-Path $Ucrt64 'bin\g++.exe'
$Cmake  = Join-Path $Ucrt64 'bin\cmake.exe'
$Ninja  = Join-Path $Ucrt64 'bin\ninja.exe'

# --- Paths ---
$Root = $PSScriptRoot
$Deps = Join-Path $Root 'deps'

function Assert-Tool($path, $name) {
    if (-not (Test-Path $path)) {
        throw "$name not found at $path. Install the MSYS2 UCRT64 toolchain first (see win-build.md)."
    }
    Write-Host "  [ok] $name -> $path"
}

# ---------------------------------------------------------------------------
# Check toolchain
# ---------------------------------------------------------------------------
Write-Host '== Checking UCRT64 toolchain =='
if (-not (Test-Path $Ucrt64)) {
    throw "UCRT64 not found at $Ucrt64. Install MSYS2 and the UCRT64 toolchain first."
}
Assert-Tool $Gcc   'gcc'
Assert-Tool $Gpp   'g++'
Assert-Tool $Cmake 'cmake'
Assert-Tool $Ninja 'ninja'

# Put UCRT64 bin first on PATH for this process (gcc, g++, cmake, ninja, and the
# static compression archives' tooling). Set CC/CXX so CMake picks GCC, not MSVC.
$env:PATH = "$Ucrt64\bin;$env:PATH"
$env:CC   = $Gcc
$env:CXX  = $Gpp

# ---------------------------------------------------------------------------
# Reset ./deps
# ---------------------------------------------------------------------------
Write-Host ''
Write-Host '== Resetting deps directory =='
if (Test-Path $Deps) {
    Write-Host "  removing existing $Deps"
    Remove-Item -Recurse -Force $Deps
}
New-Item -ItemType Directory -Path $Deps | Out-Null
Write-Host "  created $Deps"

# ---------------------------------------------------------------------------
# Clone RocksDB into ./deps/rocksdb
# ---------------------------------------------------------------------------
Write-Host ''
Write-Host "== Cloning RocksDB $RocksDbTag =="
$RocksDbSrc = Join-Path $Deps 'rocksdb'
git clone --depth 1 --branch $RocksDbTag https://github.com/facebook/rocksdb.git $RocksDbSrc
if ($LASTEXITCODE -ne 0) { throw 'git clone rocksdb failed' }

# ---------------------------------------------------------------------------
# Build RocksDB (static, with compression libs)
# ---------------------------------------------------------------------------
Write-Host ''
Write-Host '== Building RocksDB (static) =='
$RocksDbBuild = Join-Path $RocksDbSrc 'build'

& $Cmake -S $RocksDbSrc -B $RocksDbBuild -G Ninja `
    -DCMAKE_BUILD_TYPE=Release `
    "-DCMAKE_POLICY_VERSION_MINIMUM=3.5" `
    -DCMAKE_C_COMPILER="$Ucrt64/bin/gcc.exe" `
    -DCMAKE_CXX_COMPILER="$Ucrt64/bin/g++.exe" `
    -DROCKSDB_BUILD_SHARED=OFF `
    -DWITH_GFLAGS=OFF `
    -DWITH_TESTS=OFF `
    -DWITH_BENCHMARK_TOOLS=OFF `
    -DWITH_TOOLS=OFF `
    -DWITH_CORE_TOOLS=OFF `
    -DWITH_TRACE_TOOLS=OFF `
    -DWITH_LZ4=ON `
    -DWITH_ZSTD=ON `
    -DWITH_ZLIB=ON `
    -DWITH_SNAPPY=ON `
    -DWITH_BZ2=ON `
    -DWITH_LIBURING=OFF `
    -DPORTABLE=ON `
    -DFAIL_ON_WARNINGS=OFF
if ($LASTEXITCODE -ne 0) { throw 'rocksdb cmake configure failed' }

& $Cmake --build $RocksDbBuild --config Release --target rocksdb
if ($LASTEXITCODE -ne 0) { throw 'rocksdb build failed' }

# Copy the archive next to include/ so CGO link paths are simple
# (-L deps/rocksdb -I deps/rocksdb/include).
Copy-Item (Join-Path $RocksDbBuild 'librocksdb.a') (Join-Path $RocksDbSrc 'librocksdb.a') -Force

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
Write-Host ''
Write-Host '== Done =='
Write-Host "  rocksdb lib: $RocksDbSrc\librocksdb.a"
Write-Host "  rocksdb inc: $RocksDbSrc\include"
Write-Host ''
Write-Host 'Next: build & test grocksdb with  .\build-grocksdb.ps1'
