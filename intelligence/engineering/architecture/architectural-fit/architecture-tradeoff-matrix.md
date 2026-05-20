# Architecture Tradeoff Matrix

**Status**: `candidate-intelligence`

| Strategy | 適合 | 成本 | 避免情境 |
| --- | --- | --- | --- |
| CRUD | 低規則、短流程、admin/internal tool | domain model 薄弱 | 高 invariant domain |
| Vertical Slice | 功能快速演化、團隊小 | 橫向共享規則可能重複 | 強一致 domain core |
| Simple Service Layer | 中低複雜度、交付速度重要 | service 可能變 transaction script | 長期高 domain complexity |
| DDD Lite | 中度 domain language / invariant | 需維護 glossary 與 boundary | 純 CRUD |
| Full DDD | 高 domain complexity、高 invariant、長期核心 | 建模成本高 | MVP / low complexity |
| Event-driven | 跨 context coordination、audit / replay | debugging 與一致性成本 | 無事件業務價值 |
| Microservices | 獨立 deployment、scaling、ownership | 運維與分散式成本 | 只是想拆資料夾 |

## 使用方式

架構建議應引用本矩陣說明採用與拒絕理由，並回到 `architecture-fit-analysis.md` 補 evidence。
