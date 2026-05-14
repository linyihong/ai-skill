> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`workflow/software-delivery/execution-flow.md`](../../../../workflow/software-delivery/execution-flow.md)

### 2026-05-07 - Performance Test Release Gate

Status: promoted

#### One-line Summary

功能正確不代表效能可上線；會影響容量、延遲或資源的變更需要效能預算與對應測試。

#### Human Explanation

AI 產生或快速修改的程式常能通過功能測試，卻在資料量、併發、批次、快取、資料庫、外部 API 呼叫或長時間運行下暴露效能問題。單元測試通常看不出 P95/P99 延遲、吞吐量、錯誤率、資源使用與長時間退化。開發流程要把效能測試視為 release gate，而不是出事後才補的檢查。

#### Trigger

使用者提供一張效能測試知識圖，重點包含 AI code 常見效能陷阱、load / stress / spike / soak 四類測試、P95/P99、TPS/RPS、錯誤率、資源使用率，以及把小型壓測放進 CI/CD。

#### Evidence

- Tool: User-provided visual reference.
- Sanitized excerpt: Performance testing should cover normal load, saturation, sudden spikes, and long-running stability; averages alone are insufficient.
- Evidence path: Keep the original image in the conversation or project evidence if needed; reusable skill only stores the generalized rule.

#### Generalized Lesson

任何可能影響延遲、吞吐量、資源使用、啟動、背景作業、資料庫、批次、快取、併發、重試或外部呼叫量的變更，都要在開發前定義效能預算與測試類型。最小可行測試可以是 CI smoke check，但 release 判斷必須能說明指標、環境、資料量、baseline 與結果。

#### Agent Action

下次規劃 app/API/SDK/backend/tooling 變更時：

1. 在 change intake 問「這個變更會不會影響效能或資源」。
2. 若會，要求 P95/P99 延遲、吞吐量、錯誤率與資源預算。
3. 依風險選擇 load、stress、spike、soak 或 CI smoke performance check。
4. 在 template / checklist / release note 中記錄 runner、環境、資料量、baseline、結果與 owner。
5. 不用平均值 alone 宣稱效能足夠；平均值只能作為輔助背景。

#### Goal / Action / Validation

- Goal: 避免功能正確但效能不穩的變更被誤判為可上線。
- Action: 將效能預算與 load / stress / spike / soak 測試納入 planning、test strategy、checklist 與 template。
- Validation or reference source: 每個 performance-sensitive change 有明確 test type、P95/P99、throughput、error rate、resource usage、baseline 或 release-gate evidence；若不適用，必須寫明原因。

#### Applies When

- 變更可能影響 user-visible latency、API/SDK throughput、queue/job volume、database access、caching、batch processing、startup/background work、external API fan-out、memory/CPU/disk/network 使用量。
- AI 生成或大幅改寫程式碼，且資料量或併發風險不明。
- CI/CD、pre-release、nightly 或 on-demand performance gates 需要定義。

#### Does Not Apply When

- 純文件或純 UI copy 變更，且明確不影響 runtime path。
- 小型內部重命名已由 existing regression tests 證明無行為或資源路徑變化。

#### Validation

- `WORKFLOW.md` 包含 performance test gate。
- `process/README.md`、`CHECKLIST.md`、`templates/initial-development-docs.md` 能引導 agent 產出效能預算與測試證據。
- `SKILL.md` 和 `README.md` 的入口與 linked updates 能讓下一個 agent 找到這條規則。

#### Promotion Target

- `WORKFLOW.md`
- `process/README.md`
- `CHECKLIST.md`
- `templates/initial-development-docs.md`
- `SKILL.md`

#### Required Linked Updates

- 已同步更新 `WORKFLOW.md`、`process/README.md`、`CHECKLIST.md`、`templates/initial-development-docs.md`、`templates/README.md`、`DOCUMENTATION.md`、`SKILL.md` 與 `README.md`。
- 已依 `reusable-guidance-boundary.md` 檢查：skill 只保留 generalized lesson，原始圖片作為使用者提供的參考來源，不複製私有專案證據。
