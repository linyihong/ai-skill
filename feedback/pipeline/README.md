# Feedback Promotion Pipeline

`feedback/pipeline/` 將 feedback lesson 從 `skills/*/feedback_history/` 的原始觀察，透過機器可讀的 scoring、workflow 與 lifecycle automation，推進到 `workflow/`、`intelligence/`、`shared-rules/`、`memory/` 或 runtime surfaces。

## 為什麼需要 Pipeline

目前 feedback_history 已有大量 lesson（僅 `apk-analysis` 就有 **86 條**，跨 6 個分類），但：

1. **無自動 promotion** — lesson 寫完就留在 feedback_history，不會自動推進到 workflow/intelligence/shared-rules。
2. **無優先級** — 86 條 lesson 中，哪些該先 promote、哪些該 archive，沒有機器可讀的判斷。
3. **無生命週期** — lesson 不會自動過期、降級或歸檔，冷資料持續累積。
4. **無 scoring** — 無法量化 lesson 的品質、適用性與 promotion 急迫性。

Pipeline 解決這些問題：

```
feedback_history/  ──→  Promotion Engine  ──→  Promotion Workflow  ──→  Target Layer
     (86 lessons)         (score + decide)        (executable steps)       (workflow/
                                                                           intelligence/
                                                                           shared-rules/
                                                                           memory/)
                               │
                               ▼
                        Lifecycle Automation
                        (auto-archive cold,
                         auto-downgrade stale)
```

## Pipeline 架構

```text
feedback/pipeline/
  README.md                  ← 本檔：pipeline 概覽
  promotion-engine.yaml      ← Promotion scoring & decision rules
  promotion-workflow.yaml    ← Executable promotion workflow steps
  lifecycle-automation.yaml  ← Auto-archive, auto-downgrade rules
```

## 與既有層的關係

| Pipeline 元件 | 依賴的既有層 | 關係 |
|-------------|------------|------|
| promotion-engine | `skills/*/feedback_history/`, `shared-rules/feedback-lessons.md` | 讀取 lesson 內容做 scoring |
| promotion-workflow | `governance/lifecycle/README.md`, `governance/validation/README.md` | 遵循 lifecycle states 與 validation gates |
| lifecycle-automation | `governance/lifecycle/README.md`, `knowledge/runtime/sqlite/` | 使用 SQLite index 做 cold data lookup |

## 使用方式

Pipeline 不是自動化服務，而是 Agent 在每個 close-loop 階段遵循的**執行模型**：

1. **Session 結束時**：Agent 檢查是否有新 lesson 需要評估 promotion。
2. **Promotion 評估**：Agent 使用 `promotion-engine.yaml` 的 scoring 邏輯決定 lesson 的 promotion 優先級。
3. **Promotion 執行**：Agent 遵循 `promotion-workflow.yaml` 的步驟執行 promotion。
4. **Lifecycle 維護**：Agent 使用 `lifecycle-automation.yaml` 的規則自動歸檔冷 lesson、降級過期 lesson。

---

← [回到 Feedback](../README.md)
