# Multi-Layer Enforcement（多層強制執行）

**Status**: `candidate-intelligence`
**Source**: Bootstrap 三層架構 empirical 測試（2026-05-27）

## 原則

**To reliably enforce an agent obligation, combine context injection (data accuracy), prohibition prompting (awareness), and mechanical gate (failsafe) — each layer addresses a distinct failure mode that the others cannot prevent.**

要可靠地強制 agent 遵守某個 obligation，需要三層協同：context 注入（資料精確性）、禁止語氣 prompt（意識觸發）、機械關卡（兜底）— 每一層防止其他層無法防止的不同失效模式。

## 三層設計模式

| 層 | 機制 | 防止的失效模式 | 缺少此層的症狀 |
|---|---|---|---|
| **L1 資料注入層** | 執行期把預算好的精確值注入 context | Agent 編造數值（看起來正確但實際上是猜測） | Agent 輸出的數值與 source of truth 不符 |
| **L2 禁止語氣層** | 文件最前段使用禁止語氣 + 顯式排序 + execution prohibition | Agent 在「任務看起來簡單」時跳過 obligation | 簡單任務被直接回答，obligation 被忽略 |
| **L3 機械關卡層** | PreToolUse / PreAction hook 以 non-zero exit code 攔截 | Agent 偶發性繞過 prompt-based 規則 | 即使有 L2，仍有一定比例的繞過行為 |

## 為什麼單層不夠

| 單層方案 | 不夠的原因 |
|---------|-----------|
| 只有 L1（資料注入）| Agent 有資料但不一定先輸出；任務看起來簡單時仍會跳過 |
| 只有 L2（禁止 prompt）| Agent 沒有精確資料，被迫猜測數值；無機械兜底 |
| 只有 L3（機械關卡）| 關卡觸發前 agent 已可能輸出錯誤；no data = 猜值 |
| L1 + L2（無機械關卡）| Agent 偶發性繞過仍然存在，prompt-based 強度有上限 |
| L2 + L3（無資料注入）| Agent 有意識遵守但數值不精確（猜測） |

## Prompt 強度位階

實測 agent 對以下 prompt 模式的遵守強度（遞增）：

| 強度 | Prompt 形式 | 遵守可靠度 |
|------|------------|-----------|
| 低 | 一般敘述（「請先做 X」） | 容易跳過 |
| 中低 | IMPORTANT block | 仍可能跳過 |
| 中高 | 禁止語氣 + 顯式排序 + execution prohibition | 顯著提升 |
| 高 | 中高 + 機械關卡兜底 | 最可靠 |

## Primacy Effect（位置效應）

**放在文件最前面比放在中段或末尾顯著有效。**

- 當 obligation 文字在 context 最前方時，agent 更難「忽略」它
- 被其他 context 包圍的規則較易被 recency bias 或 rule-overload 覆蓋
- → 高優先級 obligation 的 prompt text 應出現在所有相關文件的最前段

## 適用時機

- 設計 agent 需要強制遵守的 session-level obligation（如 Bootstrap Receipt、Cognitive Mode 報告）
- 設計任何「agent 可能因為任務看起來簡單而跳過」的流程
- 評估現有 enforcement 機制是否有覆蓋所有三層

## 不適用時機

- 一般 task instruction（不需要 session-level enforcement）
- 低風險的 best practice 提醒（提示即可，不需機械關卡）

## 實作注意

具體工具的 hook 設定格式、matcher 字串、已知限制等**實作細節**因工具而異，放在各工具 adapter 文件中（如 `ai-tools/agent/claude.md`），本 atom 只保留工具中立的設計原則。

## 相關 atoms

- [`cognitive-boundaries.md`](cognitive-boundaries.md) — agent 無法可靠地自我檢測邊界，需外部關卡
- [`rule-overload.md`](rule-overload.md) — 規則過多時 agent 選最近/最具體的，不是最重要的
- [`stateless-validation-necessity.md`](stateless-validation-necessity.md) — AI 決策路徑驗證必須是 stateless 的

---

← [回到 agent-architecture/](README.md)
