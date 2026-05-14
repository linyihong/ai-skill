# 知識生命週期（Knowledge Lifecycle）

`governance/lifecycle/` 定義知識如何從既有來源檔案遷移到新的 AI-native 分層，同時不破壞舊 skill 入口。

## 真相來源規則（Source Of Truth Rule）

在遷移被明確提升（promote）並驗證之前，既有的 `skills/`、`shared-rules/`、`ai-tools/` 和 `scripts/` 檔案仍然是可執行行為的真相來源。

`analysis/`、`workflow/`、`intelligence/`、`runtime/`、`memory/`、`feedback/`、`models/`、`governance/`、`knowledge/` 和 `metadata/` 下的新分層檔案可以作為：

- 路由表面（Routing surfaces）
- 候選地圖（Candidate maps）
- 提升目標（Promotion targets）
- 中繼資料或摘要表面（Metadata or summary surfaces）
- 治理與執行期設計（Governance and runtime design）

它們在提升之前不得靜默取代舊 skill 行為。

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
| `source-of-truth` | 當前標準行為在此 | 既有的 skill/shared rule/tool/script 檔案 | 將較新的地圖視為覆蓋 |
| `candidate-map` | 從當前來源到未來分層目的地的地圖 | 擁有權邊界、來源到目標對應表、相容性說明 | 大量內容遷移或行為變更 |
| `candidate-atom` | 提議的 Knowledge Atom 或摘要 | 中繼資料、摘要、來源連結、驗證標準 | 未經使用就標記為穩定 |
| `validated-atom` | 已成功使用或審查的候選項目 | 路由中繼資料、摘要、檢查清單、驗證證據 | 移除舊入口 |
| `promoted` | 新分層成為支援的參考路徑 | 舊入口連結到提升後的 atom，索引在需要時路由到兩者 | 未經棄用就刪除相容路徑 |
| `deprecated` | 舊路徑正在退役，有替代方案 | 棄用說明、替代連結、驗證記錄 | 破壞既有連結或工具載入 |

## 冷資料歸檔（Cold Data Archive）

冷資料歸檔是生命週期/執行期壓縮策略，不是把知識搬離 Markdown。當 lesson、summary 或 graph 數量成長到 agent 每次都需要掃大量冷資料時，應先建立 generated summary、SQLite / FTS lookup cache 或 report view，讓 agent 先查候選 source，再按需讀全文。

### 觸發門檻（Trigger Thresholds）

符合任一條件時，應啟動冷資料歸檔 / generated lookup 檢查：

| 觸發條件 | 動作 |
| --- | --- |
| 單一 domain 的 `feedback/history/<domain>/` 超過約 50 條 lesson，或單一 category 超過約 20 條 | 產生或更新 category index、summary rows 與 SQLite / FTS index |
| Agent 為了找 lesson 需要讀多個 `feedback_history` README 或大量 lesson 全文 | 先用 SQLite / FTS query 找候選 `source_path`，再讀 1-3 個 canonical files |
| 某批 lesson 長期未修改，但仍可能被 routing / promotion 使用 | 標為 cold lookup candidate，保留 Markdown source，產生短 summary / tags / validation signal |
| feedback lesson 被 promotion 到 `workflow/`、`intelligence/`、`shared-rules/` 或 runtime route | 保留原 lesson，更新 summary / graph / registry / SQLite lookup |
| 查詢成本高於判斷成本，例如只需要知道「有哪些相關 lesson」 | 使用 generated report / SQLite index，不讀全文 |

### 規則（Rules）

- Canonical source 仍是 `feedback/history/*/*.md`、`shared-rules/`、`knowledge/summaries/`、`knowledge/graphs/` 與 routing registry。
- SQLite / FTS、generated summaries、runtime reports 都是 generated lookup views；可刪除、可重建，不作唯一來源。
- 冷資料歸檔不等於棄用。冷資料仍可能有效，只是預設不讀全文。
- 需要修改、promotion、debug、failure learning 或高信心判斷時，必須回到 canonical Markdown / YAML。
- 若 generated lookup 與 source 不一致，降級 lookup confidence 並重新產生，不要改 source 去迎合 cache。

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

## 技能仍在變更時的救援策略（Rescue Strategy While Skills Still Change）

當舊 skill 在遷移進行中被更新時，**不應再更新舊的 `skills/<name>/` 來源**，因為該 skill 正在被遷移到新分層，舊路徑應視為凍結（frozen）。

正確的救援流程：

1. **直接更新新分層的目標檔案** — 找出該 skill 對應的 `workflow/`、`analysis/`、`intelligence/` 目標，將變更內容寫入新分層。
2. **檢查 candidate map 或 promoted atom** — 確認是否有 map 或 atom 參考了變更的章節。
3. **同步更新 map、metadata、summary 或 index** — 如果有參考，在同一變更中更新。
4. **記錄不需要連動更新** — 如果沒有參考，在最終驗證中記錄即可。
5. **不要複製到舊路徑** — 除非該 atom 尚未被 promote，否則不要將新內容複製回 `skills/`。
6. **標註舊檔案為 deprecated** — 如果該 skill 的所有內容都已遷移完成，在舊檔案開頭加入 `# Deprecated — see <new path>` 標註。

## 刪除規則（Deletion Rule）

在 candidate-map 或 candidate-atom 階段不要刪除或移動舊 skill 檔案。只有在以下條件滿足時才能考慮刪除：

- 替代路徑已被提升。
- 既有的 tool adapter 仍可載入該 skill，或有文件化的替代方案。
- 連結、indexes、summaries 和 metadata 已更新。
- 存在棄用說明和復原路徑。

## 技能棄用時間表（Skills Deprecation Timeline）

舊 `skills/` 目錄的內容不會一次性刪除，而是依以下階段逐步 deprecate：

| 階段 | 條件 | 動作 | 狀態 |
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
