# Knowledge Lifecycle

`governance/lifecycle/` defines how knowledge moves from existing source files into the new AI-native layers without breaking old skill entrypoints.

## Source Of Truth Rule

Until a migration is explicitly promoted and validated, existing `skills/`, `shared-rules/`, `ai-tools/`, and `scripts/` files remain the source of truth for executable behavior.

New layer files under `analysis/`, `workflow/`, `intelligence/`, `runtime/`, `memory/`, `feedback/`, `models/`, `governance/`, `knowledge/`, and `metadata/` may act as:

- Routing surfaces.
- Candidate maps.
- Promotion targets.
- Metadata or summary surfaces.
- Governance and runtime design.

They must not silently replace old skill behavior before promotion.

## Durable Goal Boundary

Long-term lifecycle states belong in durable planning files, not in `.agent-goals/`.

| Goal type | Durable location |
| --- | --- |
| Repository roadmap, phase, migration sequence | `architecture/` |
| Layer responsibility, candidate destinations, promotion targets | Layer README files |
| Knowledge lifecycle, validation and deprecation rules | `governance/` |
| Routing, metadata and atom discovery | `knowledge/`, `metadata/`, `runtime/` |
| Active implementation work for the current conversation | `.agent-goals/` only until completion |

Before deleting an active `.agent-goals/` entry, confirm that any remaining roadmap, lifecycle, migration, promotion, deprecation or follow-up state has been written to the durable location above.

## Lifecycle States

| State | Meaning | Allowed content | Not allowed |
| --- | --- | --- | --- |
| `source-of-truth` | Current canonical behavior lives here. | Existing skill/shared rule/tool/script files. | Treating a newer map as override. |
| `candidate-map` | A map from current sources to future layer destinations. | Ownership boundaries, source-to-target tables, compatibility notes. | Bulk content migration or behavior changes. |
| `candidate-atom` | A proposed Knowledge Atom or summary. | Metadata, summary, links to source, validation criteria. | Marking as stable without use. |
| `validated-atom` | A candidate used or reviewed successfully. | Routing metadata, summary, checklist, validation evidence. | Removing old entrypoint. |
| `promoted` | New layer becomes a supported reference path. | Old entrypoint links to promoted atom, index routes to both when needed. | Deleting compatibility path without deprecation. |
| `deprecated` | Old path is being retired with a replacement. | Deprecation note, replacement link, validation record. | Breaking existing links or tool loading. |

## Cold Data Archive

Cold data archive 是 lifecycle / runtime compression 策略，不是把知識搬離 Markdown。當 lesson、summary 或 graph 數量成長到 agent 每次都需要掃大量冷資料時，應先建立 generated summary、SQLite / FTS lookup cache 或 report view，讓 agent 先查候選 source，再按需讀全文。

### Trigger Thresholds

符合任一條件時，應啟動 cold-data archive / generated lookup 檢查：

| Trigger | Action |
| --- | --- |
| 單一 skill 的 `feedback_history/` 超過約 50 條 lesson，或單一 category 超過約 20 條。 | 產生或更新 category index、summary rows 與 SQLite / FTS index。 |
| Agent 為了找 lesson 需要讀多個 `feedback_history` README 或大量 lesson 全文。 | 先用 SQLite / FTS query 找候選 `source_path`，再讀 1-3 個 canonical files。 |
| 某批 lesson 長期未修改，但仍可能被 routing / promotion 使用。 | 標為 cold lookup candidate，保留 Markdown source，產生短 summary / tags / validation signal。 |
| feedback lesson 被 promotion 到 `workflow/`、`intelligence/`、`shared-rules/` 或 runtime route。 | 保留原 lesson，更新 summary / graph / registry / SQLite lookup。 |
| 查詢成本高於判斷成本，例如只需要知道「有哪些相關 lesson」。 | 使用 generated report / SQLite index，不讀全文。 |

### Rules

- Canonical source 仍是 `skills/*/feedback_history/*.md`、`shared-rules/`、`knowledge/summaries/`、`knowledge/graphs/` 與 routing registry。
- SQLite / FTS、generated summaries、runtime reports 都是 generated lookup views；可刪除、可重建，不作唯一來源。
- Cold archive 不等於 deprecated。冷資料仍可能有效，只是預設不讀全文。
- 需要修改、promotion、debug、failure learning 或高信心判斷時，必須回到 canonical Markdown / YAML。
- 若 generated lookup 與 source 不一致，降級 lookup confidence 並重新產生，不要改 source 去迎合 cache。

### Validation

Cold-data archive 啟用或更新後，至少驗證：

1. Generated lookup 可從 clean checkout 重建。
2. Generated DB 或 cache 不進 git，除非另有明確 governance 決策。
3. Query result 包含 `source_path`、summary、status / confidence 與 validation signal。
4. Query result 只作 candidate list，不跳過 source-of-truth gate。
5. Runtime reports、SQLite counts 與 canonical summary / graph / registry counts 可交叉檢查。

## Promotion Gates

A candidate can be promoted only when all gates pass:

1. The old source path remains reachable or has a redirect note.
2. `metadata/schema.md` metadata exists for the promoted atom or surface.
3. `knowledge/indexes/README.md` routes relevant task intents.
4. The owning layer README links the new path.
5. Validation is documented in `governance/validation/`.
6. Diff review confirms no project-specific evidence, secrets, local absolute paths, or tool mirror paths were introduced.
7. Any durable roadmap or lifecycle state has been updated outside `.agent-goals/`.
8. Commit, push, readback, and clean status have completed.

## Update Strategy While Skills Still Change

When an old skill is updated before migration:

1. Update the old `skills/<name>/` source first.
2. Check whether a candidate map or promoted atom references the changed section.
3. If yes, update the map, metadata, summary, or index in the same change.
4. If no, record that no linked update was needed in the final validation.
5. Do not copy new skill text into a new layer unless the change is an explicit atom promotion.

## Deletion Rule

Do not delete or move old skill files during candidate-map or candidate-atom phases. Deletion can be considered only after:

- The replacement path is promoted.
- Existing tool adapters can still load the skill or have documented replacements.
- Links, indexes, summaries, and metadata are updated.
- A deprecation note and rollback path exist.

## Skills Deprecation Timeline

舊 `skills/` 目錄的內容不會一次性刪除，而是依以下階段逐步 deprecate：

| Phase | Condition | Action | Status |
| --- | --- | --- | --- |
| **Phase A** | 新分層已建立，舊入口仍 active | 不刪除。舊 `skills/` 維持 source of truth，新分層作為 reference / routing / promotion surface。 | ✅ 已完成 |
| **Phase B** | 所有 `techniques/` 已完成 decomposition（workflow → `analysis/`，intelligence → `intelligence/`），且 pilot 驗證通過 | 舊 technique 檔案標註 `# Deprecated — see <new path>`，但保留檔案。`skills/` 仍可被 tool adapter 載入。 | ✅ 已完成（2026-05-12） |
| **Phase C** | Pipeline 可自動從舊 skill 提取 intelligence atoms，且驗證通過 | 可開始刪除已完全覆蓋的舊 technique 檔案。刪除前需確認：<br>1. `analysis/` 有對應 workflow<br>2. `intelligence/` 有對應 atoms<br>3. `knowledge/indexes/` 可 route 到新路徑<br>4. `knowledge/runtime/routing-registry.yaml` 已更新<br>5. 無任何 tool adapter 依賴該舊路徑 | ✅ 已完成（2026-05-12） |
| **Phase D** | 所有 skill 內容已完全遷移到新分層 | 可刪除整個 `skills/` 目錄，或保留為唯讀 archive。刪除前需通過完整的 deprecation checklist（見 `architecture/ai-native-knowledge-operating-system.md` 的 Phase 3 Deprecation Checklist）。 | ⏳ 待執行 |

### 判斷是否可刪除單一舊檔案的檢查清單

- [ ] 新分層有對應的 workflow / analysis / intelligence 檔案
- [ ] 舊檔案的 `# Intelligence Extracted` 或 `# Deprecated` 標註已存在
- [ ] `knowledge/indexes/README.md` 可 route 到新路徑
- [ ] `knowledge/runtime/routing-registry.yaml` 包含新路徑
- [ ] 所有 `knowledge/summaries/` 和 `knowledge/graphs/` 已更新
- [ ] 無 tool adapter 或 `.claude/settings.json` 依賴該舊路徑
- [ ] 刪除後可 rollback（git revert）
