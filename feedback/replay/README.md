# Experience Replay

`feedback/replay/` 定義「經驗重播」的系統設計。本目錄保存如何從過往 session、failure 與成功經驗中提取可重複使用的教訓，讓 agent 不需要每次都從零開始學習。

## 核心責任

- Session 經驗的回顧與結構化。
- Failure 經驗的 replay 流程（從 incident 到 generalized lesson）。
- 成功模式的 replay 流程（從成功實作到可重複 pattern）。
- Replay 的觸發條件（何時應該 replay、何時不該）。
- Replay 結果的儲存與索引（哪些進 feedback_history、哪些進 intelligence）。

## 核心原則

1. **Replay 不是重複執行**。Replay 是回顧過往經驗，提取可泛化的教訓，不是重新執行舊 session。
2. **Failure replay 優先於 success replay**。失敗經驗通常有更高的學習價值，且更容易泛化為 prevention gate。
3. **Replay 結果必須泛化**。包含專案特定 raw evidence 的 replay 結果不應 promotion；只有去除專案細節後的 generalized lesson 才能進入 `intelligence/` 或 `workflow/`。
4. **Replay 有成本**。每次 replay 消耗 token 與時間，應設定觸發門檻避免過度 replay。

## Replay 觸發條件

| 條件 | 說明 | 優先級 |
| --- | --- | --- |
| 同一類型 failure 發生 ≥2 次 | 表示不是一次性問題，需要 generalized prevention。 | high |
| Session 結束時有未解決的 blocker | Blocker 可能隱含流程或知識缺口。 | high |
| 新 skill 或 workflow 首次使用後 | 首次使用通常會發現流程缺口或模糊地帶。 | medium |
| 定期回顧（每 N 個 session） | 定期 replay 可發現漸進式退化。 | low |
| Agent 偵測到 context pollution 或 circuit breaker 觸發 | 系統層級問題可能需要 replay 來調整 runtime 配置。 | medium |

## Replay 流程

```
1. 觸發 replay
   ├─ 自動觸發（failure 重複、circuit breaker 觸發）
   └─ 手動觸發（session 結束回顧、定期維護）

2. 收集 replay 素材
   ├─ Session summary（memory/summary/）
   ├─ Failure record（enforcement/failure-patterns/ 或 feedback_history/）
   ├─ Decision record（memory/decision/ 或 decisions/）
   └─ Context health / circuit breaker 記錄（runtime/health/、runtime/guards/）

3. 分析模式
   ├─ 這是已知 failure pattern 的變體？
   ├─ 這是新的 failure pattern？
   ├─ 這是流程缺口（workflow 需要更新）？
   ├─ 這是知識缺口（intelligence 需要新增）？
   └─ 這是工具或設定問題（ai-tools/ 需要更新）？

4. 產出 replay 結果
   ├─ Generalized lesson（寫入 feedback_history/ 或 intelligence/）
   ├─ Workflow 更新建議（寫入 workflow/）
   ├─ Prevention gate 建議（寫入 enforcement/failure-patterns/）
   └─ Runtime 配置建議（寫入 runtime/ 或 tools/）

5. 決定 promotion 路徑
   ├─ 立即 promotion（高信心、已驗證）
   ├─ Candidate（需要更多驗證）
   └─ 記錄但不 promotion（低信心、一次性問題）
```

## Replay 結果格式

每次 replay 應產出：

```yaml
replay_id: <YYYY-MM-DD>-<slug>
trigger: <failure_repeat | session_end | periodic | circuit_breaker>
source_sessions:
  - <session_id_or_path>
pattern_type: <known_variant | new_pattern | workflow_gap | knowledge_gap | tool_issue>
generalized_lesson: |
  <去敏後的泛化教訓，不含專案特定 raw evidence>
promotion_candidate:
  target: <intelligence/ | workflow/ | enforcement/ | runtime/ | tools/>
  confidence: <high | medium | low>
  validation: <需要哪些驗證才能 promotion>
linked_updates:
  - <受影響的文件路徑>
```

## 與其他層的關係

- `feedback/extraction/`：Replay 產出的 generalized lesson 可進入 extraction 流程，進一步提煉為 intelligence atom。
- `feedback/refinement/`：Replay 發現的 workflow 缺口可進入 refinement 流程。
- `feedback/promotion/`：Replay 結果的 promotion 路徑由 promotion pipeline 管理。
- `memory/summary/`：Session summary 是 replay 的重要素材來源。
- `memory/decision/`：Decision record 可輔助 replay 分析。
- `enforcement/failure-patterns/`：Replay 發現的新 failure pattern 可寫入此處。
- `intelligence/`：Replay 產出的工程智慧可 promotion 到 intelligence。
- `workflow/`：Replay 發現的流程缺口可 promotion 到 workflow。
