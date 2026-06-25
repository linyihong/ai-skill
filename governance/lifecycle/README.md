# 知識生命週期（Knowledge Lifecycle）

`governance/lifecycle/` 定義知識如何從既有來源檔案遷移到新的 AI-native 分層，並在舊入口完成遷移後安全退役。

## 真相來源規則（Source Of Truth Rule）

在遷移被明確提升（promote）並驗證之前，既有的 `enforcement/`、`ai-tools/` 和 `scripts/` 檔案仍然是可執行行為的真相來源。舊 `skills/` scaffold 已退役；執行流程入口改由 `workflow/`，分析方法改由 `analysis/`，判斷智慧改由 `intelligence/` 承接。

`analysis/`、`workflow/`、`intelligence/`、`runtime/`、`memory/`、`feedback/`、`models/`、`governance/`、`knowledge/` 和 `metadata/` 下的新分層檔案可以作為：

- 路由表面（Routing surfaces）
- 候選地圖（Candidate maps）
- 提升目標（Promotion targets）
- 中繼資料或摘要表面（Metadata or summary surfaces）
- 治理與執行期設計（Governance and runtime design）

它們在提升之前不得靜默取代舊 skill 行為。

## Framework 改動前需求拷問

修改 Ai-skill framework、governance、runtime、workflow、metadata、validation、schema、generated artifact 或 tool adapter 前，必須先執行 [`../../workflow/software-delivery/requirements/pre-build-interrogation.md`](../../workflow/software-delivery/requirements/pre-build-interrogation.md)。這是 plan / implementation 前的 blocking gate，用來避免尚未釐清 source-of-truth 就開始 migration。

最低檢查：

| 面向 | 必須回答 |
| --- | --- |
| Goal / scope | 這次 framework 改動要防止哪個 failure，哪些 layer 在範圍內，哪些不做？ |
| Canonical source | 真正要改的是 owner Markdown/YAML、SQLite canonical document、compiler source、generated report 還是 tool adapter？ |
| Projection boundary | `runtime.db`、`generated_surfaces`、SQLite index、reports 或 tool config 是 source 還是 projection？ |
| Duplication risk | 是否會留下兩份 rule body、兩條 activation path、mirror、compatibility table 或 stale generated output？ |
| Close-loop | 哪些 README、routing registry、contract inventory、validation scenario、runtime compile / refresh / validate 或 tool sync 必須同步？ |

若任一答案未知且會影響 execution semantics，先問使用者或停在 planning，不得直接進入 implementation。

## Executable Contract Inventory

流程或 gate 是否需要 companion YAML，以 [`executable-contract-boundary.md`](executable-contract-boundary.md) 與 [`executable-contract-inventory.yaml`](executable-contract-inventory.yaml) 為準。Inventory 使用：

- `contract_exists`：已有 owner-layer executable YAML contract。
- `contract_required`：含 ordered steps、required reads、blocking gates、required evidence、failure action 或 final report，需補 companion YAML。
- `markdown_only`：哲學、背景、tradeoff、設計理由或索引，不投影 runtime。
- `not_applicable`：template、example、deprecated stub 或非 owner source。

Lifecycle executable contracts include [`directory-structure-governance.yaml`](directory-structure-governance.yaml), [`system-upgrade-governance.yaml`](system-upgrade-governance.yaml), [`knowledge-update-flow.yaml`](knowledge-update-flow.yaml), [`executable-contract-boundary.yaml`](executable-contract-boundary.yaml), and [`executable-contract-inventory.yaml`](executable-contract-inventory.yaml).

## Governance Pattern Template（治理模式正向模板）

[`governance-pattern-template.md`](governance-pattern-template.md) 是新增 **mechanical governance 子系統**（schema 檢查、reference resolvability、drift 偵測、coverage gate 等可由程式碼偵測違規的規則）時的正向模板。它是 `enforcement/failure-patterns/` 的建設性對偶：failure pattern 捕捉 anti-pattern，本模板捕捉應有的正向形狀。

核心命題：治理模式不是固定的六步 pipeline，而是 **invariant core + justified omissions**：

- **Invariant core（缺一即缺陷）**：Observation → Registry → Executor → Validation。
- **Conditional（僅在 predicate 成立時省略，且須記錄理由）**：Rule（near-universal；唯 pure structural invariant 可省）、Projection（direct-consumption executor 可省）。

證據紀錄（7 個樣本、兩條 falsifiable predicate）保留在 [`governance-pattern-library-draft.md`](governance-pattern-library-draft.md)：讀模板取契約，讀 draft 取「為何某步 invariant / 為何可省」的證明。兩個尚在 incubation 的 sibling family（Reference Integrity、Failure Authority）仍記在 draft，各有獨立 N≥5 gate，未 promote。

## 持久目標邊界（Durable Goal Boundary）

長期生命週期狀態應放在持久規劃檔案中，而不是 `.agent-goals/`。

| 目標類型 | 持久位置 |
| --- | --- |
| 儲存庫路線圖、階段、遷移順序 | `architecture/` |
| 分層職責、候選目的地、提升目標 | 各分層 README 檔案 |
| 知識生命週期、驗證與棄用規則 | `governance/` |
| 路由、中繼資料與 atom 發現 | `knowledge/`、`metadata/`、`runtime/` |
| 當前對話的實作工作 | `.agent-goals/`（僅到完成為止） |

刪除 active `.agent-goals/` 條目前，請確認任何剩餘的路線圖、生命週期、遷移、提升、棄用或後續狀態已寫入上述持久位置。

## 生命週期狀態（Lifecycle States）

| 狀態 | 意義 | 允許的內容 | 不允許的內容 |
| --- | --- | --- | --- |
| `source-of-truth` | 當前標準行為在此 | 既有的 shared rule/tool/script 檔案與已 promoted 的 workflow/analysis/intelligence source | 將較新的地圖視為覆蓋 |
| `candidate-map` | 從當前來源到未來分層目的地的地圖 | 擁有權邊界、來源到目標對應表、相容性說明 | 大量內容遷移或行為變更 |
| `candidate-atom` | 提議的 Knowledge Atom 或摘要 | 中繼資料、摘要、來源連結、驗證標準 | 未經使用就標記為穩定 |
| `validated-atom` | 已成功使用或審查的候選項目 | 路由中繼資料、摘要、檢查清單、驗證證據 | 移除舊入口 |
| `promoted` | 新分層成為支援的參考路徑 | 舊入口連結到提升後的 atom，索引在需要時路由到兩者 | 未經棄用就刪除相容路徑 |
| `deprecated` | 舊路徑正在退役，有替代方案 | 棄用說明、替代連結、驗證記錄 | 未標明替代入口就破壞仍 active 的工具載入 |

## 冷資料歸檔（Cold Data Archive）

冷資料歸檔是生命週期/執行期壓縮策略，不是把知識搬離 Markdown。當 lesson、summary 或 graph 數量成長到 agent 每次都需要掃大量冷資料時，應先建立 generated summary、SQLite / FTS lookup cache 或 report view，讓 agent 先查候選 source，再按需讀全文。

### 觸發門檻（Trigger Thresholds）

符合任一條件時，應啟動冷資料歸檔 / generated lookup 檢查：

| 觸發條件 | 動作 |
| --- | --- |
| 單一 domain 的 `feedback/history/<domain>/` 超過約 50 條 lesson，或單一 category 超過約 20 條 | 產生或更新 category index、summary rows 與 SQLite / FTS index |
| Agent 為了找 lesson 需要讀多個 `feedback_history` README 或大量 lesson 全文 | 先用 SQLite / FTS query 找候選 `source_path`，再讀 1-3 個 canonical files |
| 某批 lesson 長期未修改，但仍可能被 routing / promotion 使用 | 標為 cold lookup candidate，保留 Markdown source，產生短 summary / tags / validation signal |
| feedback lesson 被 promotion 到 `workflow/`、`intelligence/`、`enforcement/` 或 runtime route | 保留原 lesson，更新 summary / graph / registry / SQLite lookup |
| 查詢成本高於判斷成本，例如只需要知道「有哪些相關 lesson」 | 使用 generated report / SQLite index，不讀全文 |

### 規則（Rules）

- Canonical source 仍是 `feedback/history/*/*.md`、`enforcement/`、`knowledge/summaries/`、`knowledge/graphs/` 與 routing registry。
- SQLite / FTS、generated summaries、runtime reports 都是 generated lookup views；可刪除、可重建，不作唯一來源。
- 冷資料歸檔不等於棄用。冷資料仍可能有效，只是預設不讀全文。
- 需要修改、promotion、debug、failure learning 或高信心判斷時，必須回到 canonical Markdown / YAML。
- 若 generated lookup 與 source 不一致，降級 lookup confidence 並重新產生，不要改 source 去迎合 cache。

## Memory Promotion / Pruning Boundary

Memory 是 historical replay archive，不是 active execution state。`memory/working/` 的內容只有在 compression、qualification、abstraction 與 contamination check 後，才可 promotion 到 `memory/summary/`、`memory/episodic/`、`memory/project/`、`memory/decision/` 或 `memory/failure/`。

Promotion 到更穩定 layer 時遵守：

| Destination | Gate |
| --- | --- |
| `knowledge/` | 已抽象成 reusable navigation / summary / graph，且不依賴單一 incident。 |
| `intelligence/` | 已抽象成可重用 reasoning / heuristic。 |
| `workflow/` | 已成為可重複流程或 artifact gate。 |
| `enforcement/` | 有 recurring validated failure，且需要可執行 policy。 |

禁止 promotion：

- Raw transcript。
- Temporary blocker。
- Old `.agent-goals/` owner、lock、next action 或 active blocker。
- Unstable execution graph。
- Unresolved contradiction。
- Project-secret / private evidence。

Memory pruning 時若仍需保留經驗，只保留 generalized lesson、compatibility scope 與 source revalidation note。

### 驗證（Validation）

冷資料歸檔啟用或更新後，至少驗證：

1. Generated lookup 可從 clean checkout 重建。
2. Generated DB 或 cache 不進 git，除非另有明確 governance 決策。
3. Query result 包含 `source_path`、summary、status / confidence 與 validation signal。
4. Query result 只作 candidate list，不跳過 source-of-truth gate。
5. Runtime reports、SQLite counts 與 canonical summary / graph / registry counts 可交叉檢查。

## 目錄結構治理（Directory Structure Governance）

新增或改名目錄前，應執行 [`directory-structure-governance.md`](directory-structure-governance.md) 定義的 5 步驟 Checkpoint：

1. **Name Conflict Check** — 新目錄名稱是否與其他層同名？例如 `intelligence/engineering/analysis/` 與根目錄 `analysis/` 同名，但內容本質不同。
2. **Boundary Clarity Check** — 目錄邊界是否清晰？Agent 能否明確判斷某份文件該放哪？
3. **Inertial Naming Check** — 命名是否只是沿用舊技能名稱（如 `apk-analysis` → `analysis`），沒有反映內容本質？
4. **Path Depth Check** — 目錄深度是否合理？是否過度嵌套？
5. **Global Reference Impact Assessment** — 改名會影響多少外部檔案？是否有完整的 linked updates 計畫？

此流程應在知識更新流程（[`knowledge-update-flow.md`](knowledge-update-flow.md)）Step 1 觸發。

## 提升關卡（Promotion Gates）

候選項目只有在所有關卡通過時才能提升：

1. 舊來源路徑仍然可達，或有重新導向說明。
2. 提升後的 atom 或 surface 有 `metadata/schema.md` 中繼資料。
3. `knowledge/indexes/README.md` 可路由相關任務意圖。
4. 所屬分層的 README 連結了新路徑。
5. 驗證記錄在 `governance/validation/` 中。
6. Diff review 確認沒有引入專案特定證據、機密、本機絕對路徑或工具 mirror 路徑。
7. 任何持久的 roadmap 或 lifecycle 狀態已在 `.agent-goals/` 之外更新。
8. 已完成 commit、push、readback 和 clean status。

## 舊技能入口退役後的救援策略（Rescue Strategy After Skills Retirement）

舊 `skills/` scaffold 已退役。若後續工作仍指向舊 skill 路徑，**不應重建或更新 `skills/<name>/` 來源**；應直接定位對應的新分層 source。

正確的救援流程：

1. **直接更新新分層的目標檔案** — 找出該能力對應的 `workflow/`、`analysis/`、`intelligence/` 目標，將變更內容寫入新分層。
2. **檢查 candidate map 或 promoted atom** — 確認是否有 map 或 atom 參考了變更的章節。
3. **同步更新 map、metadata、summary 或 index** — 如果有參考，在同一變更中更新。
4. **記錄不需要連動更新** — 如果沒有參考，在最終驗證中記錄即可。
5. **不要複製到舊路徑** — 除非該 atom 尚未被 promote，否則不要將新內容複製回 `skills/`。
6. **不要恢復舊 scaffold** — 若需要相容說明，更新 routing / README / lifecycle 文件，而不是新增 `skills/` 檔案。

## 刪除規則（Deletion Rule）

在 candidate-map 或 candidate-atom 階段不要刪除或移動舊 source。只有在以下條件滿足時才能考慮刪除：

- 替代路徑已被提升。
- 既有的 tool adapter 有文件化的替代方案。
- 連結、indexes、summaries 和 metadata 已更新。
- 存在棄用說明和復原路徑。

## 技能棄用時間表（Skills Deprecation Timeline）

舊 `skills/` 目錄已依以下階段完成 deprecation：

| 階段 | 條件 | 動作 | 狀態 |
| --- | --- | --- | --- |
| **Phase A** | 新分層已建立，舊入口仍 active | 不刪除。舊 `skills/` 維持 source of truth，新分層作為 reference / routing / promotion surface。 | ✅ 已完成 |
| **Phase B** | 所有 `techniques/` 已完成 decomposition（workflow → `analysis/`，intelligence → `intelligence/`），且 pilot 驗證通過 | 舊 technique 檔案標註 `# Deprecated — see <new path>`，但保留檔案。`skills/` 仍可被 tool adapter 載入。 | ✅ 已完成（2026-05-12） |
| **Phase C** | Pipeline 可自動從舊 skill 提取 intelligence atoms，且驗證通過 | 可開始刪除已完全覆蓋的舊 technique 檔案。刪除前需確認：<br>1. `analysis/` 有對應 workflow<br>2. `intelligence/` 有對應 atoms<br>3. `knowledge/indexes/` 可 route 到新路徑<br>4. `knowledge/runtime/routing-registry.yaml` 已更新<br>5. 無任何 tool adapter 依賴該舊路徑 | ✅ 已完成（2026-05-12） |
| **Phase D** | 所有 skill 內容已完全遷移到新分層 | 刪除整個 `skills/` scaffold，active entrypoint 改由 `workflow/`、`analysis/`、`intelligence/`、`feedback/history/` 與 routing registry 承接。 | ✅ 已完成（2026-05-20） |

### 判斷是否可刪除單一舊檔案的檢查清單

- [ ] 新分層有對應的 workflow / analysis / intelligence 檔案
- [ ] 舊檔案的 `# Intelligence Extracted` 或 `# Deprecated` 標註已存在
- [ ] `knowledge/indexes/README.md` 可 route 到新路徑
- [ ] `knowledge/runtime/routing-registry.yaml` 包含新路徑
- [ ] 所有 `knowledge/summaries/` 和 `knowledge/graphs/` 已更新
- [ ] 無 tool adapter 或 `.claude/settings.json` 依賴該舊路徑
- [ ] 刪除後可 rollback（git revert）

## Enforcement Rule Deprecation

Enforcement rules（`enforcement/` 下的規則）有專屬的 deprecation 流程，作為通用 lifecycle 的補充。

### 狀態定義

| 狀態 | 意義 | 檔案位置 | Metadata status |
|------|------|---------|----------------|
| `active` | 當前有效規則 | `enforcement/<rule>.md` | `validated` 或 `stable` |
| `deprecated` | 即將移除，有替代方案 | `enforcement/<rule>.md`（標記 notice） | `deprecated` |
| `removed` | 已搬移至 deprecated/ | `enforcement/deprecated/<rule>.md` | `deprecated`（source_path 更新） |

### Deprecation 流程（4 階段）

```
標記（Mark）→ 公告（Announce）→ 緩衝期（Buffer）→ 搬移（Move）
```

詳細步驟請參閱 [`enforcement/deprecated/README.md`](../../enforcement/deprecated/README.md)。

### 與通用 Lifecycle 的關係

| 通用 Lifecycle 狀態 | Enforcement Rule 對應 |
|-------------------|---------------------|
| `source-of-truth` | `active` — 規則在 `enforcement/` 中有效 |
| `deprecated` | `deprecated` — 規則已標記 deprecation notice |
| （無直接對應） | `removed` — 規則已搬移至 `enforcement/deprecated/` |

### Rule 版本追蹤

Enforcement rules 的版本變更應記錄在 metadata 中：

```yaml
# metadata/rules/<rule-name>.yaml
status: deprecated
replaces: enforcement.<new-rule-id>
deprecation_date: "2026-05-14"
removal_date: "2026-08-14"
```

- `replaces` 欄位指向取代它的新 rule（使用 `enforcement.<rule-id>` 格式）
- `deprecation_date` 為標記 deprecation 的日期
- `removal_date` 為預計搬移至 `enforcement/deprecated/` 的日期（預設 3 個月緩衝期）

### 檢查清單

Deprecate 一條 enforcement rule 前，確認：

- [ ] Metadata 已更新（`status: deprecated`、`replaces`、`deprecation_date`）
- [ ] 原始 rule 檔案已加入 deprecation notice
- [ ] `enforcement/README.md` 規則索引已更新
- [ ] `runtime/runtime.db` 已移除或註解該 rule
- [ ] `knowledge/graphs/rules/` 中該 rule 的 graph 記錄已更新
- [ ] Linked updates 已完成（`enforcement/linked-updates.md`）
- [ ] 取代規則至少已通過 1 次成功使用驗證
- [ ] 緩衝期（3 個月）已設定
