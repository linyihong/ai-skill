# Reference Decomposition

**Status**: `candidate-intelligence`

## 定義

Reference decomposition 是一種 reasoning pattern：當任務依據 reference material（截圖、設計稿、既有系統、API 範例、競品等）時，必須先把隱含資訊拆成 **explicit specification**，再進入實作。

它回答的問題不是「能不能仿照這張圖寫 CSS」，而是「實作者與驗收者是否共用同一份可檢查的規格」。

## Anti-pattern

```text
Reference → Implement
```

實作者補腦、驗收者對原圖 → 反覆修改。

## Correct pattern

```text
Reference
→ Decompose (inventory + attributes + constraints)
→ Explicit Spec (incl. product mapping + acceptance)
→ Implement
→ Verify against spec
```

## 六步 workflow

1. **Reference inventory** — 列出 reference 中有哪些區塊/元素；標註資產路徑與觀察條件（viewport、主題、語系）。
2. **Attribute extraction** — 每個元素抽可觀察屬性（尺寸、色碼、字級、間距、資產來源等）。UI 維度在 **專案 visual spec 模板** 填寫，不在此檔寫死。
3. **Constraint extraction** — 硬性規則（全寬、無縫隙、分隔線行為等），寫成可勾選語句。
4. **Product mapping** — Reference 項目對應產品功能；標註 **Deferred / out of scope**。
5. **Acceptance checklist** — 驗收勾選項；含部署/靜態資源等環境檢查（若適用）。
6. **Implementation** — 僅在 spec 可被 review 後動手；實作後把 code path 記回 spec（optional Implementation Notes）。

## 判斷規則

**Mandatory**（應先 decomposition）：

- 新 screen / 新 surface 對齊 reference 圖或設計稿
- 使用者或需求方要求「照圖 / 照這個風格」
- Reference 變更導致視覺或契約對齊任務

**可跳過**（不必完整 spec）：

- 單行 typo、明確 bugfix 且範圍已用文字說死
- 純邏輯/後端變更，無 reference 輸入

## Non-goals

- 不建立 Reference taxonomy / ontology（見 [`reference-types.md`](reference-types.md)）
- 不把專案 UI 色碼/token 寫進 canonical intelligence
- 不在知識落地當天同步 validation scenario、stop hook 或 enforcement governance（先證明 spec 降 rework）

## Related Reasoning Families

This reasoning family may belong to a broader **Model Before Action** family.

```text
Observation → Explicit Model → Execution
```

| Family | Shape |
|--------|-------|
| Validation Reasoning | Claim → Validation model → Verification |
| Failure Authority | Failure → Authority model → Action |
| Evidence chain validation | State claim → Evidence model → Depth check |
| **Reference Decomposition** | Reference → Specification model → Implementation |

Current evidence:

- Validation Reasoning
- Failure Authority
- Reference Decomposition

Further observation required. Do **not** create `execution/model-before-action/` directory or move files until multiple pilots show the umbrella is stable.

## 相關知識

- [`reference-types.md`](reference-types.md) — illustrative hints only
- [`../validation-reasoning/README.md`](../validation-reasoning/README.md)
- Project Visual Reference Workflow — `visual-reference-spec.template.md` + per-screen `*-spec.md`
