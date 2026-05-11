# Governance

`governance/` 負責「知識治理與系統維護」。本層保存 cleanup、splitting、lifecycle、validation 與 dependency maintenance 的架構設計，支援知識長期可維護。

## 放什麼

- Knowledge lifecycle、deprecation、archive 與 cleanup strategy。
- Duplicate detection、splitting rules 與 ownership boundary。
- Validation gate、dependency maintenance 與 linked update strategy 的架構化設計。
- 知識治理流程如何協調 `shared-rules/`、`skills/`、`knowledge/` 與 `metadata/`。

## 不放什麼

- 目前可執行的 dependency reading、linked updates 與 close-loop 規則正文；放到 `shared-rules/`。
- Tool-specific sync 或 hook 操作；放到 `ai-tools/` 或 scripts 文件。
- Active conversation goal state；放到 `.agent-goals/`。
- 單一 skill 的 checklist 正文；保留在 `skills/` 或後續依遷移策略拆分。

## 與既有層的關係

- `shared-rules/` 仍是可執行 governance policy 的 source of truth。
- `scripts/` 提供 goal ledger 與 close-loop helper；本層描述治理架構，不取代腳本文件。
- `metadata/` 提供治理可讀取的 ranking、confidence、compatibility 與 lifecycle 控制資料。
- `knowledge/` 的 atom、index、summary 與 graph 需要由本層定義維護責任。

## 第一批候選遷移來源

- `shared-rules/document-sizing.md`
- `shared-rules/linked-updates.md`
- `shared-rules/dependency-reading.md`
- `architecture/ai-native-knowledge-operating-system.md` 的 deprecation checklist
