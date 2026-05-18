# Runtime Decision Recording

`runtime/decisions/` 定義 **何時、寫到哪、用什麼格式** 記錄鎖定的設計決策。Canonical 正文在 `decisions/`（ADR）與 `memory/decision/`（session）；本目錄提供 **machine-readable 路由** 供 close-loop 與 agent 查詢。

## 入口

- **[`decision-recording.yaml`](./decision-recording.yaml)** — tier 分類、路徑、close-loop 自問、registry 連結
- [`decisions/README.md`](../../decisions/README.md) — ADR 生命週期與錯誤查詢索引
- [`memory/decision/README.md`](../../memory/decision/README.md) — session 級決策

## 與 session lifecycle 的關係

[`../pipeline/session-lifecycle.yaml`](../pipeline/session-lifecycle.yaml) 的 **close-loop** 階段含 `decision-recording-checkpoint`，依本目錄 yaml 分流寫入 ADR 或 session/project 決策檔。

## Inbound References

- `route.runtime.decision-recording` — [`routing-registry.yaml`](../../knowledge/runtime/routing-registry.yaml)
- `route.decisions.adr` — 查詢歷史 ADR
- [`governance/lifecycle/knowledge-update-flow.md`](../../governance/lifecycle/knowledge-update-flow.md) Step 1 — 決策紀錄自問
