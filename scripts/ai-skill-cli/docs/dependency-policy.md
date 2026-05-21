# Dependency Policy：Ai-skill CLI Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## 原則

`ai-skill` CLI 的核心目標是降低使用者端 runtime 安裝成本。核心命令應優先使用 Go standard library 或 pure Go dependency，避免要求使用者預先安裝 Ruby、Python、sqlite3 CLI、pip、gem、C compiler 或 POSIX shell。

## Dependency 分類

| 類型 | 政策 | 目前狀態 |
| --- | --- | --- |
| Go standard library | 優先使用 | JSON、flag、path、process、filesystem checks 已使用 standard library |
| Pure Go dependency | 可使用；需記錄用途與替代方案 | SQLite 採用 `modernc.org/sqlite` |
| CGO dependency | 預設不使用；若要使用，必須有 ADR、CI matrix、Windows 安裝成本與 fallback | `mattn/go-sqlite3` 不作為預設 |
| 外部 desktop binary | 只在使用者環境語意不可替代時允許 | Git 保持 external dependency |
| Shell / Ruby / Python | 只能是 wrapper-mode 過渡依賴；不得成為長期核心依賴 | Runtime compiler / validators Phase 3 先 wrapper |
| Tool-specific API / path | 放在 tool adapter，不進通用 CLI 預設 | Cursor / Roo Code 設定 helper 維持 tool-specific |

## SQLite 決策

預設 SQLite engine 使用 `modernc.org/sqlite`，理由：

- pure Go，符合單一 binary 與跨平台部署目標。
- 不要求使用者安裝 `sqlite3` CLI、SQLite development headers 或 C compiler。
- Windows、macOS、Linux 的 build friction 低於 CGO SQLite。
- Phase 1 / Phase 3 已用 `doctor --check-runtime` 建立 in-memory query 與 temporary file-backed create / insert / query / integrity proof，並可用同一 driver 檢查 `runtime.db` integrity。

`mattn/go-sqlite3` 不作為預設，因為它需要 CGO 與平台 C toolchain。只有在 `modernc.org/sqlite` 出現無法接受的 compatibility 或 performance 缺口時，才可另開 ADR 評估 CGO 例外。

## Git 外部依賴

Git 不包進 binary。需要 commit、push、hooks、credential helper、SSH key、GPG signing、LFS 或 submodule 語意的命令，預設呼叫使用者本機 Git。

必要行為：

- `doctor --require-git` 缺 Git 時回傳 `missing_dependency` / `missing_git`。
- `close-loop`、linked-update、hook install 與任何 commit / push 命令缺 Git 時必須阻斷。
- 不得用 Go git library 假裝完成 commit / push 行為，除非另有明確 contract 與 parity tests。

## Wrapper Mode 限制

Phase 3 仍可用 Go CLI 包裝既有 Ruby compiler / validators，但 wrapper mode 必須：

- 清楚回報仍依賴 Ruby / Python / 外部工具的命令。
- 在 dependency 缺失時回傳 `missing_dependency`，不得 partial success。
- 保留 parity fixture；native replacement 通過前不得刪除舊 compiler / validator。
- 最終 closure 仍以刪除 replacement 範圍內的舊 script 為預設。
