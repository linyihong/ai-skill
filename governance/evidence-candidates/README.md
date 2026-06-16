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
  scanner-v0.py                    # ASSEMBLER（stateless；非 ai-skill CLI、非 hook）
  README.md
```

## scanner-v0（Assembler）

`scanner-v0.py` 是 standalone stateless **assembler**，不接 ai-skill CLI / hook / runtime。
input：artifact + **明確 criteria_hits**（`criteria_source.actor` 標於 scanner 外）+ rule registry；
output：`candidate{status: create}` 寫入 `inbox/`。
做 schema validate / pointer resolve / dedupe（content-hash id，idempotent）/ invariant check / persist；
**不做** infer / match / classify / score / rank / accept / expire。
真正的 match-by-pattern 是 Phase 2（`criteria_source.actor: matcher-v2`），scanner 永不自產 criteria。
用法：`python3 scanner-v0.py < artifact.json`。

## Invariants

1. **`evidence-candidates/` 不可成為 plan authority**：`evidence-rules/*.pointer.yaml` 只放 registry
   pointer（`plan_ref` + `schema` + `source.embedded_section`），rule 定義留在 plan。否則會長出
   `plan.md → rule copy → registry` 雙 source。
2. **candidate 不可指向 candidate**：`candidate.source` 只能 reference 原始 artifact，鏈為
   `artifact → candidate → plan evidence`。
3. **scanner MUST be stateless**：掃當前 artifact → 產 candidate → 結束；不記狀態、不增量 cache、不長 ontology。
4. **inbox 永不 commit**：candidate 是本機 observation state（gitignored），accept 才寫回 plan。
5. **`criteria_hits` MUST originate outside scanner**：candidate 必帶 `criteria_source.actor`（`human` 等）；
   scanner 自身被拒。scanner v0 是 assembler，不自產 criteria（matcher 是 Phase 2）。

## 狀態

Phase 1 完成：1A（schema + attach + criteria）、1B（round-trip + lifecycle proof）、1C（assembler
scanner-v0）。scanner 為 assembler：candidate 由人工標 `criteria_hits` 後經 scanner 組裝進 `inbox/`；
真正的 auto-detection（matcher）與 Phase 2 accumulation runtime 皆 deferred。詳見設計 plan。
