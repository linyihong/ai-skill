### 2026-04-30 - APK metadata：`aapt` 不在 PATH 時走 SDK build-tools；launcher 用 `resolve-activity`

Status: validated

#### One-line Summary

`apkanalyzer` 或環境找不到 build-tools 時，改用 `$ANDROID_HOME/build-tools/<version>/aapt`；badging 若沒有 launcher 行，用 `cmd package resolve-activity` 取得 `am start -n` 所需的完整 component。

#### Human Explanation

許多機器只把 `adb`（platform-tools）放在 PATH，但沒有 `aapt`。`apkanalyzer` 仍會嘗試在 SDK 內解析 `aapt`，若本機 SDK 配置不完整，可能拋出「Cannot locate latest build tools」之類錯誤。此時不必急著裝 jadx：只要已安裝 Android SDK，`build-tools` 目錄裡通常已有對應版本的 `aapt`，直接呼叫即可取得 `package`、`versionName`、`permissions`、`native-code` 等盤點資訊。

另外，`aapt dump badging` 的輸出有時**沒有** `launchable-activity:` 行（多重 activity、工具版本或 manifest 複雜度都可能造成）。這不代表無法冷啟動：在**已安裝該 package 的裝置**上，`adb shell cmd package resolve-activity --brief <package>` 常能給出預設 launcher 的 `package/class`，可用於腳本化的 `am start -n`，也比依賴 `monkey` 更穩定。

#### Trigger

- `aapt` / `aapt2`：command not found。
- `apkanalyzer`：報錯無法定位 build tools / `aapt`。
- `aapt dump badging app.apk | grep launchable-activity` 無結果，但需要自動化啟動 App。

#### Evidence

- Tool：`apkanalyzer`（依賴 SDK 內部 `aapt` 解析）。
- Tool：`aapt dump badging`（部分 APK 無 `launchable-activity` 行）。
- Tool：`adb shell cmd package resolve-activity --brief <package>`（回傳形如 `pkg/component.name`）。
- Sanitized excerpt：`IllegalStateException: Cannot locate latest build tools`（`apkanalyzer`）；`resolve-activity` 成功回傳預設 activity component。

#### Generalized Lesson

1. **Metadata fallback**：PATH 無 `aapt` 時，優先嘗試 `$ANDROID_HOME/build-tools/<任意已裝版本>/aapt`（版本目錄可用 `ls` 選最新穩定版），再考慮安裝獨立工具鏈。
2. **Launcher fallback**：需要可重現的 `am start` 時，若 badging 無 launcher，可在裝置上用 `cmd package resolve-activity --brief`；靜態-only 環境則仍可用 `aapt dump xmltree AndroidManifest.xml` 搭配 intent-filter 判讀（較費工）。

#### Agent Action

- 盤點 APK 時若 `aapt` 不在 PATH，先檢查 `$ANDROID_HOME/build-tools`，用該路徑下的 `aapt dump badging`。
- 不要僅因 `apkanalyzer` 失敗就判定「無法讀 manifest metadata」。
- 撰寫冷啟動／代理測試腳本時，優先用 `resolve-activity` 或明確的 `-n pkg/activity`；`monkey` 僅作備用。
- 技能文件與 FEEDBACK 中**不要**寫入使用者本機絕對路徑；用 `$ANDROID_HOME` 與占位符描述。

#### Applies When

- 已安裝 Android SDK（含 `build-tools`），且需要快速 badging。
- 裝置上已安裝目標 App，需要 launcher component 做自動化。

#### Does Not Apply When

- 完全沒有 Android SDK／沒有 `build-tools`（需改用人類安裝 SDK、或改用其他解析器）。
- 僅有 APK 檔、無法 adb 到裝置：launcher 需靠 xmltree／反編譯推斷。

#### Validation

- `aapt ... dump badging` 輸出含 `package:` 與 `versionName`，且與裝置 `dumpsys package` 一致（同一簽署／同一版）。
- `am start -n pkg/activity` 能啟動 App 到預設桌面入口。

#### Promotion Target

- `TOOLS.md`（命令模板、常見失敗判讀）
