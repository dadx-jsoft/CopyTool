# CopyTool
Copy all files from a folder to another folder without subfolder (flattened).

```
CopyTool/
├── cmd/copytool/     Go GUI entry
├── internal/copy/    copy logic + tests
├── java/             original Swing app (archive)
├── scripts/          release / macOS build
├── Icon.png
├── FyneApp.toml
└── go.mod
```

## Download

Latest release: **[v1.1.0](https://github.com/daduong-zen8labs/CopyTool/releases/tag/v1.1.0)** · [All releases](https://github.com/daduong-zen8labs/CopyTool/releases)

### Portable (run directly, no install)

| OS | Download |
|---|---|
| Windows | [CopyTool-windows-amd64.exe](https://github.com/daduong-zen8labs/CopyTool/releases/download/v1.1.0/CopyTool-windows-amd64.exe) |
| Linux | [CopyTool-linux-amd64](https://github.com/daduong-zen8labs/CopyTool/releases/download/v1.1.0/CopyTool-linux-amd64) |
| macOS | Build on a Mac with `scripts/build-darwin.sh` (not in this release) |

```bash
# Linux
chmod +x ./CopyTool-linux-amd64
./CopyTool-linux-amd64
```

### Install packages

| OS | Download | How |
|---|---|---|
| Windows | [CopyTool-windows-amd64.zip](https://github.com/daduong-zen8labs/CopyTool/releases/download/v1.1.0/CopyTool-windows-amd64.zip) | Unzip → run `CopyTool.exe` |
| Linux | [CopyTool-linux-amd64.tar.xz](https://github.com/daduong-zen8labs/CopyTool/releases/download/v1.1.0/CopyTool-linux-amd64.tar.xz) | `tar -xJf ...` then `sudo cp -a usr/local/* /usr/local/` |
| macOS | — | Run `./scripts/build-darwin.sh` on a Mac → unzip `.app.zip` → open `CopyTool.app` |

## Features

- **All files** — copy every file (flattened)
- **By extension** — only matching types, e.g. `jpg, png, pdf`
- **Exclude extension** — all except listed types, e.g. `tmp, log`
- Path: paste into From/To, or **Browse** (native folder dialog on Windows / macOS / Linux)

## Build from source

Requires **Go 1.22+** and a **C compiler** (Fyne uses CGO):

| OS | C compiler |
|---|---|
| Windows | MinGW / TDM-GCC / LLVM (`gcc` or `clang`) |
| macOS | Xcode Command Line Tools (`xcode-select --install`) |
| Linux | `build-essential` (or equivalent `gcc`) |

```bash
go test ./...
go run ./cmd/copytool
```

### One-shot release (Windows host + Docker)

```powershell
go install github.com/fyne-io/fyne-cross@latest
go install fyne.io/tools/cmd/fyne@latest
powershell -ExecutionPolicy Bypass -File .\scripts\build-release.ps1
```

Produces `dist/portable/*` and `dist/install/*` for Windows and Linux. macOS is skipped unless `COPYTOOL_MACOSX_SDK` points to a valid `MacOSX*.sdk`.

### macOS on a Mac

```bash
chmod +x ./scripts/build-darwin.sh
./scripts/build-darwin.sh
```

### Native single-OS build

```bash
# Windows
go build -o CopyTool.exe ./cmd/copytool

# macOS / Linux
go build -o CopyTool ./cmd/copytool
```

## Java (original)

Sources live in `java/`. Build jar there, then use Launch4j to create a .exe file.
