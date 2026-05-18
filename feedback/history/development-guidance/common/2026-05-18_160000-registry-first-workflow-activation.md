> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-18 — Workflow activation 應 registry-first，勿每個 workflow 一列 activation-table

Status: promoted

#### One-line Summary

Workflow 觸發條件與依賴應寫在 `routing-registry.yaml` 各 `route.workflow.*` 的 `activation_triggers` / `required_dependencies`；`activation-table` 只保留 **#27 通用閘門** + §Workflow Discovery SOP。

#### Human Explanation

若為每個 workflow（開發、APK、greenfield、travel…）在 activation-table 各開一列（#27、#28、#29…），會與 `workflow-routing.md`、`routing-registry.yaml` 三重維護，且新增 workflow 時表格無限膨脹。

正確分工：

- **#27**：是否進入 workflow 編排世界（通用閘門）。
- **registry**：各 workflow 的觸發與必讀依賴。
- **workflow-routing.md**：多 route 同時命中時的歧義裁決。

#### Trigger

- 討論「察覺開發要不要強制 routing discovery」時，在 activation-table 為 software-delivery / apk-analysis 各寫專向列。
- 使用者指出 workflow 變多後不應持續新增 activation 編號。

#### Evidence

- Sanitized excerpt: `activation-table #27+#28` duplicated triggers already present under `route.workflow.software-delivery.activation_triggers`.
- Evidence path: `<AI_SKILL_REPO>/knowledge/runtime/routing-registry.yaml`, `runtime/router/activation-table.md`.

#### Generalized Lesson

1. 新增 `route.workflow.*` 時，在 registry 補 `activation_triggers` 與 `required_dependencies`。
2. 不要為每個 workflow 新增 activation-table 專向列。
3. 多 route 命中時用 `workflow-routing.md` §歧義裁決，不靠再加 activation 列。
4. `activation_table_ref` 等 registry → table 反向引用應避免。

#### Agent Action

- 命中 workflow 任務時：先 #27 / §Discovery → 掃描 registry `activation_triggers` → 載入選定 route 的 `primary_source` + `required_dependencies`。
- 勿只載入單一 intelligence 檔（如 docs-first）就寫可觀察產品碼。
- 擴充 workflow 時改 registry，並更新 workflow-routing 歧義表（若需要）。

#### Applies When

- 任何需要 `route.workflow.*` 編排的任務。
- 設計或修改 activation-table、routing-registry、workflow-routing 時。

#### Promotion Target

- `governance/lifecycle/routing-philosophy.md`
- `runtime/router/activation-table.md`
- `workflow/workflow-routing.md`
- `knowledge/runtime/routing-registry.yaml`
- `enforcement/dependency-reading.md`
