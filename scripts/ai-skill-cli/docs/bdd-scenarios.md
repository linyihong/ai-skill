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

## 場景：不安全 repo 狀態阻斷 commit

**Given** repository 處於 merge、rebase 或 cherry-pick 狀態
**When** 執行 `ai-skill close-loop --commit`
**Then** 命令以 `unsafe_repo_state` 結束
**And** 不建立 commit。

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
- Mobile support 決策不承諾 iOS 任意 native binary。

## 驗證目標

- **證明類型**：fixture-backed automated tests
- **測試 / fixture / checklist**：[`test-fixture-plan.md`](test-fixture-plan.md)
- **限制**：這些場景尚未證明 performance、release signing、mobile app feasibility 或完整 compiler parity。

## 回歸範圍

- [ ] 既有 shell script 行為必須保留到 replacement 完成驗證。
- [ ] 需要新增缺 Git、不安全 repo、dry-run、舊腳本 parity、runtime assertion 與 iOS 不支援決策測試。
- [ ] 需要測試資料 / fixtures：temporary repo、PATH isolation、fake home、runtime source fixture、legacy script parity fixture。
