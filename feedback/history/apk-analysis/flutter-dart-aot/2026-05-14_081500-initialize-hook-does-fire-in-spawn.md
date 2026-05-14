> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-14 - Frida 實戰：初始化函數在 attach 模式不觸發不代表 offset 錯誤——改用 spawn 模式驗證

Status: validated

#### One-line Summary

當 Frida hook 在 attach 模式（`frida -U <package>`）下從未觸發時，**不要假設 offset 錯誤**——初始化函數（如 static initializer、singleton constructor）只在 App 啟動階段執行，必須使用 spawn 模式（`frida -U -f <package>`）才能捕獲。Spawn 模式確認 hook 會觸發後，再回頭檢查 attach 模式的限制。

#### Human Explanation

在分析某 App 的 `SkyShieldIntegration.initialize` 函數時，Frida hook 在 attach 模式下從未觸發。這導致了兩個錯誤假設：

1. **「offset 錯誤」**——以為 `0xe5b458` 不是正確的函數入口
2. **「該函數從未被呼叫」**——以為 App 不使用這個初始化路徑

實際上，`SkyShieldIntegration.initialize` 是一個 **static initializer**，只在 App 冷啟動時執行一次。Attach 模式是在 App 已經啟動後才附加 hook，因此錯過了初始化階段。

改用 spawn 模式（`frida -U -f <package> -l script.js`）後，hook 正確觸發——`onEnter` 和 `onLeave` 都被呼叫。進一步分析顯示，initialize 前後的 singleton 狀態沒有變化——這不是 hook 沒觸發，而是 initialize 的副作用在 heap 物件中，不在 singleton 的前 128 bytes 內。

**關鍵教訓**：初始化函數的 hook 必須用 spawn 模式。Attach 模式不觸發時，先試 spawn 模式再下結論。

#### Trigger

- `Interceptor.attach` 在 attach 模式下從未觸發
- 目標函數是初始化相關（static initializer、singleton constructor、`JNI_OnLoad`、library constructor）
- 文件或註解說「該函數 never fires」
- 嘗試多個 offset 都無效

#### Evidence

- Tool: Frida hook script comparison（attach vs spawn）
- Sanitized excerpt:
  - Attach mode: `hookInitialize` at offset `0xe5b458` — never fires
  - Spawn mode: `hookInitialize` at offset `0xe5b458` — `onEnter` and `onLeave` both fire
  - Singleton first 128 bytes are identical before and after initialize
- Evidence path: `<PROJECT_ROOT>/capture/frida_capture_iv_20260514.log`

#### Generalized Lesson

1. **初始化函數必須用 spawn 模式驗證**——如果目標函數在 App 啟動時執行（static initializer、singleton constructor、native library constructor），attach 模式永遠無法捕獲。
2. **「Hook 不觸發」不等於「offset 錯誤」**——先用 spawn 模式確認 offset 是否正確，再判斷是 offset 問題還是模式問題。
3. **「Hook 不觸發」也不等於「函數從未被呼叫」**——函數可能在 App 啟動階段執行完畢，attach 時已經結束。
4. **Spawn 模式的缺點**：會重新啟動 App，可能觸發登入、session recovery 或 rate limit。但這是驗證初始化 hook 的唯一可靠方式。
5. **混合策略**：先用 spawn 模式確認初始化 hook 正確，再用 attach 模式進行後續分析。

#### Agent Action

部署 Frida hook 時：

1. **先判斷目標函數的執行時機**——如果是初始化階段，使用 spawn 模式
2. **如果 attach 模式 hook 不觸發，先試 spawn 模式**——不要急著假設 offset 錯誤
3. **記錄測試過的模式**——在文件或註解中標明「tested with attach: no fire, tested with spawn: fires」
4. **如果 spawn 模式也不觸發**——再懷疑 offset 錯誤或函數未被編譯

#### Goal / Action / Validation

- Goal: 正確驗證初始化函數的 Frida hook
- Action: 先用 spawn 模式確認 hook 會觸發，再用 attach 模式
- Validation or reference source: 同一 hook 在 spawn 模式下觸發但在 attach 模式下不觸發

#### Applies When

- 需要 hook App 啟動階段的初始化函數
- Frida attach 模式的 hook 從未觸發
- 文件或註解說「該函數 never fires」但沒有說明測試模式
- 懷疑 offset 錯誤但沒有用 spawn 模式驗證過

#### Does Not Apply When

- 目標函數在使用者操作階段執行（attach 模式即可）
- 已經用 spawn 模式驗證過且確認不觸發
- 使用 Frida Gadget（嵌入模式，初始化時機不同）

#### Validation

- Spawn 模式成功觸發 `hookInitialize`（`onEnter` + `onLeave`）
- Attach 模式同一 hook 不觸發
- 確認函數在 App 啟動階段執行完畢

#### Promotion Target

- `intelligence/engineering/analytical-reasoning/heuristics/` — 更新現有 heuristic：「Frida spawn vs attach init timing」（與 `2026-05-14_073700-frida-spawn-vs-attach-init-timing-no-buffer.md` 合併）
- `workflow/apk-analysis/execution-flow.md` — 新增步驟：「初始化函數先用 spawn 模式驗證」

#### Required Linked Updates

- `feedback/history/apk-analysis/common/2026-05-14_073700-frida-spawn-vs-attach-init-timing-no-buffer.md` — 此 lesson 與該篇高度重疊，考慮合併或交叉引用
