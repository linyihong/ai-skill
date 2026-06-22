> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Partial standalone SDK when HMAC sign is offline but requestTime/session still native

Status: candidate

#### One-line Summary

自訂簽名 API 的 standalone client **不必全有或全無**：當 `sign = HMAC-SHA256(mode_key, canonical)` 已離線驗證後，可實作 **hybrid signer**（host 算 `sign` + Frida/device relay `requestTime` 與 live `session`）；`141700` 的「blocked until RE or relay」應細化為 capability 矩陣，而非單一 boolean。

#### Human Explanation

v-api 類請求常同時需要：`canonical`（可從 interceptor 規則重建）、`requestTime`（獨立 native 加密、常雙次呼叫）、`session`（登入態）、`sign`（canonical 的 MAC）。完整離線 SDK 需四者皆可離線；實務上 sign 先突破時，仍可用最小 relay：`getSystemLaunchTime(apkPath)` + `session` 從運行中 app 讀取，其餘在 host 完成。SDK readiness 表應分列：canonical ✅、sign ✅、requestTime ❌、session ⚠️。避免在 sign 已解後仍宣稱「完全不能離線」而延遲 hybrid 測試。

#### Trigger

- `141700` 場景：custom `sign` + `requestTime` headers
- HMAC sign verified offline against Frida capture
- `requestTime` still `native` with per-call varying blob
- User asks for SDK / curl / Python client progress

#### Evidence

- Tool: offline `hmac_sign` + Frida RPC `getrequesttime` + hybrid POST script
- Sanitized excerpt: same canonical → host sign matches JNI; POST still needs fresh requestTime from device
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`、`<PROJECT_ROOT>/docs/domain-baseline.md` SDK table

#### Generalized Lesson

```text
Custom-sign SDK readiness (update per breakthrough):
  canonical layout     → document from interceptor / capture
  sign algorithm+key   → offline when HMAC verified (142300)
  requestTime          → native relay until RE (142100 double-call)
  session/auth         → capture or login RE
Ship hybrid client when sign+canonical offline; do not block on requestTime RE alone.
```

#### Agent Action

1. 更新 project `domain-baseline.md` SDK 矩陣為分列狀態。
2. 提供 `hybrid_*` script 範式（Frida fields + offline sign）。
3. 每輪結束檢查：是否有新可離線欄位 → 更新矩陣 + 是否需新 Ai-skill 範式（本條）。

#### Goal / Action / Validation

- Goal: 進度可 incremental 交付，不因單一 native 欄位卡住已可驗證部分。
- Action: hybrid POST dry-run + live test when app on device.
- Validation: offline sign matches JNI; live POST 200 on free endpoint (optional).

#### Applies When

- Multi-field custom signing on mobile API
- At least one field still native-only

#### Does Not Apply When

- All sign inputs offline including requestTime
- Server accepts unsigned dev endpoints

#### Validation

- SDK table shows per-field status; hybrid script exists in project

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §SDK readiness matrix

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- `141700` 可交叉引用本條作 partial-unblock 細化
- 已依 sanitization / reusable-guidance-boundary 自查
