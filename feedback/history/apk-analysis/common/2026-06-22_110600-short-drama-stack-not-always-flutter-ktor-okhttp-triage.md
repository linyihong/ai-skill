> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Vertical-video category does not imply Flutter stack

Status: candidate

#### One-line Summary

同為短劇／豎屏視頻品類，技術棧可能是 Java/Kotlin + Ktor/OkHttp + 原生播放器 SDK，而非 Flutter Dart AOT；triage 必須先查 `libapp.so`，不可因品類或上一個專案經驗就套用 Blutter 主線。

#### Human Explanation

若 workspace 內已有 Flutter 短視頻分析先例，agent 容易對**同品類新 target** 產生 path bias，直接上 Blutter / Dart Frida。實務上也可能為原生多 DEX + Ktor Client + OkHttp3 + 商用播放器 native。錯誤主線會導致 MITM/Java hook 策略選錯、浪費 Blutter 時間、漏掉 Ktor pipeline hook 點。

#### Trigger

- 任務描述含 short-form video / episode unlock / vertical feed，且 agent 未驗證就假設 Flutter。
- `unzip -l` **無** `libapp.so` / `libflutter.so`，但有大量 `classes*.dex` 與播放器/downloader native。
- dex strings 出現 `io/ktor`、`okhttp3`、`retrofit2` 等 Java/Kotlin HTTP 棧，而非 Dart AOT 特徵。

#### Evidence

- Tool: `unzip -l`, `aapt dump badging`, dex strings grep（類別名、HTTP client、路徑 prefix 模式）。
- Sanitized excerpt: 多個 `classes*.dex`；native 含 player/downloader 類 lib；HTTP 棧為 Ktor + OkHttp；API 路徑在 dex 中以 `/api/<segment>/` 形式出現（**具體 path/host 留 project docs**）。
- Evidence path: stack 摘要與 path inventory 留在 `<PROJECT_ROOT>/docs/`；**本 lesson 不含** publisher host、完整 endpoint 列表、schema、token 或 class 真名。

#### Generalized Lesson

**Stack triage 決策樹（品類無關）：**

```text
libapp.so 存在？
  是 → Flutter/Dart AOT 主線（blutter + Dart Frida）
  否 → Java/Kotlin 主線
        ├─ okhttp3 / ktor / retrofit → Java/Ktor MITM + hook
        ├─ player/downloader native libs → 媒體：控制面 API vs CDN 分線
        └─ protect/armor 類 native → 預留 pinning / anti-tamper 排查
```

**品類 ≠ stack**；上一個專案的 Flutter 結論不可移植。

#### Agent Action

1. 靜態 triage 前：`unzip -l | grep -E 'libapp|libflutter'` — 結果決定主線。
2. 原生主線：jadx + OkHttp/Ktor Frida；path 線索用 dex strings 的 **prefix 模式**，完整 route 寫 project API docs。
3. 寫入 Ai-skill 時只保留決策樹與工具選擇；**不得**寫入 target 專屬 host、service 名、簽章材料或 capture 片段。

#### Goal / Action / Validation

- Goal: 縮短 stack triage 誤判時間。
- Action: `analysis/apk/traffic-triage.md` 強調「libapp 檢查優先於品類假設」。
- Validation: 無 libapp 時 Blutter 被 explicit skip；Java hook 或 MITM 有命中即驗證主線。

#### Applies When

- 新 short-form video / episode-gated content APK 分析。
- Workspace 內存在 Flutter 先例（path bias 風險）。

#### Does Not Apply When

- 已確認 `libapp.so` + `libflutter.so`。
- Unity / React Native 等（各自 triage 表）。

#### Validation

- 報告明確寫 Flutter yes/no 及**檔案列表依據**，不是品類推測。
- Ai-skill lesson diff 通過去敏自查（無 host、無 package 真值、無本機路徑）。

#### Promotion Target

- `analysis/apk/traffic-triage.md` §主線選擇
- `workflow/apk-analysis/execution-flow.md` §2（Dart AOT 步驟前先確認 libapp）

#### Required Linked Updates

- 依 [`linked-updates.md`](../../../../enforcement/linked-updates.md)：`feedback/history/apk-analysis/README.md` 索引待追加。
- 已依 [`sanitization.md`](../../../../enforcement/sanitization.md) 與 [`reusable-guidance-boundary.md`](../../../../enforcement/reusable-guidance-boundary.md) 檢查。
