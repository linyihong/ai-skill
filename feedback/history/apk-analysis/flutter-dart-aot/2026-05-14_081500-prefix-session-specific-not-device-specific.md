> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-14 - 逆向加密 token 時：固定前綴的變異邊界要跨 session 比對，不要假設是 device-specific

Status: validated

#### One-line Summary

逆向加密 token 時，若所有 ciphertext 共享一個固定前綴，**不要假設它是 device-specific（不變的裝置指紋）**——它可能是 **session-specific**（每次 App 啟動重新產生）。跨 session 比對是唯一可靠的判斷方式。在同一個 session 內看起來「固定」的前綴，在不同 session 間可能在特定 bytes 位置變化。

#### Human Explanation

在分析某 App 的 `x-aspnet-version` token 時，所有 ciphertext（32 decoded bytes）共享一個 **15-byte 固定前綴**。在同一個 session 內捕獲的數十個 token 中，前綴完全一致。這導致了錯誤假設：「前綴是 device-specific 的，永遠不變」。

但跨 session 比對（不同時間啟動 App）顯示：
- Session A 前綴：`s0jOc/2OJwu/l5ecDFfl`（bytes 13-14 = `0x97 0xe6`）
- Session B 前綴：`s0jOc/2OJwu/l5e1FTfl`（bytes 13-14 = `0xd7 0x53`）

前綴的 **bytes 13-14 在不同 session 間變化**，表示前綴不是 device-specific 的固定值，而是每次 App 啟動時重新產生的 session-specific 值。

**關鍵教訓**：跨 session 比對是判斷前綴是否固定的唯一可靠方式。同一個 session 內的「固定」可能只是假象。

#### Trigger

- 所有 ciphertext 共享一個固定前綴
- 在同一個 session 內前綴完全一致
- 文件或分析假設「前綴是 device-specific」
- 沒有做過跨 session 比對

#### Evidence

- Tool: Frida capture + manual cross-session comparison
- Sanitized excerpt:
  - Session A (2026-05-14 capture): prefix = `s0jOc/2OJwu/l5ecDFfl` (bytes 13-14 = `0x97 0xe6`)
  - Session B (2026-05-14 capture, different launch): prefix = `s0jOc/2OJwu/l5e1FTfl` (bytes 13-14 = `0xd7 0x53`)
  - Within each session: prefix is stable across all calls
- Evidence path: `<PROJECT_ROOT>/capture/frida_capture_iv_20260514.log`

#### Generalized Lesson

1. **跨 session 比對是必要步驟**——不要只依賴同一個 session 內的觀察。至少比對 2-3 個不同 session 的資料。
2. **「固定」在 session 內 ≠ 跨 session 固定**——前綴可能在 session 內穩定，但在不同啟動之間變化。
3. **前綴的變異邊界可能很小**——不是整個前綴都變，可能只有特定 bytes（如 bytes 13-14）變化。需要逐 byte 比對。
4. **Session-specific 前綴的來源**——可能是 random seed、timestamp-based nonce、或 session key 的一部分。
5. **對 self-generation 的影響**：如果前綴是 session-specific，self-generation 需要能夠重現前綴的產生邏輯，不能 hardcode。

#### Agent Action

分析加密 token 的固定前綴時：

1. **在同一個 session 內收集多個樣本**——確認 session 內是否真的固定
2. **跨 session 比對**——至少比對 2-3 個不同啟動的 session
3. **逐 byte 比對**——找出哪些 bytes 固定、哪些變化
4. **記錄變化模式**——bytes 13-14 變化？bytes 0-12 固定？這暗示了前綴的結構
5. **不要假設 device-specific**——除非有跨 session 證據支持

#### Goal / Action / Validation

- Goal: 正確判斷加密 token 前綴的變異邊界
- Action: 跨 session 比對前綴，逐 byte 分析變化模式
- Validation or reference source: 跨 session 比對顯示 bytes 13-14 變化，bytes 0-12 固定

#### Applies When

- 分析加密 token 或 ciphertext 的固定前綴
- 所有樣本來自同一個 session
- 文件假設前綴是 device-specific 但沒有跨 session 證據
- 需要實作 self-generation（前綴的變異邊界影響實作策略）

#### Does Not Apply When

- 已經有跨 session 證據確認前綴完全固定
- 前綴長度為 0（沒有固定前綴）
- 分析的是靜態資料（非動態產生的 token）

#### Validation

- Session A 和 Session B 的前綴在 bytes 13-14 不同
- 同一個 session 內所有 token 的前綴完全一致
- 前綴的 bytes 0-12 跨 session 固定

#### Promotion Target

- `intelligence/engineering/analytical-reasoning/heuristics/` — 新增 heuristic：「加密 token 前綴的跨 session 比對」
- `workflow/apk-analysis/execution-flow.md` — 新增步驟：「分析固定前綴時必須跨 session 比對」

#### Required Linked Updates

- 無需連動更新；這是新 lesson。
