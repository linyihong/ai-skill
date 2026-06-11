# Sanitization Leak on Canonical Write（把專案私有 token 寫進可重用層）

Status: validated（2026-06-06 commit 214a415 + sibling plan v1–v4 累積 leak）
Class: `reusable-guidance-boundary-violation` / `sanitization-miss`

## Trigger

在對 **shared / reusable 層**（`plans/`、`enforcement/`、`governance/`、`workflow/`、`knowledge/`、`runtime/`、`validation/`）做 canonical write 時，把某個專案在自身 metadata 宣告為 private 的 entity token（codename / client / customer / product / individual）寫入內容。

具體觸發訊號：

- 某 downstream / project-local 的具體名稱、客戶名、產品 codename 出現在 shared-layer plan / rule / workflow 文件
- `enforcement/sanitization.md` + `enforcement/reusable-guidance-boundary.md` 明文禁止的「project incident 具體名稱」進入可重用文件
- 同一個私有 token 在 project-local（`.agent-goals/`，`shared_layer:false`）合法，但被複製進 shared layer
- email / 電話 / OS 絕對路徑 / 憑證前綴等 PII / secret pattern 出現在 staged 可重用內容

## Failure Mode

「私有 token 進入可重用層」→ 該 token 隨 framework 散佈給所有 downstream consumer。導致：

1. **Silent leak**：`sanitization.md` 是 prose rule，無 mechanical executor 時純靠 agent 自律；自律失效即 leak
2. **延遲偵測**：要等使用者手動 inspect 才發現（214a415 即如此，修正於 728282c）
3. **每次只修個案**：發現後只 sanitize 當前這幾處，不防下一次
4. **跨 project 汙染**：framework 端散佈了 A 專案的私有名稱給 B 專案

## Empirical Evidence

| # | 事件 | 說明 | 處置 |
|---|---|---|---|
| 1 | commit 214a415（2026-06-06） | 某 downstream project label 七處寫入 `plans/active/` shared layer（一份 v0 草稿） | 使用者手動指認 → sanitize 於 728282c |
| 2 | sibling plan `2026-05-31-2000` v1–v4 | allowlist 路線成熟前的 cold-start window 累積 leak | 由本 metadata-derived 路線取代 |

## Risk

- **可重用性信任崩塌**：consumer 發現 framework 文件含他人私有名稱
- **Reusable-guidance-boundary 形同虛設**：規則寫了卻擋不住
- **Sanitization rule_class 假性 coverage**：registry 顯示有規則，實際無 enforcement

## Executor Binding（mechanical）

| 層 | Artifact |
|---|---|
| Rule (prose) | `enforcement/sanitization.md` + `enforcement/reusable-guidance-boundary.md` |
| Forbidden tokens（metadata-derived） | 各 project `.ai-skill-project.yaml` `private_entities[*].match_tokens` → `runtime.db.derived_match_tokens` |
| Pattern（regex） | `runtime/sanitization-patterns.yaml` → `runtime.db.sanitization_patterns` |
| Scope classification | `runtime/repository-topology.yaml`（shared-layer 判定） |
| Executor | `scripts/ai-skill-cli/internal/app/sanitization_scan.go::validateSanitizationStagedContent`（pre-commit，block）|
| Registry | `enforcement/enforcement-registry.yaml` `rule_classes[sanitization]` |

## 與 Failure Authority 的關係

Scanner 對 staged 可重用內容 block；但「malformed metadata 是否有資格阻塞 compile」由 Failure Authority invariant 治理（見 [`governance/lifecycle/governance-pattern-library-draft.md`](../../governance/lifecycle/governance-pattern-library-draft.md) §Failure Authority）。本 leak 屬 shared-layer authoritative 範圍，block 合理。

## Validation Scenarios

- `validation/scenarios/runtime/sanitization-metadata-derived-fail-v1.yaml`（214a415 reconstruction）
- `validation/scenarios/runtime/sanitization-metadata-derived-pass-v1.yaml`
- `validation/scenarios/runtime/sanitization-placeholder-allowed-v1.yaml`
