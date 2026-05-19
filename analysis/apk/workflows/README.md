# APK 分析工作流程

`analysis/apk/workflows/` 存放 APK 分析中單一技術路線的操作流程與證據取得步驟。

命名上的 `workflows/` 是歷史相容路徑：本目錄屬於 `analysis/apk/`，負責技術方法內部的 procedure；端到端 agent orchestration、任務分解、交接與完成門檻仍由 `workflow/apk-analysis/` 負責。

## Scope

本目錄負責 **HOW TO COLLECT EVIDENCE** — analysis-local procedure：

- Frida hook 操作流程（command、setup、adb、proxy）
- 代理導流設定步驟
- 靜態分析操作步驟
- 媒體串流驗證操作步驟
- 常見 dump 方法與命令模板

本目錄不負責：

- **HOW TO THINK** 決策智慧（何時該用哪個技術）；使用 `intelligence/engineering/analytical-reasoning/`
- 端到端 agent 執行流程；使用 `workflow/apk-analysis/`
- 工具選擇與失敗判讀；使用 `analysis/apk/tools-and-failures.md`
- Raw case logs、App 專屬 class dump、host、endpoint、token 或一次性 reverse 過程；保留在業務專案 evidence。

## 與其他層的關係

- `analysis/apk/traffic-triage.md` 決定流量分流路線，本目錄提供路線內的具體操作
- `analysis/apk/tools-and-failures.md` 提供工具選擇與命令模板，本目錄提供完整操作流程
- `intelligence/engineering/analytical-reasoning/heuristics/` 提供「何時該用哪個流程」的判斷
- `workflow/apk-analysis/` 提供端到端 agent 執行流程，本目錄提供單一技術的深度操作

## 目前 workflows

| Workflow | 說明 | 來源 |
|----------|------|------|
| [`frida-hook-flow.md`](frida-hook-flow.md) | Frida Hook 操作流程 — 6 步驟（確認 Flutter/Dart AOT、搜尋關鍵字、Hook Request Options、Hook Response Decode/Decrypt、Dart String Decoding、對齊與去敏） | `skills/apk-analysis/techniques/flutter-dart-aot/`（已刪除，見 `frida-hook-flow.md`） |
| [`http-api-documentation-flow.md`](http-api-documentation-flow.md) | HTTP API 文件化操作流程 — 7 步驟（API Entry → Group Index → Per-API Detail → Coverage/Gap Matrix → SDK Mapping → Finish Gate → UI Automation） | `skills/apk-analysis/techniques/http-api/`（已刪除，見 `http-api-documentation-flow.md`） |
| [`local-proxy-hook-flow.md`](local-proxy-hook-flow.md) | Local Proxy Hook 操作流程 — 6 步驟（確認證據 → 識別 Handler → Hook → Cast Netty → 去敏 → 歸因） | `skills/apk-analysis/techniques/local-proxy/`（已刪除，見 `local-proxy-hook-flow.md`） |
| [`media-hls-analysis-flow.md`](media-hls-analysis-flow.md) | Media/HLS 分析操作流程 — 7 步驟（分離控制面/資料面 → Playlist → Key → Segments → 合併 → 容器驗證） | `skills/apk-analysis/techniques/media-hls/`（已刪除，見 `media-hls-analysis-flow.md`） |

## Generated Surface

本目錄的 4 個操作流程已由 [`runtime/compiler/compiler-engine.rb`](../../runtime/compiler/compiler-engine.rb) 編譯為 [`generated-phases.yaml`](generated-phases.yaml)，包含：

- 各 workflow 的步驟定義（step title + source line）
- 前置準備項目（prerequisites）
- 成功產出格式（output format）

### 位置說明

`generated-phases.yaml` 放在本目錄而非 `runtime/generated/`，因為 APK workflow 屬於
領域操作技巧（analysis 層），不是系統 runtime 的一部分。遵循三層分離原則：
領域層的 generated surface 跟隨 source 目錄。

相關連結：
- [`workflow/apk-analysis/execution-flow.md`](../../workflow/apk-analysis/execution-flow.md) — 端到端 agent 執行流程
- [`runtime/runtime.db`](../../runtime/runtime.db) — 系統層 generated surface（SQLite，取代 `runtime/generated/`）
- [`runtime/routing/`](../../runtime/routing/) — Runtime routing 流程
