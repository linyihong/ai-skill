> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Frida Java.enumerateLoadedClasses causes script load timeout on large apps

Status: candidate

#### One-line Summary

大型商業 APK 上 Frida 腳本若在 `Java.perform` 內 **`Java.enumerateLoadedClasses`** 全量掃描找 obfuscated 類名，常導致 **script load timeout**；sign/crypto RE 應改用 **DEX 靜態定位類名 +  targeted `Java.use('fully.qualified.Name')`** hook。

#### Human Explanation

enumerate 在數萬 class 上迴圈，載入腳本階段即超時，表現為 Frida 無 hook 輸出。靜態 androguard/jadx 已能從 `Interceptor`、`sign` literal、util 方法名反查 R8 後類名（如單字母 package）。probe 腳本應只 hook 3–5 個已知類，必要時延遲 hook OkHttp chain。

#### Trigger

- Frida attach 後無 `[INIT]` log、script failed to load / timeout
- 腳本含 `enumerateLoadedClasses` / `enumerateClassLoaders` 全掃
- 目標 class 名已可從 DEX 取得

#### Evidence

- Tool: Frida + androguard DEX class search
- Sanitized excerpt: 移除 enumeration 後 probe 正常輸出 SIGN_IN/SIGN_OUT
- Evidence path: `<PROJECT_ROOT>/scripts/frida/probe_sign.js`

#### Generalized Lesson

```text
Frida hook strategy (large app):
  1. Static DEX → resolve obfuscated class names first
  2. Java.use('exact.name') for interceptor, util, okhttp chain
  3. Avoid enumerateLoadedClasses in initial probe
  4. Optional: enumerate only in secondary script after narrow filter
```

#### Agent Action

1. Probe 腳本模板預設 targeted hooks。
2. 若需動態發現，用 DEX strings / androguard 預過濾關鍵字再 enumerate 子集。

#### Goal / Action / Validation

- Goal: sign RE probe 第一輪就能出 log。
- Validation: script loads <10s; INIT lines present。

#### Applies When

- Large multi-dex commercial APK Frida probes
- R8 obfuscated packages

#### Does Not Apply When

- Small app / early bootstrap where enumerate is fast
- Frida `-f` spawn with minimal class load (still prefer targeted)

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §Frida probe guardrails

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
