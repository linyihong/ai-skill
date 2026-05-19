> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-18 — Flutter + Java 混合架構 App 的安全架構總覽

#### One-line Summary

目標 App 是一個 Flutter (Dart AOT) + Java/Kotlin 混合架構的 Android 應用，使用多層安全防護：內建防機器人 SDK（自訂 TLS + AES-CTR token）、第三方閘道（TLS 指紋 + HTTP 標頭驗證）、雙層 AES-CBC 請求/回應加密、以及三種不同的請求簽名標頭格式（標準、訪客登入、minimal）。

#### Human Explanation

這類 App 的架構和安全設計非常複雜，分析時需要同時理解多個層面。以下是完整的架構圖和安全機制分析，幫助未來更快定位問題。

#### Trigger

當分析一個 Flutter + Java 混合架構的 Android App，且遇到以下情況時：
- API 請求被阻擋（HTTP 400/551）
- 需要理解 App 的加密/簽名機制
- 需要區分不同 API 的請求簽名標頭格式
- 遇到第三方防機器人系統

#### Evidence

從原始碼、Frida capture、文件分析得出的完整架構：

**1. 技術框架**

| 層級 | 技術 | 說明 |
|------|------|------|
| UI 框架 | **Flutter (Dart)** | 所有 UI 使用 Flutter 實作，Dart AOT 編譯到 `libapp.so` |
| 原生橋接 | **Java/Kotlin** | 防機器人 SDK、ProxyServer、部分 native 功能用 Java |
| HTTP 客戶端 | **Dart `dart:io` HttpClient** | 主要 API 呼叫使用 Dart 層級的 HttpClient |
| 本地 Proxy | **Netty (Java)** | ProxyServer 使用 Netty 作為透明轉發代理 |
| 序列化 | **Dart `jsonDecode`** | 所有 API 回應先經 Dart 層 jsonDecode 再傳給 UI |
| 建置工具 | **Gradle + Maven** | App 本體用 Gradle，SDK 專案用 Maven |

**2. 安全架構（4 層防護）**

```
┌─────────────────────────────────────────────────────┐
│ 第 1 層：第三方閘道（TLS 指紋 JA3/JA3S）             │
│  - 阻擋標準 Java/OkHttp TLS                         │
│  - 只信任 Dart `dart:io` HttpClient 的 TLS 指紋      │
│  - 或內建 ProxyServer 的自訂 TLS                     │
├─────────────────────────────────────────────────────┤
│ 第 2 層：第三方閘道 HTTP 標頭驗證                    │
│  - 檢查請求簽名標頭是否存在                          │
│  - 檢查 User-Agent 是否為 `Dart/<version> (dart:io)` │
│  - 阻擋 HTTP/2 升級嘗試                              │
├─────────────────────────────────────────────────────┤
│ 第 3 層：請求簽名標頭（AES-CBC）                     │
│  - 三種不同格式（見下方）                             │
│  - 金鑰：`<encryption-key>`（硬編碼）                │
│  - 輸出：16-char prefix + Base64(ciphertext)         │
├─────────────────────────────────────────────────────┤
│ 第 4 層：回應加密（AES-CBC）                         │
│  - 相同金鑰 `<encryption-key>`                       │
│  - 格式：16-char prefix + Base64(ciphertext)         │
│  - 部分 API 使用 dual-payload 格式（兩個加密 payload 串接）│
└─────────────────────────────────────────────────────┘
```

**3. 請求簽名標頭的三種格式**

| API 類型 | 明文格式 | 核心欄位 | 使用場景 |
|----------|---------|---------|---------|
| **標準** | `FUCKYOU\|ts\|rand\|apiPath\|\|service&uid=...&token=...&package_name=...&version=...&device=...` | uid, token, package_name, version, device | 一般已登入 API 呼叫 |
| **訪客登入** | `FUCKYOU\|ts\|rand\|apiPath\|\|service=Login.guestLogin&package_name=...&deviceInfo=...&rtm=...&_sign=...` | deviceInfo, rtm, _sign（無 uid/token） | 訪客登入（未登入狀態） |
| **Minimal** | `FUCKYOU\|ts\|rand\|apiPath\|\|service` | 只有 service 名稱 | 唯讀查詢 API（不需要身份驗證） |

**4. 內建防機器人 SDK**

目標 App 內建防機器人 SDK，提供：
- **自訂 token**：AES-CTR 加密，8-byte key + 8-byte seed（zero-padded to 16 bytes），輸出 44 Base64 chars（32 decoded bytes）
- **輔助 token**：與自訂 token 相同機制
- **ProxyServer**（特定埠）：透明轉發代理，使用自訂 TLS 實作繞過第三方閘道
- **簽名參數**：訪客登入請求中的簽名，需進一步分析

**5. 回應解密**

所有 API 回應都可能被 AES-CBC 加密：
- 單一 payload：`16-char prefix + Base64(ciphertext)`
- Dual payload：`prefix1 + Base64(ciphertext1) + prefix2 + Base64(ciphertext2)`
- 解密方式：prefix 作為 IV（`IV.fromUtf8(prefix)`），金鑰與請求簽名相同

**6. 特定內容模組 API 特性**

特定內容模組 API 與一般 API 有顯著不同：
- 使用 `serviceHash`（如 `<service-hash>`）而非 service name
- 部分 API 使用 minimal 簽名格式（無 uid/token）
- 回應可能是 dual-payload 加密
- 需要先調用列表 API 獲取內容 ID，再用該 ID 調用詳情 API
- 播放 API 需要多個參數

#### Generalized Lesson

分析 Flutter + Java 混合架構的 Android App 時，應按以下順序建立心智模型：

1. **先識別框架**：檢查 `libapp.so`（Dart AOT）、`libflutter.so`、Java classes.dex
2. **再識別安全層**：
   - 是否有防機器人 SDK
   - 請求簽名機制（自訂標頭、token）
   - 回應加密
3. **區分 API 類型**：不同 API 可能有不同的簽名格式和加密方式
4. **建立 API 呼叫圖**：列表 → 詳情 → 播放的依賴關係
5. **使用 Frida 驗證**：hook `jsonDecode`、`encrypt`、socket 操作來驗證假設

#### Agent Action

分析此類 App 時：

1. 先讀取 `libapp.so` 的 Dart AOT 符號（使用 unflutter/blutter）
2. 檢查 Java 類別中是否有防機器人 SDK 相關套件
3. 用 Frida hook `jsonDecode` 快速取得 API 回應結構
4. 用 Frida hook `encrypt` 函數取得加密參數（金鑰、IV、明文格式）
5. 區分不同 API 的簽名格式（不要假設所有 API 使用相同格式）
6. 注意回應加密的存在——不要被 Base64 字串誤導為 plaintext

#### Goal / Action / Validation

- **Goal**：快速建立 App 的完整架構和安全模型
- **Action**：按上述步驟依序分析
- **Validation**：能夠成功發送 API 請求並正確解密回應

#### Applies When

- 分析 Flutter + Java 混合架構的 Android App
- 遇到多層安全防護（防機器人 + 請求簽名 + 回應加密）
- 需要區分不同 API 的認證/簽名格式

#### Does Not Apply When

- 純原生 Android App（無 Flutter）
- 無自訂安全機制的簡單 API

#### Validation

- 能正確生成三種請求簽名標頭格式並被伺服器接受
- 能正確解密單一 payload 和 dual-payload 回應
- 能透過裝置 proxy 成功發送請求（繞過第三方閘道 TLS 檢查）
- 內容模組 API 的列表 → 詳情 → 播放流程完整可執行

#### Promotion Target

`intelligence/engineering/analytical-reasoning/heuristics/flutter-java-hybrid-security-layer-mapping.md`（✅ 已 promote）

#### Required Linked Updates

- `workflow/apk-analysis/execution-flow.md`：加入 App 架構分析步驟
- `intelligence/engineering/analytical-reasoning/heuristics/`：加入 Flutter+Java 混合架構分析啟發式
