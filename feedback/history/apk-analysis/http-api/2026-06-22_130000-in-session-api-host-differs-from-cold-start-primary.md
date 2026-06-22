> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - In-session API host may differ from cold-start primary host

Status: candidate

#### One-line Summary

短劇／內容 App 的**業務 REST 路徑**可能在冷啟動走 dex 上的 primary API host，進入播放或 in-session 後改走 **CDN edge（如 CloudFront）** 同一 path prefix；只 filter primary host 會漏掉 heartBeat、getChapterContent 等播放期 API。

#### Human Explanation

靜態 triage 常只記一個 `v-api.*` 或 publisher primary host。動態上 client 可能在 init/config 後把 Retrofit base URL 切到 CDN 域名，**path 與 header 家族不變**。若 Frida/MITM 只 match primary host，W1 有流量、W3 播放期卻「0 條業務 API」——實際是 host 變了，不是 hook 失效。

#### Trigger

- W1 cold-start capture 在 primary host 有 `/api/` POST
- 同 session 播放／切集後，primary host 無新 `/api/` 請求
- Capture log 出現 `*.cloudfront.net`（或其它 CDN）+ 相同 `/api/video/` 或 `/api/ms/` path

#### Evidence

- Tool: Frida OkHttp chain hook，對比 cold spawn vs attach 播放窗口
- Sanitized excerpt: 冷啟動 host A 有 splash/config；播放期 host B 有 heartBeat / chapter content，path 前綴一致
- Evidence path: `<PROJECT_ROOT>/api/dynamic-*.md`；**本 lesson 不含** 具體 distribution、book/chapter id、raw log

#### Generalized Lesson

**Host 排查（與 pinning tier 並用）：**

```text
同一 session 內 business /api/ 流量為 0？
  1. 導流 / pinning（見 local-proxy pinning tier lesson）
  2. Host filter 是否只含 static primary？
     → 放寬：*.cloudfront.net、*.akamaized.net 等 + path 含 /api/
  3. 分窗口記錄：cold-start hosts vs in-session hosts
```

報告應列 **host × phase** 矩陣，不是單一 canonical host。

#### Agent Action

1. Frida `interest()` / MITM 過濾：path 優先（`/api/`）或 CDN + `/api/`，勿只 match dex primary FQDN。
2. Capture plan：W1 cold + W3 attach 播放各跑一窗，合併 path inventory。
3. Ai-skill 只寫決策規則；FQDN 真值留 project inventory。

#### Goal / Action / Validation

- Goal: 避免「primary 有 API、播放期以為沒 API」的 false negative。
- Action: `workflow/apk-analysis/execution-flow.md` 動態段補 host failover 檢查。
- Validation: 播放窗口在 CDN host 上命中 ≥1 條與 W1 同 path family 的 POST。

#### Applies When

- DEX 有 primary API host + 多 CDN 域名
- 播放／heartbeat／chapter 類 path 在靜態存在但 attach 期 primary 無流量

#### Does Not Apply When

- 全 session 僅 single host（無 failover 證據）
- 流量走 gRPC/WebSocket 非 REST path

#### Validation

- 同一 log 內列出 cold vs in-session 的 distinct host（path 去重）
- 放寬 host filter 後播放 API 計數 > 0

#### Promotion Target

- `analysis/apk/traffic-triage.md` §host inventory
- `workflow/apk-analysis/execution-flow.md` §W1/W3 windows

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 交叉引用 `local-proxy/2026-06-22_120100-mitm-cdn-visible-primary-api-invisible-pinning-tier.md`
