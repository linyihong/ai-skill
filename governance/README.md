# Governance

`governance/` 負責「知識治理與系統維護」。本層保存 cleanup、splitting、lifecycle、validation 與 dependency maintenance 的架構設計，支援知識長期可維護。

## 目前入口

- [`contributing.md`](contributing.md)：人類維護者與 PR 的驗證流程入口（指令索引與權威文件連結）。
- [`lifecycle/`](lifecycle/README.md)：定義舊 `skills/` 與新分層之間的 lifecycle、promotion gates 與 source-of-truth 保護。
  - [`directory-structure-governance.md`](lifecycle/directory-structure-governance.md)：目錄結構治理流程 — 新增/改名目錄前的 5 步驟 Checkpoint（名稱衝突、邊界清晰度、慣性命名、路徑深度、全域引用影響）。
- [`validation/`](validation/README.md)：定義新分層變更的 validation gates、migration checklist 與 pass / block rules。
- [`cleanup/`](cleanup/README.md)：定義知識重複偵測、拆分規則與所有權邊界。
- [`dependency/`](dependency/README.md)：定義知識依賴圖的維護規則與連動更新流程。
- [`document-sizing.md`](document-sizing.md)：文件大小與拆分原則 — 定義文件拆分門檻、決策流程、拆分後必做事項與建議結構。

## 放什麼

- Knowledge lifecycle、deprecation、archive 與 cleanup strategy。
- Duplicate detection、splitting rules 與 ownership boundary。
- Validation gate、dependency maintenance 與 linked update strategy 的架構化設計。
- 知識治理流程如何協調 `enforcement/`、`skills/`、`knowledge/` 與 `metadata/`。

## 不放什麼

- 目前可執行的 dependency reading、linked updates 與 close-loop 規則正文；放到 `enforcement/`。
- Tool-specific sync 或 hook 操作；放到 `ai-tools/` 或 scripts 文件。
- Active conversation goal state；放到 `.agent-goals/`。
- 單一 skill 的 checklist 正文；保留在 `skills/` 或後續依遷移策略拆分。

## 誰會參考這裡（Inbound References）

- [`route.governance.durable-goal-boundary`](../knowledge/runtime/routing-registry.yaml:127) — required_dependencies 引用 `governance/lifecycle/README.md`
- [`route.feedback.promotion-pipeline`](../knowledge/runtime/routing-registry.yaml:285) — required_dependencies 引用 `governance/lifecycle/README.md` 與 `governance/validation/README.md`
- [`enforcement/failure-learning-system.md`](../enforcement/failure-learning-system.md) — 引用 governance 的 lifecycle 規則
- [`plans/archived/next-stage-upgrade-plan.md`](../plans/archived/next-stage-upgrade-plan.md) — 引用 `governance/lifecycle/` 的 Skills Deprecation Timeline

## 與既有層的關係

- `enforcement/` 仍是可執行 governance policy 的 source of truth（例如可重用文件的語言與用語見 [`enforcement/neutral-language.md`](../enforcement/neutral-language.md)；`governance/validation` 以 **Neutral language** gate 掛鉤該規則，而非在 `governance/` 重複條文）。
- `scripts/` 提供 goal ledger 與 close-loop helper；本層描述治理架構，不取代腳本文件。
- `metadata/` 提供治理可讀取的 ranking、confidence、compatibility 與 lifecycle 控制資料。
- `knowledge/` 的 atom、index、summary 與 graph 需要由本層定義維護責任。

## 第一批候選遷移來源

- `enforcement/linked-updates.md`
- `enforcement/dependency-reading.md`
- `architecture/ai-native-knowledge-operating-system.md` 的 deprecation checklist
