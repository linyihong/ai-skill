# 去敏與占位符

可重用文件（含 `feedback_history` lesson）**不要**包含：

- 完整 `Authorization` / token、session cookie、可識別特定使用者的實體裝置識別。
- AES/HMAC/簽名密鑰（**除非**教學用、合成、可公開測試向量）。
- 未去敏的 raw request/response。
- 本機**真實**絕對路徑、使用者帳號名稱、私用工作目錄、git clone 實體路徑。
- Project incident 的具體 app / project 名稱、endpoint、host、payload fragment、sample ID、class/test 名稱、live run 結果或環境 quirks；這些依 [`reusable-guidance-boundary.md`](reusable-guidance-boundary.md) 留在專案文件。

一律改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 等占位符（見本庫 [README.md](../README.md)）。

← [回到共用規則索引](README.md)
