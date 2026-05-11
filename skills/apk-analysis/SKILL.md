---
name: apk-analysis
description: Guides authorized APK traffic analysis, dynamic capture, Flutter/Dart AOT inspection, response decoding, evidence documentation, and reusable learning updates. Use when analyzing APK network behavior, reverse engineering API flows with permission, planning Frida/pcap/Proxyman/blutter workflows, or turning newly discovered APK analysis techniques into reusable docs.
---

# APK Analysis

Use this skill for authorized APK analysis only. The goal is to recover how an app communicates, captures, decodes, and validates data in a reproducible way, then turn findings into documentation and reusable tests without leaking secrets or target-specific private details.

**Shared policy (do not duplicate in every lesson):** read [`shared-rules` index](../../shared-rules/README.md) once (or [`feedback-lessons.md`](../../shared-rules/feedback-lessons.md) for feedback-specific rules), apply [`dependency-reading.md`](../../shared-rules/dependency-reading.md) when this skill or a related rule/template/lesson has changed, apply [`neutral-language.md`](../../shared-rules/neutral-language.md) when naming or summarizing docs, and apply [`goal-action-validation.md`](../../shared-rules/goal-action-validation.md) so important conclusions include a goal, action, validation, or reference source. When reloading this skill after an update or when the user asks whether shared rules were read, create a dependency read ledger before continuing project analysis: required files, files read, `not applicable` missing files (for example no `CHECKLIST.md`), blocked items, and validation method. Per-technique files in `feedback_history/` should **reference** shared rules, not paste the full common rules.

**Cross-skill references:** follow [`cross-skill-references.md`](../../shared-rules/cross-skill-references.md). `apk-analysis` may reference another skill when analysis output must be consumed by that skill, but it must name the trigger, handoff artifact, ownership boundary, and sanitization boundary instead of copying the target skill's full workflow.

## Quick Start

1. Confirm scope and authorization:
   - Identify the APK, version, device/emulator, architecture, and allowed analysis actions.
   - Do not collect or publish credentials, full tokens, private keys, personal data, or unrelated third-party traffic.
2. Separate method from target facts:
   - Put reusable techniques in this skill folder.
   - Put target-specific endpoints, hosts, schemas, and findings in the project’s API/reference docs.
   - Keep target-specific analysis documents scoped to the observed behavior of the APK under analysis. Do not mix in implementation plans, product design, self-built app architecture, clone behavior, or facts from other apps. If findings need to become development guidance, put them in a clearly separate handoff / guidance document and label them as derived guidance.
3. Start with network path triage:
   - Check localhost traffic.
   - Check whole-device pcap for external TLS/SNI/host timing.
   - Check whether system proxy/MITM is actually used.
   - Check Java HTTP hooks only if evidence suggests Java networking.
   - For Flutter apps, inspect Dart AOT/native paths early.
4. Route into a technique category only after evidence points there:
   - Use [`techniques/`](techniques/) as the category index.
   - Do not read every category by default; Flutter/Dart, local proxy, media, and HTTP API docs are separate routes.
   - If the APK type is still unknown, stay in the common workflow until traffic/runtime evidence narrows it.
5. Build a UI architecture map when the device/app can be operated:
   - Start lightweight: capture only enough sanitized screenshots/UI hierarchy to understand major tabs, drawers, routes, and key screens.
   - Mark which screens are scrollable and which visible elements are clickable entry points.
   - Document how to reach each important screen, including entry state, tap/swipe steps, expected destination, and reusable operation id.
   - When the user asks to analyze a named page/tab/module, create or update a project-level page map artifact (for example `docs/UI架構地圖/<page>.md`) instead of leaving the UI-to-API findings only inside API docs, tool docs, or chat.
   - Before reporting the analysis complete, force a UI Behavior backfill: update the project's UI Behavior entry/index and the relevant page-level map with observed App actions, visible sort labels, tap/swipe steps, data source mapping, evidence, and unknowns. If UI behavior was not captured, explicitly mark `needs capture` / `Trigger confidence: low` in the project docs rather than omitting the section.
   - Keep operation maps scoped to in-app pages; if a step opens another app, system screen, browser, or external intent, document the transition instead of treating it as an app screen.
   - For key flows, optionally create a small replayable app-operation script so API capture can be repeated with stable timing.
   - Adapt the order if screenshots or device control make the app slow: solve core API/decode first, then bind important APIs back to UI actions.
   - Record the exact UI path and action window only for flows that need API attribution.
6. Prefer high-semantic hooks:
   - Request object hooks are better than raw socket hooks.
   - Response interceptor/decode hooks are better than trying to reconstruct TLS bytes.
7. Convert dynamic results into durable assets:
   - UI architecture map and operation-to-API matrix.
   - Redacted HTTP/API docs with headers, request fields, response fields, and per-field meaning/type notes.
   - **Domain / runtime baseline（實取數據門檻）:** API 條目與 schema **不足以**讓 SDK、client 或回放工具自動連線。必須在**專案**內另有一份 **domain／執行環境基線**（見 [`DOCUMENTATION.md`](DOCUMENTATION.md) § *Domain／執行環境基線*）：host／path family 與環境維度（去敏占位）、TLS／代理／內嵌通道假設、**session／登入／裝置身分**對列表資料的依存、必填 **opaque query**（例如 session 標量）從哪些操作或 response 衍生、簽章或 gateway 前置條件（不留 secret）、分頁**地面真相**（明確旗標 vs 長度推斷）、限流／錯誤恢復、`service`/`serviceHash` 與公開端點的對應表（值可 placeholder）。並以表格 **連回** API Catalog 條目與 UI operation id。可新建 `docs/domain-baseline.md`、`api/domain-environment.md` 或附在現有 inventory；**不放進**可重用 skill 正文。
   - **Development readiness gate:** if the next step is SDK/client/app/tool development that must talk to a real service, replay a feature, or run live integration, do **not** start implementation from API shape alone. First verify the project docs contain the minimum runnable factors: endpoint/path family, route/service mapping or placeholder strategy, session/bootstrap dependency, opaque parameter source and lifetime, authorized identity material self-generation audit, signing/gateway prerequisites, response decrypt/unwrap boundary, pagination truth, error/session recovery, and a replay checklist. Missing factors must be written as blockers or explicitly scoped out; only offline parser, fixture, mock, or docs-only work may proceed from a skeleton baseline.
   - **Authorized identity material self-generation audit:** when device, install, account, session seed, vendor attestation, or server-issued session material is required for live access, the analysis must answer two separate questions for each key group: **can an SDK/tool generate or initialize it without the target app runtime?** and **if yes, what sanitized generation recipe, lifecycle, reset/cooldown, and validation matrix prove that?** If the answer is no or unknown, name the provider boundary (`caller-provided`, `server-issued`, `trusted-bridge`, `private-adapter-required`, or `unknown`) and treat it as a live-development blocker unless explicitly scoped to a private adapter.
   - **API Catalog / API list finish gate:** when the task needs an API reference, SDK/client input, mock API, contract test, or rebuildable feature, create or update a project-level API Catalog with total API entry, grouped indexes, per-API detail files, coverage/gap status, UI/API mapping, SDK/client field usage, evidence, validation, and open questions. Do not leave confirmed APIs only in schema catalogs, correlation tables, hook logs, or chat.
   - Feature reconstruction handoff: capability, behavior scenarios, candidate domain concepts, API/interface contracts, state/error handling, data lifecycle, fixtures, and open questions.
   - **Feature handoff finish gate:** once a named feature/module is understood well enough to mark core flows `Confirmed` or to answer implementation-style questions, create or update a project-level feature handoff / architecture document in the same session. Do not leave the understanding only in API tables, hook logs, chat, or page maps.
   - Redacted request/response samples.
   - Offline decoders or fixtures.
   - API/schema docs.
   - Contract tests where the project has an SDK or client implementation.
8. **Automatic skill feedback (agents)**:
   - Whenever you learn a **new reusable** technique, failure pattern, or validation rule during analysis, **write it into this skill in the same session**—do **not** wait for the user to say「記得回饋」.
   - Treat a user's improvement idea, reviewer comment, or question such as「這個能不能加入技巧」/「為什麼沒有觸發閉環」as an immediate feedback-trigger check. First decide whether the idea generalizes beyond the current APK; if yes, open an Ai-skill writeback transaction before continuing long project work.
   - Do not require the idea to be fully validated before writing it down. If it is plausible but not yet replay-proven, create a `candidate` or `experimental` lesson with clear validation criteria and defer promotion until evidence exists.
   - Add **one new Markdown file** under the matching [`feedback_history/<category>/`](feedback_history/) folder using [shared-rules/feedback-lessons.md](../../shared-rules/feedback-lessons.md) naming + template (generalized, sanitized, with evidence and applicability). Use `feedback_history/common/` for cross-category rules. Optionally add a row to [`feedback_history/README.md`](feedback_history/README.md).
   - If the lesson is already credible from evidence in this session, also patch [TOOLS.md](TOOLS.md), [WORKFLOW.md](WORKFLOW.md), or [DOCUMENTATION.md](DOCUMENTATION.md) as appropriate; label uncertain items `experimental` in that lesson file instead of promoting prematurely.
   - If the reusable lesson is about how to build **your own future apps** more safely, write the development guidance in [`app-development-guidance`](../app-development-guidance/) and keep only the APK-analysis method here.
   - Target-specific hosts, endpoints, tokens, or one-off product conclusions stay in the **project** docs, not in reusable skill files.

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

## Required Output Style

When reporting analysis progress, include:

- Goal, action, and validation or reference source for the current work unit.
- What was tested.
- What evidence was observed.
- What was ruled out.
- What remains unknown.
- The next lowest-risk action.

When documenting a new finding, include:

- Trigger or UI path.
- Tool and command summary.
- Evidence file path or sanitized excerpt.
- Feature/capability mapping and operation id when the finding supports functional reconstruction.
- Page-level UI map path when the task targets a named page/tab/module and UI/API mapping was established.
- **Domain／執行環境基線**更新點（若影響 host 族、session、opaque 欄位、分頁或簽章前置）：指向專案內 baseline 檔或小節錨點。
- Generalized lesson.
- Follow-up validation.

## Safety and Sanitization

Never write raw secrets into reusable skill docs:

- Full Authorization tokens.
- Session cookies.
- Device identifiers that belong to a real user/device.
- AES/HMAC secrets unless they are synthetic examples.
- Private hostnames not meant for the reusable method guide.
- Personal user data.

Use placeholders:

```text
<package-name>
<device-serial>
<api-host>
<proxy-host>:<proxy-port>
<token-redacted>
<secret-redacted>
```

## Feedback Loop

If analysis discovers a new reusable idea, or a user/reviewer suggests an improvement that may generalize:

1. Create **`feedback_history/<category>/YYYY-MM-DD_HHMMSS-<slug>.md`** as a dated lesson **proactively** (same session as the discovery unless blocked by missing evidence). Use `common/` when the lesson is cross-category. Follow [feedback-lessons.md](../../shared-rules/feedback-lessons.md) naming rules (`HHMMSS` = local 24h time).
2. Generalize it so it is not tied to one APK.
3. Add evidence and validation criteria.
4. Promote it into `WORKFLOW.md`, `TOOLS.md`, or `DOCUMENTATION.md` only after it has been validated or is clearly labeled as experimental in the lesson file.

**Root-cause check when feedback did not trigger:** before patching only the missed lesson, identify whether the trigger was too implicit, the idea was wrongly treated as project-only, validation uncertainty blocked a `candidate` lesson, or the Ai-skill writeback transaction was not opened. Strengthen the correct trigger/checklist text in this skill or shared rules, then complete sync/commit/push before returning to long-running APK analysis.

Do not silently overwrite prior lesson files. Add new files or add a short deprecation note in an older file pointing to the replacement.

**Agent checklist before ending an APK-analysis task:** Did any **new generalized lesson** emerge, including from user suggestions or review feedback? If yes → **`feedback_history/<category>/`** or **`feedback_history/common/`** has a new file (minimum); optional promotion to `TOOLS.md` / `WORKFLOW.md` / `DOCUMENTATION.md` / `techniques/<category>/` when justified. If nothing new → no forced entry, but answer must state why the idea stayed project-specific or unvalidated when the user explicitly asked about skill feedback.

**Git（本 repository）：**若在**同一工作區**修改了 `<AI_SKILL_REPO>` 底下的 `skills/apk-analysis/`、`shared-rules/` 等（含新建 `feedback_history/` 檔），**除非使用者明講不要提交**，否則在結束任務前**必須**於 `<AI_SKILL_REPO>` 根目錄執行 `git status`，將相關變更 **`git add` → `git commit`**（訊息清楚）→ **`git push`**；需要權限時**必須**向使用者申請（例如 git_write／網路）。第一次寫入時先依 [`dependency-reading.md`](../../shared-rules/dependency-reading.md) 開啟 writeback transaction；僅重載或重讀任一工具/skill **不會**自動完成這一步。若本輪明確使用或更新本機 tool sync / bundle mirror，才依 [`ai-tools/`](../../ai-tools/README.md) 執行對應同步；reference-only 不需要同步。
