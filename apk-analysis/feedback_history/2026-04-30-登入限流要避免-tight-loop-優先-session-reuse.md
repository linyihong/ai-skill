### 2026-04-30 - 登入限流要避免 tight-loop，優先 session reuse

Status: validated

#### One-line Summary

遇到 login too frequently，不要盲目旋轉單一參數；先重用 session 並記錄風控維度。

#### Human Explanation

登入限流可能不是單一 request 欄位造成，而是伺服器用 device、User-Agent、IP、時間窗、App fingerprint、帳號狀態等多維度計算。一直換 device id 或 tight-loop login 可能讓問題更嚴重。穩定做法是同一輪測試重用 session，記錄每次登入嘗試時間與參數，必要時用 device/session pool。

#### Trigger

- API 回 login too frequently。
- 多個測試每次都重新登入。
- 改某個 device/body 欄位後結果不穩定。

#### Evidence

授權測試中曾做過參數 probe，無法把限流穩定歸因於單一欄位；session reuse 明顯降低重複登入風險。

#### Generalized Lesson

登入流程測試要有節流與重用策略。不要為每個測試方法重新登入；不要在沒有證據時假設旋轉單一 device 欄位即可繞過限制。

#### Agent Action

設計 live integration 或 runner 時，優先共用 session/context，記錄 login attempt metadata。遇到限流時先停止 tight-loop，再分析時間窗與風控維度。

#### Promotion Target

- `WORKFLOW.md`
- `DOCUMENTATION.md`
