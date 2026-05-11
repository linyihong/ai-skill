# Claude Code 配置與發現系統

本目錄 (`.claude/`) 提供 Claude Code 的配置和規則發現機制，讓 AI 自動找到並應用 Ai-skill Knowledge Operating System 的規則。

**重要：** 本目錄是 Claude Code 工具特定的配置層。通用的 Claude 工具使用說明見 [`../ai-tools/claude.md`](../ai-tools/claude.md)。

## 📋 文件清單

| 檔案 | 用途 |
|------|------|
| `settings.json` | Claude Code 主配置（bootstrap、工作流、token 優化、git 規則、發現設定） |
| `discovery-map.json` | 架構層級發現地圖 - 指引 AI 在不同情境下讀什麼文件 |
| `README.md` | 本檔案 - 發現機制入口 |

## 🔗 快速導航

### 工具層面（通用）
→ 詳見 [`../ai-tools/claude.md`](../ai-tools/claude.md)
- Claude 如何讀取 shared rules、skill 入口、依賴文件
- 建議提示詞與使用注意
- Tool adapter 與 skill-specific 策略

### 配置層面（專案特定）
→ 本目錄提供
- Default Bootstrap 規則自動加載
- 架構層級發現機制
- Token 優化與工作流偏好

---

## 🔍 規則發現機制

### 核心概念

`.claude/discovery-map.json` 提供**情境化發現規則**，讓 AI 自動判斷該讀什麼文件：

```json
{
  "when_starting_task": "讀 knowledge/indexes/README.md",
  "when_defining_rules": "檢查 shared-rules/linked-updates.md",
  "when_stuck_or_uncertain": "查詢 shared-rules/failure-learning-system.md"
}
```

### 架構層級快速參考

`discovery-map.json` 索引了 7 個關鍵層級：

| 層級 | 檔案 | 何時讀 |
|------|------|-------|
| **knowledge** | `knowledge/indexes/README.md` | 規劃新任務 |
| **metadata** | `metadata/schema.md` | 建立知識原子 |
| **governance** | `governance/lifecycle/README.md` | 評估推廣/淘汰 |
| **intelligence** | `intelligence/README.md` | 架構決策 |
| **workflow** | `workflow/README.md` | 複雜任務規劃 |
| **feedback** | `feedback/promotion/README.md` | 提取和推廣教訓 |
| **memory** | `memory/README.md` | 長期記憶管理 |

---

## ⚙️ 配置詳情

### 1. Default Bootstrap（啟動時自動讀取）

AI 啟動時必讀的 13 個規則檔案，見 `settings.json` 中的 `default_bootstrap` 陣列。

### 2. 工作流優化

```json
{
  "workflow_preferences": {
    "response_style": "concise",      // 簡潔回應
    "avoid_repetition": true,         // 避免重複
    "batch_operations": true          // 合併操作
  }
}
```

### 3. Token 優化

```json
{
  "token_optimization": {
    "use_limited_file_reads": true,   // 只讀需要部分
    "batch_git_commands": true,       // 合併 git 命令
    "skip_redundant_explanations": true
  }
}
```

### 4. 權限白名單

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

---

## 📖 新 AI 啟動清單

進入專案時應做：

- [ ] 讀 `settings.json` 了解基本配置
- [ ] 讀 `discovery-map.json` 了解層級結構
- [ ] 讀主 `README.md` 了解知識庫架構
- [ ] 讀 [`../ai-tools/claude.md`](../ai-tools/claude.md) 了解 Claude 工具使用規範
- [ ] 根據任務，使用 `discovery_rules` 查詢該讀什麼

---

## 🎯 設計原則

1. **單一真相來源（DRY）**
   - 通用規則在 `ai-tools/claude.md`
   - 工具特定配置在 `.claude/`
   - 不重複內容

2. **自動發現**
   - `discovery-map.json` 指導不同情景
   - AI 自動查詢而非手動指導

3. **可擴展性**
   - 新層級直接加入 `discovery-map.json`
   - 新規則透過 `discovery_rules` 自動應用

---

## 📚 相關文件

- **通用 Claude 說明**: [`../ai-tools/claude.md`](../ai-tools/claude.md)
- **規則索引**: `../shared-rules/README.md`
- **任務路由**: `../knowledge/indexes/README.md`
- **檔案依賴**: `../shared-rules/dependency-reading.md`
- **規則優先序**: `../shared-rules/rule-weight.md`

---

**最後更新**: 2026-05-11  
**設計意圖**: Claude Code 配置入口 + 規則發現機制  
**維護者**: Ai-skill Knowledge Operating System

https://claude.ai/code/session_013aHoFHESgu26J2kPSH82rj
