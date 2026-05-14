# APK 分析方法

`analysis/apk/` 是可重用的 APK 觀測與拆解方法的候選目錄。在 pilot 階段，活躍 skill 仍為 `skills/apk-analysis/`；本目錄提供以參考優先的分析層，不搬移既有檔案。

## 目錄結構

```
analysis/apk/
├── README.md                       # 本文件
├── traffic-triage.md               # 流量分流與路線選擇
├── tools-and-failures.md           # 工具選擇、失敗判讀、命令模板
├── techniques/                     # 保留中：舊 technique 混合層（逐步拆分中）
└── workflows/                      # 操作流程（HOW TO DO）
```

## 範圍（Scope）

本層負責：

- 流量與執行路徑分流（traffic / runtime path triage）。
- 以證據優先選擇路線：pcap、MITM、Java hooks、native hooks、Flutter / Dart AOT、local proxy、media、offline decoding。
- 從動態捕獲中萃取模式，轉化為可重用的分析方法。
- 在撰寫 workflow 或工程結論之前，進行分析類別路由。
- **HOW TO DO** 執行知識：workflow、命令、設定、追蹤、hook 步驟、dump 方法（`workflows/`）。

本層不負責：

- 工具 skill 觸發文字；保留在 `skills/apk-analysis/SKILL.md`。
- 端到端的 agent 執行流程；使用 `workflow/apk-analysis/`。
- **HOW TO THINK** 決策智慧：heuristics、anti-patterns、failure learning、signal detection；使用 `intelligence/engineering/analytical-reasoning/`。
- 特定目標的 API host、endpoint、token、原始樣本或執行證據；保留在專案文件中。

## 目前來源參考（Current Source References）

| 主題 | 目前來源 | Pilot 目標狀態 |
| --- | --- | --- |
| 通用流量／執行路徑分流 | `../../skills/apk-analysis/WORKFLOW.md` | ✅ 已萃取至 `traffic-triage.md` |
| 工具選擇與失敗判讀 | `../../skills/apk-analysis/TOOLS.md` | ✅ 已萃取至 `tools-and-failures.md` |
| 媒體驗證工具 | `../../skills/apk-analysis/TOOLS.md` | ✅ 已萃取至 `tools-and-failures.md` |
| 自動化腳本安全邊界 | `../../skills/apk-analysis/TOOLS.md` | ✅ 已萃取至 `tools-and-failures.md` |
| Flutter / Dart AOT 方法 | `../../skills/apk-analysis/techniques/flutter-dart-aot/`（已刪除） | ✅ workflow → `workflows/frida-hook-flow.md`, intelligence → `intelligence/engineering/analytical-reasoning/` |
| HTTP API 文件方法 | `../../skills/apk-analysis/techniques/http-api/`（已刪除） | ✅ workflow → `workflows/http-api-documentation-flow.md`, intelligence → `intelligence/engineering/analytical-reasoning/heuristics/api-documentation-completeness.md` |
| Local proxy / loopback 方法 | `../../skills/apk-analysis/techniques/local-proxy/`（已刪除） | ✅ workflow → `workflows/local-proxy-hook-flow.md`, intelligence → `intelligence/engineering/analytical-reasoning/` |
| Media / HLS 方法 | `../../skills/apk-analysis/techniques/media-hls/`（已刪除） | ✅ workflow → `workflows/media-hls-analysis-flow.md`, intelligence → `intelligence/engineering/analytical-reasoning/signals/media-type-detection.md` |

## 閱讀順序（Read Order）

1. 當工具需要活躍 skill 觸發時，先讀 `../../skills/apk-analysis/SKILL.md`。
2. 用本文件了解分析層的邊界。
3. 讀 `traffic-triage.md` 了解流量／執行路徑分流。
4. 讀 `tools-and-failures.md` 了解工具選擇、失敗判讀與命令模板。
5. 證據確定路線後，讀 `workflows/` 了解 HOW TO DO 執行步驟。
6. 讀 `intelligence/engineering/analytical-reasoning/` 了解 HOW TO THINK 決策指引。

## 遷移備註（Migration Notes）

- 流量分流、工具選擇、失敗判讀、媒體驗證與安全邊界已從 `skills/apk-analysis/` 萃取至本目錄。
- 原始 skill 檔案仍保有權威內容；本目錄提供以參考優先的視角。
- 當 technique 被萃取時，使用 `../../metadata/schema.md` 建立 metadata，並更新 `../../knowledge/indexes/README.md`。
- 保留來自 `skills/apk-analysis/` 的連結，直到新路徑經過實際使用驗證。
- `workflows/` 是 HOW TO DO 執行知識的新目錄。舊的 `techniques/` 將逐步拆分為 `workflows/` + `intelligence/` + `techniques-archive/`。
