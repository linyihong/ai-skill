# 決策效率

本規則用來在不過度消耗 context、工具或使用者注意力的前提下，選擇下一個有用行動。目標是在降低無關閱讀、廣泛探索、重複工作與 token-heavy context 的同時，維持決策品質。

本規則泛化了 APK analysis 等技術 skill 的決策路由模式：先界定目前未知，再比較證據路徑，選擇最高收益路線，只載入該路線需要的文件。

## 核心規則

做實質工作前，用一句話說明目前決策點：

```text
Current unknown: <what must be learned or decided next>
```

接著依下列標準選擇下一步：

| 標準 | 問題 |
| --- | --- |
| Time to evidence | 哪條路最快回答目前未知？ |
| Semantic distance | 哪個來源最接近真正決策，而不是只看雜訊症狀？ |
| Safety / reversibility | 哪個行動最不具破壞性、最容易回退？ |
| Validation signal | 哪條路能提供清楚的 pass/fail 或信心更新？ |
| Context cost | 哪些檔案/工具是真的需要，哪些可以等？ |
| User value | 哪個結果最能推進使用者目標或移除 blocker？ |

優先選 evidence-to-cost ratio 最高的路線，而不是第一個出現在 checklist 裡的路線。

## Context Loading

分層載入 context：

1. **Bootstrap：**讀取 shared-rule bootstrap set。
2. **Task frame：**讀取使用者要求、active `.agent-goals/` entry，以及直接相關的開啟檔案。
3. **Skill entry：**讀取符合任務的 `SKILL.md` 與其 routing guidance。
4. **Route-specific docs：**只讀目前路線需要的 workflow/tools/docs 分類。
5. **Deep references：**只有證據顯示必要時，才讀 examples、techniques、feedback lessons 或 source files。

不要預設讀完每個分類或每個 technique。若需要廣泛讀取，說明為什麼需要廣泛 context，以及它支援哪個決策。

## Decision Routing

把 workflow 當作 routing aid，不是僵硬腳本。當 workflow 有多個分支時：

- 從最高層 triage 開始。
- 證據清楚指向某分支時就停止擴散。
- 只讀該分支的詳細文件。
- 其他分支保留為 fallback，不放進 active context。
- 證據推翻目前分支時再重新路由。

如果某個行動已回答決策點，不要只因為有工具或 checklist 就繼續跑更廣或更低層的檢查。

## Token And Noise Reduction

降低 context 與輸出負擔：

- 先摘要大型證據，再決定是否展開。
- 先讀索引，再讀子檔。
- 已知名稱用 exact search；廣泛探索才用 semantic search。
- raw logs、大 payload、screenshots、generated dumps 放在專案 artifacts，回答中引用路徑或去敏摘錄。
- 記錄 open questions，而不是為了填補推測缺口讀無關檔案。
- 可重用但冗長的材料，依 [`document-sizing.md`](document-sizing.md) 移到聚焦子檔。

不得用 token reduction 當作跳過 required dependencies 的理由。若 [`dependency-reading.md`](dependency-reading.md) 要求某依賴，必須讀取，或標為 blocked / not applicable。

## Stop Conditions

遇到下列情況時，停止目前路線並重新評估：

- 目前路線已產生足夠證據回答未知。
- 路線產生雜訊但沒有提高信心。
- 路線相較可用替代方案變得破壞性高、不穩或太慢。
- 出現更高語意距離的來源。
- 使用者優先順序或 blocker 改變目標。
- 另一個 active goal/owner/lock 讓平行工作不安全。

停止時，若該決策影響未來工作，需在回答、active goal 或 document TODO 記錄原因。

## Output Shape

重要決策回報：

```text
Current unknown:
Options considered:
Chosen next action:
Why this is the highest-yield path:
What was deferred:
Validation signal:
```

保持簡短。目的在於讓路由選擇可稽核，不是產生冗長推理 dump。

## 與其他規則的關係

- Required dependency scope 依 [`dependency-reading.md`](dependency-reading.md)。
- 決策指引變長或變成路線專屬時，依 [`document-sizing.md`](document-sizing.md)。
- 目標/執行/驗證閉環依 [`goal-action-validation.md`](goal-action-validation.md)。
- 路線變更產生新 goal、blocker 或 handoff 時，依 [`conversation-goal-ledger.md`](conversation-goal-ledger.md)。

← [Back to shared rules index](README.md)
