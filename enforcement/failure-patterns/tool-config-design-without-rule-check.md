# Tool Config Design Without Existing Rule Check（工具配置設計漏讀現有規則）

Status: candidate
Class: `tool-strategy-gap`

## Trigger

當 agent 被要求為工具（如 Claude Code、Cursor 等）設計新配置、適配層或自動化機制，但沒有先檢查 `ai-tools/<tool>.md` 中是否已有設計規則或邊界定義時，使用此 pattern。

## Failure Mode

Agent 基於假設或記憶提出工具配置設計，而不是從 `ai-tools/<tool>.md` 中的現有規則推導。結果可能導致：

1. 重複現有規則（違反 DRY 原則）
2. 忽視工具層責任邊界
3. 建立與既有設計衝突的配置
4. 在不同工具文檔間產生不一致

## Risk

- **知識碎片化**: 同一個設計在多個地方各寫一遍，後續維護時同步成本高
- **邊界混淆**: 不知道某些配置應該放在 `ai-tools/` 還是 `.claude/` 或 skill adapter
- **決策錯誤**: 提出的設計可能與工具文檔中明確的規則衝突
- **漏讀成本**: 需要事後重新調整而非一開始就正確

## Required Agent Action

提議任何工具配置、適配、自動化或發現機制時：

1. **檢查現有規則**
   - 讀 `ai-tools/<tool>.md`
   - 查找「配置」、「自動化」、「發現」、「適配」等相關段落

2. **理解責任邊界**
   - `ai-tools/<tool>.md` = 工具層通用規範
   - `.claude/` = Claude Code 特定配置
   - `skills/<name>/tool-adapters/<tool>.md` = Skill 特定適配
   - 避免重複內容

3. **根據規則設計**
   - 設計應該遵循或擴展現有規則，而不是憑假設
   - 若發現空缺，明確指出並記錄

4. **標注邊界**
   - 新文檔應明確指向既有規則
   - 避免成為「另一份真相來源」

## Prevention Gate

設計工具配置前，agent 必須能回答：

| Check | Required answer |
|-------|-----------------|
| 現有規則 | `ai-tools/<tool>.md` 中有哪些現有設計規則？ |
| 責任邊界 | 這個配置應該放在哪一層？為什麼？ |
| DRY 檢查 | 是否會重複 `ai-tools/<tool>.md` 的內容？ |
| 衝突檢查 | 新設計是否違反或擴展現有規則？ |
| 指向方式 | 新文檔如何指向既有規則而避免重複？ |

若無法清楚回答，先讀 `ai-tools/<tool>.md` 的完整內容。

## Validation

符合下列條件時，此 pattern 已被驗證：

- Agent 在提建議前讀過 `ai-tools/<tool>.md`
- 設計遵循現有規則邊界，或明確指出需要新規則
- 避免重複現有內容（DRY 原則）
- 新文檔（如 `.claude/README.md`）明確指向並索引現有規則
- 最終設計清楚標示各層級責任和指向關係

## Linked Rules

- [`dependency-reading.md`](../dependency-reading.md) - 讀取相關文件時必須系統性檢查依賴
- [`tool-neutral-documentation.md`](../tool-neutral-documentation.md) - 通用規則保持工具中立，工具差異放在 `ai-tools/`
- [`linked-updates.md`](../linked-updates.md) - 工具文檔更新時需要同步
- [`../ai-tools/README.md`](../../ai-tools/README.md) - AI 工具使用說明的責任邊界

## 相關文獻

這個失效源於對 `ai-tools/claude.md` 的規則漏讀：

- **第 5-13 行**: 「自動配置」- 說明 `.claude/settings.json` 的角色
- **第 87-95 行**: 「與 Tool Adapter 的關係」- 明確定義不同層級的責任

## Linked Validation Scenarios

- `validate_directory_structure` — 檢查 `ai-tools/` 中各工具文件的 README 是否列出所有子檔案，防止工具配置設計漏讀現有規則

← [Back to failure patterns](README.md)
