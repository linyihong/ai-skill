# BDD-lite 場景：Ai-skill CLI Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## 需求連結

- **來源**：Cross-platform Go script runtime plan
- **Actor / system role**：Ai-skill maintainer、agent、CI runner、desktop contributor、mobile control-plane user
- **行為邊界**：CLI 命令行為、依賴偵測、runtime 驗證、安全 close-loop
- **模糊項處置**：draft；場景完成審查前不得開始實作

## 場景：缺 Git 時阻斷 close-loop

**Given** 在桌面平台執行 `ai-skill close-loop --commit`
**When** PATH 中沒有可用的 `git`
**Then** 命令以 `missing_dependency` 結束
**And** 沒有檔案被 staged、committed、pushed 或 modified
**And** 輸出提示使用者安裝 Git。

## 場景：Doctor 回報 Git 需求

**Given** 執行 `ai-skill doctor --require-git --json`
**When** Git 缺失
**Then** JSON 輸出包含 `error.code = "missing_git"`
**And** exit code 在 Windows、macOS 與 Linux 上保持穩定。

## 場景：Runtime Compile 驗證 generated surface

**Given** runtime source file 已被修改
**When** `ai-skill runtime compile --assert-source <path> --assert-keyword <keyword>` 完成
**Then** `runtime/runtime.db` 的預期 generated surface 包含 source path 與 keyword
**And** keyword 缺失時 validation 失敗。

## 場景：Dry-run 防止副作用

**Given** 某個命令可寫入檔案或執行 git 操作
**When** 以 `--dry-run` 呼叫
**Then** 輸出規劃動作
**And** 不修改 tracked file、untracked file、git index、commit、hook、mirror 或 runtime DB。

## 場景：舊腳本 Parity 需有測試證據

**Given** `script-parity-inventory.md` 中某個舊入口標記為 `native target` 或 `wrapper first`
**When** 對應 `ai-skill` 命令進入實作
**Then** 該列必須有最低測試證據或 fixture 名稱
**And** 高風險路徑必須連到 `test-fixture-plan.md` 中的 fixture
**And** 若 parity 被標為 `deferred` 或 `tool-specific`，必須說明為何不阻擋目前 Phase。

## 場景：刪除 Legacy Surface 前先有 Go Migration Map

**Given** agent 準備刪除或移植舊 shell / Ruby / Python surface
**When** 對應 Go CLI command 或 package 尚未在 `legacy-to-go-migration-map.md` 記錄
**Then** 不得刪除 legacy surface
**And** 必須先補上 old surface、new owner、source-of-truth、validation evidence
**And** 再同步 `script-parity-inventory.md` 與 `legacy-script-disposition.md`。

## 場景：新增 Automation 必須先進 Go CLI

**Given** maintainer 或 agent 需要新增 repository automation
**When** automation 不是 Git hook adapter 或短期 binary bootstrap wrapper
**Then** 不得新增長期 `.sh`、`.rb` 或 `.py` entrypoint
**And** 必須先更新 `command-contract.md` 並實作對應 Go CLI command
**And** 若保留 legacy shell，shell 不得新增新功能，只能等待 Go write-mode parity 後刪除。

## 場景：Shell 入口完成 parity 後必須刪除

**Given** `scripts/init-new-project.sh`、`scripts/agent-goals.sh` 或 `scripts/ai-skill-close-loop.sh` 已有 Go write-mode parity
**When** 對應 fixture、文件與 release gate 通過
**Then** 舊 `.sh` 檔案必須刪除
**And** active docs 必須改指向 `ai-skill` Go CLI
**And** final grep 不得在 active docs 中保留舊 shell 執行指令。

## 場景：Copilot 新 Session 需要 Guided Startup Prompt

**Given** Copilot custom instructions 不能保證在每個新 session 自動完成 bootstrap
**When** 使用者執行 `ai-skill copilot start --project <project>`
**Then** 命令輸出可貼到 Copilot 第一則訊息的 bootstrap prompt
**And** prompt 指向 `<AI_SKILL_REPO>/CORE_BOOTSTRAP.md`、`<AI_SKILL_REPO>/runtime/core-bootstrap.yaml` 與 Copilot adapter
**And** 命令不修改檔案、不啟動 editor、不假裝 Copilot 有 hard enforcement。

## 場景：Git hook logic 必須進 Go

**Given** Git hook 需要 pre-commit 或 post-commit 行為
**When** hook 需要判斷 staged files、runtime compile、runtime validate 或 Cursor sync 狀態
**Then** 判斷邏輯必須在 `ai-skill hooks run ...` Go command
**And** hook file 若保留，只能作 repo-local binary adapter
**And** 不得在 hook shell 中新增 reusable automation logic。

## 場景：不安全 repo 狀態阻斷 commit

**Given** repository 處於 merge、rebase 或 cherry-pick 狀態
**When** 執行 `ai-skill close-loop --commit`
**Then** 命令以 `unsafe_repo_state` 結束
**And** 不建立 commit。

## 場景：Glossary entry 通過 schema

**Given** `knowledge/glossary/ai-skill.md` 包含一個 H2 heading（snake_case term）緊接 YAML code block，且 YAML 內 `term` / `status` / `meaning` / `affects` / `owner-layer` 皆合法
**When** 執行 `ai-skill glossary validate --json`
**Then** exit code 為 `0`
**And** JSON `status` 為 `success`
**And** `checks` 回報 entries / aliases / relations 計數
**And** `knowledge/glossary/` 內任何檔案沒有被修改。

## 場景：Glossary Retro-Own 自動強制 — happy path

**Given** staged diff 動到 `runtime/cognitive-modes-discovery.yaml`（加新 signal），且 `knowledge/glossary/ai-skill.md` 也在 staged 清單
**When** 執行 `ai-skill hooks run commit-msg`
**Then** `glossaryRetroOwn` validator 回傳空字串（pass）
**And** commit-msg hook exit code 為 `0`
**And** commit 順利建立。

## 場景：Glossary Retro-Own 阻斷 — framework surface 改動但無 glossary stage

**Given** staged diff 只動到 `runtime/cognitive-modes-discovery.yaml` 或 `runtime/economics/*.yaml` 或 `ecosystem/*.yaml`，`knowledge/glossary/ai-skill.md` **不在** staged 清單，且 commit message body 沒含 `[skip-glossary-retro-own]`
**When** 執行 `ai-skill hooks run commit-msg`
**Then** `glossaryRetroOwn` validator 回傳非空 error 訊息
**And** commit-msg hook exit code 為 `30`
**And** error 訊息引用 `runtime/cli-modification-policy.yaml` `gate.glossary.retro_own_required` 與 upstream plan。

## 場景：Glossary Retro-Own opt-out — 純 refactor / typo 變更

**Given** staged diff 動到 framework cognitive vocabulary surface，但變更為純 comment / typo / refactor（未引入新 term），commit message body 含獨立一行 `[skip-glossary-retro-own]`
**When** 執行 `ai-skill hooks run commit-msg`
**Then** `glossaryRetroOwn` validator 回傳空字串（pass）
**And** commit 順利建立。

## 場景：Glossary entry schema violation 阻斷

**Given** `knowledge/glossary/ai-skill.md` 內一個 entry 缺少 `owner-layer`、或 `term` 為 kebab-case、或 `aliases:` 內字串等於另一 entry 的 `term`、或 `excludes:` 引用不存在的 term、或 `introduced-by` 為 commit SHA
**When** 執行 `ai-skill glossary validate --json`
**Then** exit code 為 `30`（`validation_failed`）
**And** JSON `error.code = "validation_failed"`
**And** `checks` 列出每個 violation 的 entry path、term、rule id 與 remediation hint
**And** 命令不修改任何 `knowledge/glossary/` 內檔案。

## 場景：iOS native binary 不支援

**Given** 使用者要求在 iOS 上以下載的 native binary 執行 `ai-skill`
**When** 評估平台支援
**Then** 結果為不支援
**And** 建議選項為 App-contained runtime、Browser/WASM 或 SSH remote runner。

## 驗收條件

- 缺 Git 時不能產生半套 close-loop。
- Dry-run 命令不得修改 file system、git index、commits、hooks、mirrors 或 runtime DB。
- Runtime compile 能以 content assertions 證明 source-to-DB propagation。
- 每個被取代的舊腳本能力都有 parity disposition 與最低測試證據。
- 每個被刪除的 legacy surface 都能從 migration map 反查新的 Go owner 與 source-of-truth。
- 新增 repository automation 預設走 Go CLI，不新增長期 shell / Ruby / Python entrypoint。
- Mobile support 決策不承諾 iOS 任意 native binary。

## 驗證目標

- **證明類型**：fixture-backed automated tests
- **測試 / fixture / checklist**：[`test-fixture-plan.md`](test-fixture-plan.md)
- **限制**：這些場景尚未證明 performance、release signing、mobile app feasibility 或完整 compiler parity。

## 回歸範圍

- [ ] 既有 shell script 行為必須保留到 replacement 完成驗證。
- [ ] 需要新增缺 Git、不安全 repo、dry-run、舊腳本 parity、runtime assertion 與 iOS 不支援決策測試。
- [ ] 需要測試資料 / fixtures：temporary repo、PATH isolation、fake home、runtime source fixture、legacy script parity fixture。
