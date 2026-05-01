# APK Analysis Technique Categories

Use this directory to keep technique-specific guidance out of the default workflow. The top-level `WORKFLOW.md`, `TOOLS.md`, and `DOCUMENTATION.md` stay as common routing and documentation rules. Category folders hold deeper guidance for a specific app/runtime/API family.

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
- Keep one-off or not-yet-promoted discoveries in `feedback_history/`; promote them into a category only after they are credible.
- Cross-cutting rules that apply to all APKs stay in top-level docs.
- When a new family appears, create `techniques/<category>/README.md` and link it here instead of expanding unrelated categories.
