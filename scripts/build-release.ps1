# Build portable binaries + install packages.
# Requires: Go, Docker Desktop, fyne-cross.
# Optional for macOS cross-build from Windows/Linux:
#   set COPYTOOL_MACOSX_SDK to a MacOSX*.sdk folder (see README).
#
# Usage (from repo root):
#   powershell -ExecutionPolicy Bypass -File .\scripts\build-release.ps1

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root
. (Join-Path $PSScriptRoot "build-flags.ps1")

$env:PATH = "$(go env GOPATH)\bin;$env:PATH"

if (-not (Get-Command docker -ErrorAction SilentlyContinue)) {
    throw "Docker is required for fyne-cross"
}
docker info | Out-Null
if (-not (Get-Command fyne-cross -ErrorAction SilentlyContinue)) {
    throw "fyne-cross not found. Run: go install github.com/fyne-io/fyne-cross@latest"
}

$Dist = Join-Path $Root "dist"
$Portable = Join-Path $Dist "portable"
$Install = Join-Path $Dist "install"

Remove-Item -Recurse -Force $Dist -ErrorAction SilentlyContinue
New-Item -ItemType Directory -Force -Path $Portable, $Install | Out-Null

$pkg = "./cmd/copytool"
$common = @(
    "-name", "CopyTool",
    "-app-id", "com.vn.copytool",
    "-icon", "Icon.png"
)

Write-Host "==> fyne-cross windows amd64"
& fyne-cross windows @common -arch amd64 $pkg
if ($LASTEXITCODE -ne 0) { throw "windows build failed" }

Write-Host "==> fyne-cross linux amd64"
& fyne-cross linux @common -arch amd64 $pkg
if ($LASTEXITCODE -ne 0) { throw "linux build failed" }

$sdk = $env:COPYTOOL_MACOSX_SDK
if (-not $sdk) {
    $defaultSdk = Join-Path $Root "third_party\macos-sdk\MacOSX11.3.sdk"
    if (Test-Path $defaultSdk) { $sdk = $defaultSdk }
}

if ($sdk -and (Test-Path $sdk)) {
    Write-Host "==> fyne-cross darwin amd64,arm64 (SDK: $sdk)"
    & fyne-cross darwin @common -arch "amd64,arm64" -macosx-sdk-path $sdk $pkg
    if ($LASTEXITCODE -ne 0) {
        Write-Warning "darwin cross-build failed. Build on a Mac instead: ./scripts/build-darwin.sh"
    }
} else {
    Write-Warning "Skipping darwin. Set COPYTOOL_MACOSX_SDK or run scripts/build-darwin.sh on macOS."
}

# Also keep a stripped local Windows portable (GUI subsystem, no console)
Write-Host "==> local windows portable"
go build -ldflags="$WindowsLdflags" -o (Join-Path $Portable "CopyTool-windows-amd64.exe") ./cmd/copytool

function Copy-IfExists([string]$Src, [string]$Dest) {
    if (Test-Path $Src) {
        New-Item -ItemType Directory -Force -Path (Split-Path $Dest) | Out-Null
        Copy-Item $Src $Dest -Force
        Write-Host "  + $Dest"
    }
}

$Bin = Join-Path $Root "fyne-cross\bin"
$Pkg = Join-Path $Root "fyne-cross\dist"

# Prefer fyne-cross windows binary if present
Copy-IfExists (Join-Path $Bin "windows-amd64\CopyTool.exe") (Join-Path $Portable "CopyTool-windows-amd64.exe")
Copy-IfExists (Join-Path $Bin "linux-amd64\CopyTool") (Join-Path $Portable "CopyTool-linux-amd64")
Copy-IfExists (Join-Path $Bin "linux-amd64\copytool") (Join-Path $Portable "CopyTool-linux-amd64")
Copy-IfExists (Join-Path $Bin "darwin-amd64\CopyTool") (Join-Path $Portable "CopyTool-darwin-amd64")
Copy-IfExists (Join-Path $Bin "darwin-arm64\CopyTool") (Join-Path $Portable "CopyTool-darwin-arm64")

Copy-IfExists (Join-Path $Pkg "windows-amd64\CopyTool.zip") (Join-Path $Install "CopyTool-windows-amd64.zip")
Copy-IfExists (Join-Path $Pkg "linux-amd64\CopyTool.tar.xz") (Join-Path $Install "CopyTool-linux-amd64.tar.xz")
Copy-IfExists (Join-Path $Pkg "darwin-amd64\CopyTool.app.zip") (Join-Path $Install "CopyTool-darwin-amd64.app.zip")
Copy-IfExists (Join-Path $Pkg "darwin-arm64\CopyTool.app.zip") (Join-Path $Install "CopyTool-darwin-arm64.app.zip")

# Darwin .app folders -> zip if zip missing
Get-ChildItem (Join-Path $Pkg "darwin-*") -Directory -ErrorAction SilentlyContinue | ForEach-Object {
    $app = Get-ChildItem $_.FullName -Filter "*.app" -Directory -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($app) {
        $zipPath = Join-Path $Install ("CopyTool-" + $_.Name + ".app.zip")
        if (-not (Test-Path $zipPath)) {
            Compress-Archive -Path $app.FullName -DestinationPath $zipPath -Force
            Write-Host "  + $zipPath"
        }
    }
}

# Write short usage next to artifacts
@"
Portable: run the binary for your OS (no install).
  Windows: CopyTool-windows-amd64.exe
  Linux:   chmod +x CopyTool-linux-amd64 && ./CopyTool-linux-amd64
  macOS:   chmod +x CopyTool-darwin-arm64 (Apple Silicon) or CopyTool-darwin-amd64 (Intel)

Install packages (dist/install):
  Windows: unzip CopyTool-windows-amd64.zip and run CopyTool.exe
  Linux:   tar -xJf CopyTool-linux-amd64.tar.xz && sudo cp -a usr/local/* /usr/local/
  macOS:   unzip *.app.zip and open CopyTool.app (unsigned: right-click > Open)
"@ | Set-Content -Encoding UTF8 (Join-Path $Dist "HOW_TO_RUN.txt")

Write-Host ""
Write-Host "Done."
Get-ChildItem -Recurse $Dist | Select-Object FullName, Length | Format-Table -AutoSize
