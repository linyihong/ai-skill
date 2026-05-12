# APK Analysis Workflows

`analysis/apk/workflows/` 存放 APK 分析的操作流程與執行步驟。

## Scope

本目錄負責 **HOW TO DO** — execution knowledge：

- Frida hook 操作流程（command、setup、adb、proxy）
- 代理導流設定步驟
- 靜態分析操作步驟
- 媒體串流驗證操作步驟
- 常見 dump 方法與命令模板

本目錄不負責：

- **HOW TO THINK** 決策智慧（何時該用哪個技術）；使用 `intelligence/engineering/apk-analysis/`
- 端到端 agent 執行流程；使用 `workflow/apk-analysis/`
- 工具選擇與失敗判讀；使用 `analysis/apk/tools-and-failures.md`

## 與其他層的關係

- `analysis/apk/traffic-triage.md` 決定流量分流路線，本目錄提供路線內的具體操作
- `analysis/apk/tools-and-failures.md` 提供工具選擇與命令模板，本目錄提供完整操作流程
- `intelligence/engineering/apk-analysis/heuristics/` 提供「何時該用哪個流程」的判斷
- `workflow/apk-analysis/` 提供端到端 agent 執行流程，本目錄提供單一技術的深度操作

## 目前 workflows

（pilot 階段逐步建立）
