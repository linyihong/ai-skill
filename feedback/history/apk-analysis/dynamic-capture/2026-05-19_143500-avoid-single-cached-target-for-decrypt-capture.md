> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-19 - Avoid Single Cached Target For Decrypt Capture

Status: candidate

#### One-line Summary

解密或影像 hook 被單一目標 cache 卡住時，改抓同一資源族群的其他新樣本，先建立乾淨 input/output 對照，再回頭驗證原目標。

#### Human Explanation

動態分析常會為了追蹤一個具名資源而鎖定單一 URL、hash 或 UI card。若該資源已經進入 App 的 memory cache、disk cache 或 Flutter image cache，後續操作可能只重用已解密結果，不再觸發 HTTP、decryptor 或 image loader。此時繼續強迫同一目標重跑，容易導致錯誤 UI 操作、加重 hook 或誤判 decrypt path。

較穩定的做法是把目標從「某一個已快取樣本」提升成「同一格式、同一 loader、同一 path family 的任意新樣本」。單一樣本仍然很有價值，但它的角色應是定位點：用來找到正確 UI route、loader branch、path family、hook offset 或 cache key 形狀；真正建立 decrypt fixture 時，不必拘泥於定位點本身。先抓到多組新的 encrypted body、decrypt input/output、memory payload 或 final image，再用離線腳本還原共通格式；原本的具名樣本只保留作最後 parity 驗證。

#### Trigger

- 指定資源已在 UI 可見，但 `loadImage`、`getCachedImageData`、decryptor 或 HTTP hook 沒有新事件。
- 小幅滑動或切 tab 仍只命中 memory/disk cache。
- Agent 為了追單一目標開始加入不可靠座標操作、重啟、清 cache 或 heavy hook。
- 同一頁面有多個同格式資源可供取樣，例如多張封面、avatar、富文本圖片或 media item。

#### Evidence

- Tool: Frida attach/spawn、image loader hook、decryptor hook、HTTP body download、offline byte/magic probe。
- Sanitized excerpt: 把 hook 條件從單一 path hash 改成同一資源 path family 後，短窗口內取得多組 `decrypt` / `memory-cache` fixture；其中 final image payload 可由 wrapper 後方 magic bytes 定位。
- Evidence path: `<PROJECT_ROOT>/capture/` 中的 private fixture 與 log；lesson 不包含 package name、host、auth query、sample id、raw bytes 或本機絕對路徑。

#### Generalized Lesson

1. **先分清目標層級**：若任務是理解 decrypt algorithm，樣本不必一開始就是業務指定那一筆；只要同 path family、同 loader、同 decrypt branch，就能先解共通格式。
2. **避免單點 cache 迷思**：單一目標越常被操作，越可能被 App cache 吃掉；繼續執著會降低 evidence-to-cost ratio。
3. **定位點不等於證據來源**：指定樣本可用來定位正確邊界，但 fixture 可以來自同邊界下更乾淨、尚未快取、較容易重放的樣本。
4. **改用族群式 arm 條件**：把 hook 從 `targetHash == X` 放寬到 `path contains /resource-family/`，檔名再以 path hash 區分樣本。
5. **先拿對照，再回原目標**：離線復現應先用最乾淨的新樣本完成，再用原業務目標做 parity check。
6. **UI 由人控或 checkpoint 控**：若 automation 已經跑偏，停止自動 tap；只 attach hook，讓人工或已驗證 checkpoint 操作目標列表。

#### Agent Action

1. 當單一目標 hook 沒事件時，先問：這是「指定樣本問題」還是「同族格式問題」？
2. 把指定樣本當定位點，先用它確認 UI route、loader branch、path family、hook offset 或 cache key 形狀。
3. 若目標是格式或 decrypt path，將 hook filter 放寬到同一資源族群。
4. 用 manifest 記錄 path hash、event type、byte length、magic offset 與檔名。
5. 過濾樣本：挑 final payload 有 image magic，且同 hash 有 encrypted/decrypt/cache 對照的組合。
6. 離線復現成功後，再回到原指定目標驗證是否同算法。

#### Goal / Action / Validation

- Goal: 避免已快取的單一樣本阻塞 decrypt 或 image capture 分析。
- Action: 先用單一目標定位正確邊界，再從同族資源 sampling 建立多組 sanitized fixture。
- Validation or reference source: 至少一組新樣本有 encrypted input、intermediate/decrypt evidence、final image payload；原指定樣本作為最後 parity，而非唯一 capture 來源。

#### Applies When

- 資源可由多個同格式樣本代表，且任一樣本都能說明共通 decrypt/container format。
- App 有明顯 memory/disk/image cache，單一目標已經被多次載入。
- 需要先理解 algorithm，再做 SDK/client implementation。

#### Does Not Apply When

- 任務要求驗證單一資源的業務語義、授權狀態或內容完整性。
- 不同樣本可能走不同 decrypt key、不同 path family 或不同 loader branch。
- 授權範圍只允許分析指定樣本，不允許擴大到同頁其他資源。

#### Validation

- Hook log 顯示同族 filter 命中多個 path hash。
- Fixture manifest 可把每個 hash 對應到 event 類型與 byte evidence。
- Offline probe 能從 final payload 定位標準 magic bytes，且不需要依賴單一已快取目標。

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md`（capture window / sample targeting）
- `analysis/apk/workflows/frida-hook-flow.md`（image/decrypt hook sampling）

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 更新 `dynamic-capture` 數量。
- `feedback/extraction/apk-analysis-index.md` 加入本 lesson 的 candidate row。
- 已檢查 `enforcement/failure-patterns/skill-local-feedback-bypass.md` 與 `correction-loop-bypass.md`；本輪的 agent failure 已有可用 pattern，不新增重複 failure pattern。
