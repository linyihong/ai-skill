# Source Triage Heuristic（旅行資訊來源分級經驗法則）

**Status**: `candidate-intelligence`
**Source**: [`workflow/travel-planning/execution-flow.md`](../../workflow/travel-planning/execution-flow.md), [`skills/travel-planning/WORKFLOW.md`](../../skills/travel-planning/WORKFLOW.md)

## 原則

**If a travel claim is not backed by a current official source, label it as unconfirmed and provide a safer backup.**

如果一個旅行宣稱沒有當前官方來源支持，標記為未確認並提供更安全的備案。

## 為什麼

1. **旅行資訊的時效性極高** — 營業時間、票價、路線、住宿 availability 可能在短時間內變化。
2. **社群來源（部落格、影片、論壇）的資訊可能過時或個人化** — 作者的經驗不一定適用於所有旅客。
3. **官方來源雖然不一定最方便，但是最可靠的 baseline** — 官方網站、營運商時刻表、旅遊局頁面是可信度最高的來源。
4. **未經 triage 的資訊會導致行程不可執行** — 到了現場才發現景點休館、交通停駛、住宿客滿。

## 何時適用

- 任何旅行規劃的資訊收集階段。
- 使用者提供未附來源的建議或宣稱時。
- 從社群媒體、部落格、影片中獲取靈感時。
- 比較不同來源的價格或時間時。

## 何時不適用

- 使用者明確表示「這是個人經驗，不需要驗證」。
- 資訊來自使用者自己的 recent experience（例如「我上週才去過」）。
- 通用知識（例如「東京車站位於東京站」）不需要來源。

## 決策流程

```text
有旅行宣稱？
  ├── 這是什麼類型的資訊？
  │     ├── 營業時間/最後入場 → 官方設施頁面、官方 SNS、預約頁面
  │     ├── 套裝行程/模型路線 → 旅行社行程頁面、官方旅遊模型路線
  │     ├── 地點身份 → Google Maps 連結、官方地址、Mapcode
  │     ├── 交通時刻/票價 → 營運商時刻表、官方票價表
  │     ├── 住宿 → 官方預約平台、旅遊局頁面
  │     ├── 道路狀況/通行費 → 道路管理機構、高速公路營運商
  │     └── 靈感/發現 → 地圖、社群、部落格（標記 needs confirmation）
  ├── 來源是官方且當前？
  │     ├── 是 → 可採用
  │     └── 否 → 標記 needs confirmation，提供 safer backup
  └── 記錄到 Source Table
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「部落格說這個景點很棒，直接排進行程」 | 先確認營業時間、門票、交通方式，再排入 |
| 「Google Maps 顯示有營業，應該沒問題」 | Google Maps 資訊可能過時；以官方網站或預約頁面為準 |
| 「這個套裝行程的價格看起來合理」 | 確認價格 basis（每人/每團體/每車）、包含/不包含項目 |
| 「大家都說這條路好開」 | 確認季節性封閉、路面狀況、通行費、停車場 availability |

## Token Impact

避免因使用未驗證的來源導致行程不可執行。一個未驗證的景點可能導致當天行程全部需要重排，浪費 30-60 分鐘的重新規劃時間。

---

← [回到 intelligence/travel/](README.md)
