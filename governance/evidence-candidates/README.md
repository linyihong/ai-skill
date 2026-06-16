# Evidence Candidates（observation infrastructure）

跨 plan 的**候選證據索引層**。解「人工記不住哪個案例該回流到哪個 plan」的**索引問題**。

設計來源（canonical）：[`plans/active/2026-06-16-1131-evidence-candidate-system.md`](../../plans/active/2026-06-16-1131-evidence-candidate-system.md)。

## 這是什麼、不是什麼

| 是 | 不是 |
|---|---|
| observation infrastructure（owned by governance） | runtime（不接入 runtime） |
| 被動 Detection + 人工 Capture | Enforcement / Promotion（人工 + plan gate + maintainer 才有判定權） |
| registry / index | plan authority（rule 定義不在這裡，見下方 invariant） |
| 跨 session 可回收觀察層 | memory（memory 只存 work-continuation pointer） |

## Ownership 鏈

```text
plan        → defines rule    （rule 定義 owner = 各 plan 的「Evidence Rule」section）
governance  → indexes rule    （本目錄 evidence-rules/，pointer only）
candidate   → references rule  （inbox/ 的 candidate）
evidence    → enters plan      （accept 後寫回 plan）
```

## 目錄

```text
governance/evidence-candidates/
  schema/                          # 只定 shape，不寫 criteria
    candidate-schema.md
    evidence-rule-schema.md
  evidence-rules/                  # POINTER ONLY — 不承載規則內容
    economics.pointer.yaml
    governance-pattern.pointer.yaml
    interaction-hazard.pointer.yaml
  inbox/                           # gitignored — candidate 本機暫存，永不 commit
  README.md
```

## Invariants

1. **`evidence-candidates/` 不可成為 plan authority**：`evidence-rules/*.pointer.yaml` 只放 registry
   pointer（`plan_ref` + `schema` + `source.embedded_section`），rule 定義留在 plan。否則會長出
   `plan.md → rule copy → registry` 雙 source。
2. **candidate 不可指向 candidate**：`candidate.source` 只能 reference 原始 artifact，鏈為
   `artifact → candidate → plan evidence`。
3. **scanner MUST be stateless**：掃當前 artifact → 產 candidate → 結束；不記狀態、不增量 cache。
4. **inbox 永不 commit**：candidate 是本機 observation state（gitignored），accept 才寫回 plan。

## 狀態

Phase 1A — contract establishment（skeleton）。scanner（Phase 1C）尚未實作；目前無自動掃描，
candidate 由人工放入 `inbox/`。詳見設計 plan。
