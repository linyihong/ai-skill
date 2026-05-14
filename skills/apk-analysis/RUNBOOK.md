# 新 APK 專案使用 Runbook

這份 runbook 給新專案第一天使用。目標是讓 AI 或人類拿到一個新 APK 後，知道怎麼套用 `apk-analysis` 的分析方法，並把新專案學到的通用技巧回饋回這個知識庫。

canonical repository 是你本機 clone 的 `Ai-skill` repository。下文用 `<AI_SKILL_REPO>` 表示該路徑：

```text
<AI_SKILL_REPO>
```

完成新專案分析後，請在 `<AI_SKILL_REPO>` commit 並 push。

## 放置位置

### 0. 這台機器上的共用 repo

優先讀：

```text
<AI_SKILL_REPO>/runtime/onboarding/apk-analysis-setup.md    # 開場提示詞與放置方式
<AI_SKILL_REPO>/workflow/apk-analysis/execution-flow.md     # 分析執行流程
<AI_SKILL_REPO>/workflow/apk-analysis/artifact-gates.md     # 產出格式與完成定義
<AI_SKILL_REPO>/analysis/apk/workflows/                     # 操作細節流程（依分類）
<AI_SKILL_REPO>/intelligence/engineering/analytical-reasoning/      # 決策智慧（heuristics / anti-patterns / signals / failure）
```

新技巧回饋到：

```text
<AI_SKILL_REPO>/shared-rules/feedback-lessons.md                          # 檔名規則與 lesson 模板（全庫共用）
<AI_SKILL_REPO>/feedback/history/apk-analysis/<category>/                 # 每一條 lesson 獨立檔（YYYY-MM-DD_HHMMSS-<slug>.md）
```

驗證後同步更新：

```text
<AI_SKILL_REPO>/workflow/apk-analysis/execution-flow.md
<AI_SKILL_REPO>/workflow/apk-analysis/artifact-gates.md
<AI_SKILL_REPO>/analysis/apk/workflows/<workflow>.md
<AI_SKILL_REPO>/intelligence/engineering/analytical-reasoning/<atom-type>/<name>.md
```

### 1. 當作專案文件使用

將以下路徑加入專案的 AI 工具設定（如 `.roomodes`、`CLAUDE.md`、`.cursorrules`）：

```text
<AI_SKILL_REPO>/runtime/onboarding/apk-analysis-setup.md
<AI_SKILL_REPO>/workflow/apk-analysis/execution-flow.md
<AI_SKILL_REPO>/workflow/apk-analysis/artifact-gates.md
<AI_SKILL_REPO>/shared-rules/feedback-lessons.md
```

使用時要明確告訴 AI：

```text
請先閱讀 <AI_SKILL_REPO>/runtime/onboarding/apk-analysis-setup.md，
並依照 workflow/apk-analysis/execution-flow.md 與 artifact-gates.md 分析這個 APK。
如果過程中學到可重用技巧，請依 shared-rules/feedback-lessons.md 在 feedback/history/apk-analysis/<category>/ 新增檔案回饋；跨分類用 feedback/history/apk-analysis/common/。
```

### 2. 當作工具可讀 skill 使用

將 `workflow/apk-analysis/`、`analysis/apk/`、`intelligence/engineering/analytical-reasoning/` 與 `shared-rules/` 成對部署到你使用的 AI / agent 工具，或在工具提示中明確指定 `<AI_SKILL_REPO>` 內的路徑。工具專屬部署方式請看 [`../../ai-tools/`](../../ai-tools/README.md)，不要把工具專屬路徑寫進通用流程。

## 開場提示詞

新專案開始時，可以直接給 AI：

```text
請使用 apk-analysis 分析方法協助分析這個 APK。

目標：
- 先判斷核心流量走哪一層。
- 不要一開始假設是 pinning。
- 先做 localhost / pcap / proxy / Java hook / native or Flutter 判斷。
- 所有 token、device id、私密 host、個資都要去敏。
- 如果發現新技巧，請**主動**在 <AI_SKILL_REPO>/feedback/history/apk-analysis/<category>/ **新增 lesson 檔**（跨分類用 common/；格式見 <AI_SKILL_REPO>/shared-rules/feedback-lessons.md；不要等使用者說「記得回饋」），且要讓人類也看得懂；可同步更新 workflow/apk-analysis/ 或 intelligence/engineering/analytical-reasoning/ 如已驗證。

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

5. 分類路由（見 `intelligence/engineering/analytical-reasoning/evidence-first-routing.md`）：
   - 先用共通流程判斷 runtime / traffic family。
   - 證據指向 Flutter/Dart AOT 才讀 `analysis/apk/workflows/frida-hook-flow.md`。
   - 證據指向 local proxy / loopback 才讀 `analysis/apk/workflows/local-proxy-hook-flow.md`。
   - 目標是 HTTP API 文件化才讀 `analysis/apk/workflows/http-api-documentation-flow.md`。
   - 目標是 HLS/media 才讀 `analysis/apk/workflows/media-hls-analysis-flow.md`。
   - 不要在分類未明時一次讀完所有 workflow folders。

6. 文件化：
   - 失敗路徑也要寫
   - 成功證據要可重現
   - API 結論和方法論要分開
   - 去敏後才保存樣本
   - 若下一步是 SDK/client/app tool/live integration 開發，先補 project-level domain/runtime baseline 的最小可跑因素；若需要 device/install/account/session/vendor/server-issued material，補 authorized identity material self-generation audit，回答能否自生成、怎麼生成或由誰提供、reset/cooldown/error 如何驗證；缺項要列 blocker，不要直接進 live-facing code

7. 回饋 skill：
   - 新技巧寫入 `feedback/history/apk-analysis/<category>/` 或 `feedback/history/apk-analysis/common/`（見 `shared-rules/feedback-lessons.md` 命名規則）
   - 使用者或 reviewer 提出的可泛化操作優化、失敗模式或驗證規則，也算回饋 trigger；不要等到完全驗證才記錄，可先標 `candidate` / `experimental`
   - 已驗證技巧再同步進 `workflow/apk-analysis/execution-flow.md`、`artifact-gates.md` 或 `intelligence/engineering/analytical-reasoning/` 對應分類

## 回饋規則

新專案得到的經驗可以回饋，但要先判斷是不是通用技巧。

應該回饋：

- 新的抓包判斷順序。
- 新的 Frida / proxy / pcap 失敗模式。
- 新的 Flutter / native / Java stack 判讀方式。
- 新的去敏與 fixture 沉澱方式。
- 新的媒體或解密驗證方法。
- 使用者或 reviewer 指出的可泛化 UI/API 操作優化、停止條件、驗證門檻或閉環缺口。

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
如果使用者或 reviewer 提出一個可能跨 APK 重用的操作改進、失敗模式或驗證規則，請先做泛化判斷；能泛化就立刻進 feedback loop。
請依 <AI_SKILL_REPO>/shared-rules/feedback-lessons.md，在同一輪對話於 feedback/history/apk-analysis/<category>/ 或 feedback/history/apk-analysis/common/ **新增一個** lesson 檔。
如果 lesson 尚未完全驗證，先標 candidate / experimental 並寫清楚 validation criteria；如果已經被驗證，也請同步更新 workflow/apk-analysis/execution-flow.md、artifact-gates.md 或 intelligence/engineering/analytical-reasoning/ 對應分類。
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
- 若目標包含 SDK/client/replay/live integration：domain/runtime baseline 已回答最小可跑因素，或缺口已列 blocker / scoped out；僅 skeleton 時不得宣稱可開始 live-facing 開發。
- 是否有新 lesson 回饋到 skill。
