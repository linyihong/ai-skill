# Session-Scoped Implicit State

**Status**: `candidate-intelligence`
**Source**: 經典 Web Monolith 模式（ASP.NET Session / Rails session / Flask g）於 agent runtime、distributed system、replayable execution 場景的災難化。

## 反模式

把「當下使用者」「目前 phase」「current tenant」這類執行上下文塞進 session-scoped service（或 thread-local、ambient context），讓下游函式不需要顯式接收參數，直接從注入的 service 讀。

對單一進程的 stateful Web App 還能跑；一旦執行單元可能被 **fork、replay、跨 agent 並發、或從 SQLite 重建狀態**，就會變成不可觀察、不可測試、不可重放的隱形依賴。

## 訊號

- Service / controller 簽章看不到 `userId`、`phaseId`、`tenantId`，但函式行為依賴它們。
- 測試需要 mock session middleware 才能跑單一函式。
- 同樣的 input 在不同 caller 下產生不同 output，差異藏在 ambient context。
- 嘗試把邏輯抽出來給排程器 / background worker / agent 用時，需要「偽造一個 session」。
- Bug 報告長這樣：「我重跑同樣的步驟，得到不同的結果。」

## 為什麼對 agent runtime 特別致命

Agent 執行天然就是 replayable、fork-able、可被多個 caller 觸發的。Session-scoped state 在這個場景會把三件事一起破壞：

1. **Replay 不確定**：cache 重放時 ambient context 已經消失，函式拿到 stale 或 null。
2. **並發 corruption**：兩個 agent 共用同一個 ambient slot，互相覆蓋。
3. **觀察性破洞**：trace 只看得到函式呼叫，看不到「當時的隱式狀態是什麼」，事後無法歸因。

## 判斷邊界

| 情境 | Session-scoped 可接受 | 必須 explicit context |
|------|---------------------|----------------------|
| 單一進程、單一使用者請求生命週期 | ✅ | — |
| Background job / scheduler | ❌ | ✅ |
| Multi-agent runtime | ❌ | ✅ |
| Event replay / time-travel debug | ❌ | ✅ |
| 函式可能被測試獨立呼叫 | ❌ | ✅ |

判斷句：**state 越隱式，replay / fork / 並發越痛苦。** 一旦執行單元的觸發方式 > 1 種，session-scoped 就該升級成 explicit `Context` 參數。

## 修正

- 把 ambient context 物化成 `ExecutionContext` / `RunContext` 物件，沿呼叫鏈顯式傳遞。
- SQLite-as-canonical 的 runtime（如本專案 `runtime/runtime.db`）尤其要避免：phase / obligation 應從 context 物件讀，而不是由 service 「偷偷查 DB 當前 phase」。
- 若無法一次重構，至少在進入函式時 **snapshot 一次** ambient state 成 local immutable，後續邏輯只讀 snapshot。

## 常見誤用

| 誤用 | 正確 |
|------|------|
| `ISessionService.CurrentUser` 散佈在 domain layer | Domain 函式接收 `User` 參數 |
| Agent 從 `runtime.db` 直接查 `current_phase` | Agent 接收 `RunContext { phase_id, ... }` |
| 用 `AsyncLocal<T>` / thread-local 在 async pipeline 攜帶狀態 | 顯式 context 傳遞，或 structured concurrency |

## Token Impact

避免 replay debug 時無法重現的 ghost bug。一個被 ambient state 污染的 agent runtime，事後歸因成本通常是 explicit context 重構成本的 5–10 倍。

---

← [回到 engineering/anti-patterns/](README.md)
