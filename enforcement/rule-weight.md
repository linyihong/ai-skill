# 規則權重與衝突優先序

本規則定義 agent 在多個 instructions、documents、tool adapters 或 compatibility paths 看似衝突時，應如何判斷規則權重。

## 核心原則

不是每條規則權重都相同。規則競合時，選擇直接控制目前風險的最高權重規則；低權重規則只能在不削弱高權重要求時使用。

Rule weight 不是看哪個檔案最近被讀到，而是看該規則控制的 user-facing 或 repository risk。

## 權重順序

| 權重 | 規則類型 | 範例 | Agent 行為 |
| --- | --- | --- | --- |
| P0 | Safety、authorization、secrets、data loss、destructive actions、source-of-truth integrity | Authorization scope、sanitization、no secrets、未明確同意不得 destructive git、source vs mirror boundaries | 不得因方便、速度或工具限制而繞過。若 blocked，停止並詢問。 |
| P1 | 使用者最新明確要求與 active goal closure | 最新 user instruction、accepted plan、active `.agent-goals/` goal、completion criteria | 除非與 P0 衝突，否則驅動目前任務。被 redirect 時更新或暫停舊 goal。 |
| P1 | Canonical repository writeback 與 validation gates | Dependency reading、linked updates、diff review、commit/push/readback、clean status | 宣稱 repository changes 完成前必須滿足。Tool reload 或 local sync 不能取代。 |
| P2 | Cross-repo operating policy | Tool-neutral documentation、failure learning、goal/action/validation、document sizing（含 token 成本管理）、neutral language | 一致套用，但輸出長度可依任務調整。 |
| P2 | Workflow-specific execution flows 與 checklists | `workflow/<domain>/execution-flow.md`、`artifact-gates.md`、相關 README / checklist | 在 enforcement operating rules 之後套用；除非 enforcement rule 明確 delegating，否則不可覆蓋 enforcement rules。 |
| P3 | Tool adapter 與 compatibility guidance | `ai-tools/`、`.cursor/rules/`、sync scripts、symlink/bundle/copy snapshot details | 只在 active tool 或 compatibility path 適用。不得重新定義 canonical source，也不得把 optional sync 變成 universal。 |
| P3 | Efficiency 與 style preferences | Decision efficiency、output shape、formatting preferences、optional cleanup | 只在高權重要求滿足後最佳化。 |

## 衝突規則

1. **高權重勝出。**低權重 instruction 不得削弱高權重 safety、source、validation 或 user-goal requirement。
2. **同權重時具體勝過泛用。**Task-specific accepted plan 可以細化 general workflow，但不能跳過 required validation。
3. **最新使用者要求勝過 stale context。**若最新 user message 重新導向任務，更新 goal state 並跟隨新要求，除非違反 P0。
4. **Tool adapters 不定義 source truth。**工具文件可說明工具如何讀取或同步內容，但 canonical edits 仍發生在 `<AI_SKILL_REPO>`。
5. **Efficiency 不能跳過 required dependencies。**Context-saving 與 decision-efficiency rules 可選擇閱讀順序，但不可移除 required reads。
6. **Compatibility layers 保持條件式。**`symlink`、`bundle`、`copy snapshot` 只在有意使用時適用；`reference-first` 仍是預設。

## 不確定時

若 agent 無法判斷哪條規則權重較高：

1. 用一句話說明衝突。
2. 辨識每條規則控制的 risk。
3. 選擇保留 safety、canonical source、validation 與最新 user goal 的路線。
4. 只有當衝突影響 scope、permissions、destructive action 或互斥 outcomes 時才詢問使用者。

## 常見範例

| # | 情境 | 較高權重規則 | 正確行動 | Conflict Matrix |
| --- | --- | --- | --- | --- |
| 1 | Tool adapter 說要 sync local bundle，但目前 setup 是 reference-first。 | P1 canonical writeback 高於 P3 conditional compatibility。 | 不要預設跑 bundle sync。Commit/push/readback canonical repo，並標記 tool sync not applicable。 | [`conflict-001`](../metadata/rules/conflict-matrix.yaml#L56) |
| 2 | Decision efficiency 建議少讀 context，但 `dependency-reading.md` 要求讀特定依賴。 | P1 validation gate 高於 P3 efficiency。 | 讀 required dependency，或標成 blocked/not applicable；不能為了速度跳過。 | [`conflict-002`](../metadata/rules/conflict-matrix.yaml#L67) |
| 3 | 舊 `.agent-goals/` entry 指向舊工作，但最新 user message 重新導向任務。 | P1 latest user request 高於 stale goal context。 | 更新、暫停或完成舊 goal，並跟隨最新要求。 | [`conflict-003`](../metadata/rules/conflict-matrix.yaml#L78) |
| 4 | Skill workflow 建議 shortcut，但 enforcement rules 要求 sanitization 或 source/mirror checks。 | P0 safety/source integrity 高於 P2 skill workflow。 | 先套用 enforcement rule，再調整 skill workflow。 | [`conflict-004`](../metadata/rules/conflict-matrix.yaml#L89) |
| 5 | 使用者要求 destructive git action，但 repository rules 要求 explicit confirmation。 | P0 destructive-action safety 在確認前高於 P1 user goal。 | 先明確詢問確認並說明風險。 | [`conflict-005`](../metadata/rules/conflict-matrix.yaml#L100) |
| 6 | Compatibility script 可執行，但沒有 active workflow 依賴 native scan 或 local mirrors。 | P3 compatibility 保持條件式。 | 不使用 script；記錄 reference-first 已足夠。 | [`conflict-006`](../metadata/rules/conflict-matrix.yaml#L111) |
| 7 | 文件超過 300 行且混合多主題，但 agent 為了省事不想拆分。 | P2 document sizing 高於 P3 efficiency。 | 必須拆分，並更新索引和連結。拆分本身是 P2，但拆分後的連動更新是 P1。 | [`conflict-007`](../metadata/rules/conflict-matrix.yaml#L122) |
| 8 | Token 成本太高想合併小檔案，但合併後檔案超過 300 行且主題不一致。 | P2 document sizing（單一目的原則）高於 P3 token optimization。 | 不要合併。改用 summary layer 或 better routing 降低 token 成本。 | [`conflict-008`](../metadata/rules/conflict-matrix.yaml#L133) |

## Conflict Matrix

本文件的衝突規則與範例已結構化為 machine-readable 格式，存放於：

[`metadata/rules/conflict-matrix.yaml`](../metadata/rules/conflict-matrix.yaml)

該檔案包含：
- **6 條衝突規則**（對應本文件第 23-31 行）
- **12 個衝突配對**（8 個來自本文件範例 + 4 個額外配對）
- **權重對照表**（P0/P1/P2/P3 各包含哪些規則）
- 每個配對的 `rule_id` 直接對應 `metadata/rules/*.yaml`

Agent 在遇到規則衝突時，應優先查詢 conflict matrix 的 machine-readable 配對，
再回本文件查閱完整說明。

## 驗證

關閉涉及規則衝突的工作前，確認：

- 最終行動沒有繞過 P0 safety/source/secret rules。
- 最新 user request 已反映在 active goal 或 final answer。
- Required dependency reading 與 linked updates 已完成，或明確標示 not applicable。
- Tool-specific sync 或 compatibility behavior 只在有意使用時套用。
- Final response 回報實際 validation，而不是假設已合規。

← [Back to enforcement index](README.md)
