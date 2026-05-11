# Runtime

`runtime/` 負責「AI 系統如何運作」。本層保存 dynamic loading、context routing、context pruning、agent coordination 與 orchestration 的設計，不取代可執行 shared rules。

## 目前入口

- [`routing/`](routing/README.md)：定義 task intent → knowledge index → metadata → source-of-truth gate 的 context loading 流程。

## 放什麼

- Context routing 與 dynamic loading 規則的設計。
- Runtime orchestration、scheduler、coordination 與 context pruning pattern。
- 如何利用 metadata 選擇 rules、workflow、knowledge atoms 與 model profile。
- AI-native Knowledge Operating System 的 runtime architecture 草案。

## 不放什麼

- 目前必須執行的共用政策；放到 `shared-rules/`。
- 單一工具的 hook、setting、UI 或 mirror sync 細節；放到 `ai-tools/`。
- Skill-specific workflow 全文；放到 `workflow/` 或保留在 `skills/`。
- 長期記憶內容本身；放到 `memory/` 或 `knowledge/`。

## 與既有層的關係

- `shared-rules/` 是目前可執行 policy layer；本層先描述 runtime design。
- `metadata/` 提供 runtime selection 與 ranking 的控制資料。
- `knowledge/` 提供可被 runtime 找到的 atom、index、summary 與 graph。
- `ai-tools/` 記錄各工具如何實作或近似 runtime 行為。

## 第一批候選遷移來源

- `architecture/ai-native-knowledge-operating-system.md` 的 reference-first 與 compatibility inventory
- `architecture/next-stage-upgrade-plan.md` 的 runtime / context routing 章節
- `shared-rules/decision-efficiency.md` 中可抽成 runtime routing design 的概念
