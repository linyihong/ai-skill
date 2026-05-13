#!/usr/bin/env python3
"""
Roo Code 全域 Custom Instructions 寫入腳本

使用方法：
  1. 先關閉 VS Code（重要！）
  2. 執行此腳本：python3 scripts/set-roo-global-custom-instructions.py
  3. 重新開啟 VS Code

注意：VS Code 開啟時執行此腳本無效，因為 VS Code 的 extension host
會管理 state.vscdb，直接寫入會被 VS Code 的記憶體狀態覆寫。
"""

import json
import sqlite3
import os
import subprocess
import sys

DB_PATH = os.path.expanduser(
    "~/Library/Application Support/Code/User/globalStorage/state.vscdb"
)

# 要寫入的全域 Custom Instructions 內容
CUSTOM_INSTRUCTIONS = """你是一個 AI 助手，遵循以下啟動流程：

## 啟動流程

1. 首先讀取 CORE_BOOTSTRAP.md 了解核心啟動流程
2. 依照 CORE_BOOTSTRAP.md 的指示載入依賴文件（dependency-reading.md）
3. 依照 dependency-reading.md 的規則載入 shared-rules 目錄下的必要規則
4. 根據任務類型載入對應的 workflow、intelligence、analysis 文件

## 語言偏好

- 預設使用英文回應
- 如果使用者使用中文提問，則用中文回應
- 如果使用者切換語言，跟隨其切換
- 所有輸出（包含技術分析、表格欄位、章節標題）都必須與使用者當前語言一致

## 對話目標閉環

- 每個對話必須有明確的目標
- 完成目標後使用 attempt_completion 呈現結果
- 不要在結果結尾提出問題或要求進一步協助

## 專案設定

此 Ai-skill 專案位於使用者的 Documents/Ai-skill 目錄。
當使用者在此專案中工作時，遵循專案內的 .roomodes 和 shared-rules 設定。
當使用者在其他專案中工作時，此全域設定確保啟動流程仍然有效。
"""


def is_vscode_running():
    """檢查 VS Code 是否正在執行"""
    try:
        result = subprocess.run(
            ["pgrep", "-f", "Visual Studio Code"],
            capture_output=True,
            text=True,
            timeout=5,
        )
        return result.returncode == 0 and len(result.stdout.strip()) > 0
    except (subprocess.TimeoutExpired, FileNotFoundError):
        return False


def main():
    # 步驟 1：檢查 VS Code 是否正在執行
    if is_vscode_running():
        print("⚠️  VS Code 正在執行中！")
        print("   直接寫入 state.vscdb 會被 VS Code 的記憶體狀態覆寫。")
        print()
        print("   請先關閉 VS Code（Cmd+Q），然後重新執行此腳本。")
        print("   關閉後執行：python3 scripts/set-roo-global-custom-instructions.py")
        sys.exit(1)

    # 步驟 2：確認資料庫檔案存在
    if not os.path.exists(DB_PATH):
        print(f"❌ 找不到資料庫檔案：{DB_PATH}")
        print("   請確認 Roo Code 擴充功能已安裝並至少執行過一次。")
        sys.exit(1)

    print(f"📁 資料庫路徑：{DB_PATH}")

    # 步驟 3：讀取現有 JSON
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()

    cursor.execute(
        "SELECT value FROM ItemTable WHERE key = 'RooVeterinaryInc.roo-cline'"
    )
    row = cursor.fetchone()

    if not row:
        print("❌ 找不到 Roo Code 設定資料。")
        print("   請確認 Roo Code 擴充功能已至少執行過一次。")
        conn.close()
        sys.exit(1)

    data = json.loads(row[0])
    print(f"📊 讀取到 {len(data)} 個設定鍵")

    # 步驟 4：寫入 customInstructions
    old_value = data.get("customInstructions", None)
    data["customInstructions"] = CUSTOM_INSTRUCTIONS

    new_value = json.dumps(data, ensure_ascii=False)
    cursor.execute(
        "UPDATE ItemTable SET value = ? WHERE key = 'RooVeterinaryInc.roo-cline'",
        (new_value,),
    )
    conn.commit()

    # 步驟 5：強制 WAL checkpoint
    cursor.execute("PRAGMA wal_checkpoint(TRUNCATE);")
    conn.close()

    # 步驟 6：驗證寫入
    conn2 = sqlite3.connect(DB_PATH)
    cursor2 = conn2.cursor()
    cursor2.execute(
        "SELECT value FROM ItemTable WHERE key = 'RooVeterinaryInc.roo-cline'"
    )
    row2 = cursor2.fetchone()
    data2 = json.loads(row2[0])
    has_ci = "customInstructions" in data2
    conn2.close()

    if has_ci:
        ci_len = len(data2["customInstructions"])
        print(f"✅ 寫入成功！customInstructions 長度：{ci_len} 字元")
        if old_value:
            print(f"   舊值長度：{len(old_value)} 字元")
        print()
        print("🚀 現在可以重新開啟 VS Code 了。")
        print("   開啟後，Roo Code 在任何專案中都會載入全域 Custom Instructions。")
    else:
        print("❌ 寫入失敗！customInstructions 欄位不存在。")
        print("   請檢查資料庫權限或重新執行。")
        sys.exit(1)


if __name__ == "__main__":
    main()
