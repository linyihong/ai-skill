> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - R8-obfuscated OkHttp Request breaks RealCall hooks; use RealInterceptorChain.proceed

Status: candidate

#### One-line Summary

R8/ProGuard 可能把 `okhttp3.Request` 混淆成短類名（如 `okhttp3.o0`）；hook `RealCall.execute` / `request()` 常報 `not a function` 或 overload 不匹配；應 probe `RealInterceptorChain` 的 `proceed` 簽名並 hook 混淆後的 Request 類型。

#### Human Explanation

Java 主線 OkHttp 捕獲常從 `okhttp3.RealCall` 或 `okhttp3.Request` 入手。混淆後 public API 名稱可能消失：`proceed()` 只接受 `okhttp3.<short>` 而非 `okhttp3.Request`，`RealCall.request()` 在 Frida 中也可能不可用。若 hook 只裝上 banner、從不印 URL，先懷疑 **overload 錯誤**，不是「app 沒發請求」。

#### Trigger

- Frida hook `RealCall.execute` / `enqueue` 報 `TypeError: not a function` 或 `request is not a function`
- `Java.use('okhttp3.Request')` 存在但 `proceed('okhttp3.Request')` overload 失敗；錯誤訊息列出實際參數型別（如 `okhttp3.o0`）
- `RealInterceptorChain` hook 已安裝但 0 條業務 URL

#### Evidence

- Tool: Frida attach/spawn + `getDeclaredMethods()` on `okhttp3.internal.http.RealInterceptorChain`
- Sanitized excerpt: `proceed` 簽名為 `proceed(okhttp3.<obfuscated>)`；`request()` 回傳同混淆型別；改 hook 該 overload 後業務 API URL 出現
- Evidence path: hook 腳本與 capture 摘要留在 `<PROJECT_ROOT>/scripts/frida/` 與 `<PROJECT_ROOT>/capture/`（gitignore 視專案設定）

#### Generalized Lesson

**OkHttp 混淆捕獲流程：**

```text
1. Java.use('okhttp3.internal.http.RealInterceptorChain')
2. chain.class.getDeclaredMethods() → 找 proceed(...) 的參數型別 T
3. hook proceed.overload(T) → 記錄 request.toString() 或 request.a().toString()（依混淆後方法名 probe）
4. 勿假設 T === 'okhttp3.Request'
```

與 `RealCall` hook 並用時，以 **chain proceed 命中** 為準；RealCall 為輔助。

#### Agent Action

1. INIT 階段若 `proceed('okhttp3.Request')` 失敗，立即 probe 實際 overload，不要只 hook shadow/ad SDK okhttp。
2. 用 10–30 秒 cold-start spawn 驗證：INIT banner 之外至少 1 條目標 host 前綴 URL。
3. 寫入 Ai-skill 只保留 probe/hook 模式；混淆類名、endpoint、header 真值留 project docs。

#### Goal / Action / Validation

- Goal: 縮短「hook 裝了但 0 流量」的排查時間。
- Action: `analysis/apk/traffic-triage.md` Java hook 小節補「混淆 Request 型別 probe」。
- Validation: 修正 overload 後同場景業務 URL 計數 > 0。

#### Applies When

- Release APK 含 OkHttp3 且啟用 R8/ProGuard
- Frida Java hook 主線 MITM 輔助捕獲

#### Does Not Apply When

- 未混淆 debug build（`okhttp3.Request` 仍可用）
- 流量完全不走 OkHttp（需改 native / Cronet / 自訂 stack triage）

#### Validation

- `getDeclaredMethods()` 列出 `proceed` 實參型別
- Hook 後 cold-start 窗口內有業務 API 日誌

#### Promotion Target

- `analysis/apk/traffic-triage.md` §Java OkHttp hook
- `workflow/apk-analysis/execution-flow.md` §動態捕獲

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
