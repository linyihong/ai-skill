# App Development Guidance Feedback History Index

> ⚠️ **Lesson 已搬遷**：所有 lesson 實體已移至 [`feedback/history/development-guidance/`](../history/development-guidance/README.md)。本索引保留作為 extraction 追蹤用途，來源路徑 `skills/app-development-guidance/feedback_history/` 已刪除。

本索引列出原 `skills/app-development-guidance/feedback_history/`（已刪除）下所有 lessons，依其 `Promotion Target` 分類到對應的目標層。此索引讓 `feedback/` 層可發現 lessons，並追蹤哪些 lessons 已被提取到目標層。

## 索引說明

| 欄位 | 說明 |
| --- | --- |
| **目標層** | lesson 的 Promotion Target 對應到的新架構層 |
| **來源檔案** | `skills/app-development-guidance/feedback_history/` 下的原始 lesson 檔案 |
| **Status** | lesson 的成熟度狀態（promoted / validated / candidate） |
| **提取狀態** | 是否已提取到目標層（✅ 已提取 / ⬜ 未提取 / 🔄 部分提取） |
| **目標檔案** | 提取後的目標檔案路徑 |

---

## 1. 目標層：`workflow/software-delivery/`

這類 lessons 的 Promotion Target 包含 `WORKFLOW.md` 或 `process/README.md`，適合提取到 `workflow/software-delivery/`。

### common/

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-05-05_194400-contract-first-development-flow.md` | promoted | Contract-first development flow | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-05_200500-existing-project-doc-backfill-bdd-required.md` | promoted | Existing project doc backfill requires complete BDD | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-05_201000-missing-requirements-block-development.md` | promoted | Missing requirements block development | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-06_081600-change-intake-before-code.md` | promoted | Change intake before code | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-06_082000-separate-regression-from-new-code-validation.md` | promoted | Separate regression from new code validation | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-06_083000-embedded-hardware-product-flow.md` | promoted | Embedded hardware product flow | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-06_083200-implemented-first-contract-governance.md` | promoted | Implemented-first contract governance | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-06_103200-product-brief-validation-gate.md` | promoted | Product Brief validation gate | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-06_150000-same-session-doc-sync-after-code-fix.md` | promoted | Same-session doc sync after code fix | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-07_122800-performance-test-release-gate.md` | promoted | Performance test release gate | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-07_152100-private-live-adapter-smoke-gate.md` | candidate | Private live adapter smoke gate | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-07_153600-schema-derived-synthetic-fixtures.md` | candidate | Schema-derived synthetic fixtures | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-07_154400-analysis-sdk-contract-drift-gate.md` | candidate | Analysis-to-SDK contract drift gate | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-07_160200-media-metadata-private-decrypt-boundary.md` | candidate | Media metadata private decrypt boundary | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-11_093100-public-mirror-drift-gate.md` | candidate | Public mirror drift gate | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-11_093200-session-login-concurrency-matrix.md` | candidate | Session login concurrency matrix | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-12_013400-private-adapter-inside-test-module.md` | candidate | Private adapter inside test module | ✅ | `workflow/software-delivery/execution-flow.md` |
| `common/2026-05-12_043500-dart-io-httpclient-bypasses-java-frida-hooks.md` | candidate | Dart `dart:io` HttpClient Bypasses Java-Level Frida Hooks | ✅ | `workflow/software-delivery/execution-flow.md` |

---

## 2. 目標層：`analysis/development-guidance/`

這類 lessons 的 Promotion Target 包含 `implementation/`、`templates/` 或 `DOCUMENTATION.md`，適合提取到 `analysis/development-guidance/`。

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-05-05_194400-contract-first-development-flow.md` | promoted | Contract-first development flow | ✅ | `analysis/development-guidance/implementation-catalog.md` |
| `common/2026-05-05_200500-existing-project-doc-backfill-bdd-required.md` | promoted | Existing project doc backfill requires complete BDD | ✅ | `analysis/development-guidance/implementation-catalog.md` |
| `common/2026-05-06_083200-implemented-first-contract-governance.md` | promoted | Implemented-first contract governance | ✅ | `analysis/repo/contract-governance.md` |
| `common/2026-05-06_103200-product-brief-validation-gate.md` | promoted | Product Brief validation gate | ✅ | `analysis/repo/` |
| `common/2026-05-07_153600-schema-derived-synthetic-fixtures.md` | candidate | Schema-derived synthetic fixtures | ✅ | `analysis/development-guidance/` |
| `common/2026-05-07_160200-media-metadata-private-decrypt-boundary.md` | candidate | Media metadata private decrypt boundary | ✅ | `analysis/development-guidance/` |

---

## 3. 目標層：`workflow/software-delivery/artifact-gates.md`（`DOCUMENTATION.md` 對應）

這類 lessons 的 Promotion Target 包含 `DOCUMENTATION.md` 或 `CHECKLIST.md`，適合提取到 `workflow/software-delivery/artifact-gates.md`。

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-05-05_194400-contract-first-development-flow.md` | promoted | Contract-first development flow | ✅ | `workflow/software-delivery/artifact-gates.md` |
| `common/2026-05-06_081600-change-intake-before-code.md` | promoted | Change intake before code | ✅ | `workflow/software-delivery/artifact-gates.md` |
| `common/2026-05-07_081800-keep-project-incidents-out-of-skills.md` | promoted | Keep project incidents out of skills | ✅ | `workflow/software-delivery/artifact-gates.md` |

---

## 4. 目標層：`analysis/development-guidance/controls-catalog.md`

這類 lessons 的 Promotion Target 包含 `controls/`。

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `controls/2026-05-01_142100-client-encrypted-header-not-boundary.md` | promoted | Client encrypted header is not a security boundary | ✅ | `analysis/development-guidance/controls-catalog.md` |

---

## 5. 目標層：`enforcement/`

這類 lessons 影響全庫共用規則。

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-05-07_081800-keep-project-incidents-out-of-skills.md` | promoted | Keep project incidents out of skills | ✅ | `enforcement/reusable-guidance-boundary.md` |
| `common/2026-05-05_194400-contract-first-development-flow.md` | promoted | Contract-first development flow | ✅ | `enforcement/linked-updates.md` |

---

## 6. 目標層：`workflow/software-delivery/development-process.md`

這類 lessons 的 Promotion Target 包含 `process/README.md`。

| 來源檔案 | Status | 標題 | 提取狀態 | 目標檔案 |
| --- | --- | --- | --- | --- |
| `common/2026-05-05_194400-contract-first-development-flow.md` | promoted | Contract-first development flow | ✅ | `workflow/software-delivery/development-process.md` |
| `common/2026-05-05_200500-existing-project-doc-backfill-bdd-required.md` | promoted | Existing project doc backfill requires complete BDD | ✅ | `workflow/software-delivery/development-process.md` |
| `common/2026-05-05_201000-missing-requirements-block-development.md` | promoted | Missing requirements block development | ✅ | `workflow/software-delivery/development-process.md` |
| `common/2026-05-06_081600-change-intake-before-code.md` | promoted | Change intake before code | ✅ | `workflow/software-delivery/development-process.md` |
| `common/2026-05-06_082000-separate-regression-from-new-code-validation.md` | promoted | Separate regression from new code validation | ✅ | `workflow/software-delivery/development-process.md` |
| `common/2026-05-06_083000-embedded-hardware-product-flow.md` | promoted | Embedded hardware product flow | ✅ | `workflow/software-delivery/development-process.md` |
| `common/2026-05-06_083200-implemented-first-contract-governance.md` | promoted | Implemented-first contract governance | ✅ | `workflow/software-delivery/development-process.md` |
| `common/2026-05-06_103200-product-brief-validation-gate.md` | promoted | Product Brief validation gate | ✅ | `workflow/software-delivery/development-process.md` |
| `common/2026-05-06_150000-same-session-doc-sync-after-code-fix.md` | promoted | Same-session doc sync after code fix | ✅ | `workflow/software-delivery/development-process.md` |
| `common/2026-05-07_122800-performance-test-release-gate.md` | promoted | Performance test release gate | ✅ | `workflow/software-delivery/development-process.md` |

---

## 統計摘要

| 目標層 | lessons 數量 | 已提取 | 未提取 |
| --- | ---: | ---: | ---: |
| `workflow/software-delivery/` | 18 | 18 | 0 |
| `analysis/development-guidance/` | 6 | 6 | 0 |
| `workflow/software-delivery/artifact-gates.md` | 3 | 3 | 0 |
| `analysis/development-guidance/controls-catalog.md` | 1 | 1 | 0 |
| `enforcement/` | 2 | 2 | 0 |
| `workflow/software-delivery/development-process.md` | 10 | 10 | 0 |
| **總計** | **40** | **40** | **0** |

---

## 相容性說明

- ✅ `skills/app-development-guidance/feedback_history/` 已於 2026-05-13 刪除，所有 lesson 已搬遷至 `feedback/history/development-guidance/`。
- 此索引僅供 `feedback/` 層發現 lessons，不改變 lesson 的 storage location。
- 提取到目標層時，在原始 lesson 檔案開頭加入 `# Extracted — See <target path>` 標記。
- 此索引應隨 lessons 的新增或提取狀態變更而更新。
