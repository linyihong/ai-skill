# Highest Leverage APK Analysis Path

本文件把 `apk-analysis` 中已驗證到工作流層的「最高收益路線」整理成 engineering intelligence。它是決策判斷，不取代 `skills/apk-analysis/SKILL.md`、`WORKFLOW.md` 或原始 feedback lesson。

## Intelligence Status

| 欄位 | 值 |
| --- | --- |
| Status | `candidate-intelligence` |
| Source skill | `skills/apk-analysis/` |
| Source lesson | `skills/apk-analysis/feedback_history/common/2026-05-07_131000-highest-leverage-analysis-path.md` |
| Current workflow anchor | `skills/apk-analysis/WORKFLOW.md` § `2. 選擇主線` |
| Runtime route | `route.intelligence.apk-highest-leverage-path` |

## Decision Rule

每個 APK 分析 checkpoint 先界定目前未知，再選擇最能快速產生可驗證證據的路線。工作流是 routing aid，不是固定順序。

比較路線時使用下列問題：

| Criterion | Question |
| --- | --- |
| Time to evidence | 哪條路最快回答目前未知？ |
| Semantic proximity | 哪條路最接近業務物件、API contract、decode boundary 或 UI attribution？ |
| Safety / reversibility | 哪條路最不會寫入狀態、觸發風控或污染 evidence？ |
| Validation clarity | 哪條路能產生明確 pass / fail / blocker？ |
| Boundary preservation | 是否保留 App-owned session、signing、decrypt、gateway 或 UI operation path？ |

## Preferred Pattern

1. 用一句話寫出目前 unknown。
2. 列出 2-4 條已被目前 evidence 支援的 route。
3. 選擇 evidence-to-cost ratio 最高的主線。
4. 記錄 fallback route 與切換條件。
5. 如果較慢 route 仍需要做 attribution 或 edge-case confirmation，明確標成 deferred，而不是丟棄。

## Examples Of Route Choice

| Situation | Higher-leverage route | Deferred / fallback route |
| --- | --- | --- |
| 已有 app-owned read-only API boundary，可保留 session / signing / decrypt | 短窗參數覆寫或 API-first replay | 長 UI scroll，只保留作 UI attribution |
| 已找到高語意 request / response / decode hook | 高語意 hook 取得 schema / key set / timing | socket bytes 或 broad native hooks |
| Java HTTP broad hooks 在使用者操作下仍無業務 host | Flutter / Dart AOT、native connect、pcap SNI 或有效 MITM path | 繼續擴大 Java hooks |
| UI behavior 本身是問題 | UI capture / bounded gesture / package + feature context guard | 直接 API replay 只能當輔助 |

## Validation Signal

此 intelligence 可用於任務時，agent 應能反查：

- 已描述目前 unknown。
- 已比較至少兩條可行 route。
- 選定 route 有明確 validation signal。
- 延後 route 有原因與恢復條件。
- 若 attribution 重要，結果會回到 UI / operation evidence 對照。

## Boundaries

- 不可用本文件跳過授權、去敏、source-of-truth 或 skill dependency reads。
- 不可把 faster route 當成允許猜測 secrets、繞過 scope 或執行未授權寫入。
- 不可把本文件當作完整 APK workflow；實作分析仍從 `skills/apk-analysis/SKILL.md` 與 `WORKFLOW.md` 進入。
