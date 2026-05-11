# Claude Code Configuration & Discovery System

本目錄包含 Claude Code 的配置和規則發現系統，讓 AI 自動找到並應用 Ai-skill Knowledge Operating System 的規則。

## 📋 文件清單

| 檔案 | 用途 |
|------|------|
| `settings.json` | Claude Code 主配置（bootstrap、工作流、token 優化、git 規則） |
| `discovery-map.json` | 架構層級發現地圖 - 指引 AI 在不同情境下讀什麼文件 |
| `README.md` | 本檔案 - 說明發現機制 |

## 🔍 發現機制如何工作

### 問題：為什麼需要發現機制？

1. **舊方式的問題**
   - AI 看到文檔後不知道要讀什麼規則
   - 每次都要手動告訴 AI "讀這個文件"
   - 規則分散，難以自動發現

2. **新方式的優勢**
   - AI 啟動時自動讀 `settings.json` 和 `discovery-map.json`
   - 遇到特定情境時，自動查詢該讀什麼文件
   - 其他 AI 也能利用同一套發現系統

### 如何使用

#### 步驟 1：引導 AI 讀取配置（已自動化）
```bash
# Claude Code 啟動時自動讀取
.claude/settings.json
.claude/discovery-map.json
```

#### 步驟 2：AI 根據情境查詢
```
情境：規劃新任務
├─ 讀取：knowledge/indexes/README.md（任務路由）
├─ 找到：Primary source 和 Related sources
└─ 執行：任務，同時讀取依賴文件
```

#### 步驟 3：AI 自動驗證規則
```
創建新規則時：
├─ 查詢：discovery_rules.when_defining_rules
├─ 檢查：shared-rules/linked-updates.md
├─ 驗證：metadata/schema.md
└─ 確認：governance/validation 要求
```

## 📌 關鍵概念

### Default Bootstrap（預載規則）

AI 啟動時必讀的 12 個規則檔案：

```
1. README.md                          # 知識庫架構
2. shared-rules/README.md             # 規則索引
3. shared-rules/dependency-reading.md # 依賴讀取規則
4. shared-rules/linked-updates.md     # 連動更新規則
5. shared-rules/conversation-goal-ledger.md
6. shared-rules/tool-neutral-documentation.md
7. shared-rules/rule-weight.md        # 規則優先序
8. shared-rules/decision-efficiency.md
9. shared-rules/failure-learning-system.md
10. shared-rules/document-todo-list.md
11. shared-rules/document-sizing.md
12. shared-rules/goal-action-validation.md
13. shared-rules/neutral-language.md
```

### 架構層級發現

`discovery-map.json` 包含 7 個重要層級：

| 層級 | 檔案 | 何時讀取 |
|------|------|---------|
| **knowledge** | knowledge/indexes/README.md | 規劃新任務 |
| **metadata** | metadata/schema.md | 建立知識原子 |
| **governance** | governance/lifecycle/README.md | 評估推廣或淘汰 |
| **intelligence** | intelligence/README.md | 做架構決策 |
| **workflow** | workflow/README.md | 規劃複雜任務 |
| **feedback** | feedback/promotion/README.md | 提取和推廣教訓 |
| **memory** | memory/README.md | 存儲/檢索長期記憶 |

### 發現規則

根據情境自動選擇讀取檔案：

```json
{
  "when_starting_task": "讀 knowledge/indexes/README.md",
  "when_defining_rules": "檢查 shared-rules/linked-updates.md",
  "when_evaluating_promotion": "讀 governance/lifecycle/README.md",
  "when_stuck_or_uncertain": "查詢 shared-rules/failure-learning-system.md"
}
```

## ✅ 工作流優化

### Token 優化配置

```json
{
  "workflow_preferences": {
    "response_style": "concise",        // 簡潔回應
    "avoid_repetition": true,           // 避免重複
    "batch_operations": true            // 合併操作
  },
  "token_optimization": {
    "use_limited_file_reads": true,     // 只讀需要的部分
    "batch_git_commands": true,         // 合併 git 命令
    "skip_redundant_explanations": true // 跳過冗餘解釋
  }
}
```

### 權限白名單

預配置常用命令，避免權限提示：

```json
{
  "permissions": {
    "allow": [
      "Bash(git log *)",
      "Bash(git diff *)",
      "Bash(find *)",
      "Bash(grep *)",
      "Bash(jq *)"
    ]
  }
}
```

## 🎯 給新 AI 的指引

### 啟動檢查清單

新 AI 進入時應做：

- [ ] 讀 `settings.json` 了解基本配置
- [ ] 讀 `discovery-map.json` 了解層級結構
- [ ] 讀主 `README.md` 了解知識庫架構
- [ ] 根據任務，使用 `discovery_rules` 查詢該讀什麼
- [ ] 避免重複讀檔 - 使用 "Wasted call" 快取

### 常見情境查詢

| 我要做... | 查詢 | 讀... |
|---------|------|-------|
| 規劃新任務 | `discovery_rules.when_starting_task` | knowledge/indexes/ |
| 定義新規則 | `discovery_rules.when_defining_rules` | shared-rules/linked-updates.md |
| 不知道該怎辦 | `discovery_rules.when_stuck_or_uncertain` | 多個來源 |
| 做架構決策 | `discovery_rules.when_making_architecture_decision` | intelligence/ |

## 🚀 為什麼這樣做

### 解決的問題

1. **自動發現**
   - AI 不用問「我該讀什麼？」
   - 自動根據情境找答案

2. **減少重複勞動**
   - 每個進來的 AI 都能遵循同一套發現流程
   - 不用每次都手動指導

3. **智能探索**
   - 規則本身是自文檔化的
   - AI 可以「聰明地去找」

4. **節省 Token**
   - 避免不必要的對話確認
   - 直接根據規則行動

## 📚 相關檔案

- **主配置**: `.claude/settings.json`
- **發現地圖**: `.claude/discovery-map.json`
- **規則入口**: `shared-rules/README.md`
- **任務路由**: `knowledge/indexes/README.md`
- **驗證規則**: `shared-rules/rule-weight.md`

---

**最後更新**: 2026-05-11  
**作者意圖**: 讓 AI 自動化「聰明地去找規則」的過程

https://claude.ai/code/session_013aHoFHESgu26J2kPSH82rj
