# Rule Overload（規則超載）

**Status**: `candidate-intelligence`
**Source**: 本系統實際運作觀察

## 原則

**When too many rules compete for attention, agent follows the most recently loaded or most concrete rule, not the most important one.**

當太多規則競爭注意力時，agent 會遵循最近載入或最具體的規則，而不是最重要的規則。

## 為什麼

1. **規則不是權重驅動的** — Agent 沒有內建的規則優先級排序器。它依賴 recency 和 specificity 來決定哪條規則適用。
2. **最近載入的規則有優勢** — 如果 `dependency-reading.md` 在第 10 步載入，`document-sizing.md` 在第 50 步載入，後者可能覆蓋前者的行為。
3. **具體規則勝過抽象規則** — 「檔案超過 300 行要拆分」比「保持文件簡潔」更有行為影響力，即使後者更重要。
4. **規則之間的真實衝突很少被檢測** — Agent 通常不會主動檢查規則是否矛盾，而是選擇最近看到的那條。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **規則跳躍** | Agent 在不同 session 對相同情境採用不同規則 | 高 |
| **忽略核心規則** | Agent 遵守了文件拆分規則，卻違反了「不刪除舊檔案」的更高優先級規則 | 高 |
| **規則選擇不一致** | 同一個 session 中，前半段遵守規則 A，後半段遵守規則 B | 中 |
| **過度遵守次要規則** | Agent 花大量時間滿足一個低優先級規則的格式要求，卻忽略了主要目標 | 中 |

## 預防方式

1. **規則分層** — 將規則分為 always-apply（Core Bootstrap）和 lazy-load（條件觸發），減少同時競爭的規則數量
2. **明確優先級** — 使用 `rule-weight.md` 標記規則權重，讓 agent 有判斷依據
3. **減少 lazy-load 規則** — 每個 domain 的 lazy-load 規則不超過 5-7 條
4. **規則去重** — 定期檢查規則是否重複或矛盾
5. **測試規則衝突** — 使用 `validation/scenarios/` 測試在特定 signals 下 agent 是否選擇正確規則

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 把所有規則都設為 always-apply | 規則越多，agent 越難區分優先級 |
| 依賴 agent 自行判斷規則優先級 | Agent 沒有可靠的優先級排序能力 |
| 規則之間有隱含矛盾而不處理 | Agent 不會主動發現矛盾 |

## 相關 atoms

- [`attention-budgeting.md`](attention-budgeting.md) — 注意力預算管理
- [`cognitive-boundaries.md`](cognitive-boundaries.md) — 認知邊界
- [`context-collapse.md`](context-collapse.md) — 上下文崩塌

## Token Impact

Rule overload 導致 agent 花費 token 在規則選擇和切換上，而不是實際任務執行。每條多餘的 lazy-load 規則約消耗 200-500 token 的載入成本，10 條就是 2K-5K token。

---

← [回到 agent-architecture/](README.md)
