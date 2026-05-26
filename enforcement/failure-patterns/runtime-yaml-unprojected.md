# Runtime YAML Unprojected（`runtime/*.yaml` 沒設 projection）

Status: validated
Class: `source-of-truth-duplication` / `governance-drift`

## Trigger

當新建或修改 `runtime/*.yaml` 但**沒設** `runtime_projection.enabled: true` 或缺 `target_key`，使用此 pattern。

具體訊號：

- Staged `runtime/<x>.yaml` 但 grep 不到 `runtime_projection:` 或 `enabled: true`
- Staged `runtime/<x>.yaml` 但 grep 不到 `target_key:`
- Compiler silent-skip 該 YAML（generated_surfaces 沒對應 row）

## Failure Mode

Compiler 自動 walk `runtime/`，但只 project 有 `enabled: true` 的 YAML。沒 enabled 的 YAML **被 silent skip** — 沒有任何錯誤訊息，agent / validator 也無法從 runtime.db 查到。導致：

1. **規則寫了但不生效** — agent 以為加了 contract，實際 generated_surfaces 沒 row
2. **Drift accumulates** — 多個 runtime/*.yaml 之中可能混入 unprojected ones，靠 audit 才發現
3. **Validator 無法 individual check** — 缺 enabled 的 YAML 無法被 per-obligation dispatcher 載入

本 session 真實案例：使用者提醒「runtime 裡的 YAML 都要進 runtime.db」才執行 audit；雖然 11/11 當前 PASS，但**規則本身沒機械強制**過。

## Required Agent Action

建立 / 修改 `runtime/*.yaml` 時：

1. 必含：
   ```yaml
   runtime_projection:
     enabled: true
     target_key: runtime.<unique-key>
     surface: generated_surfaces
   ```
2. Compile 後 verify：`sqlite3 runtime/runtime.db "SELECT 1 FROM generated_surfaces WHERE target_key='runtime.<key>';"`
3. 若是 plan-declared 外放→收斂模式（暫不 project），plan 必含 §Deferred Runtime Projection 宣告 (a) 不 project 的 reason、(b) 預定 project 的 phase；commit 用 `[skip-runtime-yaml-projection]` opt-out

## Prevention Gate

Commit-msg validator: `validateRuntimeYamlProjects` in `hooks.go`。Read 每個 staged `runtime/*.yaml`，缺 enabled / target_key → block。

Plan template rule: [`plans/README.md`](../../plans/README.md) §Plan 模板必填章節 row `Deferred Runtime Projection`。

Canonical contract: [`governance/lifecycle/system-upgrade-governance.yaml`](../../governance/lifecycle/system-upgrade-governance.yaml) §conditional_sections。

## Validation

符合下列條件即此 pattern 已被防止：

- 每個 `runtime/*.yaml` 都 projected（generated_surfaces 有對應 row）
- 例外案例都有 plan §Deferred Runtime Projection 宣告 + commit opt-out

## Source

- 2026-05-26 session：使用者提醒「runtime YAML 都要進 runtime.db」rule 沒機械強制；當前 audit PASS 但需要 enforcement 防未來。Phase 6 加入 `validateRuntimeYamlProjects` + plan template §Deferred Runtime Projection 規則。

## Related

- [`governance/lifecycle/system-upgrade-governance.yaml`](../../governance/lifecycle/system-upgrade-governance.yaml) — conditional_sections canonical
- [`plans/README.md`](../../plans/README.md) — plan template rule
- [`cli-doc-drift.md`](cli-doc-drift.md) — 同類 source-of-truth-duplication

## Linked Validation Scenarios

- `runtime-yaml-projects-v1`

← [Back to failure patterns](README.md)
