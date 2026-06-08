# Evidence Collapse Point

**Status**: `candidate-intelligence`

## 定義

Evidence collapse point 是診斷推理模式：在 evidence chain 中，第一個失去證據一致性的節點。

它用於從「驗證不足」進一步定位「失敗從哪一段開始」。

## 例子

```text
Login ✓
Cookie ✓
API ✓
DB ✗
SSR ✗
UI ✗
```

Collapse point 是 `DB`。後面的 SSR/UI 失敗可能只是下游結果，不應先從 UI 猜修。

另一個例子：

```text
Producer ✓
Topic ✓
Consumer ✗
DB ✗
Business Effect ✗
```

Collapse point 是 `Consumer`。

## 操作步驟

1. 依 [`evidence-chain-validation.md`](evidence-chain-validation.md) 列出 chain。
2. 逐段標記 evidence 狀態：`confirmed`、`missing`、`contradicted`、`not-observed`。
3. 找出第一個 `missing` 或 `contradicted` 的必要節點。
4. 將後續 failure 視為下游症狀，除非有獨立證據顯示多點故障。
5. 下一步 debugging 優先針對 collapse point 補 instrumentation、readback 或 integration proof。

## 邊界

Collapse point 不等於 root cause。它是 evidence chain 中第一個無法支持 claim 的位置。Root cause 仍可能在上游設計、資料契約、部署設定或外部依賴。

## 相關知識

- [`evidence-chain-validation.md`](evidence-chain-validation.md)
- [`state-visibility-gap.md`](state-visibility-gap.md)
- [`evidence-model.md`](evidence-model.md)
