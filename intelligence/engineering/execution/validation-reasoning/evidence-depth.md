# Evidence Depth

**Status**: `candidate-intelligence`

## 定義

Evidence depth 是依風險選擇最低證據層級的決策規則。它不只處理測試，也包含 production inspection、log analysis、metric validation、replay validation 與外部確認。

## 層級

```yaml
evidence_depth:
  - component
  - service
  - adapter
  - api
  - live_system
  - independent_observation
```

| Depth | 證明範圍 |
| --- | --- |
| `component` | 單一函式、元件或純邏輯成立 |
| `service` | use case / domain service 層的規則成立 |
| `adapter` | boundary wrapper 能呼叫或轉換成功 |
| `api` | request/response contract 與入口路徑成立 |
| `live_system` | 真實產品路徑跨身份、狀態、持久化、UI 或外部 boundary 成立 |
| `independent_observation` | 由外部系統、獨立記錄、使用者可觀察結果或第二資料源確認 |

## Risk Mapping

```yaml
risk_to_minimum_evidence_depth:
  low:
    - component

  medium:
    - api

  high:
    - live_system

  critical:
    - live_system
    - independent_observation
```

`critical` 適用於 payment、email、external API、storage、queue、entitlement 或其他 proxy success 容易誤導的路徑。

## 決策規則

- Simple pure logic 可以停在 `component` 或 `service`。
- API/schema compatibility 至少需要 `api` 或 contract evidence。
- 命中 [`state-visibility-gap.md`](state-visibility-gap.md) 的高風險流程，最低為 `live_system`。
- 若 side effect 發生在外部系統、非同步 consumer、金流、寄信、儲存或權限授予，增加 `independent_observation`。
- 多個低層證據不能自動相加成高層證據。三個 adapter tests 仍然不是 live system proof。

## 完成宣告

Agent 在宣稱完成時應說明：

1. Claim 對應的風險等級。
2. 已取得的 evidence depth。
3. 尚未覆蓋的 propagation segment。
4. 若低於最低 depth，縮小完成宣告或標記 blocked。

## 相關知識

- [`evidence-model.md`](evidence-model.md)
- [`evidence-chain-validation.md`](evidence-chain-validation.md)
- [`state-visibility-gap.md`](state-visibility-gap.md)
