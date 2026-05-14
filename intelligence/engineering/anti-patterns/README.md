# Engineering Anti-Patterns

放**常見錯誤設計**。

## 核心

AI 自動避雷。

## 範例內容

- `generic-repository-overuse.md` — Generic repositories often hide query intent and reduce performance visibility.
- `shared-database-microservices.md` — Shared DB across services creates hidden coupling.
- `god-service-pattern.md` — Large service classes signal missing domain boundaries.

## 與根目錄 `anti-patterns/` 的差異

| 位置 | 內容 |
|------|------|
| [`anti-patterns/`](../../anti-patterns/) | Agent 操作層的 anti-patterns（context explosion、recursive tool loop、hallucination loop 等） |
| `intelligence/engineering/anti-patterns/` | 工程領域的 anti-patterns（設計錯誤、架構錯誤等） |

## 與其他層的關係

- Agent 操作層的 anti-patterns → [`anti-patterns/`](../../anti-patterns/)
- 具體的 failure pattern policy → `enforcement/failure-patterns/`
