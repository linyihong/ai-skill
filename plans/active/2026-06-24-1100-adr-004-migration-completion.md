---
id: 2026-06-24-1100-adr-004-migration-completion
plan_kind: main
status: draft
execution_status: active
owner: linyihong
created: 2026-06-24
parent: null
baseline_ref: 2026-06-23-1500-adr-004-migration-drift-diagnosis
---

# ADR-004 Migration Completion Plan

> **這是 implementation plan，不是 diagnosis。** 診斷已於
> [`2026-06-23-1500-adr-004-migration-drift-diagnosis.md`](2026-06-23-1500-adr-004-migration-drift-diagnosis.md)
> 凍結（`status: frozen`）。本文件把那份 frozen diagnosis 當 **immutable baseline**，
> **不重述**證據、loop 圖或失效模式命名；只引用其錨點，並把 P0-A/B/C 與 Contract Layer
> 展開為可執行 phase。命名刻意是 *Completion*（讓 migration 真正完成）而非 *Repair*（修 bug）。

## 0. Context

### Reference Diagnosis (immutable baseline)
- Baseline：`2026-06-23-1500-adr-004-migration-drift-diagnosis`（frozen）。
- 引用而非複製：Layer 狀態表、Drift Propagation Loop、Semantic Validation Drift /
  Canonical Knowledge Propagation Failure、Canonical-Path Derivation Invariant、Contract
  Ownership 鏈，全部以 baseline 為準；本文件若與 baseline 衝突，以 baseline 為真。
- 關鍵錨點（取自 baseline「證據錨點」段，勿在此重抄行號內容）：
  `runtime.go:1213` 死 glob、`runtime.go:1876-1882` masking gate、`runtime_test.go:244,625,640`
  old-world fixtures、`routing-registry.yaml` `route.feedback.history`（Path 1 正確）。

### Out of Scope
- [ ] 重新論證 diagnosis（已 frozen）。
- [ ] 修改 `constitution/ADR-004`（憲法層正確，禁改——見 Rollback）。
- [ ] feedback lesson 內容/品質本身。
- [ ] Path 1（reference-first routing）行為變更——它正確且 load-bearing。
- [ ] <TODO: 其餘 out-of-scope 待 Phase 0 inventory 後補足>

### Success Definition
- 兩條 discovery path 一致：`Path 1 (direct read) == Path 2 (runtime query)` 對 lesson 同集合。
- canonical path 只有單一 owner（registry），無任何 consumer hardcode 路徑。
- Exit Criteria 的 derive 鏈五層全綠（見文末 Exit Criteria）。

---

## Phase 0 — Migration Inventory (Step 0)

**目標**：完成 99-file classification + intent map（baseline 已點名但未展開）。

**產出**：
- [ ] Consumer Matrix（每個 consumer：path、讀法、Path 1/2、是否 load-bearing）
- [ ] Canonical Path Usage Table（誰宣告 path、誰消費 path）
- [ ] Class A/B/C/D 分類（A=doc 文字、D=code/fixture… <TODO: 鎖定分類定義>）

**Gate（未過不得進 Phase 1）**：
- [ ] 無 Unknown Consumer（每個觸及 feedback path 的點都已分類）
- [ ] 每個 consumer 有 owner

---

## Phase 1 — Observability Restoration (P0-A)

> **原則：Truth before repair.**
> observability 不是修復的前置步驟，而是**獨立可交付能力**——故獨立成 phase。
> 對應 baseline 的「P0-A 不保證 repair」：本 phase 成功 = 量測誠實，**不**等於 lesson 已索引。

**交付**：
- [ ] provenance health-check（取代 `runtime.go:1879` 的 `MATCH 'feedback'` token 檢查）
- [ ] `source_path` assertions（≥1 atom 且其 source_path 在 canonical lesson sink 之下）
- [ ] index coverage metrics（可量測 Path 2 的真實 lesson 覆蓋率）

**成功條件**：
- [ ] Path 2 壞 → health-check 必須**紅**（先證明它能偵測現況的空集合）
- [ ] Path 2 修 → health-check 必須**綠**

**禁止**：
- 不允許在本 phase 修改 indexer（observability 與 repair 必須可分離驗證）。

---

## Phase 2 — Runtime Repair (P0-B)

> **原則：Registry-derived indexing only.**

**交付**：
- [ ] index repoint：`skills/*/feedback_history` → canonical sink（路徑由 registry 提供）
- [ ] registry consumption：indexer 從 routing-registry 取 path，不 hardcode
- [ ] compatibility notes（舊路徑消費者/快取的過渡說明）

**成功條件**：
- [ ] `Path 1 == Path 2`（同一 lesson 集合）
- [ ] Phase 1 的 health-check 由紅轉綠

**禁止**：
- 禁止任何 hardcoded canonical path（違反 Canonical-Path Derivation Invariant）。

---

## Phase 3 — Reality Alignment (P0-C)

> **原則：Fixture must model reality.**

**交付**：
- [ ] fixture migration：`runtime_test.go` seed 改為 canonical sink
- [ ] negative tests：舊路徑**不**被索引
- [ ] deleted-world lock：機械化阻止 fixture 再教 deleted world

**成功條件**：
- [ ] 舊路徑（`skills/*/feedback_history`）若存在於 fixture/索引 → 測試 **fail**

---

## Phase 4 — Contract Extraction

把 baseline 的 Contract Ownership 鏈從文字落為可執行 artifact。

**交付**：
- [ ] `knowledge/runtime/contracts/<name>.yaml`（machine-readable contract）
- [ ] contract → routing-registry 的 single-direction materialize 連結（registry = 唯一 path owner）

**成功條件**：以下全部 **derive**（不得平行宣告 path）：
- [ ] runtime
- [ ] validator
- [ ] docs
- [ ] tests

---

## Rollback

只允許：
- [ ] disable enforcement（停用新 validator / health-check，回到觀測前狀態）

禁止：
- rollback constitution（`ADR-004` 是正確來源，回退它＝回到 drift 世界觀）

---

## Exit Criteria

**Invariant（任一條斷裂 → 不可 close）**：

```
Constitution
  ⊨ Contract
  ⊨ Registry
  ⊨ Runtime
  ⊨ Validation
```

- [ ] 五層 derive 鏈全綠
- [ ] Canonical-Path Derivation Invariant 無違反（無繞過 registry 的 path 宣告）
- [ ] `Path 1 == Path 2`
- [ ] baseline diagnosis 仍 frozen、未被本計畫修改

---

## Document TODO

| 項目 | 狀態 |
| --- | --- |
| Phase 0 分類定義（Class A/B/C/D 邊界） | pending |
| Out of Scope 補足 | pending |
| contract 命名（`<name>.yaml`） | pending |
| 各 phase 的 validator 實作細節 | pending（進 phase 時展開） |
