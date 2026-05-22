# ADR-005: Memory Architecture（6 子層記憶模型）

## Status

**Accepted**

## Framework Generation

- **世代分類**：Gen 3 確立
- **當前世代文件**：[`architecture/ai-native-cognitive-execution-system.md`](../architecture/ai-native-cognitive-execution-system.md)
- **適用狀態**：6 子層 memory 模型（working / summary / episodic / project / decision / failure）在 Gen 3 完全有效，並由 `memory/retrieval-governance/` 補上 activation threshold 與 budget 治理。

## Context

AI agent 需要記憶來維持跨 session 的上下文連續性。最初只有一個 `memory/` 目錄，所有類型的記憶混在一起，導致：

1. **Session-local vs long-term 不分**：working memory（可丟棄）與 decision memory（不可變）放在同一層。
2. **無 recall 機制**：沒有 episodic memory，agent 無法 recall 過去在類似情境中學到的經驗。
3. **無專案脈絡**：沒有 project memory，agent 每次處理同一專案都從零開始。
4. **無失效學習**：沒有 failure memory，同樣的錯誤可能重複發生。

需要一個記憶架構，讓不同類型的記憶有各自的 lifecycle、storage 策略與使用方式。

## Decision

採用 **6 子層記憶模型**，將記憶分為三組：

### 短期記憶（Session-local）

| 子層 | 用途 | Lifecycle | Token 限制 |
| --- | --- | --- | --- |
| `working/` | 進行中 task 狀態 | session 結束可丟棄 | 無硬限制 |
| `summary/` | 壓縮 session 歷史 | session 結束時建立 | ≤500 tokens |

### 中期記憶（跨 Session）

| 子層 | 用途 | Lifecycle | Token 限制 |
| --- | --- | --- | --- |
| `episodic/` | 情境經驗 recall | 情境驅動，可 archive | ≤300 tokens |
| `project/` | 專案脈絡保持 | active → archived | ≤500 tokens |

### 長期記憶（不可變 / 抽象化）

| 子層 | 用途 | Lifecycle | Token 限制 |
| --- | --- | --- | --- |
| `decision/` | 輕量 ADR | immutable, numbered | ≤500 tokens |
| `failure/` | 抽象化失效模式 | active → monitored → resolved → archived | ≤400 tokens |

### 核心設計原則

1. **Session boundary**：Working memory 只存活於目前 session；session 結束時重要內容 promotion 到 summary 或 decision。
2. **Recall 不是 re-execution**：Episodic memory 記錄「當時發生了什麼」，不是「如何重新執行」。
3. **抽象化邊界**：Failure memory 不保存專案私有 raw evidence，只保存可泛化的失效模式。
4. **不重複 ADR**：Project memory 只保留架構決策的摘要與連結，完整決策放在 `memory/decision/` 或 `constitution/`。
5. **Token-aware**：每個記憶 record 有明確的 token 上限，避免記憶層級 context pollution。

### 與外部層的關係

- `feedback/replay/` 使用 `memory/episodic/` 作為分析原料。
- `intelligence/engineering/failure/` 從 `memory/failure/` 抽象化 intelligence atom。
- `enforcement/failure-patterns/` 從 `memory/failure/` promotion 可執行規則。
- `governance/lifecycle/` 管理所有記憶子層的 lifecycle（active → archived）。

## Consequences

### 正面

- **精準記憶管理**：不同類型的記憶有各自的 lifecycle 與使用策略。
- **Token 效率**：每種記憶有明確的 token 上限，避免記憶層級爆炸。
- **Recall 能力**：Episodic memory 讓 agent 能在類似情境中快速 recall 過往經驗。
- **失效學習**：Failure memory 提供結構化的失效記錄與 prevention 策略。

### 負面

- **6 子層管理成本**：需要維護 6 個子目錄的索引與 lifecycle。
- **Promotion 路徑複雜**：記憶需要在不同子層之間 promotion（working → summary → decision），需要明確的規則。
- **初期內容稀疏**：episodic/project/failure 在初期 record 較少，效益不明顯。

## Alternatives Considered

- **單一 memory 層**：所有記憶混在一起。但 session-local 與 long-term 的 lifecycle 衝突。不採用。
- **3 子層（working/summary/decision）**：缺少 episodic/project/failure，無法支援 recall、專案脈絡與失效學習。不採用。
- **記憶全部放在 SQLite**：使用資料庫取代檔案系統。但失去 Markdown 的可讀性與 git 的版本管理。不採用。

## Related

- [`memory/working/README.md`](../memory/working/README.md) — 工作記憶
- [`memory/summary/README.md`](../memory/summary/README.md) — 摘要記憶
- [`memory/decision/README.md`](../memory/decision/README.md) — 決策記憶
- [`memory/episodic/README.md`](../memory/episodic/README.md) — 情境記憶
- [`memory/project/README.md`](../memory/project/README.md) — 專案記憶
- [`memory/failure/README.md`](../memory/failure/README.md) — 失效記憶
- [`feedback/replay/README.md`](../feedback/replay/README.md) — 使用 episodic memory 作為原料
- [`enforcement/failure-learning-system.md`](../enforcement/failure-learning-system.md) — 失效學習流程
