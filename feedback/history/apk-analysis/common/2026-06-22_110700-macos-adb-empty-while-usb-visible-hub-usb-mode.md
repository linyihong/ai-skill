> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - macOS adb empty while USB device visible (hub + USB mode)

Status: candidate

#### One-line Summary

macOS `system_profiler` 已看到 Android 裝置，但 `adb devices` 為空時，優先查 USB 模式（僅充電）與是否經 USB Hub 連線，而非假設是授權對話框問題。

#### Human Explanation

使用者已在手機上允許 USB 偵錯，但 adb 仍完全看不到裝置（不是 `unauthorized`，是**空列表**）。常見原因：(1) USB 用途仍為僅充電，adb 介面未 enumerate；(2) 經 USB Hub 連 host，adb interface 不穩；(3) host 端 Android File Transfer 佔用 USB。Wireless debugging（Android 11+）是可備援路線。

#### Trigger

- `adb devices` 空，但 `system_profiler SPUSBDataType` 仍列出 Android 裝置。
- 使用者稱已授權，狀態從未出現 `unauthorized`。
- `adb wait-for-device` 長時間阻塞或 protocol fault。

#### Evidence

- Tool: `adb devices -l`, `system_profiler SPUSBDataType`, `adb kill-server && adb start-server`。
- Sanitized excerpt: USB 層可見裝置；adb 空；改 USB 檔案傳輸 + 直插 host 後出現 `device`；或 wireless adb pair/connect 成功。
- Evidence path: 裝置型號、序號、hub 型號等 **incident 細節只留 `<PROJECT_ROOT>/` install/runbook**；lesson 不含 serial、本機路徑或使用者名稱。

#### Generalized Lesson

**診斷順序：**

```text
adb devices 空 + USB 可見？
  1. 通知欄 USB → 檔案傳輸 / Android Auto（非僅充電）
  2. 開發者選項 → 預設 USB 設定 → 檔案傳輸
  3. 撤銷 USB 偵錯授權 → 重插 → 解鎖螢幕再允許
  4. 避開 Hub，直插 host
  5. killall Android File Transfer* ; adb kill-server
  6. 仍失敗 → Wireless debugging: adb pair → adb connect
```

`unauthorized`（需對話框）與**空列表**（USB mode / enumerate）是不同故障模式。

#### Agent Action

1. 不要對空列表無限 `adb wait-for-device`；加 timeout 與 USB mode / hub 提示。
2. 診斷腳本放 `<PROJECT_ROOT>/scripts/`（project-local）；**Ai-skill lesson 不嵌入** script 內的真實 serial 或絕對路徑。
3. install 腳本用 polling + 上限 timeout，避免 session 掛死。

#### Goal / Action / Validation

- Goal: 縮短 APK 安裝 / capture 前置阻塞時間。
- Action: `runtime/onboarding/apk-analysis-setup.md` 補「USB visible but adb empty」分支。
- Validation: `adb devices` 顯示 `device` 後 install/pull 成功。

#### Applies When

- macOS + 實體 Android + adb 前置。
- 已聲稱授權但 adb 仍空。

#### Does Not Apply When

- 狀態為 `unauthorized` 或 `offline`。
- 純 wireless adb 且 USB 未使用。

#### Validation

- 直插或改 USB 模式後 adb 出現 `device`。
- Lesson 正文無 device serial、hub 品牌型號、本機路徑。

#### Promotion Target

- `runtime/onboarding/apk-analysis-setup.md`

#### Required Linked Updates

- 依 [`linked-updates.md`](../../../../enforcement/linked-updates.md)：`feedback/history/apk-analysis/README.md` 索引待追加。
- 已依 [`sanitization.md`](../../../../enforcement/sanitization.md) 檢查。
