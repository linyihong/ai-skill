### 2026-05-21 — 流程型治理與規則文件需要 executable YAML contract boundary

#### One-line Summary

流程、gate、activation、blocking condition、required evidence 或 failure action 若會影響 agent 執行，必須提供 owner-layer YAML contract，並投影到 `runtime.db`。

#### Human Explanation

只把規則寫在 Markdown 裡，agent 容易讀到部分段落或 sub-pipeline 後跳過 master flow。把所有 YAML 都搬進 `runtime/` 又會破壞 owner boundary。正確模式是 Markdown 保留說明，YAML contract 放在原 owner layer，runtime compiler 只投影 opt-in contract 到 `runtime.db`。

#### Trigger

使用者要求整理「enforcement / governance / workflow 的 YAML 要放哪裡、哪些要進 runtime、以後更新怎麼遵守」的框架規則。

#### Evidence

- `governance/lifecycle/knowledge-update-flow.yaml` 已能阻止 agent 用 sub-pipeline 取代 master flow。
- `runtime/README.md` 已規定 runtime internal config canonical 在 `runtime.db`，不保留 `runtime/**/*.yaml` mirror。

#### Generalized Lesson

Source ownership and runtime projection must be separate. Executable contracts stay in their owner layer and opt into runtime projection with `runtime_projection.enabled: true`.

#### Agent Action

新增或修改可執行流程文件時：

1. 判斷是否含 step、trigger、dependency、exit gate、blocking gate、required evidence 或 failure action。
2. 若有，建立 companion YAML contract。
3. 在 YAML 中設定 `runtime_projection.enabled: true`。
4. 跑 runtime compile / refresh / validate。
5. 若暫不建立 YAML，明確記錄 not applicable 或 linked-update gap。

#### Goal / Action / Validation

- Goal: 防止流程型文件只靠 Markdown 導致 agent 跳步。
- Action: 建立 executable contract boundary 與候選 inventory。
- Validation: `runtime.db.generated_surfaces` 必須包含 opt-in executable YAML contract。

#### Applies When

- governance / enforcement / workflow 文件定義可執行步驟或 blocking gate。
- metadata rule 需要被 runtime activation 或 validation 使用。
- agent 因 prose 流程太長而跳過必要步驟。

#### Does Not Apply When

- 文件只是哲學、背景、ADR、非執行型導航或歷史說明。
- YAML 是 graph、summary、validation scenario 或 metadata，且沒有 runtime execution effect。

#### Validation

- `governance/lifecycle/executable-contract-boundary.yaml` 存在。
- `runtime_projection.enabled: true` 的 contract 被 runtime compiler 投影。
- `runtime validate` 通過。

#### Promotion Target

- `governance/lifecycle/executable-contract-boundary.md`
- `governance/lifecycle/executable-contract-boundary.yaml`
- `scripts/ai-skill-cli/internal/app/runtime_compiler.go`

#### Required Linked Updates

- `governance/README.md`
- `governance/lifecycle/compiler-philosophy.md`
- `runtime/README.md`
- `knowledge/runtime/routing-registry.yaml`
- `knowledge/summaries/README.md`
- `knowledge/summaries/executable-contract-boundary.md`
- `runtime/runtime.db`
