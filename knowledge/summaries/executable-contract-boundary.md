## governance.executable-contract-boundary

| 欄位 | 值 |
| --- | --- |
| Atom ID | `governance.executable-contract-boundary` |
| Source path | `governance/lifecycle/executable-contract-boundary.md`, `governance/lifecycle/executable-contract-boundary.yaml` |
| Lifecycle | `candidate` |
| Summary | 定義 executable YAML contract 邊界：source 留在 owner layer，Markdown 解釋脈絡，YAML 承載 activation / steps / gates / evidence；會影響 execution 的 contract 以 `runtime_projection.enabled: true` opt in 到 `runtime.db generated_surfaces`。Framework contract / projection 改動前必須先做 pre-build interrogation，避免雙寫 source-of-truth。 |
| When to read | 新增或修改含 steps、activation、dependencies、exit gates、blocking gates、required evidence、failure actions 的 governance / enforcement / workflow 文件，或調整 runtime projection / YAML placement / framework source-of-truth 時。 |
| Do not use for | 不要把所有 YAML 搬到 `runtime/`；不要把 ordinary metadata、graph、validation 或 philosophy YAML 投影到 runtime，除非它明確 opt in 且通過 contract schema。 |
| Validation signal | `runtime/runtime.db.generated_surfaces` 包含 opt-in contracts；contract inventory 已標記 owner source；pre-build interrogation 已確認 canonical owner、projection boundary 與 duplication risk。 |
| Last checked | 2026-05-21 |
