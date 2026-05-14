# shared-rules/ → enforcement/ 搬遷計畫

Status: `draft`
Created: 2026-05-14
Related: [`shared-rules/README.md`](../../shared-rules/README.md), [`runtime/README.md`](../../runtime/README.md), [`governance/README.md`](../../governance/README.md)

---

## 1. 背景與動機

### 1.1 為什麼要改名

`shared-rules/` 這個名稱在架構上已不再準確：

| 面向 | 現狀問題 | 改名後 |
|------|---------|--------|
| **語意** | 「共用規則」暗示「可選的、建議性的」 | `enforcement/` 明確表示「這是必須執行的政策」 |
| **架構定位** | 名稱沒有反映它在三層模型中的角色 | `enforcement/` 是 governance → enforcement → runtime 的中間層 |
| **與 runtime/ 的邊界** | `shared-rules/` 和 `runtime/` 名稱上沒有層級關係 | `enforcement/` 和 `runtime/` 形成明確的 policy/engine 分離 |
| **未來擴充** | 名稱無法承載 rule metadata、activation condition、runtime cost 等進階概念 | `enforcement/` 可以自然擴充為完整的 runtime enforcement layer |

### 1.2 三層架構模型

```
┌─────────────────────────────────────────────────────┐
│  governance/  (Policy Architecture / WHY)           │
│  ── lifecycle, validation, cleanup, dependency      │
│  定義「知識如何被治理」的架構設計                      │
├─────────────────────────────────────────────────────┤
│  enforcement/ (Runtime Enforcement Layer / RULES)   │
│  ── rule-weight, dependency-reading, linked-updates │
│  定義「AI agent 必須遵守的執行政策」                   │
├─────────────────────────────────────────────────────┤
│  runtime/    (Runtime Engine Layer / LOADING)       │
│  ── routing, onboarding, context, health            │
│  定義「AI 系統如何運作」— dynamic loading、routing   │
└─────────────────────────────────────────────────────┘
```

### 1.3 Layer Responsibility Contract

| 層級 | 路徑 | 責任 | 不負責 |
|------|------|------|--------|
| **Policy Architecture** | [`governance/`](../../governance/README.md) | 知識治理架構設計：lifecycle、validation gate、dependency maintenance、cleanup strategy | 不存放可執行的 agent 政策規則；不定義 runtime loading 行為 |
| **Runtime Enforcement** | [`enforcement/`](../../enforcement/README.md) | AI agent 必須遵守的執行政策：rule-weight、dependency-reading、linked-updates、failure-learning、sanitization、feedback-lessons | 不設計 runtime engine；不存放 tool-specific 設定；不存放 governance 架構文件 |
| **Runtime Engine** | [`runtime/`](../../runtime/README.md) | Context routing、dynamic loading、activation rules、onboarding quickstart、orchestration | 不存放可執行 policy 正文；不存放 tool-specific hook 細節 |

---

## 2. 影響分析

### 2.1 受影響檔案分類

| 類別 | 檔案數量 | 影響程度 | 說明 |
|------|---------|---------|------|
| **enforcement/ 內部** | ~26 檔 | 🔴 高 | 所有內部相對路徑引用需更新（`../` → `../` 不變，但 `shared-rules/` → `enforcement/`） |
| **根目錄入口** | 2 檔 | 🔴 高 | `README.md`、`CORE_BOOTSTRAP.md` 直接引用 `shared-rules/` |
| **governance/** | ~5 檔 | 🔴 高 | `governance/README.md` 引用 shared-rules 作為 migration source |
| **runtime/** | ~3 檔 | 🟡 中 | `runtime/README.md` 提及 shared-rules 作為 policy layer |
| **scripts/** | 4 檔 | 🔴 高 | `validate-knowledge-runtime.rb`、`ai-skill-close-loop.sh`、`sync-cursor-bundle.sh`、`init-new-project.sh` |
| **ai-tools/** | ~4 檔 | 🟡 中 | `README.md`、`agent-onboarding.md`、`agent/roo.md`、`agent/cursor.md`、`agent/claude.md` |
| **skills/** | ~10+ 檔 | 🟡 中 | 各 SKILL.md、README.md、模板中的 shared-rules 引用 |
| **feedback/** | ~5 檔 | 🟡 中 | `README.md`、`pipeline/README.md`、`extraction/*.md`、`history/README.md` |
| **knowledge/** | ~4 檔 | 🟡 中 | `README.md`、`indexes/README.md`、`graphs/README.md`、`runtime/README.md` |
| **intelligence/** | ~3 檔 | 🟢 低 | `README.md`、`engineering/*/README.md` |
| **workflow/** | ~5 檔 | 🟢 低 | 各 domain 的 README.md |
| **analysis/** | ~2 檔 | 🟢 低 | `README.md`、`*/README.md` |
| **metadata/** | ~5 檔 | 🟢 低 | `README.md`、`rules/README.md`、`ranking/README.md` 等 |
| **models/** | ~2 檔 | 🟢 低 | `README.md`、`profiles/README.md` |
| **memory/** | ~2 檔 | 🟢 低 | `README.md`、`failure/README.md` |
| **plans/active/** | ~1 檔 | 🟢 低 | `next-stage-upgrade-plan.md` |
| **plans/archived/** | ~2 檔 | 🟢 低 | `context-cost-optimization.md`、`apk-analysis-pilot-migration.md` |
| **plans/** | ~2 檔 | 🟢 低 | `README.md`、`active/*.md` |
| **anti-patterns/** | ~1 檔 | 🟢 低 | `README.md` |
| **validation/** | ~1 檔 | 🟢 低 | `README.md` |

### 2.2 內部檔案路徑更新

`enforcement/` 內部檔案的相對路徑引用模式：

| 原始路徑 | 新路徑 | 變更類型 |
|---------|--------|---------|
| `../shared-rules/xxx.md` | `../enforcement/xxx.md` | 路徑前綴變更 |
| `shared-rules/xxx.md`（同層級引用） | `xxx.md`（不變，因同目錄） | 不變 |
| `../shared-rules/failure-patterns/xxx.md` | `../enforcement/failure-patterns/xxx.md` | 路徑前綴變更 |
| `../../shared-rules/xxx.md` | `../../enforcement/xxx.md` | 路徑前綴變更 |

### 2.3 外部檔案路徑更新

全庫搜尋 `shared-rules/` 的引用模式：

| 引用模式 | 出現位置範例 | 新模式 |
|---------|------------|--------|
| `shared-rules/README.md` | 根 README.md、CORE_BOOTSTRAP.md | `enforcement/README.md` |
| `shared-rules/xxx.md` | 各層 README、linked-updates 表格 | `enforcement/xxx.md` |
| `shared-rules/failure-patterns/` | validator、failure-patterns 索引 | `enforcement/failure-patterns/` |
| `shared-rules`（字串） | close-loop.sh owner group、sync script | `enforcement` |
| `shared-rules/`（目錄路徑） | sync-cursor-bundle.sh bundle 路徑 | `enforcement/` |

---

## 3. 搬遷策略：Phased Approach

### Phase 0：準備（前置檢查）

- [ ] 確認目前 `git status --short --branch` 乾淨
- [ ] 確認無 active lock 或 pending commit
- [ ] 讀取本計畫全文
- [ ] 讀取 [`shared-rules/failure-patterns/shared-rules-architecture-drift.md`](../../shared-rules/failure-patterns/shared-rules-architecture-drift.md) 的 Prevention Gate

### Phase 1：建立 enforcement/ 目錄（安全複製）

**目標**：建立 `enforcement/` 作為 `shared-rules/` 的完全鏡像，不刪除任何既有檔案。

- [ ] 建立 `enforcement/` 目錄
- [ ] 複製 `shared-rules/*.md` 到 `enforcement/`（17 個 .md 檔案）
- [ ] 複製 `shared-rules/failure-patterns/` 到 `enforcement/failure-patterns/`（9 個 failure pattern + README.md）
- [ ] 驗證：`diff -r shared-rules/ enforcement/` 應無差異（除 `.gitkeep` 等非 md 檔案）

**向後相容**：`shared-rules/` 保留不動，所有既有連結繼續有效。

### Phase 2：更新 enforcement/ 內部路徑

**目標**：`enforcement/` 內部的相對路徑引用全部指向正確的新路徑。

- [ ] 更新 `enforcement/README.md`：
  - 標題改為「強制執行規則（分類索引）」
  - 所有 `../shared-rules/` 改為 `../enforcement/`
  - 所有 `shared-rules/`（同層級引用）保持不變
  - 新增 Layer Responsibility Contract 說明
- [ ] 更新 `enforcement/linked-updates.md`：
  - 第 7 條的 `shared-rules/` 改為 `enforcement/`
  - 常見連動關係表中所有 `shared-rules/` 改為 `enforcement/`
  - 第 51 行「架構重構」連動關係中的 `shared-rules/` 改為 `enforcement/`
- [ ] 更新 `enforcement/dependency-reading.md`：
  - 所有 `shared-rules/` 路徑改為 `enforcement/`
- [ ] 更新 `enforcement/rule-weight.md`：
  - 結尾 `← [Back to shared rules index](README.md)` 改為 `← [Back to enforcement index](README.md)`
- [ ] 更新 `enforcement/failure-patterns/README.md`：
  - 結尾 `← [Back to shared rules index](../README.md)` 改為 `← [Back to enforcement index](../README.md)`
- [ ] 更新 `enforcement/failure-patterns/shared-rules-architecture-drift.md`：
  - 所有 `shared-rules/` 路徑改為 `enforcement/`
  - 檔名本身不改（內容描述的是架構漂移 pattern，仍適用）
- [ ] 更新 `enforcement/` 其餘檔案中的路徑引用
- [ ] 驗證：`grep -rn "shared-rules" enforcement/` — 應無 `shared-rules/` 路徑引用（保留 `shared-rules` 字串在說明文字中可接受，但路徑引用必須更新）

### Phase 3：更新全庫外部引用

**目標**：所有非 `shared-rules/` 的檔案中的 `shared-rules/` 路徑引用改為 `enforcement/`。

**根目錄：**
- [ ] 更新 `README.md`：
  - OS Layout 表格：`shared-rules/` → `enforcement/`
  - 說明文字中的 `shared-rules/` → `enforcement/`
- [ ] 更新 `CORE_BOOTSTRAP.md`：
  - 必讀規則表格：`shared-rules/rule-weight.md` → `enforcement/rule-weight.md`
  - 其餘 `shared-rules/` 引用 → `enforcement/`

**governance/:**
- [ ] 更新 `governance/README.md`：
  - 「第一批候選遷移來源」中的 `shared-rules/` → `enforcement/`
  - 與既有層關係中的 `shared-rules/` → `enforcement/`
- [ ] 更新 `governance/lifecycle/intelligence-extraction-pipeline.md`：
  - Step 7a 中的 `shared-rules/` → `enforcement/`
- [ ] 更新 `governance/dependency/README.md`（若有 shared-rules 引用）
- [ ] 更新 `governance/lifecycle/knowledge-update-flow.md`（若有 shared-rules 引用）

**runtime/:**
- [ ] 更新 `runtime/README.md`：
  - 「不放什麼」：`shared-rules/` → `enforcement/`
  - 「與既有層的關係」：`shared-rules/` → `enforcement/`
- [ ] 更新 `runtime/routing/README.md`（若有 shared-rules 引用）
- [ ] 更新 `runtime/onboarding/README.md`（若有 shared-rules 引用）

**scripts/:**
- [ ] 更新 `scripts/validate-knowledge-runtime.rb`：
  - 第 452 行：`scan_dirs = %w[intelligence workflow analysis shared-rules]` → `scan_dirs = %w[intelligence workflow analysis enforcement]`
  - 第 640 行：`patterns_dir = ROOT + "shared-rules/failure-patterns"` → `patterns_dir = ROOT + "enforcement/failure-patterns"`
  - 第 680 行註解：`shared-rules/content-layering.md` → `enforcement/content-layering.md`
- [ ] 更新 `scripts/ai-skill-close-loop.sh`：
  - 第 171 行：`shared-rules/*|README.md|.gitignore) echo "shared" ;;` → `enforcement/*|README.md|.gitignore) echo "shared" ;;`
- [ ] 更新 `scripts/sync-cursor-bundle.sh`：
  - 所有 `shared-rules` 路徑 → `enforcement`
  - BUNDLE_RULES、CURSOR_SHARED 變數
  - symlink 目標路徑
- [ ] 更新 `scripts/init-new-project.sh`：
  - Custom Instructions 中的 `shared-rules/` → `enforcement/`
  - 知識更新流程中的 `shared-rules` → `enforcement`

**ai-tools/:**
- [ ] 更新 `ai-tools/README.md`（若有 shared-rules 引用）
- [ ] 更新 `ai-tools/agent-onboarding.md`（若有 shared-rules 引用）
- [ ] 更新 `ai-tools/agent/roo.md`（若有 shared-rules 引用）
- [ ] 更新 `ai-tools/agent/cursor.md`（若有 shared-rules 引用）
- [ ] 更新 `ai-tools/agent/claude.md`（若有 shared-rules 引用）

**skills/:**
- [ ] 更新 `skills/ADDING_SKILLS.md` 中的 shared-rules 引用
- [ ] 更新 `skills/_template/SKILL.md` 中的 shared-rules 引用
- [ ] 更新各 skill 的 README.md、SKILL.md、WORKFLOW.md 中的 shared-rules 引用

**feedback/:**
- [ ] 更新 `feedback/README.md`（若有 shared-rules 引用）
- [ ] 更新 `feedback/pipeline/README.md`（若有 shared-rules 引用）
- [ ] 更新 `feedback/extraction/apk-analysis-index.md`（若有 shared-rules 引用）
- [ ] 更新 `feedback/extraction/development-guidance-index.md`（若有 shared-rules 引用）
- [ ] 更新 `feedback/history/README.md`（若有 shared-rules 引用）

**knowledge/:**
- [ ] 更新 `knowledge/README.md`（若有 shared-rules 引用）
- [ ] 更新 `knowledge/runtime/README.md`（若有 shared-rules 引用）
- [ ] 更新 `knowledge/indexes/README.md`（若有 shared-rules 引用）
- [ ] 更新 `knowledge/graphs/README.md`（若有 shared-rules 引用）

**intelligence/:**
- [ ] 更新 `intelligence/README.md`（若有 shared-rules 引用）
- [ ] 更新 `intelligence/engineering/README.md`（若有 shared-rules 引用）
- [ ] 更新 `intelligence/engineering/agent-architecture/README.md`（若有 shared-rules 引用）

**workflow/:**
- [ ] 更新 `workflow/README.md`（若有 shared-rules 引用）
- [ ] 更新各 domain 的 README.md（若有 shared-rules 引用）

**analysis/:**
- [ ] 更新 `analysis/README.md`（若有 shared-rules 引用）
- [ ] 更新 `analysis/development-guidance/README.md`（若有 shared-rules 引用）

**metadata/:**
- [ ] 更新 `metadata/README.md`（若有 shared-rules 引用）
- [ ] 更新 `metadata/rules/README.md`（若有 shared-rules 引用）
- [ ] 更新 `metadata/ranking/README.md`（若有 shared-rules 引用）

**models/:**
- [ ] 更新 `models/README.md`（若有 shared-rules 引用）
- [ ] 更新 `models/profiles/README.md`（若有 shared-rules 引用）

**memory/:**
- [ ] 更新 `memory/README.md`（若有 shared-rules 引用）
- [ ] 更新 `memory/failure/README.md`（若有 shared-rules 引用）

**architecture/:**
- [ ] 更新 `plans/active/next-stage-upgrade-plan.md`（若有 shared-rules 引用）
- [ ] 更新 `plans/archived/apk-analysis-pilot-migration.md`（若有 shared-rules 引用）

**plans/:**
- [ ] 更新 `plans/README.md`：
  - 第 28 行：`shared-rules/linked-updates.md` → `enforcement/linked-updates.md`
- [ ] 更新 `plans/active/knowledge-runtime-validation-gate.md`（若有 shared-rules 引用）

**anti-patterns/:**
- [ ] 更新 `anti-patterns/README.md`（若有 shared-rules 引用）

**validation/:**
- [ ] 更新 `validation/README.md`（若有 shared-rules 引用）

### Phase 4：更新 Validator

**目標**：`scripts/validate-knowledge-runtime.rb` 中的 `shared-rules` 路徑全部更新。

- [ ] 更新 `validate_language_consistency` 的 scan_dirs
- [ ] 更新 `validate_failure_pattern_validator_coverage` 的 patterns_dir
- [ ] 更新註解中的 shared-rules 路徑
- [ ] 執行 validator 確認無誤：`ruby scripts/refresh-knowledge-runtime.rb`

### Phase 5：驗證

**目標**：確認所有連結正確，無遺漏。

- [ ] 執行 `grep -rn "shared-rules/" --include="*.md" --include="*.rb" --include="*.sh" --include="*.yaml" .` — 應只出現：
  - `shared-rules/` 目錄本身的檔案（這些是舊檔案，將在 Phase 6 移除）
  - 說明文字中的「shared-rules」字串（非路徑引用）
- [ ] 執行 `grep -rn "enforcement/" --include="*.md" --include="*.rb" --include="*.sh" --include="*.yaml" .` — 應出現所有新路徑
- [ ] 執行 `ruby scripts/refresh-knowledge-runtime.rb` — 所有 validator 通過
- [ ] 隨機抽查 5-10 個檔案的 enforcement/ 連結，確認可正確點擊

### Phase 6：移除 shared-rules/ 並提交

**目標**：移除舊目錄，完成搬遷。

- [ ] 確認 Phase 1-5 全部完成且驗證通過
- [ ] 執行 `git rm -r shared-rules/`
- [ ] 執行 `git add -A`
- [ ] 執行 `git commit -m "refactor: rename shared-rules/ to enforcement/

- Phase 1: Create enforcement/ as mirror of shared-rules/
- Phase 2: Update enforcement/ internal path references
- Phase 3: Update all external shared-rules/ references to enforcement/
- Phase 4: Update validator paths
- Phase 5: Verify all links correct
- Phase 6: Remove shared-rules/ directory

Layer Responsibility Contract:
- governance/ (Policy Architecture / WHY)
- enforcement/ (Runtime Enforcement Layer / RULES)
- runtime/ (Runtime Engine Layer / LOADING)"`
- [ ] 執行 `git push`
- [ ] 執行 `git status --short --branch` 確認乾淨

---

## 4. Layer Responsibility Contract（完整版）

### 4.1 三層定義

| 層級 | 目錄 | 角色 | 關鍵問題 |
|------|------|------|---------|
| **Policy Architecture** | `governance/` | 定義知識治理的架構設計 | 「知識如何被治理？」 |
| **Runtime Enforcement** | `enforcement/` | 定義 AI agent 必須遵守的執行政策 | 「AI agent 必須遵守什麼規則？」 |
| **Runtime Engine** | `runtime/` | 定義 AI 系統如何動態載入與路由 | 「AI 系統如何運作？」 |

### 4.2 各層包含內容

**governance/** 包含：
- Knowledge lifecycle（生命週期狀態、提升關卡、刪除規則）
- Validation gates（migration checklist、pass/block rules）
- Cleanup strategy（重複偵測、拆分規則、所有權邊界）
- Dependency maintenance（依賴圖維護、連動更新策略）
- Intelligence extraction pipeline（內容審計、類型判斷、格式轉換）

**enforcement/** 包含：
- Core Bootstrap 規則（rule-weight、dependency-reading、conversation-goal-ledger）
- Lazy-load 規則（linked-updates、failure-learning、sanitization、feedback-lessons 等 14 條）
- Failure patterns（跨 skill 可重用的 agent failure 模式）
- Runtime Activation Model（Core Bootstrap + Lazy-load 的載入策略）

**runtime/** 包含：
- Context routing（task intent → knowledge index → metadata → source-of-truth gate）
- Onboarding quickstart（新專案或新任務的初始設定指引）
- Activation rules（lazy-load 規則的觸發條件）
- Context pruning 與 orchestration 設計

### 4.3 層間互動規則

1. **governance/ 定義 enforcement/ 的治理架構**：governance 定義 enforcement/ 中的規則如何被建立、更新、廢棄，但不定義規則的具體內容。
2. **enforcement/ 定義 runtime/ 的載入政策**：enforcement 定義哪些規則必須被載入（Core Bootstrap）以及何時載入（Lazy-load activation），但不定義如何載入。
3. **runtime/ 實作 enforcement/ 的載入要求**：runtime 根據 enforcement/ 的 activation model 實作 dynamic loading，但不定義載入什麼。
4. **跨層引用**：上層可以引用下層（governance → enforcement → runtime），下層不應反向引用上層的具體內容（但可引用上層的架構設計）。

---

## 5. 風險與緩解

| 風險 | 影響 | 緩解措施 |
|------|------|---------|
| 遺漏某個檔案的路徑更新 | 連結失效，agent 讀不到正確規則 | Phase 5 的 grep 驗證 + validator 執行 |
| 搬遷期間有其他人 commit 新 shared-rules 引用 | 搬遷後新引用指向不存在路徑 | 搬遷前確認 git 乾淨，搬遷後立即推送 |
| sync-cursor-bundle.sh 的 symlink 路徑錯誤 | Cursor bundle 失效 | Phase 4 更新 script 後手動測試一次 |
| validator 中的 shared-rules 路徑未更新完全 | validator 掃描錯誤目錄 | Phase 4 專門更新 validator，Phase 5 執行 validator 驗證 |
| 向後相容不足，既有工具仍指向 shared-rules/ | 工具找不到規則 | Phase 1 保留 shared-rules/ 直到 Phase 6，提供緩衝期 |

---

## 6. 驗證標準

### 6.1 通過條件

1. `grep -rn "shared-rules/" --include="*.md" --include="*.rb" --include="*.sh" --include="*.yaml" .` 只出現 shared-rules/ 目錄本身的檔案
2. `ruby scripts/refresh-knowledge-runtime.rb` 所有 validator 通過
3. `git status --short --branch` 乾淨
4. 隨機抽查 10 個 enforcement/ 連結，全部正確

### 6.2 不通過時的補救

若驗證發現遺漏：
1. 記錄遺漏的檔案與路徑
2. 回到對應 Phase 修正
3. 重新執行驗證
4. 若為系統性遺漏（如某類檔案全部漏掉），更新本計畫的檢查清單

---

## 7. 與既有文件的關係

- [`shared-rules/failure-patterns/shared-rules-architecture-drift.md`](../../shared-rules/failure-patterns/shared-rules-architecture-drift.md) — 本計畫直接對應此 failure pattern 的 prevention gate
- [`governance/lifecycle/intelligence-extraction-pipeline.md`](../../governance/lifecycle/intelligence-extraction-pipeline.md) Step 7a — 架構重構後的 shared-rules 同步檢查
- [`shared-rules/linked-updates.md`](../../shared-rules/linked-updates.md) — 第 51 行「架構重構」連動關係
- [`plans/README.md`](../../plans/README.md) — 本計畫完成後需更新 plans/README.md 狀態並搬移至 archived/
