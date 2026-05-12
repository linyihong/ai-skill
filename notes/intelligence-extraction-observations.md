# Intelligence Extraction Observations

記錄 Technique → Intelligence 分解過程中的觀察，作為未來建立正式 extraction pipeline 的參考。

## Pilot: flutter-dart-aot

### Extraction 過程

1. **分析原始 technique**：讀取 `skills/apk-analysis/techniques/flutter-dart-aot/README.md` 和 `analysis/apk/techniques/flutter-dart-aot.md`
2. **拆解內容**：
   - HOW TO DO（workflow）：Common Flow 的 7 個步驟、工具命令、操作順序
   - HOW TO THINK（intelligence）：When To Use 的判斷信號、Core Guidance 的策略建議、Pitfalls 的錯誤模式
3. **建立 workflow**：`analysis/apk/workflows/frida-hook-flow.md`
4. **提煉 intelligence atoms**：
   - `heuristics/hook-selection.md` — 從 Core Guidance + When To Use 提煉
   - `anti-patterns/early-hook-instability.md` — 從 Pitfalls 提煉
   - `failure/frida-spawn-race.md` — 從 feedback_history 提煉
   - `signals/flutter-dart-aot-detection.md` — 從 When To Use 提煉

### 哪些 decision 可以 atomize

- **Hook 策略選擇**：可以 atomize（有明確的決策表）
- **過早 hook 不穩定**：可以 atomize（有明確的症狀與預防方式）
- **Spawn race condition**：可以 atomize（有明確的診斷與緩解方式）
- **Flutter/Dart AOT 辨識**：可以 atomize（有明確的信號與可信度）

### 哪些 decision 不容易 atomize

- **Core Guidance 的通用建議**（如「從高語意邊界開始」）：太抽象，不容易轉成決策表
- **工具選擇**（blutter vs unflutter）：已經在 `analysis/apk/tools-and-failures.md` 中

### 哪些 intelligence 最 reusable

1. **Hook selection heuristic** — 最 reusable，因為每個 Flutter app 分析都會用到
2. **Flutter/Dart AOT detection signals** — 第二 reusable，因為決定分析路線
3. **Early hook instability** — 中等 reusable，只在 hook 設定階段需要
4. **Frida spawn race** — 較低 reusable，只在 spawn 失敗時需要

### 哪些 extraction 太細沒價值

- 沒有發現太細的 extraction。4 個 atoms 都有明確的使用情境。

### 格式觀察

- **決策表格式**（heuristics）最有效：情境 → 建議做法 → 判斷信號
- **症狀表格式**（anti-patterns/failure）也有效：症狀 → 可能原因 → 診斷方式
- **信號表格式**（signals）適合 detection：信號 → 檢查方式 → 可信度
- Token impact 標註有助於 runtime lazy-load 決策

### 與既有 intelligence atoms 的關係

```
evidence-first-routing.md  ──→  決定分析路線
       │
       ▼
signals/flutter-dart-aot-detection.md  ──→  辨識技術特徵
       │
       ▼
heuristics/hook-selection.md  ──→  選擇 hook 策略
       │
       ▼
anti-patterns/early-hook-instability.md  ──→  避免錯誤做法
       │
       ▼
failure/frida-spawn-race.md  ──→  診斷與修復失敗
```

這形成一個完整的分析流程：**路線選擇 → 技術辨識 → 策略選擇 → 避免錯誤 → 失敗診斷**。

### 下一步建議

1. 在實際 APK 分析 session 中測試這些 intelligence atoms 是否能提升 AI decision quality
2. 觀察 AI 是否開始做 decision routing（根據 signal 改變策略）
3. 如果成功，再處理下一個 technique（http-api 或 local-proxy）
4. 如果發現跨 domain 的 reusable pattern，再考慮 promotion 到更高層
