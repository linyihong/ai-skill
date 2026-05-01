# 如何新增 Skill（方法與模板）

目標：在本 repository 的 **`skills/<skill-name>/`** 新增一個可被 Cursor／Agent 讀取的技巧包，並與 **`shared-rules/`** 對齊（不重複維護共用政策）。

## 1. 要不要開新 skill？

| 情況 | 建議 |
| --- | --- |
| 內容仍是 APK／流量／Frida／Flutter AOT 等同家族 | 擴充現有 **`apk-analysis/`**，用 `feedback_history/` 與 WORKFLOW／TOOLS 收斂 |
| 新領域（例：iOS 靜態流程、另一種合規測試方法論）且會長期累積 | 新建 **`skills/<name>/`** |

Skill **資料夾名稱**建議 **kebab-case**、簡短、穩定（例：`apk-analysis`，之後如 `ios-ipa-analysis`）。

## 2. 目錄結構（建議）

最小可運作組合：

```text
skills/<skill-name>/
  SKILL.md              # 必填：YAML frontmatter + 給 Agent 的入口正文
  feedback_history/     # 強烈建議：每條 lesson 一檔（見 shared-rules/feedback-lessons.md）
```

常見擴充（依需求增量建立）：

```text
  README.md             # 給人類的導讀
  RUNBOOK.md            # 新專案第一天提示詞與路徑約定
  WORKFLOW.md           # 決策流程
  TOOLS.md              # 工具與環境
  DOCUMENTATION.md      # 產出格式
  FEEDBACK.md           # 可選：極短入口，連到 shared-rules/feedback-lessons.md
```

可直接複製範本（在 `<AI_SKILL_REPO>` 根目錄執行，將 `my-skill` 改成你的名稱）：

```bash
mkdir -p "skills/my-skill/feedback_history"
cp "skills/_template/SKILL.md" "skills/_template/FEEDBACK.md" "skills/my-skill/"
# 再編輯 skills/my-skill/SKILL.md 與 FEEDBACK.md，替換 <…> 占位符
# 可選：建立 skills/my-skill/README.md、RUNBOOK.md 等
# 可選：touch skills/my-skill/feedback_history/README.md 做索引表
```

## 3. `SKILL.md` 必填欄位（Cursor／Agent）

檔案**開頭**使用 YAML frontmatter（與 `apk-analysis` 相同風格）：

| 欄位 | 說明 |
| --- | --- |
| `name` | 與資料夾名一致或至少可對應；建議 kebab-case |
| `description` | **英文**一段話，描述何時應套用此 skill（觸發詞、場景、邊界）。寫給路由／檢索用，避免特定客戶名與機密 |

正文建議包含：**何時用／不用**、**Quick Start**、**如何回饋**（指向 `shared-rules/feedback-lessons.md`）、**預設讀哪些檔**。

完整占位稿見 **[`_template/SKILL.md`](_template/SKILL.md)**。

## 4. 與共用規則的關係

- **授權、去敏、路徑占位、feedback 檔名與模板**：一律只維護在 **`shared-rules/`**，各 skill **不要**複製長文。
- **連動更新規則**：一律只維護在 **[`shared-rules/linked-updates.md`](../shared-rules/linked-updates.md)**；新增 skill 或修改 skill 結構時，受影響的索引、入口、同步文件、分類文件**必須**同步更新或明確檢查。
- 各 skill 的 **`FEEDBACK.md`**（若需要）：維持與 [`apk-analysis/FEEDBACK.md`](apk-analysis/FEEDBACK.md) 相同模式——**幾行連結**到 [`shared-rules/feedback-lessons.md`](../shared-rules/feedback-lessons.md)。
- 每一則 lesson 頂部引用 `shared-rules`（路徑依檔案深度調整 `../../../shared-rules/...`）。

## 5. 在本庫登記與提交

1. 編輯根目錄 [**`README.md`**](../README.md) 的「現有 Skills」表格，加一行 `skills/<name>/` 與簡短說明。
2. 編輯 [**`skills/README.md`**](README.md) 的表格，同步列出。
3. 依 [`shared-rules/linked-updates.md`](../shared-rules/linked-updates.md) 檢查是否還需要同步更新 `RUNBOOK.md`、`WORKFLOW.md`、`DOCUMENTATION.md`、同步腳本或 cross-link。
4. `git add` → `commit` → `push`。

## 6. 同步到本機 `~/.cursor`（可選）

與 [`shared-rules/cursor-sync.md`](../shared-rules/cursor-sync.md) 相同精神：**`shared-rules/`** 與各 **skill** 要成對出現。

### 建議（共用資產放在 `bundles/` 並列）：`bundles/shared-rules` + `bundles/ai-skill`

避免 `~/.cursor` 底下其他規則或工具與這套資產混在一起，本機優先使用：

- **`~/.cursor/bundles/shared-rules`** → 本庫 **`shared-rules/`**（單一 symlink）
- **`~/.cursor/bundles/ai-skill/<skill>/`** → 本庫 **`skills/<skill>/`**（每個 skill 一個 symlink）
- **`~/.cursor/shared-rules`** → `bundles/shared-rules`
- **`~/.cursor/skills/<skill>`** → `bundles/ai-skill/<skill>`

一鍵同步（在本庫根目錄）：

```bash
./scripts/sync-cursor-bundle.sh
```

腳本會掃描 `skills/` 下含 **`SKILL.md`** 的目錄（略過 `_template`）。若 **`~/.cursor/shared-rules`** 已是「真實資料夾」而非 symlink，腳本會先**移到** `*.bak.<隨機>` 再建立連結。

### 簡易做法（直接連到 repo）

仍可直接連到 clone（較短，但與「bundle 隔離」精神不同）：

```bash
ln -sf "${AI_SKILL_REPO}/shared-rules" "${HOME}/.cursor/shared-rules"
ln -sf "${AI_SKILL_REPO}/skills/my-skill" "${HOME}/.cursor/skills/my-skill"
```

新增 skill 後建議在 Cursor **`⌘⇧P` → Developer: Reload Window** 重載一次。

## 7. 檢查清單（新建完成前）

- [ ] `SKILL.md` 有合法 `name` / `description` frontmatter
- [ ] 正文有連到 `shared-rules` 與 `feedback-lessons`
- [ ] 已建立 `feedback_history/`（可先要 `README.md` 索引）
- [ ] 根目錄 `README.md` 與 `skills/README.md` 已更新
- [ ] 已依 `shared-rules/linked-updates.md` 完成或明確檢查必要連動更新
- [ ] 無真實本機絕對路徑、無機密寫入將 commit 的檔案

← [回到 skills 索引](README.md)
