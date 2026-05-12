# Engineering Failure Intelligence

放**工程災難智慧**。不是 incident log，而是抽象化後的失敗模式。

## 核心

AI 的「危險雷達」。

## 範例內容

- `connection-leak-patterns.md` — Symptoms: latency spike, pool exhaustion, CPU idle but requests blocked. Common Causes: per-request connection creation.
- `distributed-lock-failure.md` — Distributed locks become dangerous when lock expiration assumptions are unstable.

## 與其他層的關係

- 具體 incident 記錄 → `feedback/` 或 project docs
- 可執行的 failure prevention policy → `shared-rules/failure-patterns/`
