# Architecture Escalation Policy

## 何時升級

Architecture issue 只有在影響交付正確性或長期決策時才升級，不因純粹偏好阻塞 delivery。

## 升級訊號

- Architecture recommendation 與 business complexity 明顯不匹配。
- Low complexity project 被升級為 Full DDD / CQRS / event sourcing。
- High invariant domain 被壓成 CRUD-only，導致 invariant 無保護。
- Bounded context mismatch 造成 repeated recovery loop 或 contract mismatch。
- Team / deployment / integration boundary 被忽略，導致 delivery plan 不可信。

## 行動

1. 暫停 architecture recommendation。
2. 補 architecture fit analysis。
3. 明確列出 mismatch evidence。
4. 提出 simplification 或 strengthening option。
5. 若仍不確定，詢問使用者選擇 business priority。
