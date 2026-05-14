# Workflow

`workflow/` 負責「AI 如何執行工作」。本層保存 agent 可照著執行的 planning flow、task decomposition、review flow、orchestration flow 與 handoff flow。

## 目前入口

- [`apk-analysis/`](apk-analysis/README.md)：`apk-analysis` pilot 的 tool-neutral workflow 候選目的地；目前仍保留 `skills/apk-analysis/SKILL.md` 作為 active skill entrypoint。
- [`software-delivery/`](software-delivery/README.md)：軟體交付的執行流程（design review、code review、release review、security review）。
- [`repo-analysis/`](repo-analysis/README.md)：Repository 分析的執行流程（new onboarding、deep codebase analysis、migration impact、tech debt assessment）。
- [`travel-planning/`](travel-planning/README.md)：旅遊規劃的執行流程（itinerary planning、transportation research、budget planning）。

## 放什麼

- Skill 或任務類型的執行流程。
- 規劃、拆解、review、驗證與交接步驟。
- 如何引用 `analysis/`、`intelligence/`、`metadata/` 來完成工作。
- 可被不同 AI tools 套用的工具中立 workflow。

## 不放什麼

- 深層分析方法全文；放到 `analysis/`。
- 工程智慧、trade-off 與 domain lesson；放到 `intelligence/`。
- 工具專屬操作細節；放到 `ai-tools/` 或 skill-local `tool-adapters/`。
- Conversation goal ledger state；放到 `.agent-goals/`。

## 與既有層的關係

- `skills/` 仍是目前主要 capability entry；本層會逐步承接 skill 中的通用執行流程。
- `shared-rules/` 仍提供 dependency reading、linked updates、validation 與 source boundary。
- `analysis/` 與 `intelligence/` 是 workflow 的知識來源，不應被大量複製。
- `runtime/` 未來可依 metadata 動態載入合適 workflow。

## 第一批候選遷移來源

- `skills/apk-analysis/WORKFLOW.md`
- `skills/app-development-guidance/WORKFLOW.md`
- `skills/travel-planning/` 中可工具中立化的 itinerary workflow
- `architecture/next-stage-upgrade-plan.md` 中 `workflow/` 的分層說明
