# Sanitization Rules Slice

> **Cognitive Slice**：`apk-sanitization`（從 [`../artifact-gates.md`](../artifact-gates.md) §12 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../../governance/cognitive-slice-taxonomy.md) §7.5）。

| slice 欄位 | 值 |
|---|---|
| `id` | `apk-sanitization` |
| `purpose` | 規範哪些值必須遮蔽、哪些可以保留，避免敏感資訊外洩 |
| `type` | `execution` |
| `tags` | artifact-gate, sanitization |
| `load_when` | 任何要對外公開或 commit 的 evidence / sample 前的去敏檢查 |
| `do_not_load_when` | 純內部探索、尚未準備輸出 / commit |
| `owner_layer` | workflow |
| `layer_justification` | 規定「哪些 token / key / 個資要遮蔽，哪些 schema / header name 可保留」的 ordering / artifact gate；通過 workflow membership test |
| `canonical_source` | 本檔（原 `artifact-gates.md` §12 Sanitization Rules） |
| `dependencies` | `apk-evidence-chain`（evidence 需先過此 slice 才能 commit） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | 任何 evidence commit / publish 前應載入本 slice |

## 12. Sanitization Rules

必須遮蔽：

- `Authorization`、cookie、session token。
- device id、install id、advertising id。
- 真實帳號、電話、email、邀請碼。
- AES/HMAC key material。
- 能直接重放付費內容或個人內容的 URL。
- 本機絕對路徑、使用者名稱、私有工作目錄、clone 位置。請改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 等 placeholder。

可以保留：

- header 名稱。
- path shape。
- query key 名稱。
- response top-level schema。
- schema-only JSON 摘要：字串 length/hash、top-level keys、key type；不要保留 value。
- 已去敏的 fixture。
- magic bytes、容器格式、演算法步驟。

---

← [回到 artifact-gates 索引](../artifact-gates.md) | [workflow/apk-analysis/](../README.md)
