# Dependency Policy：Ai-skill CLI Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## 原則

`ai-skill` CLI 的核心目標是降低使用者端 runtime 安裝成本。核心命令應優先使用 Go standard library 或 pure Go dependency，避免要求使用者預先安裝 Ruby、Python、sqlite3 CLI、pip、gem、C compiler 或 POSIX shell。

## Go-first Automation Policy

新增 repository automation 時，預設必須實作在 `scripts/ai-skill-cli/` 的 Go CLI 中，並由 repo-local binaries 發佈。不要新增新的 `.sh`、`.rb` 或 `.py` 作為長期入口。

例外只允許下列情況：

- Git hook adapter：Git 需要 hook 檔案作為觸發面，但 hook 內部應呼叫 repo-local `ai-skill` binary。
- Thin bootstrap wrapper：只負責找到 repo-local `ai-skill` binary 並轉呼叫，且必須有刪除條件。
- 暫時 retained legacy shell：只有在 Go write-mode parity 尚未完成時保留，且必須列在 parity / disposition / migration map 文件中。

任何例外都必須記錄 owner、side effects、保留原因、fixture gate 與刪除條件；不得把 shell 當成新的核心實作位置。

## Dependency 分類

| 類型 | 政策 | 目前狀態 |
| --- | --- | --- |
| Go standard library | 優先使用 | JSON、flag、path、process、filesystem checks 已使用 standard library |
| Pure Go dependency | 可使用；需記錄用途與替代方案 | SQLite 採用 `modernc.org/sqlite` |
| CGO dependency | 預設不使用；若要使用，必須有 ADR、CI matrix、Windows 安裝成本與 fallback | `mattn/go-sqlite3` 不作為預設 |
| 外部 desktop binary | 只在使用者環境語意不可替代時允許 | Git 保持 external dependency |
| Shell adapters | 不可新增為長期入口；只能是 hook adapter、thin bootstrap wrapper，或尚未完成 write-mode parity 的 retained legacy surface | Runtime validate / refresh / compile / query 與 Roo setting adapter 已是 native default；`sync-cursor-bundle.sh` 已刪，`ai-skill-close-loop.sh` 仍待 Go write-mode parity |
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

## Adapter 限制

Runtime core path 不保留 Ruby / Python adapter mode。仍存在的 shell adapters 必須：

- 清楚回報外部依賴與 side effects。
- 在 dependency 缺失時回傳 `missing_dependency`，不得 partial success。
- 保留 parity fixture；native replacement 通過前不得宣稱已完成。
- 不得新增新功能；新功能必須先進 Go CLI。
- 最終 closure 仍以刪除 replacement 範圍內的舊 script 為預設。
