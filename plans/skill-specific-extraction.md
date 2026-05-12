# Skill-Specific Intelligence Extraction（遠期規劃）

> **Priority**: P4（最低優先級）
> **Status**: pending
> **啟動條件**: 所有 technique decomposition 完成 + Intelligence Extraction Pipeline 驗證成功後
> **Durable location**: 本文件（`plans/skill-specific-extraction.md`）
> **Roadmap 引用**: [`architecture/next-stage-upgrade-plan.md`](architecture/next-stage-upgrade-plan.md) 的 Durable Roadmap Goals

---

## 背景

目前的 Intelligence Extraction 策略是 **technique decomposition**（將單一 technique 拆解為 HOW TO DO → `analysis/`、HOW TO THINK → `intelligence/`），並以 flutter-dart-aot 作為 pilot 驗證此模式。

但不同 skill 的內容結構差異極大，無法用單一 extraction pipeline 處理所有 skill：

| Skill | 內容結構 | 適合的 extraction 策略 |
| --- | --- | --- |
| `apk-analysis` | `techniques/` 混合 workflow + intelligence + tools + failure 解讀 | **Decomposition**：拆成 workflow（analysis/）+ intelligence atoms（intelligence/） |
| `app-development-guidance` | `controls/`、`platforms/`、`languages/`、`implementation/`、`checklists/`、`process/` 各自獨立 | **Catalog + Direct Promotion**：controls/platforms/languages 適合 catalog，process 適合 workflow，checklists 適合 review gates |
| `travel-planning` | `WORKFLOW.md` + `DOCUMENTATION.md` 緊密耦合 | **Direct Promotion**：workflow → `workflow/travel-planning/`，templates → `analysis/travel/` |

---

## 目標

為每個 skill 設計專屬的 extraction strategy，包含：

1. **內容結構分析** — identify 哪些部分是 HOW TO DO、哪些是 HOW TO THINK、哪些是 templates/tools/checklists
2. **拆解策略選擇** — decomposition vs. catalog vs. direct promotion
3. **舊檔案標註規則** — `# Intelligence Extracted`（部分提取）或 `# Migrated To`（完全遷移）
4. **驗證方式** — pilot-driven，每個 skill 獨立驗證

---

## 執行順序

```
Phase A（目前）
  └─ Technique Decomposition Pilot（flutter-dart-aot ✅ 已完成）
       └─ 驗證 decomposition 模式是否有效

Phase B（近期）
  └─ 其餘 3 個 techniques decomposition（http-api、local-proxy、media-hls）
       └─ 完成所有 apk-analysis techniques 的 decomposition

Phase C（中期）
  └─ Intelligence Extraction Pipeline 抽象化
       └─ 從 pilot 經驗提煉出可重複的 extraction 流程

Phase D（遠期）← 本文件
  └─ Skill-Specific Extraction Strategy
       ├─ apk-analysis（techniques 已 decomposition 完成，但 SKILL.md 仍包含 workflow + tools + documentation）
       ├─ app-development-guidance（controls/platforms/languages/implementation/checklists/process 各自獨立）
       └─ travel-planning（WORKFLOW + DOCUMENTATION 緊密耦合）
```

---

## 各 Skill 初步分析

### 1. `apk-analysis`

**已完成的 extraction**：
- `techniques/` → `analysis/apk/techniques/`（catalog）
- `techniques/flutter-dart-aot/` → decomposition pilot（workflow + 4 intelligence atoms）
- `WORKFLOW.md` → `workflow/apk-analysis/execution-flow.md`
- `TOOLS.md` → `analysis/apk/tools-and-failures.md`
- `DOCUMENTATION.md` → `workflow/apk-analysis/artifact-gates.md`
- `RUNBOOK.md` → `runtime/onboarding/apk-analysis-setup.md` + `apk-analysis-completion.md`

**剩餘內容**：
- `SKILL.md` — 仍包含 Quick Start、Default Workflow、Required Output Style、Safety、Feedback Loop
- `FEEDBACK.md` — 已整併至共用規則
- `feedback_history/` — 尚未提取到 `feedback/` 層

**建議策略**：Decomposition（與目前模式一致）

### 2. `app-development-guidance`

**已完成的 extraction**：
- `controls/` → `analysis/app-development-guidance/controls-catalog.md`（catalog）
- `implementation/` → `analysis/app-development-guidance/implementation-catalog.md`（catalog）
- `platforms/` → `analysis/app-development-guidance/platforms-catalog.md`（catalog）
- `languages/` → `analysis/app-development-guidance/languages-catalog.md`（catalog）
- `checklists/` → `workflow/app-development-guidance/review-checklists.md`（direct promotion）
- `process/README.md` → `workflow/app-development-guidance/development-process.md`（direct promotion）
- `WORKFLOW.md` → `workflow/app-development-guidance/execution-flow.md`
- `DOCUMENTATION.md` → `workflow/app-development-guidance/artifact-gates.md`

**剩餘內容**：
- `SKILL.md` — 仍包含 When To Use、Quick Start、Default Workflow、Output Style、Feedback Loop
- `CHECKLIST.md` — 尚未提取
- `feedback_history/` — 尚未提取到 `feedback/` 層

**建議策略**：Catalog + Direct Promotion（controls/platforms/languages 適合 catalog，process/checklists 適合 direct promotion）

### 3. `travel-planning`

**已完成的 extraction**：
- `WORKFLOW.md` → `workflow/travel-planning/execution-flow.md`
- `DOCUMENTATION.md` → `workflow/travel-planning/`（templates 已提取）

**剩餘內容**：
- `SKILL.md` — 仍包含 When To Use、Quick Start、Default Workflow、Output Style、Feedback Loop

**建議策略**：Direct Promotion（workflow + templates 可直接 promotion，無需 decomposition）

---

## 舊檔案標註規則

| 狀態 | 標註 | 時機 |
| --- | --- | --- |
| 部分提取 | `# Intelligence Extracted — See <new path>` | 內容已被部分提取到新分層，但舊檔案仍有未提取的內容 |
| 完全遷移 | `# Migrated To — See <new path>` | 舊檔案所有內容已被完全覆蓋 |
| 已廢棄 | `# Deprecated — See <new path>` | 舊檔案不再需要，但保留以確保 tool adapter 相容性 |

---

## 驗證方式

每個 skill 的 extraction 完成後，需通過以下驗證：

1. **內容完整性**：新分層的內容是否能完整覆蓋舊 skill 的核心功能
2. **可發現性**：`knowledge/indexes/README.md` 和 `knowledge/runtime/routing-registry.yaml` 是否能 route 到新路徑
3. **向後相容**：舊 skill 入口是否仍可被 tool adapter 載入（或已有 documented replacement）
4. **AI 決策品質**：在實際 session 中使用新分層內容，AI 的決策品質是否不低於使用舊 skill

---

## 與其他文件的關係

- [`architecture/next-stage-upgrade-plan.md`](architecture/next-stage-upgrade-plan.md) — Durable Roadmap Goals 引用本文件
- [`governance/lifecycle/README.md`](governance/lifecycle/README.md) — Skills Deprecation Timeline（Phase D 對應本文件的執行結果）
- [`notes/intelligence-extraction-observations.md`](notes/intelligence-extraction-observations.md) — Pilot extraction 過程記錄
- [`plans/technique-intelligence-pilot.md`](plans/technique-intelligence-pilot.md) — Flutter-dart-aot pilot 執行計畫
