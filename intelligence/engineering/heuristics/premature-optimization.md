# Premature Optimization Heuristic（過早最佳化經驗法則）

**Status**: `candidate-intelligence`
**Source**: 通用軟體工程經驗

## 原則

**If performance issue is not measured, optimization is likely harmful.**

未經測量的效能最佳化，大概率是有害的。

## 為什麼

1. 開發者對瓶頸的直覺經常是錯的 — 真正的瓶頸往往在預期之外的地方。
2. 最佳化通常增加程式碼複雜度、降低可讀性、引入耦合。
3. 未經測量的最佳化可能解決了不存在的問題，卻引入了真實的維護成本。

## 何時適用

- 在 code review 中看到「我覺得這樣比較快」的變更。
- 在設計階段有人提議用複雜的 cache strategy 或 data structure，理由是「以後可能會慢」。
- 在架構討論中有人主張用 microservices 或 event-driven architecture，理由是「以後 scale 需要」。

## 何時不適用

- 已知的效能問題已經過 profiling 或 load test 驗證。
- 該最佳化是架構級別的選擇，且後續無法低成本重構（例如選資料庫、選語言）。
- 該最佳化是 security 或 compliance 要求的一部分。

## 決策流程

```text
有效能疑慮？
  ├── 有 measurement / profiling data？
  │     ├── 是 → 根據數據決定最佳化策略
  │     └── 否 → 先建立 benchmark 或 profile，再決定
  │
  └── 沒有效能疑慮，但「覺得以後會慢」？
        ├── 該決策後續無法低成本改變？
        │     ├── 是 → 謹慎評估，記錄 tradeoff
        │     └── 否 → 先不做，等 measurement 出現
        └── 該決策是 security / compliance 要求？
              ├── 是 → 照做，但記錄原因
              └── 否 → 先不做
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「N+1 query 以後會是問題，先加 batch layer」 | 先確認 N+1 是否真的發生在 hot path，且延遲可被測量 |
| 「先用 Redis cache，以後再拿掉」 | Cache 一旦引入，拿掉的成本遠高於一開始不加 |
| 「用 message queue 解耦，以後 scale 方便」 | MQ 增加除錯複雜度與 eventual consistency 風險；沒有明確 producer/consumer boundary 時不要引入 |

## Token Impact

避免浪費 engineering 時間在解決不存在的問題上。一個未經測量的最佳化可能花費 2-5 天實作，卻只節省 0.1% 的執行時間。

---

← [回到 engineering/heuristics/](README.md)
