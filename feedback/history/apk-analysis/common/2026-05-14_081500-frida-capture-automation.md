> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-14 - Frida capture 應自動化：建立 reusable shell script 處理完整 lifecycle

Status: validated

#### One-line Summary

重複性的 Frida capture 操作（clear app state → spawn → capture → save log）應建立 **reusable shell script** 自動化，而不是每次手動輸入命令。自動化腳本應處理：timestamped log 檔名、adb pm clear、Frida spawn、tee 輸出、log 儲存。這減少手動錯誤、確保 capture 可重現、並釋放注意力給分析本身。

#### Human Explanation

在多次 Frida capture 操作中，每次都需要手動執行以下步驟：

1. `adb shell pm clear <package>` — 清除 App 狀態
2. `frida -U -f <package> -l <script>` — spawn 模式啟動
3. 等待 capture 完成
4. 手動儲存 log 到 capture 目錄

這個流程有幾個問題：
- 每次都要記住完整的命令（package name、script path）
- Log 檔名不一致（有時用 timestamp、有時用手動命名）
- 容易忘記先 clear app state
- 無法確保 capture 的可重現性

解決方案：建立一個 **reusable shell script** 封裝完整 lifecycle：

```bash
#!/bin/bash
# auto_capture.sh — Automated Frida capture
SCRIPT="$1"
TIMESTAMP=$(date +%s)
LOGFILE="/tmp/frida_capture_${TIMESTAMP}.log"
adb shell pm clear <package>
frida -U -f <package> -l "$SCRIPT" 2>&1 | tee "$LOGFILE"
mkdir -p capture
cp "$LOGFILE" "capture/frida_capture_$(date +%Y%m%d_%H%M%S).log"
```

#### Trigger

- 需要多次執行 Frida capture（> 3 次）
- 每次手動輸入相同的命令序列
- Log 檔名不一致或難以追蹤
- 偶爾忘記執行 `adb shell pm clear`
- 需要與其他人分享 capture 流程

#### Evidence

- Tool: Reusable shell script (`auto_capture_iv.sh`)
- Sanitized excerpt:
  - Before: 手動輸入 `adb shell pm clear <package> && frida -U -f <package> -l script.js 2>&1 | tee /tmp/frida_capture_$(date +%s).log`
  - After: `./auto_capture_iv.sh hook_capture_iv.js`
  - Result: 自動 clear → spawn → capture → save log
- Evidence path: `<PROJECT_ROOT>/<target-app>/scripts/frida/auto_capture_iv.sh`

#### Generalized Lesson

1. **重複性操作應自動化**——如果一個操作序列需要執行 > 3 次，建立 reusable script。
2. **自動化腳本的最小功能**：
   - 接受 script path 作為參數
   - 自動產生 timestamped log 檔名
   - 執行 `adb shell pm clear` 確保乾淨狀態
   - 使用 `tee` 同時顯示和儲存輸出
   - 將 log 複製到專案的 capture 目錄
3. **腳本應可重入**——每次執行產生新的 log 檔案，不覆蓋舊的。
4. **腳本應有 usage message**——顯示如何使用，減少學習成本。
5. **自動化釋放注意力**——不需要記命令時，可以專注在分析輸出上。

#### Agent Action

建立 Frida capture 自動化腳本時：

1. **建立 shell script**——放在專案的 `scripts/frida/` 目錄下
2. **腳本應接受 Frida script path 作為參數**——不要 hardcode script 名稱
3. **自動處理 timestamped log**——使用 `date +%s` 或 `date +%Y%m%d_%H%M%S`
4. **包含 adb pm clear**——確保每次 capture 從乾淨狀態開始
5. **使用 tee 同時顯示和儲存**——方便即時監控
6. **將 log 複製到專案 capture 目錄**——方便後續分析

#### Goal / Action / Validation

- Goal: 減少 Frida capture 的手動操作，確保可重現性
- Action: 建立 reusable shell script 封裝完整 lifecycle
- Validation or reference source: 腳本執行後自動產生 timestamped log 檔案

#### Applies When

- 需要多次執行 Frida capture（> 3 次）
- capture 流程包含多個步驟（clear → spawn → capture → save）
- 需要確保 capture 的可重現性
- 與其他人協作需要一致的 capture 流程

#### Does Not Apply When

- 只需要執行一次 Frida capture
- capture 流程非常簡單（如 attach 模式 + 單一 hook）
- 使用 Frida Gadget 或其他自動化工具

#### Validation

- 腳本執行後自動產生 log 檔案在 `capture/` 目錄下
- Log 檔案包含完整的 capture 輸出
- 腳本可以重複執行，每次產生新的 log 檔案

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` — 新增步驟：「建立 Frida capture 自動化腳本」
- `workflow/apk-analysis/execution-flow.md` — 在 Quick Start 中加入自動化腳本範本

#### Required Linked Updates

- `<PROJECT_ROOT>/<target-app>/scripts/frida/auto_capture_iv.sh` — 已建立
