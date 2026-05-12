# Intelligence Extraction Pipeline（智慧抽取管線）

## 為什麼需要這個 Pipeline

從 Phase 28-29 的 4 個 techniques decomposition 經驗中，我們發現：

1. **Technique 文件混合了兩種本質不同的知識**：HOW TO DO（操作流程）和 HOW TO THINK（決策智慧）
2. **直接搬運 technique 到新層沒有意義** — 需要的是拆解（decompose），不是搬運（move）
3. **每個 technique 的拆解模式不同** — 有些產生多個 intelligence atoms（flutter-dart-aot → 4 atoms），有些只產生 1-2 個（media-hls → 1 atom）
4. **Feedback history 的 extraction 模式完全不同** — 已經是「已提取產品」，只需索引 + 標註

本 pipeline 定義從「混合內容」到「分層知識」的可重複流程。

## Pipeline 概覽

```
來源內容（technique / feedback / SKILL.md）
    │
    ▼
Step 1: 內容審計（Content Audit）
    │
    ▼
Step 2: 類型判斷（Type Classification）
    ├── HOW TO DO → analysis/workflows/
    ├── HOW TO THINK → intelligence/{heuristics,anti-patterns,failure,signals}/
    ├── Execution Flow → workflow/execution-flow.md
    ├── Artifact Gate → workflow/artifact-gates.md
    └── Feedback Lesson → feedback/extraction/（索引 + 標註）
    │
    ▼
Step 3: 拆解執行（Decomposition）
    │
    ▼
Step 4: 格式轉換（Format Transformation）
    │
    ▼
Step 5: 標註來源（Source Annotation）
    │
    ▼
Step 6: 驗證（Validation）
    │
    ▼
Step 7: 更新索引（Index Update）
```

---

## Step 1: 內容審計（Content Audit）

**目標**：完整讀取來源內容，識別所有可拆解的元素。

**審計維度**：

| 維度 | 檢查項目 | 範例 |
|------|---------|------|
| 操作步驟 | 有無明確的 step-by-step 流程？ | 「步驟 1-7」在 flutter-dart-aot Common Flow |
| 判斷決策 | 有無「何時該做 X vs Y」的決策？ | 「何時用 Java hook vs Dart hook」 |
| 工具命令 | 有無具體的 CLI 命令？ | `frida -U -f com.example` |
| 失敗模式 | 有無常見錯誤與解決方式？ | 「Frida spawn 後 crash」 |
| 信號偵測 | 有無判斷「這是什麼」的線索？ | 「libapp.so 存在 → Flutter/Dart AOT」 |
| 反模式 | 有無不該做的做法？ | 「不要 broad hook 多個 Dart runtime function」 |
| 通用建議 | 有無抽象原則？ | 「從高語意邊界開始」 |
| 產出規範 | 有無 artifact 格式要求？ | 「API Catalog 必須包含 field meaning」 |

**產出**：內容審計清單，標記每個元素的類型。

---

## Step 2: 類型判斷（Type Classification）

根據審計結果，將每個元素分類到對應的目標層：

| 元素類型 | 目標層 | 格式模板 | 範例 |
|---------|--------|---------|------|
| **操作流程**（HOW TO DO） | `analysis/<domain>/workflows/<flow-name>.md` | 步驟式流程 + 命令模板 | `frida-hook-flow.md` |
| **執行流程**（Execution Flow） | `workflow/<domain>/execution-flow.md` | 階段式流程 + 完成定義 | `execution-flow.md` |
| **產出規範**（Artifact Gate） | `workflow/<domain>/artifact-gates.md` | 最低要求 + 完成門檻 | `artifact-gates.md` |
| **經驗法則**（Heuristic） | `intelligence/<domain>/heuristics/<name>.md` | 問題 → 原則 → 決策表 → Token 影響 | `hook-selection.md` |
| **反模式**（Anti-pattern） | `intelligence/<domain>/anti-patterns/<name>.md` | 症狀 → 可能原因 → 診斷方式 | `early-hook-instability.md` |
| **失敗模式**（Failure） | `intelligence/<domain>/failure/<name>.md` | 問題 → 診斷 → 緩解 → 預防 | `frida-spawn-race.md` |
| **信號偵測**（Signal） | `intelligence/<domain>/signals/<name>.md` | 信號 → 檢查方式 → 可信度 | `flutter-dart-aot-detection.md` |
| **工具參考**（Tool Reference） | `analysis/<domain>/tools-and-failures.md` | 工具表 + 失敗判讀 | `tools-and-failures.md` |
| **路線選擇**（Routing） | `intelligence/<domain>/evidence-first-routing.md` | 證據 → 路線決策表 | `evidence-first-routing.md` |
| **Feedback Lesson** | `feedback/extraction/<domain>-index.md`（索引） | 索引 + `# Extracted` 標註 | `apk-analysis-index.md` |

### 判斷流程

```
內容元素
    │
    ├── 包含「步驟 1, 2, 3...」或「先做 X，再做 Y」？
    │   └── YES → 操作流程（HOW TO DO）
    │
    ├── 包含「何時該做 X vs Y」或「如果...則...」？
    │   └── YES → 經驗法則（Heuristic）
    │
    ├── 包含「不要做 X」或「常見錯誤」？
    │   └── YES → 反模式（Anti-pattern）
    │
    ├── 包含「X 失敗了怎麼辦」或「錯誤訊息 Y」？
    │   └── YES → 失敗模式（Failure）
    │
    ├── 包含「如何判斷是 X」或「X 的徵兆是...」？
    │   └── YES → 信號偵測（Signal）
    │
    ├── 包含「產出必須包含 X」或「完成條件是...」？
    │   └── YES → 產出規範（Artifact Gate）
    │
    ├── 包含「先確認 X，再決定路線」？
    │   └── YES → 路線選擇（Routing）
    │
    └── 以上皆非 → 保留在來源，標註「未分類」
```

---

## Step 3: 拆解執行（Decomposition）

### 3.1 Technique Decomposition（技術拆解）

這是 Phase 28-29 的核心模式。將一個 technique 文件拆成 workflow + intelligence atoms。

**執行步驟**：

1. **建立 workflow 文件**：提取所有操作步驟、命令、順序到 `analysis/<domain>/workflows/<flow-name>.md`
2. **提煉 intelligence atoms**：從判斷邏輯、策略建議、錯誤模式中提煉 atoms
3. **保留通用建議**：如果某個建議太抽象（如「從高語意邊界開始」），保留在來源，不強制 atomize
4. **標註來源**：在舊 technique 檔案加入 `# Intelligence Extracted` 標記

**實際案例對照**：

| Technique | Workflow 產出 | Intelligence Atoms | 備註 |
|-----------|--------------|-------------------|------|
| flutter-dart-aot | `frida-hook-flow.md` | hook-selection（heuristic）、early-hook-instability（anti-pattern）、frida-spawn-race（failure）、flutter-dart-aot-detection（signal） | 最豐富，4 atoms |
| http-api | `http-api-documentation-flow.md` | api-documentation-completeness（heuristic） | 決策智慧較少，1 atom |
| local-proxy | `local-proxy-hook-flow.md` | local-proxy-routing-diagnosis（heuristic）、local-proxy-detection（signal） | 涉及 routing 判斷，2 atoms |
| media-hls | `media-hls-analysis-flow.md` | media-type-detection（signal） | 主要是操作流程，1 atom |

### 3.2 Feedback History Extraction（反饋歷史提取）

這是 Phase 30 的核心模式。Feedback history lessons 已經是「已提取產品」，不需要重新提取內容。

**執行步驟**：

1. **分析 Promotion Target**：讀取每個 lesson 的 `Promotion Target` 欄位
2. **建立索引**：將 lessons 分類到對應的目標層，建立 `feedback/extraction/<domain>-index.md`
3. **標註提取狀態**：在每個 lesson 檔案加入 `# Extracted — See <target path>` 標記
4. **不重新提取內容**：lesson 的內容已經在目標層存在或 lesson 本身就是最終形式

### 3.3 SKILL.md Decomposition（Phase 32 ✅ 已驗證）

將 SKILL.md 中的 Quick Start、Default Workflow、Output Style、Feedback Loop 提取到對應新層。

**執行步驟**：

1. Quick Start → `runtime/onboarding/<skill>-quickstart.md`
2. Default Workflow → 補齊 `workflow/<domain>/execution-flow.md`
3. Output Style → `workflow/<domain>/artifact-gates.md`
4. Feedback Loop → `feedback/` 層
5. SKILL.md 瘦身為純 routing 文件

**實際案例對照**：

| SKILL.md | Quick Start 產出 | Output Style 產出 | 瘦身比例 |
|----------|-----------------|-------------------|---------|
| apk-analysis | `runtime/onboarding/apk-analysis-quickstart.md`（步驟 5+7） | `workflow/apk-analysis/artifact-gates.md`（既有） | ~65%（158→55 行） |
| app-development-guidance | `runtime/onboarding/app-development-guidance-quickstart.md`（15 步驟） | `workflow/app-development-guidance/artifact-gates.md`（既有） | ~51%（132→65 行） |
| travel-planning | `runtime/onboarding/travel-planning-quickstart.md`（18 步驟） | `workflow/travel-planning/artifact-gates.md`（**新建**） | ~46%（102→55 行） |

**關鍵發現**：
- `travel-planning` 缺少 `workflow/travel-planning/artifact-gates.md`，需先建立才能提取 Output Style
- `apk-analysis` 的 Quick Start 步驟 1-4 是 routing（保留），步驟 5+7 是操作內容（提取），步驟 6 是 heuristic（保留），步驟 8 是 feedback（提取）
- 瘦身後 SKILL.md 保留：header metadata、Shared Policy、When To Use、Out Of Scope、Default Workflow（純 routing）

---

## Step 4: 格式轉換（Format Transformation）

每個 intelligence atom 類型有對應的格式模板：

### Heuristic 格式

```markdown
# <Title> Heuristic（<中文標題>）

## 問題
<一句話描述這個 heuristic 解決什麼決策問題>

## 原則
- <原則 1>
- <原則 2>

## 決策表
| 情境 | 建議做法 | 判斷信號 |
|------|---------|---------|
| <情境 A> | <做法 A> | <信號 A> |
| <情境 B> | <做法 B> | <信號 B> |

## 不建議的做法
- <做法 1>
- <做法 2>

## Token 影響
<低/中/高。此 atom 在 <情境> 中 lazy-load，約 <N> tokens。>
```

### Anti-pattern 格式

```markdown
# <Title> Anti-pattern（<中文標題>）

## 問題
<這個反模式是什麼>

## 症狀
| 症狀 | 可能原因 | 診斷方式 |
|------|---------|---------|
| <症狀 A> | <原因 A> | <診斷 A> |

## 為什麼會發生
<根本原因分析>

## 預防方式
<如何避免>

## Token 影響
<...>
```

### Failure 格式

```markdown
# <Title> Failure（<中文標題>）

## 問題
<這個失敗模式是什麼>

## 診斷
| 觀察 | 可能原因 | 確認方式 |
|------|---------|---------|
| <觀察 A> | <原因 A> | <確認 A> |

## 緩解方式
<如何修復>

## 預防方式
<如何避免再次發生>

## Token 影響
<...>
```

### Signal 格式

```markdown
# <Title> Detection Signals（<中文標題>）

## 問題
<這個 signal 幫助判斷什麼>

## 判斷信號
### 主要信號（高可信度）
| 信號 | 檢查方式 | 可信度 |
|------|---------|-------|
| <信號 A> | <檢查 A> | 高 |

### 次要信號（中等可信度）
| 信號 | 檢查方式 | 可信度 |
|------|---------|-------|
| <信號 B> | <檢查 B> | 中 |

### 排除信號
| 信號 | 意義 |
|------|------|
| <信號 C> | <排除理由> |

## 判斷流程
<如何組合這些信號做出判斷>

## Token 影響
<...>
```

### Workflow 格式

```markdown
# <Title> Flow（<中文標題>）

## 前置條件
- <條件 1>
- <條件 2>

## 步驟
### 步驟 1：<名稱>
<說明>

### 步驟 2：<名稱>
<說明>

## 成功產出格式
<產出範本>

## 注意事項
- <事項 1>
```

---

## Step 5: 標註來源（Source Annotation）

在舊來源檔案中加入 extraction 標記，確保 traceability。

### Technique 標註格式

在舊 technique 檔案的標題下方加入：

```markdown
# Intelligence Extracted

本檔案的部分內容已提取到新架構層：

- **Workflow**：參見 [`analysis/<domain>/workflows/<flow-name>.md`](<relative-path>)
- **Heuristics**：參見 [`intelligence/<domain>/heuristics/<name>.md`](<relative-path>)
- **Signals**：參見 [`intelligence/<domain>/signals/<name>.md`](<relative-path>)
- **Anti-patterns**：參見 [`intelligence/<domain>/anti-patterns/<name>.md`](<relative-path>)
- **Failures**：參見 [`intelligence/<domain>/failure/<name>.md`](<relative-path>)

本檔案保留作為相容性入口，新內容請直接寫入上述目標檔案。
```

### Feedback Lesson 標註格式

在 feedback lesson 檔案的第一行（compliance header 之後）加入：

```markdown
# Extracted — See [`<target-path>`](<relative-path>)
```

---

## Step 6: 驗證（Validation）

### 完整性檢查

| 檢查項目 | 通過條件 |
|---------|---------|
| 所有操作步驟已提取 | Workflow 文件涵蓋來源中所有 step-by-step 內容 |
| 所有決策邏輯已 atomize | 每個「如果...則...」都有對應的 heuristic 或 signal |
| 所有失敗模式已記錄 | 每個已知錯誤都有對應的 failure 或 anti-pattern atom |
| 來源已標註 | 舊檔案有 `# Intelligence Extracted` 或 `# Extracted` 標記 |
| 目標層 README 已更新 | `analysis/<domain>/workflows/README.md` 和 `intelligence/<domain>/README.md` 已加入新檔案 |

### 品質檢查

| 檢查項目 | 通過條件 |
|---------|---------|
| Atom 有邊界條件 | 每個 atom 說明「何時適用」與「何時不適用」 |
| Atom 有 Token Impact | 每個 atom 標註 token 估算 |
| Workflow 可執行 | 照著 workflow 步驟可完成操作 |
| Signal 可驗證 | 每個 signal 有具體的檢查方式 |
| 無專案特定內容 | Intelligence atoms 不含專案名稱、客戶名稱、raw evidence |

### 不強制 atomize 的情況

以下情況不強制建立 intelligence atom：

- **太抽象的通用建議**（如「從高語意邊界開始」）— 保留在來源
- **已在 tools-and-failures.md 中的工具選擇** — 不重複
- **只適用單一專案的細節** — 不適合 intelligence
- **無法泛化的觀察** — 留在 feedback history

---

## Step 7: 更新索引（Index Update）

Extraction 完成後需要更新以下文件：

| 文件 | 更新內容 |
|------|---------|
| `analysis/<domain>/workflows/README.md` | 加入新 workflow 的表格列 |
| `intelligence/<domain>/README.md` | 加入新 atoms 的表格列 |
| `knowledge/indexes/README.md` | 加入新路由列（如需要） |
| `knowledge/runtime/routing-registry.yaml` | 加入新 routing record（如需要） |
| `knowledge/graphs/` | 加入新 graph edge（如需要） |
| `architecture/next-stage-upgrade-plan.md` | 更新 Phase 完成狀態 |

### Step 7a：Shared-Rules 同步檢查（架構變更專用）

當 extraction 涉及**架構重構**（目錄重組、分層新增、路徑變更、命名變更）時，除了更新上述索引，**必須**同步檢查 `shared-rules/` 中的路徑參考是否過期。

#### 檢查範圍

| 類型 | 檢查重點 | 範例檔案 |
|------|---------|---------|
| 範例路徑 | 目錄結構範例是否仍指向舊路徑 | `shared-rules/document-sizing.md` |
| 表格內容 | 連動關係表是否仍使用舊路徑 | `shared-rules/linked-updates.md` |
| 模板 | Promotion Target 是否缺少新分層選項 | `shared-rules/feedback-lessons.md` |
| 索引 | lazy-load 表格是否指向舊檔案 | `shared-rules/README.md` |
| 教學文件 | 目錄結構建議是否仍以舊結構為主 | `skills/ADDING_SKILLS.md` |
| 規則正文 | 路徑描述是否過期 | `shared-rules/content-layering.md`、`shared-rules/tool-neutral-documentation.md` |
| 流程規則 | Context Loading 步驟是否指向舊入口 | `shared-rules/decision-efficiency.md` |
| 引用規則 | Cross-skill reference 格式是否過期 | `shared-rules/cross-skill-references.md` |

#### 更新原則

1. **保留舊路徑**作為向後相容參考，標註「舊結構，向後相容」
2. **新增新分層路徑**作為主要參考，標註「新分層（優先）」
3. 通用規則（如 `document-sizing.md`）應明確標註「可跨專案通用」
4. 使用通用格式（`<domain>`）而非硬編碼 skill 名稱

#### 驗證命令

```bash
# 確認無遺漏的舊路徑參考（預期：無，或只有 intentional 向後相容參考）
grep -rn "skills/<old-path>/" shared-rules/ --include="*.md"

# 確認新分層路徑已出現
grep -rn "workflow/<domain>/" shared-rules/ --include="*.md"
grep -rn "analysis/<domain>/" shared-rules/ --include="*.md"
grep -rn "intelligence/<domain>/" shared-rules/ --include="*.md"
```

#### 觸發條件

符合任一條件時，**必須**執行 Step 7a：

- 新增或重組 `workflow/<domain>/`、`analysis/<domain>/`、`intelligence/<domain>/` 等新分層
- 將舊 `skills/<skill-name>/` 內容提取到新分層
- 修改 `skills/` 結構後，未檢查 `shared-rules/` 中的範例路徑、表格、模板
- 只更新了「被提取的檔案本身」，但沒有更新「描述如何提取的規則」

---

## Pipeline 適用範圍

### 已驗證的模式

| 模式 | 適用來源 | Phase | 狀態 |
|------|---------|-------|------|
| Technique Decomposition | `skills/*/techniques/*/` | 28-29 | ✅ 已驗證（4 techniques） |
| Feedback History Extraction | `skills/*/feedback_history/*/` | 30 | ✅ 已驗證（101 lessons） |
| SKILL.md Decomposition | `skills/*/SKILL.md` | 32 | ✅ 已驗證（3 skills） |
| Skill-Specific Extraction | `skills/*/{CHECKLIST.md,TOOLS.md,README.md}` | 33 | ✅ 已驗證（3 files） |

### 未來可能擴充的模式

| 模式 | 適用來源 | 預計 Phase | 備註 |
|------|---------|-----------|------|
| DOCUMENTATION.md Extraction | `skills/*/DOCUMENTATION.md` | 未來 | 產出規範 → artifact-gates.md（travel-planning 已完成） |
| WORKFLOW.md Extraction | `skills/*/WORKFLOW.md` | 未來 | 執行流程 → workflow/execution-flow.md（apk-analysis + app-development-guidance + travel-planning 已完成） |
| TOOLS.md Extraction | `skills/*/TOOLS.md` | 未來 | 工具參考 → analysis/tools-and-failures.md（apk-analysis TOOLS.md 尚未提取） |
| FEEDBACK.md Extraction | `skills/*/FEEDBACK.md` | 未來 | 反饋規則 → feedback/ |

---

## 與既有文件的關係

- [`feedback/extraction/README.md`](../../feedback/extraction/README.md) — 定義 extraction 的核心責任與門檻，本 pipeline 是其具體執行流程
- [`notes/intelligence-extraction-observations.md`](../../notes/intelligence-extraction-observations.md) — 記錄 extraction 過程中的觀察，本 pipeline 從中抽象
- [`governance/lifecycle/README.md`](../../governance/lifecycle/README.md) — 定義知識生命週期與 deprecation timeline，本 pipeline 的產出會進入 lifecycle
- [`metadata/schema.md`](../../metadata/schema.md) — 定義 Knowledge Atom metadata schema，intelligence atoms 應符合其格式
- [`knowledge/indexes/README.md`](../../knowledge/indexes/README.md) — 任務路由索引，extraction 完成後需更新
- [`knowledge/runtime/routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml) — Machine-readable routing registry，extraction 完成後需更新
- [`architecture/next-stage-upgrade-plan.md`](../../architecture/next-stage-upgrade-plan.md) — 整體升級規劃，本 pipeline 是 Phase 31 的核心產出
