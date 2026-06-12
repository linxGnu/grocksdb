#requires -Version 5.1
<#
  build-grocksdb.ps1 — build & test the grocksdb CGO binding against the static
  RocksDB built by build-deps.ps1 into ./deps/rocksdb.

  Uses the MSYS2 UCRT64 GCC toolchain and links statically against
  deps/rocksdb/librocksdb.a plus the UCRT64 compression archives.

  Run from the grocksdb repo root:
      .\build-grocksdb.ps1            # compile (go build) + run tests (go test)
      .\build-grocksdb.ps1 -NoTest    # compile only, skip tests

  If PowerShell blocks the script (execution policy):
      powershell -ExecutionPolicy Bypass -File .\build-grocksdb.ps1
#>

param(
    [switch]$NoTest
)

$ErrorActionPreference = 'Stop'

# --- Toolchain (UCRT64 GCC) ---
# Detection order (same as build-deps.ps1):
#   1. $env:UCRT64      — explicit override (e.g. CI runners / msys2/setup-msys2)
#   2. $env:MSYS2_ROOT  — MSYS2 install root; UCRT64 lives at <root>\ucrt64
#   3. C:\msys64        — standard MSYS2 install root
#   4. Registry         — MSYS2 installer uninstall entry (InstallLocation)
function Find-Ucrt64 {
    if ($env:UCRT64) { return $env:UCRT64 }
    if ($env:MSYS2_ROOT) {
        $candidate = Join-Path $env:MSYS2_ROOT 'ucrt64'
        if (Test-Path $candidate) { return $candidate }
    }
    if (Test-Path 'C:\msys64\ucrt64') { return 'C:\msys64\ucrt64' }
    $uninstallKeys = @(
        'HKCU:\Software\Microsoft\Windows\CurrentVersion\Uninstall\*',
        'HKLM:\Software\Microsoft\Windows\CurrentVersion\Uninstall\*',
        'HKLM:\Software\WOW6432Node\Microsoft\Windows\CurrentVersion\Uninstall\*'
    )
    foreach ($keyPath in $uninstallKeys) {
        $entry = Get-ItemProperty -Path $keyPath -ErrorAction SilentlyContinue |
            Where-Object { $_.DisplayName -like 'MSYS2*' -and $_.InstallLocation } |
            Select-Object -First 1
        if ($entry) {
            $candidate = Join-Path $entry.InstallLocation 'ucrt64'
            if (Test-Path $candidate) { return $candidate }
        }
    }
    return 'C:\msys64\ucrt64'  # keep default so the error message is actionable
}
$Ucrt64 = Find-Ucrt64
$Gcc    = Join-Path $Ucrt64 'bin\gcc.exe'
$Gpp    = Join-Path $Ucrt64 'bin\g++.exe'

# --- Locate the RocksDB dep built by build-deps.ps1 (./deps/rocksdb) ---
$Repo    = $PSScriptRoot
$Deps    = Join-Path $Repo 'deps'
$RocksDb = Join-Path $Deps 'rocksdb'

# --- Sanity checks ---
if (-not (Test-Path $Gcc)) { throw "gcc not found at $Gcc. Install the MSYS2 UCRT64 toolchain (see win-build.md)." }
if (-not (Test-Path $Gpp)) { throw "g++ not found at $Gpp. Install the MSYS2 UCRT64 toolchain (see win-build.md)." }
if (-not (Test-Path (Join-Path $RocksDb 'librocksdb.a'))) {
    throw "librocksdb.a missing under $RocksDb. Run .\build-deps.ps1 first."
}
if (-not (Test-Path (Join-Path $RocksDb 'include\rocksdb\c.h'))) {
    throw "RocksDB C headers missing under $RocksDb\include. Run .\build-deps.ps1 first."
}

# --- Put UCRT64 bin first on PATH and set the CGO toolchain ---
$env:PATH        = "$Ucrt64\bin;$env:PATH"
$env:CGO_ENABLED = '1'
$env:CC          = $Gcc
$env:CXX         = $Gpp

# --- Include paths: RocksDB headers ---
$env:CGO_CFLAGS   = "-I$RocksDb\include"
$env:CGO_CXXFLAGS = "-I$RocksDb\include"

# --- Static link flags.
#     We build with `-tags grocksdb_no_link` so grocksdb does NOT inject its own
#     default LDFLAGS (which include `-ldl`, a POSIX-only lib that does not exist
#     on Windows/mingw). We provide the full link line ourselves here.
#
#     Order matters with GCC static libs: the consumer (-lrocksdb) comes BEFORE the
#     libraries it depends on (compression libs, then C++/runtime, then Windows
#     system libs). -lbcrypt / -ldbghelp satisfy RocksDB's BCrypt*/SymInitialize
#     references on recent Windows. ---
$env:CGO_LDFLAGS = @(
    "-L$RocksDb", "-L$Ucrt64\lib",
    '-lrocksdb',
    '-lzstd', '-llz4', '-lz', '-lsnappy', '-lbz2',
    '-lstdc++', '-lm', '-lpthread',
    '-lshlwapi', '-lrpcrt4', '-lbcrypt', '-ldbghelp',
    '-static', '-static-libgcc', '-static-libstdc++'
) -join ' '

Write-Host '== grocksdb build environment =='
Write-Host "  CC            = $env:CC"
Write-Host "  CXX           = $env:CXX"
Write-Host "  CGO_CFLAGS    = $env:CGO_CFLAGS"
Write-Host "  CGO_LDFLAGS   = $env:CGO_LDFLAGS"
Write-Host ''

# --- Ensure module deps are present ---
Write-Host '== go mod download =='
go mod download
if ($LASTEXITCODE -ne 0) { throw 'go mod download failed' }

# ---------------------------------------------------------------------------
# Compile the package (CGO). `grocksdb_no_link` => use our CGO_LDFLAGS only.
# ---------------------------------------------------------------------------
Write-Host ''
Write-Host '== Building grocksdb (go build) =='
go build -v -tags grocksdb_no_link ./...
if ($LASTEXITCODE -ne 0) { throw 'go build failed' }
Write-Host '  [ok] grocksdb compiled'

# ---------------------------------------------------------------------------
# Run tests
# ---------------------------------------------------------------------------
if ($NoTest) {
    Write-Host ''
    Write-Host '== Skipping tests (-NoTest) =='
} else {
    Write-Host ''
    Write-Host '== Running tests (go test) =='
    go test -v -tags grocksdb_no_link -count=1 ./...
    if ($LASTEXITCODE -ne 0) { throw 'go test failed' }
    Write-Host '  [ok] tests passed'
}

Write-Host ''
Write-Host '=== grocksdb built & tested against static RocksDB in .\deps\rocksdb ==='
