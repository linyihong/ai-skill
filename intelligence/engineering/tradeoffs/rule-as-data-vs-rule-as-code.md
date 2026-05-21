# Rule-as-Data vs Rule-as-Code

**Status**: `candidate-intelligence`

把業務規則 / governance policy 外部化成 data（JSON、YAML、DB row）並由通用 engine 執行，可換取「不 recompile 就能改規則」的彈性；代價是 debug、版本治理、schema 演化、engine ↔ data 升級協調的長期成本。只有當規則**變動頻率高、變動者非工程師、或規則本身就是產品**時才值得付這個代價。

## 兩條路徑

| 維度 | Rule-as-Code | Rule-as-Data |
|------|--------------|--------------|
| 表達力 | 完整語言 | 受 engine schema 限制 |
| Debug | IDE / breakpoint / stack trace | 需要 engine trace 工具，自己造 |
| 版本控制 | git diff 即語意 diff | 需要 schema-aware diff |
| 修改成本 | 改 code + deploy | 改 data（理想），但要過 schema 驗證 |
| 修改者 | 工程師 | 業務 / 規則作者 / 工程師都可 |
| 升級協調 | 單一 artifact | engine 與 data 雙頭升級，有 ordering 風險 |
| 測試 | 標準單元測試 | 需要 rule fixture + engine 雙重測試 |
| 觀察性 | log / APM 直接看 | 需要 engine 自帶 evaluation trace |

## 收益門檻

Rule-as-data 只在以下條件**至少滿足兩個**時划算：

1. **變動頻率**：規則平均每月修改 ≥ 1 次。
2. **作者非工程師**：法遵 / 風控 / 業務 / 規則委員會直接擁有規則。
3. **規則即產品**：規則本身是對外可配置的 feature（multi-tenant policy、insurance product 條款）。
4. **規模**：規則數量 ≥ 20，code path 已經難以維護。

只滿足一條時，**rule-as-code + 良好命名 + 清楚 unit test** 通常更便宜。

## 常見隱性成本（容易被低估）

- **Schema drift**：規則格式演化時，舊 data 沒一起遷移，runtime 出現 silent fallback。
- **Engine / data 升級 ordering**：新規則用到 engine 還沒實作的特性，或 engine 升級後舊規則語意改變。
- **「半 code」陷阱**：data 裡開始長出 `expression`、`condition`、`script` 欄位 → 你正在重新發明一個沒有 type system 的程式語言。
- **Debug 黑洞**：規則執行失敗時，stack trace 只指到 engine 內部，不指到出問題的 rule 行。
- **版本鎖**：規則 data 沒有版本，prod 跑 v2 規則但程式碼預期 v1，事故無法歸因。

## 對 SQLite-as-canonical governance runtime 的應用

本專案的 `obligation_ledger` / `governance_gates` / `language_policy` / `output_rules` 本質上是 rule-as-data。已經跨過門檻（規則數量多、變動頻率高、跨 phase 規則作者非單一工程師），所以這條路是對的。但必須明文承擔以下三個對應投資：

1. **Rule schema version 欄位 + migration policy**：避免 schema drift。
2. **Evaluation trace**：每次 gate 評估的 input / matched rule / output 應可 replay。
3. **Engine ↔ rule data 升級 ordering 規約**：明文寫死「engine 升級不得 break 既有 rule data 語意」或反之的策略。

缺一個，rule-as-data 就會在 18–24 個月內退化成不可維護狀態。

## Token Impact

中。在 governance / policy / multi-tenant 設計階段 lazy-load，約 800 tokens。避免 rule-as-data 路徑被默認選擇後，18–24 個月內陷入「engine ↔ data 雙頭升級失序」與「半 code data」的不可維護狀態。

## 何時退回 rule-as-code

- 規則變動頻率掉到 < 1 次 / 季。
- 規則作者只剩工程師。
- 出現 ≥ 3 次「rule data 語意 bug 比 code bug 更難 debug」的事故。

退路存在，不是失敗。

## 反訊號（rule-as-data 已經失控）

- Rule data 裡出現巢狀 `if/else` 結構。
- 為了表達某條規則，需要新增 engine feature。
- Rule 作者開始問「我可以在 rule 裡呼叫 function 嗎」。
- 修一個 rule bug 同時要改 engine + data + migration + test fixture。

任何一個出現，重新評估邊界。

---

← [回到 engineering/tradeoffs/](README.md)
