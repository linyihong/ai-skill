> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - requestTime device fingerprint is colon-formatted SHA1 of APK signing cert DER

Status: candidate

#### One-line Summary

`requestTime` plaintext suffix（colon 分隔 20-byte hex）通常是 **APK 簽章憑證 X.509 DER 的 SHA1**，不是 apk path 雜湊、android_id 或整包 apk hash。Native slow-path 從 install path 解析 **APK Signing Block** 取 cert。

#### Human Explanation

`143000` 已確認 plaintext 為 `{millis}_{fingerprint}`。進一步對照：fingerprint 等於 `SHA1(pkcs7_cert_der)` 轉大寫 colon-hex，與裝置上 `base.apk` 及專案 bundled apk 一致。Native `0x774dc` 在 zip 中搜尋 EOCD / `APK Sig Block 42`，與 Android v2 簽章區塊讀法一致。離線 SDK 可用本機 apk（同簽章）算 fingerprint，不必讀裝置上的 `/data/app/...` 路徑。

#### Trigger

- Plaintext suffix 穩定 per install、隨 apk 重裝/換簽章而變
- Hook encrypt 看到 colon-hex 20-byte 格式
- 離線用 path/md5/sha1(apk file) 對不上 suffix

#### Evidence

- Tool: Python `cryptography` PKCS7 cert parse + Frida plaintext capture
- Sanitized excerpt: cert DER SHA1 colon string == encrypt plaintext suffix
- Evidence path: `<PROJECT_ROOT>/api/signing-re.md` §requestTime

#### Generalized Lesson

```text
When requestTime/plaintext has stable colon-hex suffix:
  1. Parse META-INF/*.RSA PKCS#7 → first X.509 DER
  2. Compare SHA1(DER) formatted with colons to suffix
  3. If match: offline fingerprint needs only apk (not device path on host)
  4. Native may read Signing Block from install path — document, don't hardcode offsets in Ai-skill
```

#### Agent Action

1. 更新 project signing-re fingerprint 表。
2. 交叉引用 `143000`。
3. 提供 `apk_fingerprint.py` 類 helper。

#### Goal / Action / Validation

- Goal: 離線組 plaintext 前半段不需 Frida。
- Validation: bundled apk 與 device apk 指紋一致；offline encrypt 前綴可本地生成。

#### Applies When

- requestTime / native blob seeded by apk install path
- Suffix looks like MAC-colon hex but length is 20-byte digest

#### Does Not Apply When

- Fingerprint from hardware serial only
- Plain UUID without cert linkage

#### Validation

- Documented cert SHA1 match in project evidence

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §APK signing cert fingerprint

#### Required Linked Updates

- `feedback/history/apk-analysis/README.md` 索引追加
- `143000` cross-ref
- 已依 sanitization / reusable-guidance-boundary 自查
