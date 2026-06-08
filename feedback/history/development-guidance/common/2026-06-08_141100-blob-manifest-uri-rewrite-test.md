> 遵守 [共用規則索引](../../../enforcement/README.md)、[dependency-reading](../../../enforcement/dependency-reading.md)、[neutral-language](../../../enforcement/neutral-language.md)、[goal-action-validation](../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-08 - Blob manifest URI rewrite test

Status: validated

#### One-line Summary

把 HLS / manifest 內容轉成 `blob:` URL 前，必須先把 manifest 內會被播放器再次請求的 key、segment 或 asset URI 改成可解析的絕對 URL，並用 helper unit test 覆蓋這個 rewrite。

#### Human Explanation

播放器常見流程是先下載 manifest，修改其中的 key 或 asset URI，再用 `URL.createObjectURL(new Blob([...]))` 交給瀏覽器或播放器函式庫載入。這時 manifest 的 base URL 會變成 `blob:`，不再是原始 CDN manifest URL。

如果改寫後的 manifest 還保留 root-relative 或其他相對 URI，播放器在解析這些 URI 時可能會以 `blob:` 當 base，導致 key request 或 asset request 無法產生有效 URL。HTTP probe 只驗證原始 manifest、key route、segment 都能通，仍可能漏掉這個瀏覽器端解析錯誤。

#### Trigger

- Client code 會 fetch manifest、改寫內容，再建立 `blob:` URL 給播放器。
- Manifest 內有 `#EXT-X-KEY`、segment、subtitle、init map 或其他二次請求資源。
- Live integration test 只驗 HTTP endpoint 可達，沒有測播放器實際看到的 rewritten manifest。
- 使用者回報「單獨抓 manifest/key/segment 都成功，但瀏覽器播放仍失敗」。

#### Evidence

- Tool: agent-assisted frontend debugging
- Sanitized excerpt: A playback adapter fetched a manifest and rewrote the key URI, but the rewritten manifest was loaded from a `blob:` URL. The key URI needed to be absolute same-origin; endpoint-level tests did not catch the browser URL resolution failure.
- Evidence path: Project-specific media IDs, hosts, and live run details stay in `<PROJECT_ROOT>` project tests or incident notes; this lesson only records the generalized rewrite/test rule.

#### Generalized Lesson

Manifest rewrite tests must validate the exact text that the browser/player receives, not only the upstream HTTP chain.

| Risk | Required check |
| --- | --- |
| Rewritten key URI is relative inside a `blob:` manifest | Assert key URI is absolute, usually `window.location.origin` plus the same-origin route. |
| Segment or asset URI loses original base | Assert relative media lines are resolved against the original manifest URL before blob creation. |
| Helper behavior is hidden inside a UI component effect | Extract the pure rewrite helper and test it directly. |
| Integration test only fetches original URLs | Add a unit or adapter test for the rewritten manifest output. |

#### Agent Action

When implementing manifest-to-blob playback:

1. Extract manifest rewrite logic into a pure helper.
2. In the helper, resolve media segment and asset lines against the original manifest URL.
3. For key or same-origin control routes, produce absolute URLs using the active browser origin or an injected origin.
4. Add a unit test that sets a deterministic origin, passes a representative manifest, and asserts the rewritten key and media URLs.
5. Keep live integration tests for endpoint reachability, but do not treat them as a substitute for rewritten-manifest tests.

#### Goal / Action / Validation

- Goal: Prevent playback regressions where endpoint probes pass but the player fails after loading a blob-backed manifest.
- Action: Require a focused rewrite helper test for manifest content before `URL.createObjectURL`.
- Validation or reference source: The test fails if key URI becomes root-relative, if media segment URI is not absolute, or if the helper no longer preserves the manifest key identifier.

#### Applies When

- HLS, DASH, subtitles, media playlists, or other manifests are fetched and rewritten client-side.
- A browser or player library consumes a `blob:` URL or data URL instead of the original manifest URL.
- Rewritten manifest entries include key, token, signed URL, same-origin control route, init segment, subtitle, or media segment references.

#### Does Not Apply When

- The player loads the original network manifest directly without rewriting or blob creation.
- The manifest rewrite API can explicitly pass a stable base URL to the player and tests verify that base behavior.
- The content contains only fully qualified absolute URLs and no key/control URI rewrite.

#### Validation

The prevention worked when:

- A unit or adapter test constructs a manifest with a relative key identifier and relative segment URL.
- The rewritten manifest contains an absolute key URL rooted at the expected origin.
- The rewritten manifest contains absolute media URLs resolved against the original manifest URL.
- Removing the absolute key URL rewrite causes the test to fail before browser/manual testing.

#### Promotion Target

- `workflow/software-delivery/artifact-gates.md`
- `workflow/software-delivery/execution-flow.md`
- `analysis/development-guidance/`

#### Required Linked Updates

- Updated `feedback/history/development-guidance/README.md` to include this lesson count.
- Checked reusable guidance boundary: this lesson contains generalized manifest rewrite and test guidance only; project-specific media IDs, hosts, and live evidence remain outside Ai-skill reusable docs.
- No enforcement rule update needed yet: this is a development-guidance technique, not a cross-agent behavioral failure pattern beyond the already-recorded feedback report validation gap.
