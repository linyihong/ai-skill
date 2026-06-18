> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson。

### 2026-06-18 - Play-view KPI: SQL/API pass, DOM still wrong

Status: validated

#### One-line Summary

用户可见的播放量 / KPI 标签争议时，**API 或 listing SQL 与 DB 一致不足以结案**；须在同一 `video_id`、同一路径上增加 **browser DOM 对照**，才能区分 DB/部署问题与 keep-alive / sessionStorage 陈旧 view model。

#### Human Explanation

Integration 可证明 `portal_record_play_view` 写入正确、`view_count + play_count` 查询正确、search API 返回 display total。但若首页 continuation 从 sessionStorage 还原旧的 `viewCountText`，用户仍看到 source-only 数字。这是 **L3 投影 / client cache** 问题，不是 SQL bug。

#### Trigger

- 用户回报「看了没加播放量」「列表数字不对」但 DB 或 API 抽查已正确。
- 仅 fetch/SQL integration 通过；无 browser 或 DOM 断言。
- 页面使用 tab keep-alive、sessionStorage continuation、或 SSR 首屏 + client load-more 合并 view model。

#### Evidence

- Tool: Vidoe-Test play-view-dedup v1 (2026-06-18)
- Sanitized excerpt: Fixture video DB `24990 + 2 = 24992`; search API/HTML `24992`; home 擦边剧 load-more card showed `24990次播放` until continuation refetch on restore + pull-to-refresh cache clear.
- Evidence path: `external/Vidoe-Test/docs/plans/2026-06-18-1430-play-view-dedup-event-db.md`, `tests/integration/play-view-display-count.browser.integration.mjs` — project-local only.

#### Generalized Lesson

For **user-visible counters** tied to server KPI:

| Layer | Validates |
| --- | --- |
| DB | Counter columns / dedup tables |
| API / listing SQL | Same formula as product contract |
| Browser DOM | Same video id, same user path, rendered label |

Passing layer N does not imply layer N+1 without an explicit gate.

#### Agent Action

1. Reproduce with the **same video_id** on the path the user used (search vs category load-more vs detail).
2. Run DB + API alignment first (cheap).
3. If API correct and user report persists, run CDP/browser test: scroll to card, read overlay text, optionally seed stale sessionStorage and assert refetch.
4. Classify: **Integration** (SQL/deploy) vs **Continuation/cache** (client restore) before patching SQL again.

#### Goal / Action / Validation

- Goal: Avoid mis-diagnosing keep-alive cache as missing SQL or failed deploy.
- Action: Add a browser regression gate when KPI is rendered in DOM after async/continuation merge.
- Validation: Browser test fails if DOM shows source-only count while DB display total is higher; passes after refetch/clear path.

#### Applies When

- KPI shown on list cards with load-more, keep-alive tabs, or sessionStorage-backed continuation.
- Dispute is **display** not **write** (user saw wrong number vs count not incremented).

#### Does Not Apply When

- Bug is purely server dedup or increment (API `playViewRecorded` / DB counters wrong).
- Page is SSR-only with no client-side view model cache on that surface.

#### Promotion Target

- `workflow/software-delivery/test-strategy.md` (KPI / counter row)
- `workflow/software-delivery/validation.md` (L3 projection)
- `workflow/software-delivery/perf-governance.md` (POST burst generalization — separate lesson)

#### Required Linked Updates

- Project lesson: `external/Vidoe-Test/.ai-skill/project/feedback/play-view-display-api-vs-dom-gate.md`
- Perf second incident pointer: `governance/evidence-candidates/evidence-rules/play-view-dedup.pointer.yaml` (optional index)
