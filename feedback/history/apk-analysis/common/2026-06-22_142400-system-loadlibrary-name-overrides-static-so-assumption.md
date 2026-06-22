> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - System.loadLibrary short name overrides assumed protection .so from static grep

Status: candidate

#### One-line Summary

靜態分析若把「含 SHA256 / protection 字樣的 `.so`」當成 sign 實作庫，可能 **指錯模組**；必須以 DEX 中 util class `<clinit>` 的 **`System.loadLibrary("…")` 短名** 為準，再 `Process.findModuleByName("lib….so")` / spawn JNI hook 驗證實際載入庫。

#### Human Explanation

APK 常同時含多個 native 庫：通用加固殼、廣告 SDK、業務 crypto。`strings` / `nm` 在 `libprotect*.so` 見 SHA256 符號不代表 Java util 的 JNI 綁在那裡。實際載入鏈是 `loadLibrary("stupid")` → `libstupid.so`。錯庫會導致：RegisterNatives hook 無 util 方法、symbol offset 對不上、RE 浪費在無關 `.so`。正確 triage：(1) jadx/androgrep 找 util class static initializer；(2) 記錄 `loadLibrary` 短名；(3) Frida `Module.load` / spawn 時只 hook **該** `lib<name>.so`；(4) 靜態 `static-analysis.md` 標註「assumed vs verified load library」。

#### Trigger

- Sign RE 卡在「protection `.so` 有 SHA256 但 JNI/Java RPC 路徑對不上」
- `nm -D libprotect*.so` 有 crypto 但 `Java.use(util).sha256Encrypt` 來自另一模組
- RegisterNatives log 顯示 fn_ptr 不在先前假設的 `.so` 位址空間
- 多個 `.so` 均含 `SHA256` / `encrypt` strings

#### Evidence

- Tool: DEX `loadLibrary` grep + Frida `Process.enumerateModules()` + RegisterNatives spawn log
- Sanitized excerpt: util `<clinit>` loads `libX.so`; earlier RE targeted `libY.so` with similar strings
- Evidence path: `<PROJECT_ROOT>/docs/static-analysis.md`、`<PROJECT_ROOT>/api/signing-re.md`

#### Generalized Lesson

```text
Native crypto module identification:
  1. PRIMARY: DEX System.loadLibrary short name on the util class that owns natives
  2. VERIFY: Frida module list + RegisterNatives fn_ptr module base
  3. SECONDARY: strings/nm on candidate .so — only after (1)(2) lock module
  4. Do not equate "largest protection .so" with "sign .so"
```

#### Agent Action

1. Project `static-analysis.md` / `signing-re.md` 寫 **verified** load library 名。
2. Ai-skill 只寫判斷樹，不寫特定 `.so` 檔名。
3. 與 `142000`（dynamic JNI）並用：鎖對模組後再 spawn RegisterNatives。

#### Goal / Action / Validation

- Goal: 避免在錯誤 `.so` 上做 offset/symbol RE。
- Action: loadLibrary 優先於 strings 啟發式。
- Validation: RegisterNatives 或 JNI hook 的 fn_ptr 落在 `lib<loadLibrary>.so` 位址範圍。

#### Applies When

- Multiple native libs; util natives in DEX
- Static phase guessed protection library for crypto

#### Does Not Apply When

- Single `.so` APK or clear `Java_com_*` in app-owned library only
- Crypto entirely in Java

#### Validation

- Documented loadLibrary name + Frida module match for util JNI

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §native module triage

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
