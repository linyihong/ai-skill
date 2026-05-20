# HTTP API Documentation Flow（HTTP API 文件操作流程）

`analysis/apk/workflows/http-api-documentation-flow.md` 是從 `skills/apk-analysis/techniques/http-api/`（已刪除）拆解出的 **HOW TO DO** 操作流程。決策智慧（何時該開始、何時該停、欄位信心判斷）請見 `intelligence/engineering/analytical-reasoning/heuristics/api-documentation-completeness.md`。

> **Intelligence Extracted**
> See:
> - `intelligence/engineering/analytical-reasoning/heuristics/api-documentation-completeness.md`

## 前置準備

### 必要條件

- HTTP/HTTPS 流量已可觀測（MITM、Frida hook、local proxy handler、Flutter/Dart interceptor、pcap）
- 目標 API 的 method/path/header/request/response metadata 已可見或已解碼

### 工具

```bash
# 流量捕獲
mitmproxy -p 8080          # MITM proxy
tcpdump -i any port 443    # pcap 捕獲

# Frida hook（選擇性）
frida -U -l hook_okhttp.js -f com.target.app --no-pause
```

## 步驟 1：建立 API Entry

建立 API 目錄的入口文件，記錄：

| 項目 | 內容 |
|------|------|
| Hosts/base URLs | 所有觀測到的主機與 base path |
| Traffic families | API 分組（如 `api.*.com/v1/`、`cdn.*.com/`） |
| Wrapper/decode rules | 統一的 response wrapper 結構與解碼規則 |
| Shared headers | 所有 API 共用的 header（如 `Authorization`、`X-Device-Id`） |
| Links | 指向 coverage/gap matrix、UI map、SDK/client notes |

## 步驟 2：建立 Group Index

將 API 依 path prefix、domain、feature 或 protocol family 分組：

```markdown
| Group | Base Path | APIs | Status |
|-------|-----------|------|--------|
| User | `/api/v1/user/` | login, profile, settings | documented |
| Content | `/api/v1/content/` | feed, detail, search | partial |
| Media | `/cdn.example.com/` | image, video | needs capture |
```

每行連結到 per-API detail。

## 步驟 3：建立 Per-API Detail

對每個觀測到的 API endpoint，記錄以下欄位：

| Area | Required Notes |
|------|---------------|
| Identity | Method, host/path shape, auth conditions, evidence source, UI path if confirmed |
| Capability mapping | Feature/capability, operation id, user-visible behavior, trigger confidence, startup/preload/background or direct user action |
| Request headers | Header name, purpose, required/optional, source, sensitivity, token/sign/device/session involvement |
| Request query/body | Field type, meaning, required/optional, example shape, sensitivity, signing/encryption participation |
| Response headers | Status behavior, content type, cache/rate-limit/session headers; if invisible, state why |
| Response wrapper | `status`, `code`, `message`, `data`, `error`, and other outer fields with type and meaning |
| Inner payload | Field type, meaning, nullability, list item shape, media/source fields, derived values |
| Functional contract | Candidate domain concepts, commands/events, state impact, empty/error behavior, pagination/cache semantics, open questions |
| Validation | Replay, fixture, contract test, or hook/pcap/MITM sequence proving request/response alignment |
| Catalog status | grouped, per-API detail exists, coverage status, SDK/client mapping status when relevant |

如果 UI binding 尚未完成，寫 `UI path: unknown` 和 `Trigger confidence: low`。

## 步驟 4：建立 Coverage / Gap Matrix

追蹤哪些 API 已觀測、已 replay、已解碼、已對應 UI：

```markdown
| API | Observed | Replayed | Decoded | UI-bound | Tested | Gaps |
|-----|----------|----------|---------|----------|--------|------|
| POST /api/v1/login | ✅ | ✅ | ✅ | ✅ | ✅ | none |
| GET /api/v1/feed | ✅ | ❌ | ✅ | ❌ | ❌ | pagination params unknown |
```

## 步驟 5：建立 SDK/Client Mapping（如適用）

記錄 SDK/client 實際消耗的欄位、相容性預期、fixture/test 狀態。

## 步驟 6：執行 Catalog Finish Gate

在回報 API-list 任務完成前，檢查：

- [ ] 每個觀測或解碼的 API 都在 group index 或 coverage/gap file 中
- [ ] 高價值 API 有 per-operation detail，不只是 method/path 行
- [ ] 每個 per-operation detail 包含 request fields、response fields、field meaning、evidence、validation、open questions
- [ ] Shared headers、wrapper/decode behavior、auth/session、sensitivity rules 已文件化一次並從 API details 連結
- [ ] UI/API mapping 記錄 operation id、capture window、trigger confidence、startup/preload/background 狀態
- [ ] SDK/client/tool usage 記錄 consumed fields 和 fixture/test 狀態
- [ ] 未驗證的 API 明確標記為 `candidate`、`needs capture`、`needs replay`、`meaning unknown`、`low confidence`、`out of scope` 或 `not observed`

## 步驟 7：UI Automation For API Capture（選擇性）

對高價值流程，可用自動化腳本穩定 API capture：

1. 給每個 flow 一個穩定的 `operation_id`
2. 給每個 route 一個穩定的 `route_id`，說明如何到達目標畫面
3. Route map 限於 in-app pages；遇到 system screen、browser、payment/share sheet、third-party app、external intent 時標記為 external transition
4. 分類每個 in-app screen 為 scrollable 或非 scrollable，記錄 clickable entry points
5. 每個 script 限於一個 UI path 或 action group（如 `open-home`、`scroll-feed`）
6. 將 route recipe 拆成可重用 navigation segments，並登記到對應 UI map 文件；每段記 `segment_id`、entry checkpoint、exit checkpoint、preconditions、script/function path、selector/coordinate source、evidence
7. 長流程只保存為 segment composition（如 `launch-to-home -> home-to-detail -> detail-scroll-media`），不要只保存不可拆的單一腳本
8. 後續 capture 先組合既有 segments 到目標 checkpoint，只重測缺失或失效 segment，避免每次從頭重跑完整 navigation
9. Scrollable screens 使用 bounded sampling（top/mid/bottom）
10. Clickable screens 使用 labels、resource IDs、content descriptions、hierarchy bounds 或 verified coordinates
11. 在 operation 前後印出 UTC start/end timestamps
12. 在同一 window 中執行 pcap/MITM/Frida capture
13. 在 operation 結束時儲存一張 sanitized screenshot 或 UI hierarchy
14. 填寫 operation-to-API matrix（route id、method/path、source、response shape、trigger confidence）

## 成功產出格式

```markdown
# API Catalog

## Entry
- Base URLs: ...
- Shared headers: ...
- Wrapper: ...

## Group Index
| Group | APIs | Status |
|-------|------|--------|

## Per-API Detail: POST /api/v1/login
- Method: POST
- Path: /api/v1/login
- Auth: none (public)
- Request: { username, password }
- Response: { token, expires_in }
- Evidence: MITM capture 2026-05-12
- Validation: replayed successfully
```

## 注意事項

- Screenshots 可支援 UI trigger attribution，但不取代 HTTP header/request/response field analysis
- 如果目標是重建 feature，API docs 應準備好讓 `workflow/software-delivery/` 轉換為 BDD、Domain Model Contract、API/Interface Contract
- 不確定的 field meaning 或 domain vocabulary 標記為 `candidate`
- 避免執行 login loops、payment、destructive actions、posting、messaging、account changes
