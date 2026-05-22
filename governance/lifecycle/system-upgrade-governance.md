# 系統升級治理要則（System Upgrade Governance）

> **目的**：確保大型系統升級在計畫階段就涵蓋所有必要治理項目，避免執行階段遺漏 README 更新、索引同步、舊檔案標註等關鍵步驟。
>
> **適用範圍**：任何會改變系統名稱、架構分層、核心流程、或對外識別的重大升級。

---

## 1. 什麼是「大型系統升級」

符合以下任一條件即視為大型系統升級，**必須**遵循本治理要則：

| 條件 | 說明 | 範例 |
|------|------|------|
| 🏷️ **系統名稱變更** | 改變 repository 的定位名稱 | `Skill Repository` → `Knowledge Operating System` → `Cognitive Execution System` |
| 🏛️ **架構分層變更** | 新增/刪除/合併頂層目錄 | 新增 `runtime/`、`governance/`、`intelligence/` 層 |
| 🔄 **核心流程變更** | 改變 agent 的啟動、執行、閉環流程 | Bootstrap 流程改寫、新增 phase machine |
| 📄 **對外文件變更** | 改變 README、CORE_BOOTSTRAP、入口文件 | README 標題/描述更新、新增 Key Documents |
| 🗑️ **舊結構淘汰** | 標註 deprecated、搬遷內容、刪除舊檔案 | Skills deprecation、techniques 搬遷 |

### 非大型升級（不適用本要則）

- 單一 intelligence atom 新增
- 單一 workflow 步驟調整
- 日常 bug fix 或文件修正
- 常規知識更新（knowledge update flow）

---

## 2. 升級計畫書必須包含的檢查清單

每一份大型升級計畫書（如 `plans/active/` 下的文件）**必須**在結尾或附錄包含以下檢查清單。執行 agent 在完成所有實作後，逐項確認並標記完成。

### 2.1 計畫書本身

- [ ] **計畫書狀態標記**：draft → active → completed / abandoned
- [ ] **完成日期記錄**：記錄實際完成日期
- [ ] **偏離記錄**：如果實作偏離原始計畫，記錄偏離內容與原因
- [ ] **歸檔**：完成後移至 `plans/archived/`

### 2.2 README 更新

- [ ] **系統名稱**：`README.md` 標題是否反映新名稱？
- [ ] **系統描述**：`README.md` 副標題/描述是否反映新定位？
- [ ] **OS Layout 表格**：是否有新增/修改的層級需要更新？
- [ ] **Key Documents**：是否有新增/歸檔的計畫文件需要更新？
- [ ] **Agent 作業流程**：流程圖是否仍正確？
- [ ] **Quickstart**：啟動步驟是否仍正確？

### 2.3 架構文件更新

- [ ] **architecture/ 文件**：如有架構變更，是否更新了對應的架構文件？
- [ ] **新世代 architecture 文件**：如果本次升級涉及**系統名稱變更**（§1 第一項條件），是否在 `architecture/` 建立新文件 `architecture/<new-system-slug>.md` 作為新世代 canonical navigation 入口？同時把舊世代文件加上 historical header？（規則 6）
- [ ] **CORE_BOOTSTRAP.md**：啟動流程是否因升級而需要修改？
- [ ] **ADR 記錄**：是否建立了 Architecture Decision Record 記錄本次升級的關鍵決策？

### 2.4 索引與路由更新

- [ ] **knowledge/runtime/routing-registry.yaml**：如有 route 層級變更，是否更新？
- [ ] **knowledge/indexes/README.md**：路由索引是否反映新結構？
- [ ] **knowledge/runtime/routing-registry.yaml**：routing registry 是否新增/更新路由？
- [ ] **knowledge/graphs/**：graph records 是否反映新檔案之間的關係？

### 2.5 舊檔案處理

- [ ] **舊檔案標註**：被取代的舊檔案是否標註 `# Deprecated — see <new path>`？
- [ ] **已提取標註**：內容被提取的檔案是否標註 `# Intelligence Extracted` 或 `# Extracted — See <target>`？
- [ ] **舊檔案刪除**：如果刪除舊檔案，是否確認無殘留引用？（建議先標註 deprecated，下一個版本再刪除）
- [ ] **引用更新**：所有指向舊路徑的引用是否已更新為新路徑？

### 2.6 Runtime Surface 更新

- [ ] **SQLite runtime 重新編譯**：執行 repo-local `ai-skill runtime compile --native-compiler` 或由 pre-commit hook 重新編譯 `runtime/runtime.db`
- [ ] **phase runtime surface**：如有 phase 變更，已更新 `runtime/runtime.db` 或現存 SQLite canonical document，並確認 `runtime/runtime.db` 的 `phase_machine` / `phases` 已同步。
- [ ] **obligation runtime surface**：如有 obligation 變更，已更新 source，並確認 `runtime/runtime.db` 的 `obligation_ledger` / `obligations` 已同步。
- [ ] **blocking gate runtime surface**：如有 gate 變更，已更新 source，並確認 `runtime/runtime.db` 的 `blocking_gates` / `gates` 已同步。

### 2.7 跨層一致性檢查

- [ ] **README 一致性**：各層 README 的表格（atoms、workflows、files）與實際檔案一致
- [ ] **引用一致性**：所有跨文件引用（相對路徑連結）有效
- [ ] **命名一致性**：新檔案命名遵循 path convention（小寫、連字號）

### 2.8 閉環驗證

- [ ] **Diff Review**：執行 `git diff` 審查所有變更
- [ ] **Linked Updates**：檢查是否有遺漏的連動更新（參見 [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md)）
- [ ] **Commit**：`git add -A && git commit`
- [ ] **Push**：`git push`
- [ ] **Readback**：確認 push 成功，遠端狀態正確

---

## 3. 從三次升級提煉的強制規則

以下規則來自實際升級經驗中的遺漏教訓，**所有大型升級必須遵守**：

### 規則 1：計畫書必須包含完成檢查清單

**教訓**：`2026-05-15-0920-runtime-execution-layer-upgrade-analysis.md` 的 Phase 規劃中沒有明確列出「完成後需要更新 README」，導致執行 agent 忘記更新系統名稱。

**強制**：每一份升級計畫書必須在結尾包含 §2 的檢查清單（或子集），執行 agent 在完成所有實作後逐項確認。

### 規則 2：系統名稱變更必須寫在計畫書中

**教訓**：從 Knowledge OS → Cognitive Execution System 的名稱變更，在計畫書中只有 §一 的「現有框架摘要」提到舊名稱，沒有明確寫出「升級完成後 README 標題要改為 XXX」。

**強制**：如果升級涉及系統名稱變更，計畫書必須明確寫出：
```
## 完成條件
- [ ] README.md 標題改為「# AI-native Cognitive Execution System」
- [ ] README.md 描述改為「AI 認知執行系統 — ...」
```

### 規則 3：舊檔案必須在升級過程中標註

**教訓**：第一次升級（skill → Knowledge OS）時，舊 `skills/` 檔案沒有即時標註 `# Deprecated`，導致 agent 仍然優先讀取舊路徑。後來才透過 `primary_entrypoint` 機制修正。

**強制**：
- 內容被搬遷的舊檔案：標註 `# Deprecated — see <new path>`
- 內容被提取的舊檔案：標註 `# Intelligence Extracted` 或 `# Extracted — See <target>`
- 刪除舊檔案前：至少保留一個版本的 deprecated 過渡期

### 規則 4：索引必須在升級完成前更新

**教訓**：`knowledge/runtime/routing-registry.yaml` 和 `knowledge/indexes/README.md` 在多次升級中都曾被遺漏，直到跨層一致性檢查才補上。

**強制**：索引更新是升級的**完成條件**，不是「有時間再補」。在檢查清單中必須有對應項目。

### 規則 5：Compiler 必須在升級完成前執行

**教訓**：修改 prose 檔案後忘記執行 repo-local `ai-skill runtime compile --native-compiler` / runtime validation，導致 generated runtime surface 與 canonical source 不一致。

**強制**：任何修改 prose 檔案的升級，在 commit 前必須執行 compiler 重新編譯 generated YAML。

### 規則 6：新世代系統必須有對應的 architecture/ canonical 文件

**教訓**：Knowledge Operating System → Cognitive Execution System 的升級雖然更新了根 `README.md` 標題，但 `architecture/ai-native-knowledge-operating-system.md` 仍以舊世代名稱與內容存在，被誤認為「architecture/ 唯一文件」造成讀者混淆。直到 2026-05-22 才被發現並補建第三代 navigation 文件。

**強制**：如果升級涉及**系統名稱變更**（§1 第一項條件），計畫書必須明確包含以下項目。

> ⚠️ **時機**（2026-05-22 修訂）：依 [`decision-promotion-pipeline.md`](decision-promotion-pipeline.md) §No-Proposed-ADR Rule，下列 ADR 與 architecture/ 文件的建立時機**從升級開始時改為 plan completed 時**。升級開始時 plan 只需含「Architecture Doc Plan section」記錄將寫什麼；實際寫 architecture/<gen>.md 與 ADR 在 plan completed 後執行（ADR 直接 accepted，不經 proposed 階段）。

```
## 完成條件
- [ ] architecture/ 建立新世代 canonical navigation 文件
      檔名格式：architecture/<new-system-slug>.md
      內容：系統世代演化表、當前 canonical 入口列表、核心機制概述（指向真正 source-of-truth）、演化里程碑（引用相關 archived plans）
      不寫成自包含 spec — 真正 source-of-truth 是 YAML contracts + philosophy 文件 + runtime.db
- [ ] 舊世代文件加 historical header，指向新世代文件
- [ ] architecture/README.md 更新「目前文件」表，標清楚每世代狀態（current / historical）
- [ ] 補充「與 intelligence/engineering/architecture/ 的邊界」說明，避免 OS 架構與工程架構混淆
```

**為什麼是必要的**：
- 每世代系統有自己的 canonical 入口，新人或 agent 不需要從散落的 philosophy / YAML / plans 拼湊當前設計
- 歷史世代文件不被覆寫，提供 traceability
- `architecture/README.md` 的定位明確為「每世代 canonical 入口」，避免「永久文件」這種與升級事實矛盾的描述

### 規則 7：世代升級時必須評估每個既有 ADR 並建立雙向連結

**教訓**：Knowledge OS → Cognitive Execution System 升級後，constitution/ 內 7 個 ADR 沒有任何標註說明屬於哪個世代。ADR-003 標題仍使用「Knowledge / **Skills** / Intelligence」，但 Skills 在 Gen 3 已 deprecated；ADR-004 仍引用 `skills/*/feedback_history/` 路徑（已搬到 `feedback/history/<domain>/`）。讀者無法判斷 ADR 是否仍適用、哪些詞彙已演化。架構文件與 ADR 之間缺乏 traceability。2026-05-22 才被使用者指出並補建雙向連結。

**強制**：如果升級涉及**系統名稱變更**（§1 第一項條件）或**架構分層變更**（§1 第二項條件），計畫書必須包含以下項目：

```
## 完成條件 — ADR 與架構連動

1. 評估每個既有 ADR 的延伸狀態
   - [ ] 對 constitution/ADR-* 逐份判斷以下三種狀態之一：
       - Gen N 確立（決策起源於本世代）
       - cross-generation（決策跨世代延伸）
       - Gen N-1 / Gen N transition（詞彙或路徑演化但核心精神保留）
       - superseded（決策被新 ADR 推翻 — 走 ADR supersession 流程）

2. 在每個 ADR 加入 Framework Generation section（在 Status 後 Context 前）
   - [ ] 世代分類
   - [ ] 當前世代文件連結
   - [ ] 適用狀態簡述（仍 active / 詞彙演化 / 已 superseded）

3. 對需要詞彙或路徑演化說明的 ADR 加入 Vocabulary Evolution section
   - [ ] 不修改既有 immutable 正文
   - [ ] 列出原文用詞 → Gen N 對應 + 說明
   - [ ] 明寫「核心精神保留」+「為什麼不 supersede 本 ADR」

4. 在新世代 architecture/<gen>.md 加入「本世代相關 ADR」表
   - [ ] 列出本世代涵蓋的 ADR + 性質（Gen 確立 / cross-generation / transition）
   - [ ] 簡述每個 ADR 在本世代的角色

5. 在歷史世代 architecture/<old-gen>.md 加入「本世代相關 ADR（部分延伸到下一代）」表
   - [ ] 標清楚 Gen N 與 Gen N+1 狀態對比

6. 更新 constitution/README.md ADR 表格
   - [ ] 加入「Framework Generation」欄
   - [ ] 每個 ADR 填寫世代分類
```

**為什麼是必要的**：
- ADR 是 immutable 的，但系統會演化；沒有 Vocabulary Evolution 機制，新 agent 讀 ADR 會以為過時詞彙仍是 canonical
- Architecture ↔ ADR 雙向連結讓「知道 Gen 3 由哪些 ADR 構成」與「知道某個 ADR 屬於哪個 Gen」都能在一步內完成
- ADR supersession 是「決策被推翻」，Vocabulary Evolution 是「決策保留但標籤演化」— 兩個邊界清楚才不會誤用 supersession 流程

---

## 4. 升級流程圖

```
[提出升級想法]
     │
     ▼
[建立升級計畫書] ─── 放在 plans/active/
     │                  │
     │                  ├─ 定義 scope（§1 條件）
     │                  ├─ 包含完成檢查清單（§2）
     │                  ├─ 明確寫出系統名稱變更（如適用）
     │                  └─ 列出受影響檔案
     │
     ▼
[執行升級實作]
     │
     ├─ 建立/修改檔案
     ├─ 標註舊檔案
     ├─ 更新索引
     ├─ 更新 README
     ├─ 更新架構文件/ADR
     │
     ▼
[執行完成檢查清單] ─── 逐項確認（§2）
     │
     ▼
[Compiler + 閉環]
     ├─ 執行 compiler
     ├─ Diff review
     ├─ Linked updates 檢查
     ├─ Commit
     ├─ Push
     └─ Readback 確認
     │
     ▼
[歸檔計畫書] ─── 移至 plans/archived/
     │
     ▼
[完成]
```

---

## 5. 檢查清單範本（可直接複製到計畫書）

```markdown
## 完成檢查清單

### 計畫書本身
- [ ] 計畫書狀態標記為 completed
- [ ] 記錄完成日期
- [ ] 記錄偏離（如有）
- [ ] 歸檔至 plans/archived/

### README 更新
- [ ] 系統名稱已更新
- [ ] 系統描述已更新
- [ ] OS Layout 表格已更新
- [ ] Key Documents 已更新
- [ ] Agent 作業流程已確認
- [ ] Quickstart 已確認

### 架構文件
- [ ] architecture/ 文件已更新
- [ ] CORE_BOOTSTRAP.md 已確認
- [ ] ADR 已記錄關鍵決策

### 索引與路由
- [ ] knowledge/runtime/routing-registry.yaml 已更新
- [ ] knowledge/indexes/README.md 已更新
- [ ] routing-registry.yaml 已更新
- [ ] graph records 已更新

### 舊檔案處理
- [ ] 舊檔案已標註 deprecated/extracted
- [ ] 所有引用已更新為新路徑

### Runtime Surface
- [ ] compiler 已執行
- [ ] `runtime/runtime.db` 已同步
- [ ] phase / obligation / blocking gate source 已同步（embedded source 或現存 SQLite canonical document）

### 跨層一致性
- [ ] README 表格與實際檔案一致
- [ ] 跨文件引用有效
- [ ] 命名遵循 convention

### 閉環驗證
- [ ] Diff review 完成
- [ ] Linked updates 完成
- [ ] Commit 完成
- [ ] Push 完成
- [ ] Readback 確認
```

---

## 6. Runtime Surface

本文件已註冊為 runtime compiler 的 source，產生對應的 SQLite surface 供 agent 在 checkpoint 階段快速讀取。

| Runtime Surface | 位置 | 用途 |
|----------------|------|------|
| SQLite (runtime.db) | `runtime/runtime.db → generated_surfaces (type='system_upgrade_governance')` | Agent 在 checkpoint 階段查詢此記錄，了解升級條件、檢查清單分類與強制規則 |

### 觸發時機

系統升級治理檢查在每個 checkpoint 階段自動觸發（由 `runtime/runtime.db` 的 `phase_machine` / `blocking_gates` compiled surface 管理；source 在 [`runtime/runtime.db`](../../runtime/runtime.db)）：

1. Agent 進入 checkpoint phase
2. 查詢 `runtime.db` 的 `generated_surfaces` 表（快速路徑）
3. 檢查 `plans/active/` 中是否有大型升級計畫
4. 如有 → 逐項確認 §2 檢查清單
5. 如無 → 跳過

### 更新流程

修改本文件後，必須執行 compiler 重新編譯：

```bash
scripts/ai-skill-cli/bin/ai-skill-darwin-arm64 runtime compile --native-compiler
```

Pre-commit hook 會檢查 prose 與 runtime.db 是否一致，不一致時 block commit。

---

## 7. 與既有文件的關係

| 文件 | 關係 |
|------|------|
| `runtime/runtime.db → generated_surfaces (type='system_upgrade_governance')` | 本文件的 runtime surface，由 compiler 自動產生 |
| [`runtime/runtime.db`](../../runtime/runtime.db) | Checkpoint phase、`obligation.checkpoint.check_system_upgrade_governance`、`gate.checkpoint.system_upgrade_governance_checked` 的 compiled runtime surface |
| [`runtime/runtime.db`](../../runtime/runtime.db) | Phase / obligation / blocking gate 的 source；除非有 source restoration migration，不要引用已移除的 standalone YAML |
| [`runtime/runtime.db`](../../runtime/runtime.db) | 定義本文件到 YAML 的編譯規則 |
| [`enforcement/linked-updates.md`](../../enforcement/linked-updates.md) | 本要則的 §2.8 閉環驗證引用 linked-updates 規則 |
| [`governance/lifecycle/README.md`](../../governance/lifecycle/README.md) | 本要則是 lifecycle 治理的一部分，專注於「升級」這個特定 lifecycle 事件 |
| [`plans/README.md`](../../plans/README.md) | 本要則定義了計畫書必須包含的內容，與 plans/ 的目錄規則互補 |
| [`CORE_BOOTSTRAP.md`](../../CORE_BOOTSTRAP.md) | 升級可能影響啟動流程，需同步更新 |
| [`enforcement/dependency-reading.md`](../../enforcement/dependency-reading.md) | 升級後的 commit/push/readback 閉環遵循 dependency-reading 的 writeback transaction 規則 |
