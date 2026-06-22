> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - JNI registered dynamically — nm -D may show no Java_com_* for util natives

Status: candidate

#### One-line Summary

DEX 標記 `native` 的 util 方法（sign / encrypt / launchTime），在 protection `.so` 上 **`nm -D` 可能找不到 `Java_com_*`**（動態 `RegisterNatives`）；靜態 `.so` strings 分析不足時，應 **spawn + RegisterNatives hook** 或 **Frida `Java.use` RPC** 直接呼叫 Java native stub，而非等標準 JNI symbol。

#### Human Explanation

加固/殼（protection library）常延遲綁定 JNI，避免 IDA 直接 xref `Java_*`。attach 已運行進程時 RegisterNatives 可能已執行完畢，hook 抓不到註冊表。實務路徑：(1) Frida spawn `-f` 並 hook `art::JNI::RegisterNatives` 過濾目標 class；(2) 若 Java 層仍可調用，用 RPC export 包裝 `Util.sha256Encrypt` / `getSystemLaunchTime` 作 interim signer；(3) native RE 需從 RegisterNatives 回調的 fn pointer 反查，而非 symbol name search。

#### Trigger

- DEX: util class methods marked `native`
- `nm -D lib*.so | grep Java_com_` 無該 class 對應 export
- attach 後 hook native symbol 失敗或 `_Z6sha256` 有 hit 但與 Java 路徑對不上
- sign RE 卡在「找不到 JNI 入口」

#### Evidence

- Tool: androguard DEX native flag + `nm -D` + Frida RegisterNatives / Java.use RPC
- Sanitized excerpt: protection `.so` 有 SHA256 C++ symbols 但無 `Java_com_*` for app util; Frida Java hook on util still works
- Evidence path: `<PROJECT_ROOT>/docs/static-analysis.md`、`<PROJECT_ROOT>/api/signing-re.md`

#### Generalized Lesson

```text
Native util RE branches:
  A. nm -D has Java_com_* → static JNI RE
  B. No Java_com_* but Java.use() works → dynamic RegisterNatives
     B1. spawn + RegisterNatives hook → log (name, sig, fn_ptr)
     B2. interim: Frida RPC wrap Java native methods
     B3. trace fn_ptr in protection .so (not symbol name grep)
  attach-only RegisterNatives hook often too late → prefer spawn for B
```

#### Agent Action

1. Project static-analysis 表註明「dynamic JNI」與 protection library 名（generic）。
2. Ai-skill 提供分支判斷樹，不寫特定 `.so` 偏移。
3. 與 `141700`、`141900` 並用。

#### Goal / Action / Validation

- Goal: 避免在無 export 的 `.so` 上浪費時間 grep `Java_com_*`。
- Action: 優先 Java-layer Frida RPC；native RE 需 spawn RegisterNatives 或已載入模組 fn pointer。
- Validation: RPC 可從 host 傳 canonical 拿回 sign；或 RegisterNatives log 列出 util 方法映射。

#### Applies When

- App util natives in protection/obfuscation library
- Standard JNI symbol search fails

#### Does Not Apply When

- Clear `Java_com_*` exports in main app `.so`
- Crypto entirely in Java (no native)

#### Validation

- RegisterNatives log OR Java RPC reproduce sign for known canonical

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §native / JNI triage

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- 已依 sanitization / reusable-guidance-boundary 自查
