---
name: apk-analysis
description: Guides authorized APK traffic analysis, dynamic capture, Flutter/Dart AOT inspection, response decoding, evidence documentation, and reusable learning updates. Use when analyzing APK network behavior, reverse engineering API flows with permission, planning Frida/pcap/Proxyman/blutter workflows, or turning newly discovered APK analysis techniques into reusable docs.
---

# APK Analysis

Use this skill for authorized APK analysis only. The goal is to recover how an app communicates, captures, decodes, and validates data in a reproducible way, then turn findings into documentation and reusable tests without leaking secrets or target-specific private details.

**Shared policy (do not duplicate in every lesson):** read [`shared-rules` index](../../shared-rules/README.md) once (or [`feedback-lessons.md`](../../shared-rules/feedback-lessons.md) for feedback-specific rules), apply [`dependency-reading.md`](../../shared-rules/dependency-reading.md) when this skill or a related rule/template/lesson has changed, apply [`neutral-language.md`](../../shared-rules/neutral-language.md) when naming or summarizing docs, and apply [`goal-action-validation.md`](../../shared-rules/goal-action-validation.md) so important conclusions include a goal, action, validation, or reference source. When reloading this skill after an update or when the user asks whether shared rules were read, create a dependency read ledger before continuing project analysis: required files, files read, `not applicable` missing files (for example no `CHECKLIST.md`), blocked items, and validation method. Per-technique files in `feedback_history/` should **reference** shared rules, not paste the full common rules.

**Cross-skill references:** follow [`cross-skill-references.md`](../../shared-rules/cross-skill-references.md). `apk-analysis` may reference another skill when analysis output must be consumed by that skill, but it must name the trigger, handoff artifact, ownership boundary, and sanitization boundary instead of copying the target skill's full workflow.

## Quick Start（Routing）

1. **Confirm scope and authorization**: Identify the APK, version, device/emulator, architecture, and allowed analysis actions. Do not collect or publish credentials, full tokens, private keys, personal data, or unrelated third-party traffic.
2. **Separate method from target facts**: Reusable techniques → this skill folder. Target-specific endpoints, hosts, schemas, findings → project API/reference docs.
3. **Network path triage**: localhost → whole-device pcap → system proxy/MITM → Java HTTP hooks → Flutter Dart AOT/native paths.
4. **Route into technique category** only after evidence points there. Use [`techniques/`](techniques/) as category index. If APK type unknown, stay in common workflow.
5. **Build UI architecture map** → see [`runtime/onboarding/apk-analysis-quickstart.md`](../../runtime/onboarding/apk-analysis-quickstart.md) § 步驟 5.
6. **Prefer high-semantic hooks**: request object hooks > raw socket hooks; response interceptor/decode hooks > TLS byte reconstruction.
7. **Convert dynamic results into durable assets** → see [`runtime/onboarding/apk-analysis-quickstart.md`](../../runtime/onboarding/apk-analysis-quickstart.md) § 步驟 7.
8. **Automatic skill feedback** → see [`shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md) and [`feedback/`](../../feedback/) layer.

## Default Workflow

Read [WORKFLOW.md](WORKFLOW.md) before doing hands-on analysis. Use it as the default decision tree.

Use [TOOLS.md](TOOLS.md) when preparing an environment or choosing between adb, pcap, Proxyman/mitmproxy, Frida, jadx, apktool, blutter, or offline decoding scripts.

Use [DOCUMENTATION.md](DOCUMENTATION.md) when writing human-readable results.

Use [`techniques/`](techniques/) after common triage identifies a category. Only read the matching category folder unless the evidence crosses categories.

Use [shared-rules/feedback-lessons.md](../../shared-rules/feedback-lessons.md) for **how** to write feedback; put each lesson in the matching [`feedback_history/<category>/`](feedback_history/) folder. **Agents:** treat this as mandatory whenever such an idea appears—see **Quick Start §8** and **Feedback Loop** below.

Use [`app-development-guidance`](../app-development-guidance/) when analysis findings should become app development guidance, implementation patterns, PR/release checklists, or validation tests.

Automatically read and apply [`app-development-guidance/SKILL.md`](../app-development-guidance/SKILL.md) when the user wants APK analysis documents to produce an app-related tool, SDK, client, mock API, fixture-driven implementation, contract test, or rebuilt feature. Before drafting implementation plans, first apply the **Development readiness gate** above: if the work must reach a real backend or replay a live feature, the project analysis docs must already contain the minimum runnable factors, including the authorized identity material self-generation audit when identity/session/device material is involved, or the missing factors must be promoted to blocker questions. Do this so missing runtime, BDD, contract, error-handling, storage, security, ownership, or test questions are surfaced before development starts instead of being invented inside `apk-analysis`.

When the user wants a feature rebuilt from APK findings or wants analysis docs turned into app tools / SDK work, use this cross-skill handoff:

- Target skill: [`app-development-guidance`](../app-development-guidance/).
- Trigger: APK findings must become rebuildable app behavior, app-related tools, SDK/client behavior, API/interface contracts, implementation slices, fixtures, mocks, or tests.
- Handoff artifact: Feature Reconstruction Handoff with sanitized behavior, domain, API/interface, state/error, data lifecycle, fixture, and open-question detail.
- Ownership boundary: `apk-analysis` owns evidence recovery, traffic/UI attribution, schema notes, fixtures, and confidence labels; `app-development-guidance` owns BDD, Domain Model Contract, API / Interface Contract, Error Handling Contract, implementation guidance, checklists, and tests.
- Sanitization boundary: target-specific hosts, tokens, raw responses, accounts, and private business conclusions stay in project docs.

Use [RUNBOOK.md](RUNBOOK.md) when starting a new APK project or when the user asks how to apply this skill to another product.

## Output Style & Artifact Gates

See [`workflow/apk-analysis/artifact-gates.md`](../../workflow/apk-analysis/artifact-gates.md) for output format, quality gates, and sanitization rules.

## Feedback Loop

See [`shared-rules/feedback-lessons.md`](../../shared-rules/feedback-lessons.md) for the feedback lesson template and workflow. See [`feedback/`](../../feedback/) for the feedback promotion pipeline.

**Git（本 repository）：**若在**同一工作區**修改了 `<AI_SKILL_REPO>` 底下的 `skills/apk-analysis/`、`shared-rules/` 等（含新建 `feedback_history/` 檔），**除非使用者明講不要提交**，否則在結束任務前**必須**於 `<AI_SKILL_REPO>` 根目錄執行 `git status`，將相關變更 **`git add` → `git commit`**（訊息清楚）→ **`git push`**；需要權限時**必須**向使用者申請（例如 git_write／網路）。第一次寫入時先依 [`dependency-reading.md`](../../shared-rules/dependency-reading.md) 開啟 writeback transaction；僅重載或重讀任一工具/skill **不會**自動完成這一步。若本輪明確使用或更新本機 tool sync / bundle mirror，才依 [`ai-tools/`](../../ai-tools/README.md) 執行對應同步；reference-only 不需要同步。
