> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。
# Extracted — See [`analysis/apk/workflows/http-api-documentation-flow.md`](../../../../analysis/apk/workflows/http-api-documentation-flow.md)

### 2026-05-11 - Articles-First Live Adapter Smoke

Status: candidate

#### One-line Summary

SDK/private adapter 要驗真實 read-only 資料時，先讓一條核心 list route 通過，secondary routes 不應阻塞首個 live proof。

#### Human Explanation

APK 分析轉 SDK 時，常會同時知道列表、分類、詳情、留言、媒體等多條 route。若最小 live smoke 一開始就要求所有 route binding、所有 service、所有 decrypt/media 能力都齊全，會把「第一條正式資料鏈路是否可跑」和「全功能 parity 是否完成」混在一起。較好的閉環是先選核心 read route（通常 list/page 1），只要求這條 route 的 base endpoint、route binding、opaque/session provider、identity readiness、signing、decrypt/plaintext boundary，成功後再把分類、詳情、留言、next-page 或 media 當 secondary smoke 擴展。

#### Trigger

當分析輸出已足夠建立 SDK parser/transport skeleton，但 standalone live fetch 還需要 private adapter 或 trusted runtime bridge。

#### Evidence

- Tool: SDK live env gate / live smoke test / private adapter checklist.
- Sanitized excerpt: gate requires `BASE_URI`, one list-route binding, opaque provider, identity readiness marker, signing/private request material or trusted bridge, and decrypt/plaintext boundary; secondary route bindings are optional.
- Evidence path: project docs may store route id, key names, schema classes, pass/fail, and missing capability names only.

#### Generalized Lesson

最小 live adapter smoke 應該驗證一條最重要、read-only、可安全重跑的核心 route。Secondary routes 可用來擴大 route consistency，但不應成為第一個 runnable gate，除非使用者目標明確是「完整 route parity」而不是「先拿到一頁正式資料」。

#### Agent Action

下次把 APK findings 轉成 SDK/private adapter 時，先問當前 smoke 的 target route 是哪一條。把 gate 拆成 minimum required route 和 optional secondary routes；測試名稱、env var、文件都要反映這個邊界。若既有 smoke test 把 secondary route 當必填，先改成 optional，再跑 focused validation。

#### Goal / Action / Validation

- Goal: 降低第一條 live data proof 的阻塞面，避免因 secondary routes 缺私有材料而誤判 SDK 無法前進。
- Action: 只把核心 list route 的 private adapter 材料設為 minimum runnable requirements。
- Validation or reference source: focused live env tests should pass without secondary bindings, and optional route tests should run only when those bindings exist.

#### Applies When

- 目標是 SDK/private adapter 的第一個 read-only live proof。
- 已有 parser/model/transport skeleton，但 private base/service/sign/decrypt/session material 尚未全部齊。
- Secondary routes 不影響核心 list route 的安全讀取。

#### Does Not Apply When

- 使用者明確要求 full route parity 或 release gate。
- 核心 route 本身依賴 secondary route 先取得 token、cursor、category id 或 decrypt key。
- Read route 不是安全 read-only 操作，或缺少授權 runtime。

#### Validation

1. 缺少核心 route binding 時 fail-fast。
2. 缺少 secondary route binding 時 minimum smoke 仍 runnable。
3. Optional secondary route 只有在 binding/material 存在時才執行。
4. 文件把 minimum smoke 和 route expansion 分開。

#### Promotion Target

- `WORKFLOW.md`

#### Required Linked Updates

- 已同步 `WORKFLOW.md` 的 Articles-first live adapter smoke rule。
- 已依 reusable-guidance-boundary 檢查：具體 App/env var/service hash 留在 project docs，本 lesson 只保留通用 live-smoke gate 方法。
