# Workflow 選路（Task → Workflow Discovery）

本檔是 **Workflow 層** 的選路表，掛在 Ai-skill **既有 Routing Discovery** 之上，不另發明一套機制。

## 治理依據（必讀）

| 文件 | 角色 |
| --- | --- |
| [`governance/lifecycle/routing-philosophy.md`](../governance/lifecycle/routing-philosophy.md) | **Normative**：task intent → `knowledge/indexes` → `routing-registry` → primary source |
| [`enforcement/dependency-reading.md`](../enforcement/dependency-reading.md) | 查到 route 後 **必須** 讀 `primary_source` |
| [`knowledge/indexes/README.md`](../knowledge/indexes/README.md) | 人類可讀 task intent 表（含 development guidance 列） |
| [`runtime/router/activation-table.md`](../runtime/router/activation-table.md) **#27** | Workflow 編排通用閘門；具體觸發見 registry `activation_triggers` |

**強制規則**：命中任一 `route.workflow.*` 或 activation **#27** 時，**不得**只載入單一 intelligence 規則就寫碼；**必須**先完成 [`routing-philosophy.md`](../governance/lifecycle/routing-philosophy.md) Step 1–5，比對 registry `activation_triggers`，再依本檔 §選路表／§歧義鎖定單一 `route.workflow.*`。

## Discovery 流程（與 governance 對齊）

```text
0. Activation #27 或 registry `route.workflow.*.activation_triggers` 命中
1. routing-philosophy Step 1：分類 task intent
2. routing-philosophy Step 2：knowledge/indexes → routing-registry.yaml
3. 本檔 §選路表：在 route.workflow.* 中選 workflow（非單一 repo）
4. 載入選中 workflow 的 README.md → execution-flow.md
5. 若該 route 含 project_overlays：進入 workflow 後載入 <PROJECT_ROOT> 專案 overlay
```

**Catalog 索引（machine-readable）**：[`knowledge/runtime/routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) 內所有 `id: route.workflow.*`。

**人類索引**：[`workflow/README.md`](./README.md) §目前入口。

## 選路表（Task → Workflow）

| 任務性質 | 典型信號 | 接上的 Workflow | routing `id` |
| --- | --- | --- | --- |
| **App／SDK 開發交付** | 實作、改程式、寫技術 plan、契約、BDD、code review、Maven 模組 | [`software-delivery/`](./software-delivery/README.md) | `route.workflow.software-delivery` |
| **APK 逆向／抓包／Frida** | 分析 APK、hook、協議、解密、live capture、未授權行為取證 | [`apk-analysis/`](./apk-analysis/README.md) | `route.workflow.apk-analysis` |
| **從零新專案** | greenfield、specify/plan/tasks、無既有 codebase | [`greenfield/`](./greenfield/README.md) | `route.workflow.greenfield` |
| **Repo 結構分析** | onboarding、migration impact、tech debt、深讀 codebase（**不**以寫產品碼為主） | [`repo-analysis/`](./repo-analysis/README.md) | `route.workflow.repo-analysis`（若 routing 已登記） |
| **跨專案 agent 友善文件** | 撰寫/拆分 `docs/`、index-first、**零**可觀察產品行為變更 | [`documentation/`](./documentation/README.md) | `route.workflow.documentation-ai-native` |
| **旅遊規劃** | itinerary、交通、預算 | [`travel-planning/`](./travel-planning/README.md) | `route.workflow.travel-planning` |

### 常見歧義

| 情況 | 選哪個 |
| --- | --- |
| 在 `unwrapping` 寫 `apk-analysis-sdk` plan + 實作 | **software-delivery**（開發交付）；不是 apk-analysis |
| 對 TATA APK 做 Frida attach 抓 API | **apk-analysis** |
| 只改 `docs/plans/*.md` 且會導致之後要寫 SDK | **software-delivery**；純文件架構且明確零行為變更可考慮 **documentation** |
| 新 repo 從 spec 開始 | **greenfield** → 實作階段 often 再接 **software-delivery** |

## 與 activation-table / registry 的關係

| 機制 | 角色 |
| --- | --- |
| **activation #27** | 通用 Workflow 閘門 + §Discovery SOP |
| **`routing-registry` `activation_triggers`** | 各 `route.workflow.*` 的觸發條件與 `required_dependencies`（**registry-first**） |
| **本檔 §歧義** | 多 route 同時命中時的裁決（不新增 activation 列） |
| **activation #8** 撰寫文件 | 若僅 agent 友善文件 → **documentation**；若含契約/BDD/可觀察行為 → 走 **#27** + **software-delivery** |

## 專案 overlay（第二層，非 Workflow 進入點）

部分 repo 在 **已選定 workflow 之後** 提供薄 yaml（例如 `apk-analysis-sdk/runtime/workflow-activation.yaml`）。  
那是 **software-delivery 下的專案 gate**，不是第 6 種 workflow，也不產生 `route.project.*`。

## 驗證

- Agent 能否列舉目前所有 `route.workflow.*`？
- 對當前任務說出選中的 workflow 與理由？
- 是否已讀該 workflow 的 `execution-flow.md` 再動手？
