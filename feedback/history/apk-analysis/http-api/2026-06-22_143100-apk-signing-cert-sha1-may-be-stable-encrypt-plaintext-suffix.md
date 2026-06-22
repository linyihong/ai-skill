> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - APK signing cert SHA1 may appear as stable colon-hex encrypt suffix

Status: candidate

#### One-line Summary

加密 plaintext 裡**穩定不變的 colon-hex 後綴**，常是 **APK 簽章憑證 X.509 DER 的 SHA1**，而不是 install path、整包 apk hash 或 android_id。應先比對 cert digest，再假設其他種子。

#### Human Explanation

當 plaintext 形如 `{volatile}_{stable_colon_hex}` 且 stable 部分在重裝/換簽後才變，優先從 APK **Signing Block / META-INF PKCS#7** 取第一張 X.509 DER，格式化成大寫 colon-hex 與 suffix 比對。Native 可能從 install path 讀 apk 再解析簽章區塊；**offset、欄位名、key 值**留在專案 evidence，不寫進本 lesson。

#### Trigger

- Encrypt plaintext suffix 穩定 per install
- Suffix 為 20-byte digest 的 colon-hex 表示
- path/file hash 假設對不上

#### Evidence

- Tool: PKCS#7 cert parse + hook/strcpy plaintext capture
- Sanitized excerpt: `SHA1(cert_der)` formatted == suffix
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md`（專案具體欄位與腳本）

#### Generalized Lesson

```text
When encrypt plaintext has stable colon-hex suffix:
  1. Parse signing cert DER from APK (PKCS#7 / Signing Block path)
  2. Compare SHA1(DER) with colon formatting to suffix
  3. If match: host-side offline can use local APK with same signature
  4. Do not copy native offsets or wire field names into reusable docs
```

#### Agent Action

1. 專案 evidence 記錄完整鏈路與驗證腳本。
2. 交叉引用 `143000`（plaintext 可能不是表面 seed）。
3. Ai-skill 不寫 key、offset、endpoint。

#### Goal / Action / Validation

- Goal: 避免在 path hash 上浪費 RE 時間。
- Validation: cert SHA1 match documented in project only.

#### Applies When

- Native encrypt seeded by apk install path
- Suffix looks like MAC-colon hex (20 bytes)

#### Does Not Apply When

- Suffix from hardware serial / GAID only
- No apk signing cert linkage

#### Validation

- Project evidence documents match; reusable doc has no app-specific names

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §signing-cert fingerprint probe

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引
- 已依 sanitization / reusable-guidance-boundary 自查
