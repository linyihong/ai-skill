# /bootstrap

執行完整 bootstrap 序列。在新 session 開始時手動觸發，確保所有 obligations 載入完成。

## 執行步驟

請依序完成以下步驟，不得跳過：

1. 讀 `CORE_BOOTSTRAP.md`
2. 讀 `README.md`
3. 查詢 `runtime/runtime.db`：執行 `sqlite3 runtime/runtime.db "SELECT phase_id FROM phase_machine LIMIT 1;"` 取得 phase，並查詢 obligations 與 gates 數量
4. 讀 `enforcement/rule-weight.md`
5. 讀 `enforcement/dependency-reading.md`
6. 讀 `enforcement/conversation-goal-ledger.md`
7. 輸出 Bootstrap Receipt：
   ```
   Bootstrap: rules=✓ phase=<phase-id> obligations=<n> gates=<n>
   Active per-turn obligations: <obligation ids>
   ```
8. 輸出 Cognitive Mode 報告（compact 或 full table）

完成後回報「Bootstrap 完成」及 Receipt 內容。
