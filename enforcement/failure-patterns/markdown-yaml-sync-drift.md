# Markdown / YAML Sync Drift（改 .md 沒同步改 sibling .yaml）

Status: validated
Class: `source-of-truth-duplication` / `governance-drift`

## Trigger

當 commit stages 一個 canonical doc markdown（含 sibling `.yaml` 在同目錄同 stem，如 `governance/lifecycle/foo.md` ↔ `governance/lifecycle/foo.yaml`）但**沒同 commit stage sibling YAML**，使用此 pattern。

具體訊號：

- Staged `<path>/<stem>.md` 存在
- `<path>/<stem>.yaml` 存在於 disk
- Staged 集合**不**含 `<path>/<stem>.yaml`

## Failure Mode

Markdown 與 YAML companion 是兩個 surface 各帶角色：
- YAML = canonical executable contract（machine-readable, validator can iterate）
- Markdown = human-readable rationale + examples

改 markdown 不同步改 YAML → canonical contract drift：

1. **YAML 不知道 markdown 加了新規則** — agent / validator 依 YAML 行動，新規則沒效
2. **Reviewer 看 markdown 以為事情做了** — 實際 YAML 缺
3. **Cascade drift** — 多次 markdown-only 改後，YAML 落後好幾代

本 session 真實案例：

| Drift | 修補 commit |
|---|---|
| `command-contract.md` 漏 commit-msg / pre-push / runtime obligations | `2b106e9` |
| `system-upgrade-governance.yaml` 漏 Deferred Runtime Projection section（同 session 第二次同類失誤）| `ff941c3` |

## Required Agent Action

修改 canonical doc markdown 時：

1. 查 sibling `.yaml` 是否存在（同目錄、同 stem）
2. 存在 → 同 commit stage 兩個檔
3. Markdown 加新規則 → YAML 對應 schema field 也要加（required_sections / conditional_sections / format_template 等）
4. 純文字 / 排版調整 → 用 `[skip-markdown-yaml-sync]` opt-out + commit message 寫明 reason

## Prevention Gate

Commit-msg validator: `validateMarkdownYamlSync` in `hooks.go`。Check sibling pair existence；缺 staged → block。

Cross-path 例外（如 `plans/README.md` ↔ `governance/lifecycle/system-upgrade-governance.yaml`）目前**不**由本 validator 處理，需 reviewer 注意；Phase 7+ 可加入 explicit cross-path mapping。

## Validation

符合下列條件即此 pattern 已被防止：

- Sibling pair 修改都同 commit staged
- 跨路徑 companion（如 plans/README.md → system-upgrade-governance.yaml）由 reviewer 顯式檢查

## Source

- 2026-05-26 session：兩次同類失誤都由使用者提醒才修。Phase 6 加 `validateMarkdownYamlSync` 機械強制 sibling-pair；cross-path 列入 Phase 7 backlog。

## Related

- [`cli-doc-drift.md`](cli-doc-drift.md) — 同類 source-of-truth-duplication（impl ↔ doc）
- [`runtime-yaml-unprojected.md`](runtime-yaml-unprojected.md) — 同類 governance-drift
- [`bootstrap-entry-bloat.md`](bootstrap-entry-bloat.md) — 同類 thin-pointer 違反

## Linked Validation Scenarios

- `markdown-yaml-sync-v1`

← [Back to failure patterns](README.md)
