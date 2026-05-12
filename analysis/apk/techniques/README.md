# APK Analysis Technique Categories

`analysis/apk/techniques/` 負責 APK 分析中各技術路線的深度分析方法。本目錄保存特定 app/runtime/API 家族的技術指引，讓分析過程可根據證據選擇對應路線。

> **相容性規則**：`skills/apk-analysis/techniques/` 仍為 active skill entrypoint。本目錄為 reference target，兩者應保持同步。

## Routing Rule

Start with the common workflow until evidence identifies the route, then read only the matching category:

| Evidence / Goal | Category |
| --- | --- |
| `libflutter.so`, `libapp.so`, Dart AOT strings, Dio/interceptor clues | [`flutter-dart-aot/`](flutter-dart-aot/) |
| HTTP/HTTPS API docs, headers, request/response fields, replay fixtures | [`http-api/`](http-api/) |
| `127.0.0.1:<port>`, Netty handler, app-local proxy, loopback bridge | [`local-proxy/`](local-proxy/) |
| HLS playlist/key/segments, media containers, images, audio/video extraction | [`media-hls/`](media-hls/) |

If no category is identified yet, do not read every category. Continue the common network path triage until there is evidence.

## Category Rules

- Put stable, reusable guidance for a family in its category folder.
- Keep target-specific hosts, endpoints, secrets, schemas, and private findings in project docs, not here.
- Keep one-off or not-yet-promoted discoveries in `skills/apk-analysis/feedback_history/`; promote them into a category only after they are credible.
- Cross-cutting rules that apply to all APKs stay in `analysis/apk/` top-level docs.
- When a new family appears, create `techniques/<category>/README.md` and link it here instead of expanding unrelated categories.

## Migration Notes

- 本目錄為 Phase 19 提取產物，內容來自 `skills/apk-analysis/techniques/`。
- 舊入口 `skills/apk-analysis/techniques/` 仍為 active source of truth。
- 未來遷移完成條件：所有 technique categories 完全提取、索引更新、舊入口保留 redirect reference。
