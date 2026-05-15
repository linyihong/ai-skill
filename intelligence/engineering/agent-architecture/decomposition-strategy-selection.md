# Decomposition Strategy Selection（拆解策略選擇）

**Status**: `candidate-intelligence`
**Source**: 本系統 Phase 26-33 實際運作經驗（[`plans/archived/2026-05-12-1506-skill-specific-extraction.md`](../../../plans/archived/2026-05-12-1506-skill-specific-extraction.md)）

## 原則

**Not all content benefits from the same extraction strategy. Choose the strategy based on content structure, not convention.**

不是所有內容都適合同一種提取策略。根據內容結構選擇策略，而不是根據慣例。

## 為什麼

1. **不同 skill 的內容結構差異極大** — 本系統的 3 個 skill 各自有不同的組織方式：`apk-analysis` 用 `techniques/` 混合 workflow + intelligence，`app-development-guidance` 用 `controls/`、`platforms/`、`languages/` 各自獨立，`travel-planning` 用 `WORKFLOW.md` + `DOCUMENTATION.md` + `TOOLS.md` 各自獨立。
2. **單一策略無法處理所有情況** — Decomposition 適合混合型內容，Catalog 適合列舉型內容，Direct Promotion 適合已獨立成檔的內容。
3. **錯誤的策略選擇會增加工作量** — 對列舉型內容使用 decomposition 會產生大量無意義的 atom，對混合型內容使用 catalog 會遺失 intelligence。

## 三種策略

| 策略 | 適用場景 | 範例 | 工作量 |
|------|---------|------|--------|
| **Decomposition**（拆解） | 單一檔案混合了 workflow + intelligence + tools | `techniques/flutter-dart-aot/` → `workflows/frida-hook-flow.md` + 4 intelligence atoms | 高 |
| **Catalog**（目錄化） | 多個獨立檔案，每個檔案是列舉或參考資料 | `controls/`、`platforms/`、`languages/` → `controls-catalog.md` | 低 |
| **Direct Promotion**（直接提升） | 內容已獨立成檔，只需搬到新位置 | `WORKFLOW.md` → `workflow/<skill>/execution-flow.md` | 低 |

## 判斷流程

```
內容需要提取
    ↓
分析內容結構
    ↓
├── 單一檔案混合多種型態 → Decomposition
│   ├── 包含 workflow 部分 → 提取到 workflow/
│   ├── 包含 intelligence 部分 → 提取到 intelligence/
│   └── 包含 tools/failure 部分 → 提取到 analysis/
│
├── 多個獨立列舉檔案 → Catalog
│   ├── 內容是參考資料 → 合併為 catalog
│   └── 內容是檢查清單 → 合併為 review checklist
│
└── 內容已獨立成檔 → Direct Promotion
    ├── 操作流程 → workflow/
    ├── 分析方法 → analysis/
    └── 決策智慧 → intelligence/
```

## 症狀

| 症狀 | 說明 | 可信度 |
|------|------|--------|
| **提取後 atom 太細沒價值** | 一個 atom 只有 3-5 行內容，不值得獨立成檔 | 高 |
| **提取後遺失上下文** | Atom 單獨看有意義，但脫離原始檔案後無法理解 | 高 |
| **提取後 workflow 不完整** | 只提取了部分步驟，遺漏了關鍵的前置條件 | 中 |
| **Catalog 內容重複** | 多個 catalog 檔案包含相似的內容 | 中 |

## 不建議的做法

| 不建議 | 原因 |
|--------|------|
| 對所有內容都用同一種策略 | 不同內容結構需要不同策略 |
| 對列舉型內容使用 decomposition | 產生大量無意義的 atom |
| 對混合型內容使用 catalog | 遺失 intelligence，workflow 不完整 |
| 提取後不檢查完整性 | 可能遺漏關鍵的上下文或前置條件 |

## 相關 atoms

- [`pilot-first-validation.md`](pilot-first-validation.md) — 先驗證再抽象化（選擇策略前先用 pilot 驗證）
- [`linked-updates-completeness.md`](linked-updates-completeness.md) — 連動更新完整性（提取後需要更新所有引用）
- [`premature-optimization.md`](../heuristics/premature-optimization.md) — 過早最佳化經驗法則

## Token Impact

錯誤的策略選擇可能導致：
- Decomposition 策略用錯：浪費 3000-8000 token 建立無意義的 atom
- Catalog 策略用錯：遺失 intelligence，後續需要重新提取（5000-10000 token）
- Direct Promotion 策略用錯：內容結構不完整，需要後續調整（2000-5000 token）

正確的策略選擇能節省 50-70% 的 extraction 成本。

---

← [回到 agent-architecture/](README.md)
