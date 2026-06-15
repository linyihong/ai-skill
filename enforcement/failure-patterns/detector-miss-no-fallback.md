# Detector Miss → No Mechanical Fallback（偵測未命中即完全失效）

Status: validated
Class: `meta-governance-gap` / `routing-miss`

## Trigger

當下列事實鏈成立時使用此 pattern：

1. **存在 deterministic detector** — workflow / route activation 由 `routing-registry.yaml`
   的 `activation_triggers`（user keyword、path、context signal）機械判定
2. **detector miss 時 fail-open** — 沒有命中任何 route 時，gate 刻意不 block 工具
   （pre-Read 破環依賴原則：activation 不依賴 content）
3. **miss path 上沒有 mechanical fallback** — miss 之後，唯一的「補救」是 agent 自律
   回頭找對的 workflow / governance source
4. **自律失敗即靜默失效** — agent 憑常識把任務做完，workflow primary_source 從未 Read，
   artifact-gates / governance checklist 大量缺漏且無人偵測

具體訊號：

- detector miss（無 user keyword、無 path match）但任務其實屬於某 governed route
- agent 直接產出結果，未 Read 任何 `workflow/<domain>/` primary_source
- 缺漏要到使用者多輪追問才暴露（gate 名實不符：hit 才擋，miss 全放）

## Failure Mode

這是 [`workflow-detector-missing.md`](workflow-detector-missing.md) 的下一層失效：detector
**存在**且能讀 `activation_triggers`，但 **miss path 是設計接受的容忍範圍**，而該容忍範圍
對 cross-project + project-local-ontology 任務太寬。

| Pattern | 宣稱 | 實作落差 |
|---|---|---|
| `workflow-detector-missing` | registry 宣告 activation_triggers | runtime 沒 detector 讀它 |
| **本檔（detector-miss-no-fallback）** | **detector 存在且 hit 時機械擋住** | **miss 時 fail-open，無 mechanical fallback，靠 agent 自律** |

「detector hit 機械擋住」給人「workflow activation 已機械化」的安全感，但 mechanical 層
**只完成一半**：hit path 機械、miss path 靠人。系統每接一個 detector 訊號覆蓋不到的
downstream project（project-local 命名 / ontology），就會在 miss path 上原樣復發。

## 高階 Pattern：Half-Mechanized Gate（一半機械、一半靠自律）

本檔是一個更高階治理 pattern 的 instance：

> **一個 gate 只機械化「命中」分支、把「未命中」分支留給 agent 自律，就不是 mechanical gate，
> 而是 advisory gate 偽裝成 mechanical gate。** 真正的覆蓋率取決於最弱的分支（miss path），
> 不是最強的分支（hit path）。

緩解方向不是把 miss path 也變成 hard block（那會破壞 pre-Read 破環依賴原則 + deterministic
activation 契約），而是**在 miss path 上補一個 advisory-but-mechanical 的 fallback**：
detector miss → 機械觸發 Discovery → 產 candidate route → 注入 advisory（非 gate）。
miss 不再 = 完全失效，但 advisory 也從不偷渡成 activation。

## 採樣

| # | 採樣 | miss path 行為 |
|---|---|---|
| 1 | 2026-05-31 travel-planning（parent plan 原 incident） | detector 未觸發，無 fallback，缺漏未偵測 |
| 2 | 2026-06-05 downstream project review（parent 完成隔日復發） | detector miss（project-local 命名），agent 憑常識 review，travel-planning artifact-gates 19 項缺 7~10 項，使用者三輪追問才暴露 |

兩次為**同一結構性缺口的兩次採樣**：parent plan archive 收尾 evidence 把「detector miss
為設計接受的容忍範圍」當假設；採樣 #2 證明此假設對 cross-project 任務太寬。具體 project
artifact / 檔名 / 對話細節依 [`reusable-guidance-boundary.md`](../reusable-guidance-boundary.md)
留在原 project 文件，不入本 pattern。

## Why It Recurs

不是個人疏忽。是**漸進交付的自然產物**：

1. **hit path 先機械化** — detector + PreToolUse gate 是第一個可交付的機械單元
2. **miss path 的 fallback 標 deferred** — parent plan Phase 6「Discovery → Detector feedback
   loop」明寫 hot-hook auto-call 刻意延後
3. **延後代價不在當下兌現** — 隔天才在另一個 project 復發
4. **「detector 已上線」被誤讀為「workflow activation 已機械化」** — 完成宣稱蓋過 miss path 的洞

## Required Agent Action

- 設計或審查任何 deterministic gate 時，明確區分 **hit path 行為** 與 **miss path 行為**；
  miss path 若只剩「靠 agent 自律」，標記為 *advisory gate*，不可宣稱為 mechanical 覆蓋
- detector miss 不等於「此任務無 governance」；若 miss path 缺 fallback，視為已知缺口而非容忍範圍
- 補 miss path 的 fallback 時，維持 advisory 性質（非 block），避免破壞 pre-Read 破環依賴
  與 deterministic activation 契約 — advisory ranking 可用 scoring，activation 不可

## Prevention Gate

**現有機械防護**：

- `detector.go` + `gate.workflow.primary_source_read`（hit path：single locked route 時
  block 直到 workflow primary_source 被 Read）
- **Discovery Bridge Phase A**（miss path 的 advisory fallback，2026-06 landed）：detector
  miss → Light Discovery 產 top-3 candidate route + confidence → ≥ threshold 注入 PreToolUse
  advisory（`scripts/ai-skill-cli/internal/app/discovery.go` `RunDiscoveryBridge`；regression
  `validation/scenarios/runtime/workflow-discovery-bridge-light-v1.yaml`）

**缺口（tracked follow-up）**：

- **Phase B（Deep Discovery）deferred**：Light-only 在 2026-06-05 replay 下兩個 candidate 均
  < threshold（travel 0.272 / software-delivery 0.312），proposal 寫入但未達 advisory。
  miss path 的覆蓋率仍待 Phase B（piggyback Read 的 content scan 累積）+ Phase D 三週 empirical
  驗證。見 plan `2026-06-06-1700-workflow-activation-discovery-bridge`
- Discovery 是 **advisory，不是 mechanical gate** — 刻意不登記為 `enforcement-registry.yaml`
  的 `coverage: mechanical` rule_class（理由見該 plan Phase C + `governance/workflow-activation-engine.md`
  §Discovery Bridge）

## Validation

符合下列任一即此 pattern 已被緩解：

- detector miss 時 Discovery Bridge 機械觸發並（達 threshold 時）注入 advisory，agent 得以
  回頭 Read 正確 workflow primary_source
- miss path 行為在 gate 設計文件中明確標示為 advisory（非偽裝成 mechanical 覆蓋）
- 2026-06-05 empirical trigger replay 在 Phase A（+ 未來 Phase B）下至少寫出 candidate proposal

## Source

- 2026-06-05 session：parent plan `2026-05-31-1900-workflow-activation-engine`（archived）
  Phase 6 deferred 段的失效於完成隔日復發；補強計畫
  `2026-06-06-1700-workflow-activation-discovery-bridge`（Phase A landed，Phase C closeout 2026-06-15）

## Related

- [`workflow-detector-missing.md`](workflow-detector-missing.md) — 前一層：registry 宣告但 runtime 沒 detector
- [`rule-without-executor.md`](rule-without-executor.md) — 同公理：宣稱必須可機械驗證，否則治理只是裝飾
- [`governance/workflow-activation-engine.md`](../../governance/workflow-activation-engine.md) — §Discovery Bridge（miss path advisory fallback）
- [`governance/lifecycle/capability-discovery-philosophy.md`](../../governance/lifecycle/capability-discovery-philosophy.md) — Discovery → Detector feedback loop 哲學

← [Back to failure patterns](README.md)
