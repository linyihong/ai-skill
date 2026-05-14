> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-14 - 逆向加密 token 時：前綴在 session 內可能不是完全固定——比對 call #1 與後續 calls

Status: validated

#### One-line Summary

逆向加密 token 時，即使前綴在大部分 calls 中看起來固定，**第一個 call 的前綴可能與後續 calls 有微小差異**（如 byte 14 偏移 1）。不要只比對中間的 calls，一定要比對 **call #1** 與 **calls #25+**。這種差異暗示前綴產生機制有一個「初始化→穩定」的過程。

#### Human Explanation

在分析某 App 的 `x-aspnet-version` token 時，前 24 個 calls 的 prefix 看起來完全一致。但仔細比對後發現：

- **Call #1**: byte 14 = `0xe6`
- **Calls #25+**: byte 14 = `0xe5`

只有 1 byte 的差異（`0xe6` → `0xe5`），但這是一個重要的信號：

1. 前綴產生機制可能有一個 **warm-up 或初始化階段**
2. 第一個 call 可能使用不同的隨機種子或狀態
3. 前綴在大部分 calls 中穩定，但第一個 call 是例外

**關鍵教訓**：不要只取中間的樣本就假設「前綴完全固定」。比對第一個和最後一個 call，找出任何微小差異。這種差異可能揭示前綴的產生機制。

#### Trigger

- 前綴在大部分 calls 中看起來固定
- 只比對了中間的 calls（如 calls #5-#20）
- 沒有比對 call #1 和 calls #25+
- 差異只有 1-2 bytes，容易被忽略

#### Evidence

- Tool: Frida capture with sequential event numbering
- Sanitized excerpt:
  - Call #1: prefix byte 14 = `0xe6`
  - Calls #2-#24: prefix byte 14 = `0xe6`（與 call #1 相同）
  - Calls #25+: prefix byte 14 = `0xe5`（偏移 1）
  - All other bytes (0-13, 15): identical across all calls
- Evidence path: `<PROJECT_ROOT>/capture/frida_capture_iv_20260514.log`

#### Generalized Lesson

1. **比對第一個和最後一個 call**——不要只比對中間的樣本。第一個 call 可能使用不同的初始化狀態。
2. **逐 byte 比對，不只是整體比對**——1 byte 的差異在整體 Base64 比較中可能不明顯，但逐 byte hex 比對可以清楚看到。
3. **記錄 call 的序號**——在 Frida hook 中加上 event counter，方便追蹤哪個 call 產生哪個 token。
4. **微小差異可能是重要線索**——byte 偏移 1 可能表示 counter、timestamp 或 random seed 的初始化過程。
5. **對 self-generation 的影響**：如果第一個 call 的 prefix 不同，self-generation 需要能夠重現這種行為。

#### Agent Action

分析加密 token 的前綴時：

1. **在 Frida hook 中加入 event counter**——每個 call 都有一個序號
2. **比對 call #1 和 calls #25+**——不要只比對中間的 calls
3. **逐 byte hex 比對**——不要只比對 Base64 字串
4. **記錄任何差異**——即使只有 1 byte 的差異也要記錄
5. **考慮前綴產生的初始化過程**——差異可能來自 random seed、counter 或 timestamp 的初始化

#### Goal / Action / Validation

- Goal: 完整了解加密 token 前綴在 session 內的行為
- Action: 比對 call #1 與 calls #25+ 的逐 byte hex 值
- Validation or reference source: Call #1 byte 14 = `0xe6`, calls #25+ byte 14 = `0xe5`

#### Applies When

- 分析加密 token 的固定前綴
- 有多個 calls 的樣本（> 25 calls）
- 前綴在大部分 calls 中看起來固定
- 需要實作 self-generation（需要了解前綴的完整行為）

#### Does Not Apply When

- 只有少量 calls（< 5 calls）
- 前綴在所有 calls 中完全一致（包括 call #1）
- 分析的是靜態資料

#### Validation

- Call #1 和 calls #25+ 的 prefix 在 byte 14 不同（`0xe6` vs `0xe5`）
- 所有其他 bytes 在所有 calls 中一致
- Event counter 確認 call 序號正確

#### Promotion Target

- `intelligence/engineering/analytical-reasoning/heuristics/` — 更新 heuristic：「加密 token 前綴分析：比對 call #1 與後續 calls」
- `workflow/apk-analysis/execution-flow.md` — 新增步驟：「分析前綴時比對第一個和最後一個 call」

#### Required Linked Updates

- `feedback/history/apk-analysis/flutter-dart-aot/2026-05-14_081500-prefix-session-specific-not-device-specific.md` — 交叉引用此 lesson
