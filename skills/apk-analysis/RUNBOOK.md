# 新 APK 專案使用 Runbook

這份 runbook 給新專案第一天使用。目標是讓 AI 或人類拿到一個新 APK 後，知道怎麼套用 `apk-analysis` skill，並把新專案學到的通用技巧回饋回這個包。

canonical skill repository 是你本機 clone 的 `Ai-skill` repository。下文用 `<AI_SKILL_REPO>` 表示該路徑：

```text
<AI_SKILL_REPO>
```

如果新專案也放了一份 copy，請以 `<AI_SKILL_REPO>/skills/apk-analysis/` 作為最終回饋位置，完成後在 `<AI_SKILL_REPO>` commit 並 push。

## 放置位置

有兩種用法：

### 0. 這台機器上的共用 skill repo

優先讀：

```text
<AI_SKILL_REPO>/skills/apk-analysis/SKILL.md
<AI_SKILL_REPO>/skills/apk-analysis/RUNBOOK.md
```

新技巧回饋到：

```text
<AI_SKILL_REPO>/shared-rules/feedback-lessons.md         # 檔名規則與 lesson 模板（全庫共用）
<AI_SKILL_REPO>/skills/apk-analysis/feedback_history/   # 每一條 lesson 獨立檔（YYYY-MM-DD_HHMMSS-<slug>.md）
```

驗證後同步更新：

```text
<AI_SKILL_REPO>/skills/apk-analysis/WORKFLOW.md
<AI_SKILL_REPO>/skills/apk-analysis/TOOLS.md
<AI_SKILL_REPO>/skills/apk-analysis/DOCUMENTATION.md
```

### 1. 當作專案文件使用

放在：

```text
<AI_SKILL_REPO>/skills/apk-analysis/
```

使用時要明確告訴 AI：

```text
請先閱讀 <AI_SKILL_REPO>/skills/apk-analysis/SKILL.md，
並依照 WORKFLOW.md / TOOLS.md / DOCUMENTATION.md 分析這個 APK。
如果過程中學到可重用技巧，請依 shared-rules/feedback-lessons.md 在 feedback_history/ 新增檔案回饋。
```

### 2. 當作 Cursor project skill 使用

放在：

```text
.cursor/skills/apk-analysis/
.cursor/shared-rules/    # 請自中央庫複製整包 shared-rules（見 shared-rules/cursor-sync.md）
```

這樣 Cursor agent 較容易在你提到 APK 分析、抓包、Frida、Proxyman、Dart AOT、解密時自動套用 skill；共用規則另放在 `.cursor/shared-rules/` 與 skill 並列。

如果未來要跨專案共用，也可以放到個人技能資料夾：

```text
~/.cursor/skills/apk-analysis/
```

不要放到 Cursor 內建技能資料夾。

## 開場提示詞

新專案開始時，可以直接給 AI：

```text
請使用 apk-analysis skill 協助分析這個 APK。

目標：
- 先判斷核心流量走哪一層。
- 不要一開始假設是 pinning。
- 先做 localhost / pcap / proxy / Java hook / native or Flutter 判斷。
- 所有 token、device id、私密 host、個資都要去敏。
- 如果發現新技巧，請**主動**在 <AI_SKILL_REPO>/skills/apk-analysis/feedback_history/ **新增 lesson 檔**（格式見 <AI_SKILL_REPO>/shared-rules/feedback-lessons.md；不要等使用者說「記得回饋」），且要讓人類也看得懂；可同步更新 TOOLS/WORKFLOW 如已驗證。

目前材料：
- APK:
- version:
- package:
- device / emulator:
- root / Frida:
- proxy / MITM:
- 想先分析的功能流程:
```

## 第一輪分析順序

1. 盤點 APK：
   - package name
   - version
   - architecture
   - permissions
   - native libraries
   - Flutter / React Native / Java/Kotlin / native 初步判斷

2. 建立分析環境紀錄：
   - device / emulator
   - adb
   - Frida client/server
   - proxy tool
   - root status
   - allowed actions

3. 流量路徑判斷：
   - localhost / loopback
   - whole-device pcap
   - system proxy / MITM
   - Java HTTP stack
   - native connect trace
   - Flutter / Dart AOT if applicable

4. 找高語意 hook：
   - request options
   - response wrapper
   - response decoder / decryptor
   - token/session provider

5. 文件化：
   - 失敗路徑也要寫
   - 成功證據要可重現
   - API 結論和方法論要分開
   - 去敏後才保存樣本

6. 回饋 skill：
   - 新技巧寫入 `feedback_history/`（見 `shared-rules/feedback-lessons.md` 命名規則）
   - 已驗證技巧再同步進 `WORKFLOW.md` / `TOOLS.md` / `DOCUMENTATION.md`

## 回饋規則

新專案得到的經驗可以回饋，但要先判斷是不是通用技巧。

應該回饋：

- 新的抓包判斷順序。
- 新的 Frida / proxy / pcap 失敗模式。
- 新的 Flutter / native / Java stack 判讀方式。
- 新的去敏與 fixture 沉澱方式。
- 新的媒體或解密驗證方法。

不應該回饋：

- 特定產品的 host。
- endpoint 全名。
- token、secret、device id。
- 真實帳號或個資。
- 本機絕對路徑、使用者名稱、私有工作目錄、clone 位置。
- 只對單一產品有效的業務結論。

## 讓 AI 自動回饋的提示

在新專案中可以加這句：

```text
每次你發現新的可重用 APK 分析技巧時，請不要只在對話裡說明；也請不要等使用者提醒「記得回饋」。
請依 <AI_SKILL_REPO>/shared-rules/feedback-lessons.md，在同一輪對話於 feedback_history/ **新增一個** lesson 檔。
如果 lesson 已經被驗證，也請同步更新 WORKFLOW.md、TOOLS.md 或 DOCUMENTATION.md。
注意：lesson 內容要讓人類也能看懂，不只給 AI 看。
完成後請在 <AI_SKILL_REPO> commit 並 push。
```

## 完成定義

一次新 APK 初步分析完成時，至少應該得到：

- 流量路徑判斷結果。
- 代理 / MITM 是否可用的證據。
- Java / native / Flutter stack 判斷。
- 初步 request metadata 或下一步 hook 計畫。
- response wrapper 或解密定位計畫。
- 去敏規則與文件位置。
- 是否有新 lesson 回饋到 skill。
