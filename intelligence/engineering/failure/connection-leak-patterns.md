# Connection Leak Patterns（連線洩漏模式）

**Status**: `candidate-intelligence`
**Source**: 通用後端系統營運經驗

## 原則

**Symptoms: latency spike, pool exhaustion, CPU idle but requests blocked. Common Causes: per-request connection creation, missing close in error paths, thread-local connection accumulation.**

症狀：延遲飆升、連線池耗盡、CPU 空閒但請求阻塞。常見原因：每個請求建立新連線、錯誤路徑未關閉連線、thread-local 連線累積。

## 為什麼

1. **連線池是有限資源**：資料庫連線池、HTTP connection pool、gRPC channel 都有上限。一旦耗盡，新請求必須排隊等待，導致延遲飆升。
2. **錯誤路徑容易被忽略**：正常路徑會 close connection，但 exception / error 路徑經常忘記 close，導致連線洩漏。
3. **Thread-local 連線是隱形殺手**：Thread-local 連線在 thread pool 環境中不會自動釋放，累積到 pool 耗盡。
4. **Connection leak 的發現通常很晚**：因為 pool 耗盡需要時間累積，通常在 deploy 後幾小時到幾天才會被發現。

## 何時懷疑 Connection Leak

- **Deploy 後延遲逐漸升高**：剛 deploy 時正常，幾小時後延遲開始上升。
- **連線池監控顯示 active connections 持續增長**：從不下降。
- **CPU 使用率正常或偏低，但請求 timeout**：表示瓶頸不在 CPU，而在 I/O 或連線等待。
- **重啟後恢復正常**：重啟釋放所有連線，但問題會再次出現。

## 何時不懷疑 Connection Leak

- **延遲飆升伴隨 CPU 100%**：可能是計算瓶頸，不是連線問題。
- **連線數穩定在 pool 上限**：表示 pool size 設定過小，不是 leak。
- **單一 client 的連線數異常高**：可能是 client 端的 connection pool 設定問題。

## 診斷流程

```text
延遲飆升？
  ├── 檢查連線池監控
  │     ├── active connections 持續增長？
  │     │     ├── 是 → 可能是 connection leak
  │     │     └── 否 → 檢查其他指標
  │     └── pool 耗盡？
  │           ├── 是 → 緊急處理：重啟服務 + 縮小 pool size 暫時緩解
  │           └── 否 → 繼續監控
  ├── 檢查 error log
  │     ├── 有「connection closed」或「timeout」錯誤？
  │     │     ├── 是 → 檢查對應的程式碼路徑
  │     │     └── 否 → 繼續
  │     └── 有 exception 但沒有對應的 close/finally？
  │           ├── 是 → 修正：加入 try/finally 或 using 區塊
  │           └── 否 → 繼續
  └── 檢查 thread-local 連線使用
        ├── 有 ThreadLocal<Connection>？
        │     ├── 是 → 確認是否在 request 結束時清理
        │     └── 否 → 繼續
        └── 使用 thread pool？
              ├── 是 → thread-local 連線不會自動釋放，需手動清理
              └── 否 → 較少見
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「連線池 size 設大一點就好」 | 連線池 size 過大反而增加資料庫負載。真正的問題是 leak，不是 pool size |
| 「重啟就好了」 | 重啟只是暫時緩解。如果 root cause 沒修，問題會再次出現 |
| 「一定是資料庫問題」 | Connection leak 通常是 client 端的問題，不是資料庫 |

## Token Impact

Connection leak 可能導致服務中斷 30 分鐘到數小時（從 leak 開始到 pool 耗盡）。及早發現可以避免 P0 級 incident。每個 leak 的修復成本約 1-3 小時（診斷 + 修復 + 驗證）。

---

← [回到 engineering/failure/](README.md)
