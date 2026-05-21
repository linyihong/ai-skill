# Release Distribution Decision

## Decision

`ai-skill` 的桌面發佈採雙軌：

| Channel | 用途 | Status |
| --- | --- | --- |
| Repo-local `bin/` | clone repo 後不安裝 Go 也能直接執行核心命令 | Primary for this repository |
| GitHub Actions artifact | CI 產生 release-like binaries 與 checksums | Primary for automation |
| GitHub Releases | 未來正式版本發布 | Deferred until version cadence is stable |
| Homebrew / Scoop / winget | OS package manager distribution | Deferred; maintenance cost 暫不納入本計畫 |

## Repo-Local Binaries

Committed binaries:

- `scripts/ai-skill-cli/bin/ai-skill-darwin-arm64`
- `scripts/ai-skill-cli/bin/ai-skill-darwin-amd64`
- `scripts/ai-skill-cli/bin/ai-skill-linux-amd64`
- `scripts/ai-skill-cli/bin/ai-skill-linux-arm64`
- `scripts/ai-skill-cli/bin/ai-skill-windows-amd64.exe`
- `scripts/ai-skill-cli/bin/BUILDINFO`
- `scripts/ai-skill-cli/bin/SHA256SUMS`

Rebuild rule:

- Rebuild `bin/` only when CLI source changes under `scripts/ai-skill-cli/cmd/`, `scripts/ai-skill-cli/internal/`, `go.mod`, or `go.sum`.
- `go test ./...` verifies `BUILDINFO`, `SHA256SUMS`, and current-host binary smoke.
- If the guard fails with a newer source commit, rebuild with:

```bash
go run ./cmd/releasebuild --stable-names --version "repo-$(git rev-parse --short HEAD)" --commit "$(git rev-parse --short HEAD)" --dist bin
```

## CI Artifacts

`.github/workflows/ai-skill-cli.yml` runs:

- `go test ./...` on Windows, macOS, and Linux.
- `ai-skill version`, `doctor`, native runtime refresh, validate, and compile smoke.
- Ubuntu artifact build for Windows amd64, macOS amd64/arm64, Linux amd64/arm64.
- `SHA256SUMS` verification and artifact upload.

## Upgrade And Rollback

Upgrade:

1. Pull latest repo.
2. Run the host binary from `scripts/ai-skill-cli/bin/`.
3. Run `runtime refresh --repo <repo>` to regenerate local gitignored runtime caches.
4. Run `runtime validate --repo <repo>`.

Rollback:

1. Checkout the previous commit.
2. Use the previous committed `bin/` binary.
3. Regenerate local runtime caches with `runtime refresh`.

`knowledge/runtime/sqlite/runtime-index.sqlite` remains gitignored generated cache and is not part of rollback state.

## Deferred Package Managers

Homebrew, Scoop, and winget are deferred because they add:

- Formula / manifest maintenance.
- Signing and trust-chain decisions.
- Release version cadence requirements.
- Support burden for package-manager-specific failures.

The current supported path is repo-local binaries plus CI artifacts.
