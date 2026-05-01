> 遵守 [共用規則索引](../../../shared-rules/README.md) 與 [feedback-lessons](../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-01 - DoH 的 `dns=` 參數可作為「MITM 業務 host 空白」時的側信道

Status: validated

#### One-line Summary

當業務 HTTPS 未進電腦 MITM、但 **OkHttp** 仍會發 **DoH GET**（`application/dns-message`）時，可從 URL 的 **`dns=<base64url>`** 離線解出 **問句 QNAME**，用來確認 App 仍會解析哪些 **API 根域名**。

#### Human Explanation

「代理裡搜不到 `api.example.com`」常被誤判成沒打業務；若 client 改走 **HTTPS DNS**，MITM 日誌可能只留下對 **DoH resolver host** 的連線，業務域名藏在 **binary DNS wire** 裡。Frida 若只 hook **OkHttp**，仍可能記到完整 **`https://…/dns-query?dns=…`** 字串，這時不必等到 TLS 明文即可還原「想解析誰」。

#### Trigger

- Wi‑Fi／全域 HTTP MITM 幾乎沒有業務 API 的 **CONNECT** 或明文。
- **`getaddrinfo`／pcap** 已暗示業務樹存在，但需要與 **Java** 路徑交叉佐證。
- 側錄 log 出現 **`Accept: application/dns-message`** 與 **`dns=`** query。

#### Evidence

- Tool: Frida hook **`RealCall`**（headers-only 腳本即可保留 URL）；離線 **`base64.urlsafe_b64decode`** + DNS question walk。
- Sanitized excerpt: `GET https://<doh-resolver>/dns-query?dns=<base64url>` → QNAME **`<api-root-1>`**、**`<api-root-2>`**（A / AAAA）。
- Evidence path: `<PROJECT_ROOT>/capture/*.log`（勿提交含敏感標頭的 raw 檔）。

#### Generalized Lesson

DoH 請求的 **`dns=`** 是 **可離線解析** 的側信道；與 **「MITM 看不到業務 host」** 並不矛盾。應把 DoH 從「泛稱噪音」中拆出：**若 QNAME 指向產品 API 根，即為業務相關證據**。

#### Agent Action

1. 在側錄中搜 **`dns-query?dns=`** 或 **`application/dns-message`**。
2. 對每個獨立 **`dns=`** 參數解 base64url，解析問句 **QNAME**、**QTYPE**。
3. 把結果以 **去敏域名／角色** 寫進專案 API 筆記；完整 resolver hostname 若敏感可占位。
4. 可選：在專案提供 **`decode_capture_doh.py`** 類小工具避免手解。

#### Applies When

- App 使用 **OkHttp／Cronet** 等會暴露完整 URL 的 stack 發 DoH。
- 你需要在 **無 TLS 明文** 前提下確認「解析目標」。

#### Does Not Apply When

- DNS 走 **OS 解析器**且 HTTP 層看不到 query（僅 **pcap／getaddrinfo** 可見）。
- **`dns=`** 經額外加密或非標準封裝（需個案分析）。

#### Validation

對同一 **`dns=`** 字串重複解碼得到相同 **QNAME**；與 **pcap SNI／getaddrinfo** 清單 **方向一致**（不需逐字相同時戳）。

#### Promotion Target

- `TOOLS.md`（可選：列入「MITM 空白時的側信道檢查清單」）
- `DOCUMENTATION.md`
