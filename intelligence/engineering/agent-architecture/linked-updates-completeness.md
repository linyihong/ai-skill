# Linked Updates Completeness（連動更新完整性）

**Status**: `candidate-intelligence`
**Source**: 本系統 Phase 33-35 實際運作經驗（[`enforcement/linked-updates.md`](../../../enforcement/linked-updates.md)）

## 原則

**When modifying one file, agent must identify and update all files that reference or depend on it.**

修改一個檔案時，agent 必須找出並更新所有引用或依賴它的檔案。

## 為什麼

1. **知識庫是高度互連的** — 本系統有 15+ 層相互引用的檔案（indexes、routing-registry、Knowledge Graph、summaries、README 表格）。修改一個路徑可能影響 10+ 個檔案。
2. **Agent 傾向於只修改直接相關的檔案** — 因為讀取所有相關檔案需要額外的 context budget，agent 會下意識地忽略「看起來不太重要」的連動更新。
3. **連動更新遺漏是漸進的** — 一次遺漏一個檔案看起來沒什麼，但累積 5-10 次後，系統會產生大量 dead link 和過時引用。
4. **Linked updates 是系統一致性的基石** — 沒有完整的連動更新，索引會過時、routing 會失效、Knowledge Graph 會產生孤立節點。

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **Dead link 累積** | README 中的連結指向不存在的檔案 | 高 |
| **Index 過時** | `knowledge/indexes/README.md` 缺少新建立的檔案 | 高 |
| **Routing 失效** | `routing-registry.yaml` 缺少新 route | 高 |
| **Knowledge Graph 孤立節點** | Graph 中的節點沒有 `related_to` 邊 | 中 |
| **README 表格不完整** | 目錄結構表格缺少新建立的子目錄 | 高 |

## 常見遺漏的連動更新

| 修改類型 | 常見遺漏的更新目標 |
|---------|-------------------|
| 建立新檔案 | `knowledge/indexes/README.md`、`knowledge/runtime/routing-registry.yaml`、所屬目錄的 `README.md` |
| 刪除舊檔案 | 所有引用該檔案的 README、index、graph、summary |
| 修改路徑 | `skills-index.yaml` 的 `related` 欄位、`routing-registry.yaml` 的 route |
| 新增 intelligence atom | `intelligence/<domain>/README.md` 的 atoms 表格、Knowledge Graph |
| 修改 enforcement rules | `enforcement/README.md` 的索引、`dependency-reading.md` 的邊界 |
| 修改架構文件 | `plans/archived/2026-05-11-1112-next-stage-upgrade-plan.md` 的 Durable Roadmap Goals |

## 預防方式

1. **修改前先找出所有引用** — 使用 `grep -r` 或 `search_files` 找出所有引用目標路徑的檔案
2. **建立連動更新清單** — 在 commit 前列出所有需要更新的檔案
3. **使用 linked-updates.md 的常見連動關係表** — 根據修改類型檢查對應的更新目標
4. **Commit 前做完整性檢查** — 確認 indexes、routing-registry、Knowledge Graph、README 都已更新
5. **不要在同一個 commit 混合不相關的修改** — 連動更新應該與主要修改在同一個 commit 中

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 「之後再補」連動更新 | 之後不會補，遺漏會累積 |
| 只更新直接相關的檔案 | 間接相關的檔案（如 Knowledge Graph）同樣重要 |
| 依賴 agent 記憶哪些檔案需要更新 | Agent 會遺忘，特別是在 long session 中 |

## 相關 atoms

- [`context-collapse.md`](context-collapse.md) — 上下文崩塌（連動更新遺漏是 context collapse 的典型症狀）
- [`attention-budgeting.md`](attention-budgeting.md) — 注意力預算（連動更新需要額外的 context budget）
- [`cognitive-boundaries.md`](cognitive-boundaries.md) — 認知邊界（agent 無法可靠判斷「哪些檔案需要更新」）

## Token Impact

完整的連動更新每次需要 500-2000 token（讀取相關檔案 + 更新）。遺漏連動更新會導致後續的修復成本增加 5-10 倍（因為需要先診斷問題、找出所有遺漏、逐一修復）。

---

← [回到 agent-architecture/](README.md)
