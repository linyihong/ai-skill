# APK Analysis Artifact Gates（Thin Index）

本文件原為 575 行 monolith；自 2026-05-31 起改為 **thin index**——artifact gates 的 canonical prose 已切成 8 個 focused slices（見 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7.5）。

分析方法見 [`analysis/apk/`](../../analysis/apk/)；模板見本目錄的 templates。

---

## Cognitive Slice 導航

按任務意圖選擇要載入的 slice，避免拉整份 575 行 monolith：

| 任務意圖 | 載入 slice | load_when |
|---|---|---|
| 建立 / 更新 UI map | [`artifact-gates/ui-architecture-map.md`](artifact-gates/ui-architecture-map.md) | UI 觀察、navigation segment library、operation-to-API 對照 |
| 整理 API 文件 / SDK 對照 | [`artifact-gates/api-catalog.md`](artifact-gates/api-catalog.md) | API endpoint 文件、catalog completion gate、單支 API 詳細需求 |
| 建立 runtime baseline | [`artifact-gates/domain-runtime-baseline.md`](artifact-gates/domain-runtime-baseline.md) | development readiness gate（SDK/client/live integration 開發前） |
| 產出 feature handoff | [`artifact-gates/feature-handoff.md`](artifact-gates/feature-handoff.md) | feature 重建 / 「能不能重建」/「架構是什麼」問題 |
| 記錄分析證據 | [`artifact-gates/evidence-chain.md`](artifact-gates/evidence-chain.md) | 單次分析筆記 / 證據鏈 / 失敗 capture |
| evidence / sample 去敏 | [`artifact-gates/sanitization.md`](artifact-gates/sanitization.md) | commit / publish 前 |
| SDK live / identity self-gen audit | [`artifact-gates/self-generation-audits.md`](artifact-gates/self-generation-audits.md) | self-generation 宣稱 / device id / install id / account / session seed / attestation |
| 撰寫 dev notes / feedback / 回填 | [`artifact-gates/documentation-discipline.md`](artifact-gates/documentation-discipline.md) | post-analysis writeup / retrospective |

預設 suppress：純 reference 查閱請只載入 README，不必載入任何 slice。

---

## 舊章節 redirect（兼容外部連結）

下表保留舊 `#N-<title>` anchor 風格，每個指向新 slice 檔。若有 inbound link 用 `artifact-gates.md#某節`，可依此表更新。

| 舊節 | 新 slice |
|---|---|
| §1 UI Architecture Map | → [`artifact-gates/ui-architecture-map.md`](artifact-gates/ui-architecture-map.md) §1 |
| §2 API Catalog | → [`artifact-gates/api-catalog.md`](artifact-gates/api-catalog.md) §2 |
| §3 Domain/Runtime Baseline | → [`artifact-gates/domain-runtime-baseline.md`](artifact-gates/domain-runtime-baseline.md) §3 |
| §4 Feature Reconstruction Handoff | → [`artifact-gates/feature-handoff.md`](artifact-gates/feature-handoff.md) §4 |
| §5 單次分析筆記模板 | → [`artifact-gates/evidence-chain.md`](artifact-gates/evidence-chain.md) §5 |
| §6 證據鏈要求 | → [`artifact-gates/evidence-chain.md`](artifact-gates/evidence-chain.md) §6 |
| §7 失敗也要記錄 | → [`artifact-gates/evidence-chain.md`](artifact-gates/evidence-chain.md) §7 |
| §8 SDK Live Self-Generation Audit | → [`artifact-gates/self-generation-audits.md`](artifact-gates/self-generation-audits.md) §8 |
| §9 Authorized Identity Material Self-Generation Audit | → [`artifact-gates/self-generation-audits.md`](artifact-gates/self-generation-audits.md) §9 |
| §10 UI Architecture Map Template | → [`artifact-gates/ui-architecture-map.md`](artifact-gates/ui-architecture-map.md) §10 |
| §11 API Catalog Detail Requirements | → [`artifact-gates/api-catalog.md`](artifact-gates/api-catalog.md) §11 |
| §12 Sanitization Rules | → [`artifact-gates/sanitization.md`](artifact-gates/sanitization.md) §12 |
| §13 Developer Guidance Notes | → [`artifact-gates/documentation-discipline.md`](artifact-gates/documentation-discipline.md) §13 |
| §14 Feedback Lesson Writing Tips | → [`artifact-gates/documentation-discipline.md`](artifact-gates/documentation-discipline.md) §14 |
| §15 Backfill Rules | → [`artifact-gates/documentation-discipline.md`](artifact-gates/documentation-discipline.md) §15 |

---

## 切分依據

- Empirical trigger：[`slice-load-scenario-f-apk-analysis-probe.yaml`](../../validation/scenarios/software-delivery/slice-load-scenario-f-apk-analysis-probe.yaml) 實測真實 APK analysis 任務只用 6 個 gate / 12，inflation ratio ~1.57
- 8 slice 而非 7：[`slice-load-scenario-ag-schemes-a-vs-b.yaml`](../../validation/scenarios/software-delivery/slice-load-scenario-ag-schemes-a-vs-b.yaml) probe 證實 `documentation-discipline` 獨立載入經濟性更高
- Slice schema / 三層規則 / placement predicate / dependency budget：見 [`governance/cognitive-slice-taxonomy.md`](../../governance/cognitive-slice-taxonomy.md) §7.5

---

← [回到 workflow/apk-analysis/](README.md)
