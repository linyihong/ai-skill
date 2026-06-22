> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Encrypted requestTime — double native call, second output enters sign map

Status: candidate

#### One-line Summary

除 plain `ts` epoch 外，另有 **encrypted `requestTime` header** 時，interceptor 常 **連續兩次** 呼叫同一 native getter（seed 可能是 **APK install path** 等非 obvious 值），**第二次輸出** 寫入 header map 並納入 sign canonical；RPC relay 須 mirror 雙呼叫，不可用第一次輸出簽名。

#### Human Explanation

分析者常假設 `requestTime = f(devId)` 或 `f(ts)`。實務上 seed 可能是 `Context.getPackageCodePath()` 類路徑，且每次請求 native 輸出不同（含隨機/時間因子）。Bytecode 對同一 getter 連 call 兩次：第一次可能被丟棄或作 side effect，第二次才進 HashMap。Frida log 若見成對 `[REQTIME]`，須對照 `[SIGN_IN]` 中 `requestTime=` 用的是哪一次。離線 relay 腳本漏掉雙呼叫會導致 sign 正確格式但 server 拒絕。

#### Trigger

- Header 同時有 `ts`（plain epoch）與 `requestTime`（Base64-like blob）
- Frida：同一 intercept 內 getter 連續 hit 2 次，輸出不同
- `[SIGN_IN]` 的 `requestTime=` 等於第二次 `[REQTIME] out=`
- Seed 參數為 install path / file path 而非 device id

#### Evidence

- Tool: Frida hook getter + path provider + interceptor
- Sanitized excerpt: `in=<apk install path> out=<blob>` ×2；sign canonical 含第二次 blob
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`

#### Generalized Lesson

```text
requestTime RE:
  1. Hook native getter in + out (not only header on wire)
  2. Count calls per intercept — if 2, map second to sign canonical
  3. Identify seed (path, devId, ts, combo) from getter arg
  4. Frida RPC: call getter twice with same seed before sha256Encrypt
  5. Do not copy requestTime across requests (per-request blob)
```

#### Agent Action

1. Project docs 記 seed 類型與 double-call（無 live path）。
2. Ai-skill 只寫 pattern；不寫 blob 樣本。
3. SDK gap matrix：`requestTime` = blocking separate from `sign` crypto。

#### Goal / Action / Validation

- Goal: 避免 relay 用錯 requestTime 世代或只 call 一次 getter。
- Action: sign RPC wrapper mirror interceptor call order。
- Validation: RPC 產生的 header 組可通過至少一個 business POST。

#### Applies When

- Encrypted timestamp header + separate plain epoch
- Native getter in signing interceptor path

#### Does Not Apply When

- requestTime is plain ms/epoch string
- Single getter call observed end-to-end

#### Validation

- Pair REQTIME logs aligned with SIGN_IN requestTime field
- RPC replay succeeds on ≥1 endpoint

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §sign RE
- 與 `141800` canonical pattern 並列

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
