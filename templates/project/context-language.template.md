# Non-Canonical Project Context Language

> **Status**: non-canonical replay aid
>
> **不是** canonical source。**不是** runtime truth。**不是** architecture contract。
>
> 本檔提供建議格式，協助專案在 `<PROJECT_ROOT>/memory/project/context-language.md`
> 保存自己的詞彙脈絡，方便 agent 跨 session 對齊產品名詞與業務概念。
>
> 詞義衝突依 [`<AI_SKILL_REPO>/knowledge/glossary/README.md`](../../knowledge/glossary/README.md)
> §Vocabulary Resolution Priority 解析；本檔在 priority 中與 `memory replay` 同層，
> 永遠不能 override `knowledge/glossary/`、accepted ADR、active workflow。

## 使用範圍

| 是 | 不是 |
| --- | --- |
| 專案內部產品名詞、業務概念、模組代號的對齊備忘 | Ai-skill framework / runtime / cognitive vocabulary（去 `knowledge/glossary/ai-skill.md`） |
| 跨 session 的 replay aid（避免每次都重新解釋同一個詞） | Canonical source（任何宣稱「這就是正式定義」的位置） |
| 專案 docs 沒寫的非正式 working vocabulary | 取代 `architecture/`、`governance/`、ADR、合約 |
| 連結回專案正式 docs / Linear / Notion / API spec 的 pointer | 第二份正式文件 |

## 建議格式

```markdown
# <Project Name> Context Language（Non-Canonical）

> Replay aid only. Canonical source 見專案正式 docs。

## 產品 / 業務詞彙

### <term>
- **意義**：<一句話定義>
- **正式來源**：<連結到專案正式 docs / API spec / Linear ticket>
- **常見誤解**：<這個詞容易跟什麼混淆>
- **常見同義詞**：<別人怎麼稱呼這個概念>

### <next term>
...

## 內部代號 / 模組名

| 代號 | 對應 | 正式來源 |
| --- | --- | --- |
| `mod-x` | <某個服務> | <link> |
| `srv-y` | <某個 backend> | <link> |

## 已知歧義（disambiguation log）

| 詞 | 在 A context 意義 | 在 B context 意義 | 處理 |
| --- | --- | --- | --- |
| `order` | 客戶下單 | 排序動作 | 在 BDD scenario 中分別寫 `customer_order` / `sort_order` |

## 與 Ai-skill framework glossary 的關係

- 本檔的詞 **不得** 與 [`<AI_SKILL_REPO>/knowledge/glossary/ai-skill.md`](../../knowledge/glossary/ai-skill.md)
  的 framework term 同名。若不小心同名，project term 必須改名（priority §1 只允許 project
  覆蓋 project-local usage，不可改寫 framework term）。
- 本檔不需要通過 `ai-skill glossary validate`（validator 只掃 `<AI_SKILL_REPO>/knowledge/glossary/`，
  不掃 project repo）。但若想自願套用 schema 為內部 lint，可將 entries 寫成同樣的
  H2 + YAML block 格式作為內部約定。

## 維護建議

- 由 product owner / 主要 contributor 維護；不必每天更新
- 詞義穩定後盡早 promote 到專案正式 docs（API spec / glossary site / Notion），本檔只保留指向正式來源的 pointer
- 不寫入 secrets、tokens、host、私人路徑、incident raw evidence
- 不 commit 任何超出 `<PROJECT_ROOT>/memory/project/` 範圍的內容
