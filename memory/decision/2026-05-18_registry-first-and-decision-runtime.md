# Session-level Decision: Registry-First + Runtime Decision Recording

## Status

accepted

## Context

2026-05-18 工作 session 完成 workflow activation 重構與 SDK catalog plan 決策鎖定。使用者指出 **decisions 幾乎沒有紀錄**，事後查問題困難。

## Decision

1. **架構級**已升級為 [ADR-006](../../decisions/ADR-006-registry-first-workflow-activation.md)（registry-first workflow activation）。
2. **Runtime** 新增 [`runtime/decisions/decision-recording.yaml`](../../runtime/decisions/decision-recording.yaml)，close-loop 與 knowledge-update Step 1 共用 tier 路由。
3. **專案級** SDK catalog／cache／分頁決策寫入 `<PROJECT_ROOT>/apk-analysis-sdk/docs/decisions/2026-05-18-sdk-catalog-cache-and-pagination.md`。

## Consequences

- 新 session agent 應先查 `decisions/README.md` §錯誤查詢索引 或 `decision-recording.yaml` lookup。
- 僅寫在 plan 內文、未進 decisions 的鎖定項視為 **未閉環**。

## Related

- [ADR-006](../../decisions/ADR-006-registry-first-workflow-activation.md)
- [`runtime/decisions/README.md`](../../runtime/decisions/README.md)
- [`feedback/.../2026-05-18_160000-registry-first-workflow-activation.md`](../../feedback/history/development-guidance/common/2026-05-18_160000-registry-first-workflow-activation.md)
