# Agent Instructions

## Build

```bash
CGO_ENABLED=0 go build -o ntfs-tool .
```

`CGO_ENABLED=0` is required — the tool has no C dependencies and must be a static binary.

## Test

No test suite exists. Manual verification only: `./ntfs-tool list` (requires NTFS drive plugged in). Mount/unmount commands require `sudo`.

## Release

1. Tag and push: `git tag v0.2.0 && git push origin v0.2.0`
2. GitHub Actions runs GoReleaser (`.goreleaser.yml`) — builds darwin amd64/arm64
3. After release, manually update `Formula/ntfs-tool.rb` in the [homebrew-tap](https://github.com/neilning-xc/homebrew-tap) repo with new version + checksums

## Architecture

Small single-package Go CLI. `main.go` → `cmd/` (cobra commands) → `internal/` (disk detection, shell execution).

- `cmd/` — cobra commands: `list`, `mount`, `unmount`
- `internal/disk.go` — parses `diskutil info -plist` output to find NTFS partitions; detects FUSE framework and ntfs-3g binary locations
- `internal/shell.go` — `exec.Command` wrapper returning stdout/stderr/exit code

Dependencies: `github.com/spf13/cobra`, `howett.net/plist`

## Key behaviors

- Mount/unmount require `sudo` (checked via `os.Geteuid()`)
- FUSE detected by checking `/Library/Filesystems/macfuse.fs` or `fuse-t.fs`
- ntfs-3g found via hardcoded Homebrew paths (`/opt/homebrew/bin/`, `/usr/local/bin/`) then fallback to `which`
- All user-facing output is in Chinese (zh-CN)
- Plist parsing is macOS-specific (`diskutil` format)
