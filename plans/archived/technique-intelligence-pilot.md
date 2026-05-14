# Technique → Intelligence Pilot：flutter-dart-aot

## 策略摘要

**核心目標**：提升 AI decision quality，不是把分類變漂亮。

**兩個層的明確分工**：

| 層 | 內容 | 範例 |
|---|---|---|
| `analysis/` | **HOW TO DO** — execution knowledge | workflow、command、setup、tracing、hook 步驟、dump 方法、case study |
| `intelligence/` | **HOW TO THINK** — decision intelligence | heuristics、anti-patterns、failure learning、routing、tradeoffs、escalation、validation、signal detection |

**Technique Decomposition**（拆解，不是搬遷）：

```
舊 technique（flutter-dart-aot）
    ├── workflow/execution 部分 → analysis/apk/workflows/
    └── decision/intelligence 部分 → intelligence/{heuristics,anti-patterns,failure,signals}/
```

**舊 techniques 保留**，不刪除，但標註已提取 intelligence。

## Pilot 選擇：flutter-dart-aot

選擇原因：
- 同時包含 workflow、reverse、runtime patch、anti-debug、failure handling、signal detection
- 最適合驗證 Technique → Intelligence 是否能提升 AI decision quality

## 執行步驟

### Phase 28a：建立 intelligence 最小子目錄結構

只建立 4 個子目錄（最小必要結構）：

```
intelligence/
├── heuristics/       # 什麼 signal 要優先 hook
├── anti-patterns/    # 哪些 patch timing 很容易 crash
├── failure/          # spawn race / relocation timing / jit mismatch
└── signals/          # 如何辨識 flutter dart aot
```

各目錄建立 `README.md` 說明 scope 與 routing rules。

### Phase 28b：建立 analysis/apk/workflows/ 目錄

```
analysis/apk/
├── workflows/        # 新目錄：操作流程、command、setup
├── case-studies/     # （未來）
├── traces/           # （未來）
├── techniques-archive/  # （未來，舊 techniques 移入）
└── techniques/       # 保留，逐步拆分
```

### Phase 28c：拆解 flutter-dart-aot → workflow 部分

從 `skills/apk-analysis/techniques/flutter-dart-aot/` 和 `analysis/apk/techniques/flutter-dart-aot.md` 中提取：

- Frida hook 操作流程
- command 與 setup
- adb 與 proxy 設定
- dump 方法
- 常見操作步驟

寫入 `analysis/apk/workflows/frida-hook-flow.md`

### Phase 28d：提煉 intelligence atoms

從 flutter-dart-aot 中提煉 4 類 intelligence：

1. **`intelligence/heuristics/hook-selection.md`** — 何時該用哪種 hook 策略
2. **`intelligence/anti-patterns/early-hook-instability.md`** — 哪些 patch timing 容易 crash
3. **`intelligence/failure/frida-spawn-race.md`** — spawn race / relocation timing / jit mismatch
4. **`intelligence/signals/flutter-dart-aot-detection.md`** — 如何辨識 flutter dart aot

### Phase 28e：標註舊 technique 檔案

在 `analysis/apk/techniques/flutter-dart-aot.md` 和 `skills/apk-analysis/techniques/flutter-dart-aot/README.md` 加入：

```markdown
> **Intelligence Extracted**
> See:
> - `intelligence/heuristics/hook-selection.md`
> - `intelligence/anti-patterns/early-hook-instability.md`
> - `intelligence/failure/frida-spawn-race.md`
> - `intelligence/signals/flutter-dart-aot-detection.md`
```

### Phase 28f：建立 extraction observations

建立 `notes/intelligence-extraction-observations.md`，記錄：

- extraction 過程
- 哪些 decision 可以 atomize
- 哪些不能
- 哪些 intelligence 最 reusable
- 哪些 extraction 太細沒價值

### Phase 28g：更新架構文件

更新 `plans/active/next-stage-upgrade-plan.md`，記錄 Phase 28。

### Phase 28h：提交 + push

## 成功驗證標準

Pilot 成功 = AI 開始能做 decision routing：

- 以前：只會照流程 dump
- 現在：能根據 signal 改變策略

## 不做的範圍

- ❌ 不建立完整的 intelligence extraction pipeline / governance
- ❌ 不一次處理所有 4 個 techniques
- ❌ 不刪除舊 techniques
- ❌ 不建立完整的 intelligence 子目錄結構（只開 4 個）
