# Escalation Policy

本規則定義 execution 中段的 mismatch escalation：當實際 evidence 推翻原假設時，agent 必須停止局部 patch，重新讀 source-of-truth 並重建 execution graph。

## 核心規則

一旦出現 L3 以上 escalation trigger，不要繼續沿用舊 mental model。先進入 recovery frame：

```text
suspend assumption -> reload source-of-truth -> rebuild execution graph -> resume execution
```

Escalation 是 real-time control；它處理「現在不能繼續錯路線」。Failure learning 是 post-mortem prevention；它處理「之後如何避免重犯」。

## Trigger Classes

| 類別 | 觸發條件 | 預設層級 | 必要動作 |
| --- | --- | --- | --- |
| `repeated-failure` | 同一路徑、同一 automation、同一 checkpoint 連續失敗 2 次 | L3 | 停止 retry，reload source-of-truth |
| `user-contradiction` | 使用者指出「不是這樣」「你沒看文件」「你又在猜」 | L4 | 進入 recovery，重建 execution graph |
| `evidence-conflict` | UI / API / repo structure 與 workflow、contract 或 README 衝突 | L4 | reload workflow primary source 與 owner docs |
| `assumption-drift` | agent 開始靠猜 route、猜座標、沿用 stale checklist | L3 | suspend assumption，補 source reading |
| `source-of-truth-miss` | 重要操作前未讀 canonical workflow / UI map / contract / owner docs | L4 | rediscovery + dependency read ledger |
| `automation-drift` | 腳本反覆 patch 但 checkpoint 截圖、foreground、feature context 不對 | L4 | 禁止繼續自動操作，回到 navigation graph |

## Recovery Levels

| 層級 | 意義 | 典型動作 |
| --- | --- | --- |
| L1 | Simple retry | 單次工具失敗，可重試一次 |
| L2 | Reload local workflow | 重新讀當前 workflow / rule |
| L3 | Reload source-of-truth | 補讀 owner docs、UI map、contract、architecture |
| L4 | Rebuild execution graph | 停止執行，重建 goal -> validation chain |
| L5 | Assumption collapse + rediscovery | 原 routing 或 task framing 可能錯，回到 discovery |

## Forbidden Behaviors

進入 L3 以上 escalation 後，禁止：

- 繼續 patch automation。
- 繼續猜 UI route、API route、repo architecture。
- 不重讀 workflow primary source。
- 不重建 execution graph。
- 沿用已被反證的 checklist。
- 把 log、Frida event、partial API evidence 或 target PID 當成 UI / workflow checkpoint 成功證據。

## Required Recovery Output

Escalation 後的回覆或工作筆記必須能回答：

| 欄位 | 必填內容 |
| --- | --- |
| 原假設 | 原本以為什麼 |
| 反證 | 哪個 evidence 推翻 |
| 新 source-of-truth | 現在依據哪些文件或證據 |
| 新 execution graph | goal -> route -> dependencies -> checkpoint -> validation |
| 下一步 | 恢復執行前的最小安全行動 |
| Validation | 如何確認不再沿用錯誤 mental model |

## Source-of-Truth Reload

依任務類型重讀最小必要 source：

| 情境 | 必要重讀集合 |
| --- | --- |
| Workflow execution mismatch | workflow primary source、artifact gates、routing registry record、owner docs |
| APK navigation mismatch | `workflow/apk-analysis/execution-flow.md`、`workflow/apk-analysis/artifact-gates.md`、專案 UI map / API catalog（若存在） |
| Software delivery mismatch | owner contract、BDD / acceptance scenarios、implementation plan、test evidence |
| Repo structure mismatch | root README、layer README、routing registry、canonical source path |
| Tool / runtime mismatch | tool adapter、runtime README、relevant YAML or embedded runtime source |

如果 required source 不存在，標記 `not applicable` 或 `source missing`；不得假裝已讀。

## 與 Failure Learning 的關係

| 問題 | Escalation Policy | Failure Learning System |
| --- | --- | --- |
| 何時用 | failure 或 mismatch 正在發生 | failure 已被控制後 |
| 目的 | 停止錯路線並恢復正確 execution graph | 建立 durable prevention，避免重犯 |
| 產出 | recovery output、source reload、new execution graph | failure pattern、feedback lesson、validation scenario |
| 不可取代 | 不能當作長期 prevention | 不能取代 real-time recovery |

若 mismatch 已造成可重用失效模式，完成 recovery 後再依 [`failure-learning-system.md`](failure-learning-system.md) 分類與沉澱。

## 驗證

完成 escalation 後，確認：

- Trigger class 與 level 已命名。
- 舊假設已明確 suspend。
- Required source 已讀或標記不適用 / missing。
- 新 execution graph 已重建。
- 禁止事項沒有繼續發生。
- 若需要 durable prevention，已開啟 failure learning loop。

← [回到 enforcement index](README.md)
