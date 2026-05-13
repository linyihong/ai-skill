> 遵守 [共用規則索引](../../../shared-rules/README.md)、[dependency-reading](../../../shared-rules/dependency-reading.md)、[neutral-language](../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-13 - Incremental Frida Hook Versioning for Complex Capture

Status: candidate

#### One-line Summary

當需要捕獲多層網路架構的 App 流量時，使用「逐版增量」的 Frida hook 策略（v1 → v11），每版只新增一個 hook 層級，並在每版執行後撰寫分析文件，可以系統性地排除無效 hook 點、收斂到正確的捕獲策略。

#### Human Explanation

在逆向分析多層網路架構的 App 時，流量可能經過多個層級（Dart → OkHttp → Netty proxy → TLS），每個層級都可能使用不同的實作。如果一次加入所有 hook，會難以判斷哪個 hook 真正有效、哪個從未觸發。

**實務案例**：在分析某 App 的 guest login 流程時，使用了從 v4 到 v11 的增量 hook 策略：

| 版本 | Hook 層級 | 結果 |
|------|-----------|------|
| v4 | OkHttp, URL.openConnection | Dart 繞過 Java HTTP stack |
| v5 | java.net.Socket | Dart 繞過 Java socket |
| v6 | Netty channelRead（proxy 端） | 看到 proxy 流量，非 TLS |
| v7 | SSLSocketFactory.createSocket | 從未觸發（0 事件） |
| v8 | libc connect/send/sendto/write | 看到 proxy 明文 + TLS 加密 |
| v9 | Netty SslHandler | 從未觸發（非 Netty TLS） |
| v10 | Native SSL_write/SSL_read | 從未觸發（非標準 TLS） |
| v11 | Netty channelRead（body 捕獲） | 成功捕獲完整 HTTP body |

**關鍵教訓**：每版只改一個變數（新增一個 hook 層級），並在每版執行後撰寫分析文件。這樣可以：
1. 明確知道哪個 hook 有效、哪個無效
2. 避免「多個 hook 同時失敗」時的混淆
3. 累積的知識可以指導下一版的 hook 策略

#### Trigger

- App 使用多層網路架構（Dart + OkHttp + Netty proxy + TLS）
- 不確定流量經過哪些層級
- 需要系統性地找出有效的 hook 點
- 先前嘗試加入多個 hook 但無法判斷哪個有效

#### Evidence

- Tool: Frida（`hook_self_generation_phase1_v*.js`，v4 到 v11）
- Sanitized excerpt: 8 個版本逐步排除 Java HTTP → Java Socket → Netty proxy → SSLSocket → libc → Netty SslHandler → Native SSL → Netty channelRead
- Evidence path: `<PROJECT_ROOT>/capture/self_gen_phase1_v10_analysis.md`（Summary of All Hook Levels Tested）

#### Generalized Lesson

當需要為多層網路架構設計 Frida hook 策略時：

1. **從最上層開始**：先 hook 最高層級的 API（OkHttp、HttpURLConnection），因為最容易解析
2. **每版只加一個 hook 層級**：不要一次加入所有 hook。每版只新增一個 hook 點
3. **每版執行後撰寫分析文件**：記錄哪些 hook 觸發、哪些沒有、以及捕獲到的流量樣式
4. **根據結果決定下一版方向**：
   - 如果上層 hook 沒觸發 → 往更底層移動（如 libc socket）
   - 如果上層 hook 有觸發但資料不完整 → 在同一層級嘗試不同 hook 點
   - 如果看到加密流量 → 尋找 TLS 層級的 hook
5. **保留所有版本的分析文件**：即使某個 hook 沒觸發，這個「負面結果」也是重要的知識
6. **使用版本號命名**：`hook_xxx_v1.js`、`hook_xxx_v2.js`，方便追蹤演進

#### Agent Action

下次 agent 需要為多層網路架構設計 Frida hook 時：

1. 先繪製預期的網路架構圖（Dart → OkHttp → Proxy → TLS → Server）
2. 從最高層級開始設計 v1 hook script
3. 執行 v1 並撰寫分析文件（記錄觸發/未觸發的 hook）
4. 根據分析結果設計 v2（只新增一個 hook 層級）
5. 重複直到成功捕獲目標流量
6. 保留所有版本的分析文件作為參考

#### Goal / Action / Validation

- Goal: 讓 agent 能系統性地找出多層網路架構的有效 hook 點
- Action: 逐版增量 hook，每版只新增一個層級，每版執行後撰寫分析文件
- Validation or reference source: 實測驗證 — v4 到 v11 逐步排除無效 hook 點，最終在 v11 成功捕獲完整 HTTP body

#### Applies When

- App 使用多層網路架構（Dart + Java + Netty + TLS）
- 不確定流量經過哪些層級
- 需要系統性地找出有效的 hook 點
- 分析過程需要可重現的文件記錄

#### Does Not Apply When

- App 使用單一網路層級（如純 OkHttp）
- 已知有效的 hook 點（不需要探索）
- 只需要快速的一次性捕獲（不需要版本控制）

#### Validation

- 為每個 hook 層級建立獨立的版本
- 每版執行後檢查 hook 觸發計數
- 如果某個 hook 在 30 秒內沒有事件，記錄為「未觸發」
- 如果所有 hook 都未觸發，檢查 Frida 連線健康狀態

#### Promotion Target

- `analysis/apk/workflows/dynamic-capture-flow.md`（新分層）

#### Required Linked Updates

- 無需連動更新；此為新 lesson，尚未 promote 到正式文件
