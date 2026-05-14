# Plans（計畫目錄）

## 目錄規則

| 子目錄 | 用途 | 生命週期 |
|--------|------|---------|
| [`active/`](active/) | 進行中或待審閱的計畫（draft / in-progress） | 完成後搬移至 `archived/` |
| [`archived/`](archived/) | 已執行完成的計畫（執行結果記錄） | 永久保留，作為決策記錄 |

## 原則

1. **`active/` 只放尚未開始或正在執行的計畫** — 一旦計畫執行完畢，立即搬移至 `archived/`
2. **`archived/` 的計畫不刪除** — 作為歷史決策記錄，可供日後查閱
3. **計畫檔案命名規則**：`<slug>.md`，slug 需能反映計畫核心目標
4. **每個計畫必須在檔頭標註狀態**：`draft` / `in-progress` / `completed`
5. **計畫完成後，若從中提煉出可重用的系統經驗，應建立對應的 intelligence atom**

## Plan 完成閉環（Plan Completion Closure）

當一個 plan 的所有項目都標記為完成（`✅`）時，agent **必須**執行以下閉環檢查：

### 檢查清單

| # | 檢查項目 | 說明 |
|---|---------|------|
| 1 | **確認所有項目已完成** | 檢查 plan 中所有 task 是否都標記為 `✅`，無遺漏項目 |
| 2 | **執行 validator** | 若 plan 涉及 `knowledge/`、`validation/`、`intelligence/` 等層，執行 `ruby scripts/refresh-knowledge-runtime.rb` |
| 3 | **檢查連動更新** | 依 [`shared-rules/linked-updates.md`](../shared-rules/linked-updates.md) 檢查 plan 改動是否需要同步其他檔案 |
| 4 | **更新 plans/README.md 狀態** | 將本 plan 在[目前狀態](#目前狀態)表格中的狀態改為 `✅ completed` |
| 5 | **搬移至 archived/** | 將 plan 檔案從 `active/` 搬移至 `archived/`，檔名與內容不變 |
| 6 | **Commit & push** | 提交搬移與狀態更新，並推送 |
| 7 | **最終確認** | 執行 `git status --short --branch` 確認工作樹乾淨 |

### 不搬移的例外情況

若 plan 符合以下任一條件，可留在 `active/` 但標註 `✅ completed`：

- Plan 是**持續生效的基礎建設**（如 validation gate、pre-commit hook），未來可能擴充新 Phase
- Plan 的 scope 是 ongoing 的維護性任務，沒有明確的「完成」邊界

例外情況必須在 plan 檔頭或 `plans/README.md` 表格中說明原因。

## 目前狀態

| 檔案 | 狀態 | 說明 |
|------|------|------|
| [`active/cognitive-boundary-system.md`](active/cognitive-boundary-system.md) | ⏳ draft | Cognitive Boundary System 整合計畫，待審閱後開始實作 |
| [`active/shared-rules-to-enforcement-migration.md`](active/shared-rules-to-enforcement-migration.md) | 📄 draft | shared-rules/ → enforcement/ 搬遷計畫，含 Layer Responsibility Contract |
| [`active/enforcement-layer-enhancement.md`](active/enforcement-layer-enhancement.md) | 📄 draft | enforcement/ 後續強化計畫：Metadata Spec、Rule Graph、Activation Engine、Conflict Matrix、Deprecation Lifecycle |
| [`archived/knowledge-runtime-validation-gate.md`](archived/knowledge-runtime-validation-gate.md) | ✅ completed | Part 1: Validation Gate 已完成；Part 2: UI Operation Intelligence Extraction 已完成 |
| [`archived/technique-intelligence-pilot.md`](archived/technique-intelligence-pilot.md) | ✅ completed | Phase 28：Technique → Intelligence Pilot（flutter-dart-aot） |
| [`archived/skill-specific-extraction.md`](archived/skill-specific-extraction.md) | ✅ completed | Phase 33：Skill-Specific Intelligence Extraction |
| [`archived/ai-decision-contract-testing.md`](archived/ai-decision-contract-testing.md) | ✅ completed | AI Decision Contract Testing 框架設計與實作 |
| [`archived/context-cost-optimization.md`](archived/context-cost-optimization.md) | ✅ completed | Phase 1：Context Cost Optimization 執行計畫（原 architecture/） |
| [`archived/apk-analysis-pilot-migration.md`](archived/apk-analysis-pilot-migration.md) | ✅ completed | APK Analysis Pilot Migration 狀態圖（原 architecture/） |

## 與其他層的關係

- [`architecture/next-stage-upgrade-plan.md`](../architecture/next-stage-upgrade-plan.md) — 全局升級路線圖，`plans/` 中的計畫是路線圖的具體執行計畫
- [`governance/lifecycle/README.md`](../governance/lifecycle/README.md) — Skills Deprecation Timeline 等生命週期規則
- [`intelligence/engineering/agent-architecture/`](../intelligence/engineering/agent-architecture/) — 從已完成計畫中提煉的系統經驗結晶
