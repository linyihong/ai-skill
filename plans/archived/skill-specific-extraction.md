# Skill-Specific Intelligence Extraction（執行結果）

> **Priority**: P4（最低優先級）
> **Status**: ✅ completed（Phase 33）
> **啟動條件**: 所有 technique decomposition 完成 + Intelligence Extraction Pipeline 驗證成功後
> **Durable location**: 本文件（`plans/skill-specific-extraction.md`）
> **Roadmap 引用**: [`architecture/next-stage-upgrade-plan.md`](architecture/next-stage-upgrade-plan.md) 的 Durable Roadmap Goals

---

## 背景

目前的 Intelligence Extraction 策略是 **technique decomposition**（將單一 technique 拆解為 HOW TO DO → `analysis/`、HOW TO THINK → `intelligence/`），並以 flutter-dart-aot 作為 pilot 驗證此模式。

但不同 skill 的內容結構差異極大，無法用單一 extraction pipeline 處理所有 skill：

| Skill | 內容結構 | 使用的 extraction 策略 |
| --- | --- | --- |
| `apk-analysis` | `techniques/` 混合 workflow + intelligence + tools + failure 解讀 | **Decomposition**：拆成 workflow（analysis/）+ intelligence atoms（intelligence/） |
| `app-development-guidance` | `controls/`、`platforms/`、`languages/`、`implementation/`、`checklists/`、`process/` 各自獨立 | **Catalog + Direct Promotion**：controls/platforms/languages 適合 catalog，process 適合 workflow，checklists 適合 review gates |
| `travel-planning` | `WORKFLOW.md` + `DOCUMENTATION.md` + `TOOLS.md` + `README.md` 各自獨立 | **Direct Promotion**：workflow → `workflow/travel-planning/`，templates → `analysis/travel/` |

---

## 目標

為每個 skill 設計專屬的 extraction strategy，包含：

1. **內容結構分析** — identify 哪些部分是 HOW TO DO、哪些是 HOW TO THINK、哪些是 templates/tools/checklists
2. **拆解策略選擇** — decomposition vs. catalog vs. direct promotion
3. **舊檔案標註規則** — `# Intelligence Extracted`（部分提取）或 `# Migrated To`（完全遷移）
4. **驗證方式** — pilot-driven，每個 skill 獨立驗證

---

## 執行結果

### Phase A：Technique Decomposition Pilot（flutter-dart-aot）✅

- 完成於 Phase 28
- 驗證 decomposition 模式有效
- 產出：`analysis/apk/workflows/frida-hook-flow.md` + 4 intelligence atoms

### Phase B：其餘 3 個 techniques decomposition ✅

- 完成於 Phase 29
- http-api、local-proxy、media-hls 全部 decomposition 完成
- 產出：3 workflow files + 4 intelligence atoms

### Phase C：Intelligence Extraction Pipeline 抽象化 ✅

- 完成於 Phase 31
- 從 pilot 經驗提煉出 7-step pipeline
- 產出：`governance/lifecycle/intelligence-extraction-pipeline.md`

### Phase D：Skill-Specific Extraction Strategy ✅

- 完成於 Phase 33
- 實際執行內容見下方各 skill 詳細記錄

### Phase E：Skills Deprecation Phase C ✅

- 完成於 Phase 35（2026-05-12）
- 已刪除 10 個舊 technique 檔案：
  - `skills/apk-analysis/techniques/`：4 個子目錄（flutter-dart-aot/、http-api/、local-proxy/、media-hls/）+ 1 個 README.md
  - `analysis/apk/techniques/`：4 個 .md（flutter-dart-aot.md、http-api.md、local-proxy.md、media-hls.md）+ 1 個 README.md
- 刪除前已確認 Phase C 檢查清單 7 項條件全部滿足
- 15+ 個引用舊路徑的檔案已更新為指向新路徑

---

## 各 Skill 最終狀態

### 1. `apk-analysis` — 全部提取完成 ✅

| 原始檔案 | 目標路徑 | 策略 | 狀態 |
|---------|---------|------|------|
| `techniques/` | `analysis/apk/techniques/` | Catalog | ✅ |
| `techniques/flutter-dart-aot/` | `analysis/apk/workflows/frida-hook-flow.md` + 4 intelligence atoms | Decomposition | ✅ |
| `techniques/http-api/` | `analysis/apk/workflows/http-api-documentation-flow.md` + 2 intelligence atoms | Decomposition | ✅ |
| `techniques/local-proxy/` | `analysis/apk/workflows/local-proxy-hook-flow.md` + 2 intelligence atoms | Decomposition | ✅ |
| `techniques/media-hls/` | `analysis/apk/workflows/media-hls-analysis-flow.md` + 1 intelligence atom | Decomposition | ✅ |
| `WORKFLOW.md` | `workflow/apk-analysis/execution-flow.md` | Direct Promotion | ✅ |
| `TOOLS.md` | `analysis/apk/tools-and-failures.md` | Direct Promotion | ✅ |
| `DOCUMENTATION.md` | `workflow/apk-analysis/artifact-gates.md` | Direct Promotion | ✅ |
| `RUNBOOK.md` | `runtime/onboarding/apk-analysis-setup.md` + `apk-analysis-completion.md` | Direct Promotion | ✅ |
| `SKILL.md` | Slimmed to routing（~55 lines） | SKILL.md Decomposition | ✅ |
| `FEEDBACK.md` | 已整併至共用規則 | Consolidation | ✅ |
| `feedback_history/` | `feedback/extraction/apk-analysis-index.md` | Index + Annotate | ✅ |

**剩餘內容**：無

### 2. `app-development-guidance` — 全部提取完成 ✅

| 原始檔案 | 目標路徑 | 策略 | 狀態 |
|---------|---------|------|------|
| `controls/` | `analysis/app-development-guidance/controls-catalog.md` | Catalog | ✅ |
| `implementation/` | `analysis/app-development-guidance/implementation-catalog.md` | Catalog | ✅ |
| `platforms/` | `analysis/app-development-guidance/platforms-catalog.md` | Catalog | ✅ |
| `languages/` | `analysis/app-development-guidance/languages-catalog.md` | Catalog | ✅ |
| `checklists/` | `workflow/app-development-guidance/review-checklist.md` | Direct Promotion | ✅ |
| `process/README.md` | `workflow/app-development-guidance/development-process.md` | Direct Promotion | ✅ |
| `WORKFLOW.md` | `workflow/app-development-guidance/execution-flow.md` | Direct Promotion | ✅ |
| `DOCUMENTATION.md` | `workflow/app-development-guidance/artifact-gates.md` | Direct Promotion | ✅ |
| `SKILL.md` | Slimmed to routing（~65 lines） | SKILL.md Decomposition | ✅ |
| `CHECKLIST.md` | `workflow/app-development-guidance/review-checklist.md` | Direct Promotion | ✅ |
| `FEEDBACK.md` | 已整併至共用規則 | Consolidation | ✅ |
| `feedback_history/` | `feedback/extraction/app-development-guidance-index.md` | Index + Annotate | ✅ |

**剩餘內容**：無

### 3. `travel-planning` — 全部提取完成 ✅

| 原始檔案 | 目標路徑 | 策略 | 狀態 |
|---------|---------|------|------|
| `WORKFLOW.md` | `workflow/travel-planning/execution-flow.md` | Direct Promotion | ✅ |
| `DOCUMENTATION.md` | `workflow/travel-planning/artifact-gates.md` | Direct Promotion | ✅ |
| `TOOLS.md` | `analysis/travel/sources-and-tools.md` | Direct Promotion | ✅ |
| `README.md` | `analysis/travel/README.md` | Direct Promotion | ✅ |
| `SKILL.md` | Slimmed to routing（~55 lines） | SKILL.md Decomposition | ✅ |
| `FEEDBACK.md` | 已整併至共用規則 | Consolidation | ✅ |

**剩餘內容**：無

---

## 舊檔案標註狀態

| 檔案 | 標註 | 狀態 |
|------|------|------|
| `skills/apk-analysis/techniques/flutter-dart-aot/README.md` | `# Intelligence Extracted` | ✅ |
| `skills/apk-analysis/techniques/http-api/README.md` | `# Intelligence Extracted` | ✅ |
| `skills/apk-analysis/techniques/local-proxy/README.md` | `# Intelligence Extracted` | ✅ |
| `skills/apk-analysis/techniques/media-hls/README.md` | `# Intelligence Extracted` | ✅ |
| `analysis/apk/techniques/flutter-dart-aot.md` | `# Intelligence Extracted` | ✅ |
| `analysis/apk/techniques/http-api.md` | `# Intelligence Extracted` | ✅ |
| `analysis/apk/techniques/local-proxy.md` | `# Intelligence Extracted` | ✅ |
| `analysis/apk/techniques/media-hls.md` | `# Intelligence Extracted` | ✅ |
| `skills/app-development-guidance/CHECKLIST.md` | `# Extracted — See workflow/app-development-guidance/review-checklist.md` | ✅ |
| `skills/travel-planning/TOOLS.md` | `# Extracted — See analysis/travel/sources-and-tools.md` | ✅ |
| `skills/travel-planning/README.md` | `# Extracted — See analysis/travel/README.md` | ✅ |
| `skills/apk-analysis/feedback_history/`（101 files） | `# Extracted — See feedback/extraction/apk-analysis-index.md` | ✅ |
| `skills/app-development-guidance/feedback_history/`（40 files） | `# Extracted — See feedback/extraction/app-development-guidance-index.md` | ✅ |

---

## 驗證結果

| 驗證項目 | 結果 | 說明 |
|---------|------|------|
| 內容完整性 | ✅ Pass | 所有舊 skill 的核心功能已被新分層完整覆蓋 |
| 可發現性 | ✅ Pass | `knowledge/indexes/README.md` 和 `knowledge/runtime/routing-registry.yaml` 可 route 到新路徑 |
| 向後相容 | ✅ Pass | 舊 skill 入口已不再作為 active entrypoint，新分層已完全承接 |
| AI 決策品質 | ✅ Pass | 新分層內容的結構化程度更高，AI 決策品質不低於使用舊 skill |

---

## 與其他文件的關係

- [`architecture/next-stage-upgrade-plan.md`](architecture/next-stage-upgrade-plan.md) — Durable Roadmap Goals 引用本文件
- [`governance/lifecycle/README.md`](governance/lifecycle/README.md) — Skills Deprecation Timeline（Phase D 對應本文件的執行結果）
- [`notes/intelligence-extraction-observations.md`](notes/intelligence-extraction-observations.md) — Pilot extraction 過程記錄
- [`plans/technique-intelligence-pilot.md`](plans/technique-intelligence-pilot.md) — Flutter-dart-aot pilot 執行計畫
