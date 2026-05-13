# 框架依賴偏誤：ide/ 錯誤放在 engineering/ 下

## 發生時間

2026-05-13 14:40 (JST)

## 情境

將 `vscode-extension-global-state.md` 從 `ai-tools/ide/` 昇華到 `intelligence/` 時，我直接把 `ide/` 放在了 `intelligence/engineering/ide/` 下，沒有考慮到 `ide/`（IDE 生態系統知識）與 `engineering/`（軟體工程經驗法則）的本質差異。

## 當下使用的工具

Roo Code（Code mode），deepseek-chat model

## 分析

### 為什麼會犯這個錯

這是**框架依賴偏誤（Framework Dependency Bias）**：

1. **類比捷徑**：因為 `engineering/` 已經有 `apk-analysis/`、`app-development-guidance/` 等領域，我直覺地把 `ide/` 當作另一個領域，沒有考慮本質差異
2. **沒檢查父目錄定義**：`intelligence/engineering/README.md` 的標題是「Engineering Intelligence」，描述是「軟體工程經驗法則」、「架構決策智慧」—— `ide/`（IDE 生態系統知識）明顯不屬於
3. **沒有退一步問**：「`engineering/` 的邊界是什麼？`ide/` 真的符合嗎？」

### 為什麼沒有加入測試案例

這是**Failure-to-validator Closure**：

1. **修復心態**：專注於「把錯的改對」（搬移目錄），沒有切換到「防止再發」的預防模式
2. **沒有泛化錯誤模式**：沒有把這個具體錯誤抽象化為「新目錄必須在 intelligence/README.md 結構圖中註冊」的通用檢測規則
3. **validator 存在但沒想起來**：`scripts/validate-knowledge-runtime.rb` 已經存在，但修復時不會自動想到去擴充它

## 教訓

1. **新增目錄到 `intelligence/` 時，必須先確認分類邊界**：檢查 `intelligence/README.md` 的結構圖，確認新目錄應該放在哪一層
2. **修復錯誤後，必須加入對應的 validator 測試**：否則同樣的錯誤模式可能再次發生
3. **建立「錯誤修復檢查清單」**：□ 修復完成 □ 已泛化錯誤模式 □ 已加入 validator 測試 □ 已驗證測試有效性

## 相關文件

- Failure pattern: [`shared-rules/failure-patterns/failure-to-validator-closure.md`](../../shared-rules/failure-patterns/failure-to-validator-closure.md)
- Validator test: [`scripts/validate-knowledge-runtime.rb`](../../scripts/validate-knowledge-runtime.rb) 中的 `validate_intelligence_classification_boundary`
