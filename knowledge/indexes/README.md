# Knowledge Navigation Indexes

`knowledge/indexes/` 負責讓 agent 從任務意圖快速找到 task-relevant knowledge。這裡先保存索引格式與第一版 routing table；不搬移既有 `skills/` 或 `shared-rules/` 內容。

## 使用方式

1. 先用任務意圖對照 [任務路由索引](#任務路由索引)。
2. 讀取 `Primary source`，再依需要讀 `Related sources`。
3. 若需要 machine-readable route，讀取 [`../runtime/routing-registry.yaml`](../runtime/routing-registry.yaml)。
4. 若任務需要建立或評估 Knowledge Atom，使用 [`../../metadata/schema.md`](../../metadata/schema.md) 的欄位與 controlled values。
5. 若路由產生新候選知識，先記錄在對應 source 層或 `.agent-goals/`，不要直接搬移大量內容。

## 索引欄位

| 欄位 | 用途 |
| --- | --- |
| `Task intent` | 使用者或 agent 的工作意圖。 |
| `Primary source` | 第一個應讀的 canonical source。 |
| `Related sources` | 視任務深度補讀的來源。 |
| `Target layer` | 未來遷移或抽象化時的主要歸屬層。 |
| `When to read` | 何時載入此知識，避免無關 context。 |
| `Validation signal` | 如何確認這條路由仍然可靠。 |

## 任務路由索引

| Task intent | Primary source | Related sources | Target layer | When to read | Validation signal |
| --- | --- | --- | --- | --- | --- |
| 啟動或接手 Ai-skill work | [`../../README.md`](../../README.md) | [`../../shared-rules/README.md`](../../shared-rules/README.md), [`../../shared-rules/dependency-reading.md`](../../shared-rules/dependency-reading.md) | `runtime/`, `governance/` | 新 session、接手長任務、或使用者要求繼續 Ai-skill 升級時 | Bootstrap required set 已讀，`git status --short --branch` 已檢查 |
| 判斷規則優先序或 source/mirror 邊界 | [`../../shared-rules/rule-weight.md`](../../shared-rules/rule-weight.md) | [`../../shared-rules/dependency-reading.md`](../../shared-rules/dependency-reading.md), [`../../architecture/ai-native-knowledge-operating-system.md`](../../architecture/ai-native-knowledge-operating-system.md) | `governance/`, `runtime/` | instructions、tool adapter、local mirror 或 validation gate 看似衝突時 | 選擇保留 canonical source、validation 與最新 user goal 的路線 |
| 維護舊 skills 與新分層的 source-of-truth 邊界 | [`../../governance/lifecycle/README.md`](../../governance/lifecycle/README.md) | [`../../governance/validation/README.md`](../../governance/validation/README.md), [`../../metadata/compatibility/README.md`](../../metadata/compatibility/README.md) | `governance/`, `metadata/` | 舊 skill 仍在更新、建立 candidate map、promotion 或 deprecation 時 | Old entrypoint remains reachable; lifecycle state and validation gate are explicit |
| 決定 runtime context loading 路線 | [`../runtime/routing-registry.yaml`](../runtime/routing-registry.yaml) | [`../../runtime/routing/README.md`](../../runtime/routing/README.md), [`../../metadata/ranking/README.md`](../../metadata/ranking/README.md), [`../../metadata/confidence/README.md`](../../metadata/confidence/README.md), [`../../knowledge/summaries/README.md`](../../knowledge/summaries/README.md) | `runtime/`, `knowledge/`, `metadata/` | 需要從 task intent 選擇 primary source、related sources 或 candidate summaries 時 | Primary source, deferred sources, source-of-truth gate, and validation signal are recorded |
| 規劃下一階段 OS 分層 | [`../../architecture/next-stage-upgrade-plan.md`](../../architecture/next-stage-upgrade-plan.md) | [`../README.md`](../README.md), [`../../metadata/README.md`](../../metadata/README.md) | `knowledge/`, `metadata/`, `governance/` | 建立或調整 `analysis/`、`workflow/`、`metadata/` 等 top-level layers 時 | Roadmap status 與 root dashboard 同步 |
| 執行 APK analysis 能力 | [`../../skills/apk-analysis/SKILL.md`](../../skills/apk-analysis/SKILL.md) | [`../../architecture/apk-analysis-pilot-migration.md`](../../architecture/apk-analysis-pilot-migration.md), [`../../analysis/apk/README.md`](../../analysis/apk/README.md), [`../../workflow/apk-analysis/README.md`](../../workflow/apk-analysis/README.md), [`../../intelligence/engineering/apk-analysis/README.md`](../../intelligence/engineering/apk-analysis/README.md) | `analysis/`, `workflow/`, `intelligence/` | 使用者要求授權 APK traffic/runtime/response 分析時 | Skill dependency set 已讀，authorization 與 sanitization gate 已套用；old skill entrypoint remains active |
| 選擇 APK 分析最高收益路線 | [`../../intelligence/engineering/apk-analysis/highest-leverage-analysis-path.md`](../../intelligence/engineering/apk-analysis/highest-leverage-analysis-path.md) | [`../../skills/apk-analysis/SKILL.md`](../../skills/apk-analysis/SKILL.md), [`../../skills/apk-analysis/WORKFLOW.md`](../../skills/apk-analysis/WORKFLOW.md), [`../../knowledge/summaries/apk-highest-leverage-analysis.md`](../../knowledge/summaries/apk-highest-leverage-analysis.md), [`../../knowledge/graphs/apk-highest-leverage-analysis.yaml`](../../knowledge/graphs/apk-highest-leverage-analysis.yaml) | `intelligence/`, `runtime/` | APK 分析卡住、可選多條 route，或需要判斷 UI、API replay、hook、pcap、MITM、static xref 哪條先做時 | Current unknown、candidate routes、chosen validation signal 與 fallback route 已記錄 |
| 產出 app/API/embedded development guidance | [`../../skills/app-development-guidance/SKILL.md`](../../skills/app-development-guidance/SKILL.md) | [`../../skills/app-development-guidance/README.md`](../../skills/app-development-guidance/README.md), [`../../intelligence/README.md`](../../intelligence/README.md), [`../../workflow/README.md`](../../workflow/README.md) | `workflow/`, `intelligence/` | 使用者要把觀察、規格或反向工程成果轉成 buildable guidance 時 | BDD、contract、validation 與 checklist closure 可反查 |
| 規劃 evidence-based travel itinerary | [`../../skills/travel-planning/SKILL.md`](../../skills/travel-planning/SKILL.md) | [`../../skills/travel-planning/README.md`](../../skills/travel-planning/README.md), [`../../workflow/README.md`](../../workflow/README.md) | `workflow/`, `intelligence/` | 使用者要求旅遊路線、交通、餐飲、住宿或行程 app 欄位時 | Current source links、營業時間、交通與可行性檢查可反查 |
| 沉澱或修正 feedback lesson | [`../../shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md) | [`../../shared-rules/failure-learning-system.md`](../../shared-rules/failure-learning-system.md), [`../../feedback/README.md`](../../feedback/README.md), [`../../memory/README.md`](../../memory/README.md) | `feedback/`, `memory/`, `intelligence/` | 使用者指出可重用 lesson、agent failure 或需要 promotion 時 | Promotion target、linked updates、sanitization 與 validation 已檢查 |
| 規劃 feedback lesson promotion 或 downgrade | [`../../feedback/promotion/README.md`](../../feedback/promotion/README.md) | [`../../shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md), [`../../shared-rules/failure-learning-system.md`](../../shared-rules/failure-learning-system.md), [`../../governance/lifecycle/README.md`](../../governance/lifecycle/README.md), [`../../knowledge/summaries/feedback-promotion-pipeline.md`](../../knowledge/summaries/feedback-promotion-pipeline.md) | `feedback/`, `governance/`, `knowledge/` | Lesson 要推進到 workflow、intelligence、shared-rules、memory 或 runtime surface，或需要退回 / 降級時 | 原 lesson source 保留，最小 durable target、runtime refresh 與 close-loop validation 已記錄 |
| 設計 Knowledge Atom metadata | [`../../metadata/schema.md`](../../metadata/schema.md) | [`../../metadata/README.md`](../../metadata/README.md), [`../../architecture/next-stage-upgrade-plan.md`](../../architecture/next-stage-upgrade-plan.md), [`../README.md`](../README.md) | `metadata/`, `knowledge/` | 建立 schema、index、summary、graph 或 runtime metadata 時 | Schema 欄位能套用到第一批 atom candidates |

## 維護規則

- 新增索引列時，優先連到 canonical source，不連到 tool mirror。
- `Primary source` 必須是 agent 能直接讀的檔案，不使用只有人能理解的模糊描述。
- `Target layer` 可以指向未來遷移層，但不得暗示內容已完成搬移。
- 若新增或改動路由會影響 root dashboard、roadmap、metadata schema 或 skill entry，依 `shared-rules/linked-updates.md` 同步更新。
- 本索引只是 navigation layer；可執行政策仍以 `shared-rules/` 為準。
