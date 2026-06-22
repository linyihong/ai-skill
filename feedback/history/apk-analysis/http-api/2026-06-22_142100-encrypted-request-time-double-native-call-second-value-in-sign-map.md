> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Encrypted time header — double native call, second output enters sign map

Status: candidate

#### One-line Summary

除 plain epoch header 外，另有 **encrypted time header**（Base64-like blob）時，interceptor 常 **連續兩次** 呼叫同一 native getter（seed 可能是 **APK install path** 等非 obvious 值），**第二次輸出** 寫入 header map 並納入 sign canonical；RPC relay 須 mirror 雙呼叫，不可用第一次輸出簽名。

#### Human Explanation

分析者常假設 encrypted time = f(deviceId) 或 f(plain epoch)。實務上 seed 可能是 install path 類參數，且每次請求 native 輸出不同（含隨機/時間因子）。Bytecode 對同一 getter 連 call 兩次：第一次可能被丟棄或作 side effect，第二次才進 HashMap。Frida log 若見成對 getter hit，須對照 sign canonical 用的是哪一次輸出。離線 relay 腳本漏掉雙呼叫會導致 sign 格式正確但 server 拒絕。**wire 欄位名**留在專案 evidence。

#### Trigger

- Header 同時有 plain epoch 與 encrypted-time blob
- Frida：同一 intercept 內 getter 連續 hit 2 次，輸出不同
- Sign canonical 含 encrypted-time 欄位，且值等於第二次 getter 輸出
- Seed 參數為 install path / file path 而非 device id

#### Evidence

- Tool: Frida hook getter + path provider + interceptor
- Sanitized excerpt: `in=<apk install path> out=<blob>` ×2；sign canonical 含第二次 blob
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`

#### Generalized Lesson

```text
Encrypted-time header RE:
  1. Hook native getter in + out (not only header on wire)
  2. Count calls per intercept — if 2, map second to sign canonical
  3. Identify seed (path, devId, ts, combo) from getter arg
  4. Frida RPC: call getter twice with same seed before sign MAC
  5. Do not copy encrypted-time blob across requests (per-request)
```

#### Agent Action

1. Project docs 記 seed 類型與 double-call（無 live path）。
2. Ai-skill 只寫 pattern；不寫 blob 樣本。連續呼叫快取語意見 `142800`。
3. SDK gap matrix：encrypted-time = blocking separate from `sign` crypto。

#### Goal / Action / Validation

- Goal: 避免 relay 用錯 encrypted-time 世代或只 call 一次 getter。
- Action: sign RPC wrapper mirror interceptor call order。
- Validation: RPC 產生的 header 組可通過至少一個 business POST。

#### Applies When

- Encrypted timestamp header + separate plain epoch
- Native getter in signing interceptor path

#### Does Not Apply When

- Encrypted-time header is plain ms/epoch string (no native blob)
- Single getter call observed end-to-end

#### Validation

- Pair getter logs aligned with sign canonical encrypted-time field
- RPC replay succeeds on ≥1 endpoint

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §sign RE
- 與 `141800` canonical pattern 並列

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
