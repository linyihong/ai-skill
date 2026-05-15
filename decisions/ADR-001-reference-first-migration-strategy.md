# ADR-001: Reference-First Migration Strategy

## Status

**Accepted**

## Context

本 repository 原本是單一 `skills/` 目錄結構，所有內容（workflow、analysis methods、engineering intelligence、templates、feedback lessons）都放在 skill 目錄下。隨著內容增長，出現以下問題：

1. **跨 skill 重複**：多個 skill 需要相同的分析流程或 intelligence，但只能複製貼上。
2. **Context 浪費**：載入一個 skill 時，連帶載入不相關的 workflow、tools、feedback lessons。
3. **難以演化**：feedback lesson 只能回到原 skill，無法跨 skill 推廣。
4. **無 navigation**：沒有機制讓 agent 快速找到「哪個文件處理什麼問題」。

解決方案有兩種：

- **Migration-first**：直接將 `skills/` 內容搬到新分層，舊路徑設 redirect。
- **Reference-first**：保留 `skills/` 為 source of truth，在新分層建立 reference / summary / index，逐步提取內容。

## Decision

採用 **Reference-first** 策略。

核心原則：

1. **`skills/` 維持 source of truth**：所有現有 skill 文件保持不變，繼續可被 agent 直接讀取。
2. **新分層只建立 reference / summary / index**：`analysis/`、`workflow/`、`intelligence/` 等新目錄先建立 README 定義 scope，再從 `skills/` 提取內容到對應位置。
3. **提取後的內容標記為 `candidate-map` 或 `candidate-atom`**：不取代舊路徑，只新增 navigation entry。
4. **舊入口維持 active**：直到新分層內容經過 validation 並 promotion 為 `validated-atom` 後，才考慮 deprecate 舊路徑。
5. **相容性規則**：新分層內容不得破壞舊 `skills/` 的可讀性；graph edges 使用 `preserves_entrypoint` 而非 `replaces`。

## Consequences

### 正面

- **零風險**：舊結構完全不受影響，agent 行為不變。
- **可逐步驗證**：每個提取步驟都可獨立驗證，不影響既有功能。
- **學習曲線低**：開發者仍可從熟悉的 `skills/` 路徑讀取內容。
- **相容性明確**：graph edges 的 `preserves_entrypoint` 語意讓 runtime 知道「新路徑是補充，不是取代」。

### 負面

- **過渡期較長**：完整遷移需要多個 phase 才能完成。
- **雙重維護**：提取期間，內容同時存在於舊路徑與新路徑，需要 governance 確保同步。
- **Navigation 複雜度**：agent 需要知道「什麼時候讀舊路徑，什麼時候讀新路徑」。

## Alternatives Considered

- **Migration-first**：直接搬移內容並設 redirect。風險是搬移過程中可能破壞現有 agent 行為，且無法 rollback。不採用。
- **Big bang rewrite**：一次全部重寫新結構。風險極高，無法驗證中間狀態。不採用。

## Related

- [`plans/archived/2026-05-11-apk-analysis-pilot-migration.md`](../plans/archived/2026-05-11-apk-analysis-pilot-migration.md) — 第一個 pilot migration map
- [`governance/lifecycle/README.md`](../governance/lifecycle/README.md) — lifecycle states（candidate-map → candidate-atom → validated-atom → promoted）
- [`knowledge/graphs/README.md`](../knowledge/graphs/README.md) — graph edge types（`preserves_entrypoint`）
- [`plans/archived/2026-05-11-next-stage-upgrade-plan.md`](../plans/archived/2026-05-11-next-stage-upgrade-plan.md) — 整體升級規劃
