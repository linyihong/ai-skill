# Workflow

`workflow/` 負責「AI 如何執行工作順序」。本層保存 agent 可照著執行的 planning flow、task decomposition、review flow、orchestration flow 與 handoff flow。

判斷原則：如果內容核心是在回答「先做 A 再做 B，如何收口與交接」，通常屬於 `workflow/`。如果只是某個技術路線內部如何取得證據，優先放在對應 `analysis/<domain>/`；如果是在回答何時選哪條路或如何避免錯誤，放到 `intelligence/`。

## Workflow 選路（先做這步）

察覺任務後、動手前，先讀 **[`workflow-routing.md`](./workflow-routing.md)**：依任務性質對照選路表 → 查 [`routing-registry.yaml`](../knowledge/runtime/routing-registry.yaml) 的 `route.workflow.*` → 進入對應子目錄 README + `execution-flow.md`。

[`activation-table.md`](../runtime/router/activation-table.md) 的 **#27（Workflow 編排閘門）** 與 registry 內各 `route.workflow.*.activation_triggers` 會觸發上述 discovery；**registry-first**，不為每個 workflow 新增 activation 列。

## 目前入口

- [`apk-analysis/`](apk-analysis/README.md)：APK 分析的 tool-neutral workflow。已從舊 `skills/apk-analysis/` 遷移至本層。
- [`software-delivery/`](software-delivery/README.md)：軟體交付的執行流程（design review、code review、release review、security review）。
- [`repo-analysis/`](repo-analysis/README.md)：Repository 分析的執行流程（new onboarding、deep codebase analysis、migration impact、tech debt assessment）。選路見 [`workflow-routing.md`](./workflow-routing.md)；routing 條目以 registry 為準。
- [`travel-planning/`](travel-planning/README.md)：旅遊規劃的執行流程（itinerary planning、transportation research、budget planning）。
- [`documentation/`](documentation/README.md)：跨專案 **agent 友善文件** 的撰寫與分類流程（index-first、降低無效 token、與 `enforcement/` / `governance/` 對齊）。

## 放什麼

- Workflow 或任務類型的執行流程。
- 規劃、拆解、review、驗證與交接步驟。
- 如何引用 `analysis/`、`intelligence/`、`metadata/` 來完成工作。
- 可被不同 AI tools 套用的工具中立 workflow。

## 不放什麼

- 深層分析方法全文；放到 `analysis/`。
- 單一技術路線的命令模板、hook 細節、dump procedure；放到對應 `analysis/<domain>/`。
- 工程智慧、trade-off 與 domain lesson；放到 `intelligence/`。
- Raw logs、case dump、專案 findings 或一次性 execution transcript；留在業務專案 evidence，去敏後再進 `feedback/history/`。
- 工具專屬操作細節；放到 `ai-tools/` 或 workflow 的 tool adapter 說明。
- Conversation goal ledger state；放到 `.agent-goals/`。

## 誰會參考這裡（Inbound References）

- [`route.workflow.apk-analysis`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `workflow/apk-analysis/README.md`
- [`route.workflow.software-delivery`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `workflow/software-delivery/README.md`
- [`route.workflow.travel-planning`](../knowledge/runtime/routing-registry.yaml) — candidate_sources 引用 `workflow/travel-planning/README.md`
- [`route.workflow.documentation-ai-native`](../knowledge/runtime/routing-registry.yaml) — `primary_source` 引用 `workflow/documentation/README.md`
- [`route.intelligence.apk-highest-leverage-path`](../knowledge/runtime/routing-registry.yaml:251) — required_dependencies 引用 `workflow/apk-analysis/execution-flow.md`、`workflow/apk-analysis/artifact-gates.md`

## 與既有層的關係

- 舊 `skills/` scaffold 已退役；本層是 capability execution 的 active entrypoint。
- `enforcement/` 仍提供 dependency reading、linked updates、validation 與 source boundary。
- `analysis/` 與 `intelligence/` 是 workflow 的知識來源，不應被大量複製。
- `runtime/` 未來可依 metadata 動態載入合適 workflow。

## 遷移狀態

- `skills/apk-analysis/WORKFLOW.md` — ✅ 已遷移至 `workflow/apk-analysis/execution-flow.md`
- `skills/app-development-guidance/WORKFLOW.md` — ✅ 已遷移至 `workflow/software-delivery/execution-flow.md`
- `skills/travel-planning/WORKFLOW.md` — ✅ 已遷移至 `workflow/travel-planning/execution-flow.md`
- `plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md` 中 `workflow/` 的分層說明
