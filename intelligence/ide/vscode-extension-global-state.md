# VS Code Extension 全域設定修改方法

## 背景

VS Code Extension 的全域設定（globalState）**不是**存在檔案系統的 JSON 或 YAML 檔中，而是存在 VS Code 的 **SQLite 資料庫**中。這是因為 VS Code 的 Extension API 提供 `context.globalState` 介面，底層使用 SQLite 儲存。

## 資料庫位置

```
~/Library/Application Support/Code/User/globalStorage/state.vscdb
```

> **注意**：此路徑為 macOS 預設值。Windows 為 `%APPDATA%/Code/User/globalStorage/state.vscdb`，Linux 為 `~/.config/Code/User/globalStorage/state.vscdb`。

## 資料庫結構

| 項目 | 說明 |
|------|------|
| 資料表 | `ItemTable` |
| 欄位 | `key` (TEXT PRIMARY KEY), `value` (TEXT) |
| key 格式 | `{publisher}.{extension-name}`（如 `RooVeterinaryInc.roo-cline`） |

## 查詢方法

```bash
# 列出所有 extension 的 key
sqlite3 ~/Library/Application\ Support/Code/User/globalStorage/state.vscdb \
  "SELECT key FROM ItemTable WHERE key LIKE '%.%' ORDER BY key;"

# 查看特定 extension 的設定（JSON 格式）
sqlite3 ~/Library/Application\ Support/Code/User/globalStorage/state.vscdb \
  "SELECT value FROM ItemTable WHERE key = 'RooVeterinaryInc.roo-cline';" | \
  python3 -m json.tool --no-ensure-ascii
```

## 修改方法

### 方法一：Python 腳本（建議）

```python
import json, sqlite3

db_path = os.path.expanduser(
    "~/Library/Application Support/Code/User/globalStorage/state.vscdb"
)

conn = sqlite3.connect(db_path)
conn.execute("PRAGMA journal_mode=WAL")  # 避免鎖定問題
cursor = conn.cursor()

# 讀取
cursor.execute(
    "SELECT value FROM ItemTable WHERE key = ?",
    ("RooVeterinaryInc.roo-cline",)
)
d = json.loads(cursor.fetchone()[0])

# 修改
d['language'] = 'zh-CN'  # 範例：改語言

# 寫回
new_value = json.dumps(d, ensure_ascii=False)
cursor.execute(
    "UPDATE ItemTable SET value = ? WHERE key = ?",
    (new_value, "RooVeterinaryInc.roo-cline")
)
conn.commit()
conn.close()
```

### 方法二：直接使用 sqlite3 CLI

```bash
# 匯出 JSON
sqlite3 ~/Library/Application\ Support/Code/User/globalStorage/state.vscdb \
  "SELECT value FROM ItemTable WHERE key = 'RooVeterinaryInc.roo-cline';" \
  > /tmp/roo_state.json

# 編輯 JSON（用任何編輯器）
vim /tmp/roo_state.json

# 匯入回資料庫
sqlite3 ~/Library/Application\ Support/Code/User/globalStorage/state.vscdb \
  "UPDATE ItemTable SET value = '$(cat /tmp/roo_state.json)' WHERE key = 'RooVeterinaryInc.roo-cline';"
```

## 注意事項

1. **VS Code 必須關閉或使用 WAL mode**：如果 VS Code 正在執行，資料庫可能被鎖定。使用 `PRAGMA journal_mode=WAL` 可以避免此問題。
2. **修改後需要重啟 VS Code**：Extension 通常在啟動時讀取 globalState，修改後需要重新載入視窗（Cmd+Shift+P → "Developer: Reload Window"）才會生效。
3. **JSON 結構不可破壞**：如果寫入無效的 JSON，extension 可能無法正常啟動。建議先備份。
4. **部分欄位會被覆蓋**：某些 extension 會在執行時重新寫入特定欄位（如 Roo Code 的 `customModes` 會被 `.roomodes` 覆蓋），修改這類欄位無效。
5. **機密資訊加密儲存**：API key 等機密資訊使用 `secret://` 前綴的 key 加密儲存，無法直接讀取或修改。

## 適用場景

| 場景 | 說明 | 適用工具 |
|------|------|----------|
| 修改語言偏好 | 改變 extension 的 UI/回應語言 | Roo Code |
| 調整允許指令 | 修改指令白名單 | Roo Code |
| 重置設定 | 清除特定 extension 的所有設定 | 通用 |
| 遷移設定 | 在不同機器間複製設定 | 通用 |

## 備份與還原

```bash
# 備份
cp ~/Library/Application\ Support/Code/User/globalStorage/state.vscdb \
   ~/Library/Application\ Support/Code/User/globalStorage/state.vscdb.backup

# 還原
cp ~/Library/Application\ Support/Code/User/globalStorage/state.vscdb.backup \
   ~/Library/Application\ Support/Code/User/globalStorage/state.vscdb
```

## 與其他層的關係

- `ai-tools/agent/roo.md` 記錄 Roo Code 專屬的 key 與欄位（如 `RooVeterinaryInc.roo-cline`、`language` 欄位），以及如何與 `.roomodes` 搭配使用。
- `enforcement/failure-patterns/language-preference-drift.md` 記錄語言偏好漂移的失效模式。
- 本文件只記錄 VS Code Extension 全域設定的通用知識（資料庫位置、結構、查詢/修改方法），不包含特定工具的設定細節。

---

← [回到 Engineering Intelligence 索引](../README.md)
