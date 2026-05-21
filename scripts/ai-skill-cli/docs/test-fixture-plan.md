# 測試 fixture 計畫：Ai-skill CLI Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## Fixture 建置策略

Fixture 必須避開使用者真實 home 目錄、真實 git config、真實 Cursor bundle 與真實工作樹。所有具破壞性或可寫入的行為，都必須先從 dry-run fixture 開始驗證。

## 必要 fixtures

| Fixture | 用途 | 必要斷言 |
| --- | --- | --- |
| `fixture/temp-repo-clean` | clean repo 的 close-loop 行為 | dry-run 無實際動作，status 為 clean |
| `fixture/temp-repo-dirty-owned` | owner group 分類 | 只把預期檔案分到對應群組 |
| `fixture/temp-repo-merge-state` | 不安全 repo 狀態 | 阻斷 commit / push |
| `fixture/temp-repo-hooks-install` | hook install dry-run | source / target / conflict / force 行為可驗證，且不寫入 `.git/hooks/` |
| `fixture/missing-git-path` | PATH 中沒有 Git | `doctor` / `close-loop` 阻斷並提供安裝指引 |
| `fixture/fake-home` | home / Cursor bundle 隔離 | 不寫入真實使用者 home |
| `fixture/windows-paths` | path separator 與 drive 處理 | 正規化後路徑符合 contract |
| `fixture/cursor-mirror-copy-fallback` | Cursor mirror 策略 | Windows / 權限受限環境預設 copy fallback；symlink 只允許明確 opt-in |
| `fixture/runtime-source-change` | runtime.db assertion | source keyword 出現在 generated surface |
| `fixture/runtime-golden-generated-surfaces` | runtime refresh / compile generated surfaces | 同一份 source 可產出 runtime report、model reports、SQLite index、`runtime.db` assertions，且不寫 production artifacts |
| `fixture/runtime-refresh-ordered-steps` | runtime refresh wrapper orchestration | generator / validator steps 順序固定；任一步失敗時停止後續 steps 並回 failure |
| `fixture/runtime-generator-parity` | Ruby vs Go generator parity | Go-native generator output 必須與 Ruby stdout byte-for-byte 相同，通過前不得接 production refresh |
| `fixture/runtime-native-reports-opt-in` | opt-in native report writes | `--native-reports` 只改用 Go 寫三個 Markdown reports；預設 refresh 不切換 |
| `fixture/runtime-native-index-opt-in` | opt-in native SQLite index write | `--native-index` 只改用 Go 寫 SQLite index；validators 仍走 Ruby steps，可不帶 flag rollback |
| `fixture/runtime-sqlite-index-parity` | Ruby vs Go SQLite index parity | Ruby temp DB 與 Go temp DB 的 table counts、row content、source checksum map、FTS hits 必須一致；通過前不得切換 production refresh |
| `fixture/runtime-compiler-snapshot-parity` | runtime compiler DB parity | Ruby temp `runtime.db` snapshots 固定比對 schema row counts、generated surfaces、compiler metadata；`--native-compiler` snapshot output 先比對 generated surfaces / metadata |
| `fixture/native-sqlite-file-proof` | pure Go SQLite proof | temporary DB 可 create / insert / query / integrity check，且不依賴外部 `sqlite3` CLI |
| `fixture/runtime-db-native-validator` | native runtime.db validator | valid DB、missing required table、invalid JSON column、stale compiler metadata warning 均有固定結果 |
| `fixture/runtime-query-index` | native runtime index query | ranking、filter、limit、empty result、missing DB 均有固定結果 |
| `fixture/runtime-query-graph` | native knowledge graph query | source / target / type / keyword filters、limit、empty result、missing filter 均有固定結果 |
| `fixture/runtime-index-native-validator` | native runtime index validator | valid index、missing table、stale checksum、FTS count mismatch、git-ignore boundary 均有固定結果 |
| `fixture/release-artifacts` | cross-platform release build | Windows/macOS/Linux artifacts 與 `SHA256SUMS` 可由 Go release builder 產生 |
| `fixture/legacy-script-parity` | 舊腳本覆蓋率驗證 | 每個 native target / wrapper first 舊入口都有命令映射與測試證據 |

## 缺 Git fixture

目的：證明 Git 是外部依賴，且缺 Git 時不能產生半套 linked-update / close-loop。

設定：

- 在隔離 PATH 的環境中執行命令。
- 確認 `git` 不可被 discover。
- 使用含檔案的暫時 repo-like 目錄，但不實際執行 Git 操作。

斷言：

- `ai-skill doctor --require-git` 回傳 `missing_dependency`。
- `ai-skill close-loop --commit` 回傳 `missing_dependency`。
- 輸出包含安裝指引。
- 不修改任何檔案、index、commit、push、hook、mirror 或 runtime DB。
- Phase 2 已用 isolated PATH tests 覆蓋 `doctor --require-git`、`close-loop --dry-run` 與 `close-loop --commit` 的 `missing_git` 阻斷；`close-loop --commit` 缺 Git 必須優先回 `missing_git`，不得只回 write mode 未完成。

## Runtime DB assertion fixture

目的：證明 source 變更會進入 `runtime/runtime.db`。

設定：

- 使用含唯一 keyword 的 fixture source 內容。
- 分別以 dry-run 與真實 fixture mode 執行 compile。
- 依 source path 與 keyword 查詢 generated surface。

斷言：

- Compile 成功。
- `generated_surfaces` 包含預期 source path。
- `data` 包含預期 keyword。
- keyword 缺失時回傳 `validation_failed`。
- Phase 3 已用 `doctor --check-runtime` / Go unit tests 覆蓋 `native-sqlite-file-proof`：pure Go SQLite 可對 temporary DB 執行 create / insert / query / `PRAGMA integrity_check`，不需外部 `sqlite3` CLI。
- Phase 3 已新增 `runtime-db-native-validator` Go fixture，覆蓋 valid `runtime.db`、missing required table、invalid JSON column、stale compiler metadata warning；stale warning 不阻斷 `runtime validate` 成功狀態。
- Phase 3 已新增 `runtime-query-index` Go fixture，覆蓋 `query-runtime-index.rb` 的 ranking、layer / type / status filters、limit、empty result 與 missing DB 行為；`query-knowledge-graph.rb` fixture 待補。
- Phase 3 已新增 `runtime-query-graph` Go fixture，覆蓋 `query-knowledge-graph.rb` 的 source / target / type / keyword filters、limit、empty result 與 missing filter 行為。
- Phase 3 已新增 `runtime-index-native-validator` Go fixture，覆蓋 SQLite runtime index valid case、missing table、stale checksum、FTS count mismatch、git-ignore boundary；git-ignore boundary 以 Go 呼叫 external Git 檢查。
- Phase 3 已新增 `runtime-golden-generated-surfaces` integration fixture：以同一份 canonical source 產生 runtime report、model context report、model checklists、臨時 SQLite index、臨時 `runtime.db`，並固定檢查 report anchors、index row counts、FTS hit、`generated_surfaces` 與 compiler metadata；fixture 只寫入測試 temp dir。
- Phase 3 已新增 `runtime-refresh-ordered-steps` fixture：`runtime refresh` Go wrapper mode 依序執行 model context report、model checklists、runtime report、SQLite index、SQLite index validation、knowledge runtime validation；第一個 failing step 會阻斷並不執行後續 steps。
- Phase 3 已新增 `runtime-generator-parity` fixtures：`generate-knowledge-runtime-report.rb`、`generate-model-context-report.rb` 與 `generate-model-checklists.rb` stdout 均與各自 Go-native builder byte-for-byte 一致；這些 fixtures 只建立替換前護欄，尚不改 production refresh path。
- Phase 3 已新增 `runtime-native-reports-opt-in` fixture：`runtime refresh --native-reports` 先以 Go 寫三個 Markdown reports，再執行剩餘 Ruby steps；預設 `runtime refresh` 不切換。
- Phase 3 已新增 `runtime-native-index-opt-in` fixture：`runtime refresh --native-reports --native-index` 先以 Go 寫三個 Markdown reports 與 SQLite index，再執行 Ruby validator steps；預設 `runtime refresh` 不切換，rollback 是不帶 `--native-index` 重新執行。
- Phase 3 已新增 `runtime-sqlite-index-parity` fixture：同一份 canonical source 分別產生 Ruby / Go temporary SQLite index，並比對 atoms / sources / edges / fts row counts、row-level content、source checksum map 與 `runtime` / `feedback` / `route` FTS hit counts；另以 nested feedback-history fixture 固定 recursive feedback parity。
- Phase 3 已新增 `runtime-compiler-snapshot-parity` harness：Ruby compiler 連續輸出兩個 temporary `runtime.db`，比對 required table row counts、generated surfaces content 與 compiler/schema metadata，忽略 dynamic timestamps；`runtime compile --native-compiler --db <temp>` output 先接上 generated surfaces / metadata snapshot parity，完整 row-count parity 留給真正 port compiler slice。

## 舊腳本 parity fixture

目的：防止新 CLI 漏掉既有腳本能力。

設定：

- 讀取 [`script-parity-inventory.md`](script-parity-inventory.md)。
- 對每個標記為 `native target` 或 `wrapper first` 的舊入口建立 fixture expectation。
- 對高風險路徑建立至少一個正向或負向案例：`init-project`、`goals`、`hooks install`、`sync-cursor-bundle`、`close-loop`、`runtime refresh`、`runtime compile`、`runtime validate`、`runtime query`。

斷言：

- 每個舊入口都有目標 CLI 命令、parity 狀態與必要驗證。
- 若舊入口會寫檔、動 git、寫 runtime DB、寫 tool mirror 或寫使用者設定，必須有 dry-run 或 fake-root fixture。
- `deferred` 與 `tool-specific` 條目不得被當成已替代；輸出必須明確標示不在目前 Phase 範圍。

## Windows path fixture

目的：證明 path normalization 不假設 POSIX shell。

設定：

- 使用 Windows-style paths、drive letters、空格與 mixed separators。
- 避免 shell-specific quoting 假設。

斷言：

- Path normalization 具穩定可重現結果。
- JSON 輸出使用穩定的 path 表示。
- 命令不需要 Git Bash / WSL。
- Phase 1 已以 `internal/pathutil` unit tests 覆蓋 drive letter、UNC path、mixed separators、spaces 與 relative path normalization。

## Fake home fixture

目的：防止意外寫入真實使用者設定。

設定：

- 使用 fixture paths 覆寫 home / config roots。
- 以 dry-run 執行 `init-project`、`goals`、`sync-cursor-bundle` 與 `hooks install`。

斷言：

- fixture root 外沒有寫入。
- JSON 輸出列出規劃寫入項目。
- 權限不足時回傳 `permission_denied`。
- Phase 2 已用 `init-project --dry-run` fake project tests 覆蓋 planned actions、既有檔案 conflict、`--force`、plain output 與無寫入 assertion。
- Phase 2 已用 temporary git repo tests 覆蓋 `close-loop --dry-run` clean repo、dirty owner group、active lock、merge/rebase state、unknown path block 與 missing Git block；commit / push write mode 仍阻斷，待 parity fixture 完成後再開啟。
- Phase 2 已用 temporary git repo tests 覆蓋 `hooks install --dry-run` source / target / conflict / force、merge state warning、missing Git 與無寫入 assertion；copy / chmod write mode 仍阻斷，待 parity fixture 完成後再開啟。
- Phase 2 已用 fake target tests 覆蓋 `sync-cursor-bundle --dry-run` 明確 target、target 不可在 repo 內、copy-fallback mirror strategy、技能 mirror planning 與無寫入 assertion；write mode 仍阻斷，待 managed mirror parity 完成後再開啟。

## 產物關卡

Phase 1 Go 實作前：

- [x] Fixtures 具備穩定名稱與預期斷言。
- [x] Missing Git fixture 已列入規劃測試套件。
- [x] Runtime DB assertion fixture 已連到命令契約。
- [x] 舊腳本 parity fixture 已連到舊腳本盤點。
- [x] Windows path fixture 不要求 POSIX shell。
- [x] Fake home fixture 可防止寫入真實本機設定。
