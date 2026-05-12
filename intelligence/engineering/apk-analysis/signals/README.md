# APK Analysis Signal Detection

`intelligence/engineering/apk-analysis/signals/` 存放 APK 分析中用來辨識技術特徵的信號檢測知識。

## Scope

本目錄負責：

- Flutter/Dart AOT 的辨識信號（libapp.so、Dart snapshot、AOT 混淆特徵）
- HTTP 流量層級的辨識信號（OkHttp vs dart:io HttpClient vs WebSocket）
- 代理導流成功與否的信號
- 媒體串流類型的辨識信號（HLS vs DASH vs progressive download）

## 與其他層的關係

- `analysis/apk/traffic-triage.md` 提供流量分流流程，本目錄提供分流所需的判斷信號
- `intelligence/engineering/apk-analysis/heuristics/` 使用本目錄的信號來決定策略

## 目前 atoms

| Atom | 說明 | 來源 | 跨領域推廣 |
|------|------|------|-----------|
| [`flutter-dart-aot-detection.md`](flutter-dart-aot-detection.md) | Flutter/Dart AOT 辨識信號 — 主要/次要/排除信號表與判斷流程 | `skills/apk-analysis/techniques/flutter-dart-aot/README.md` | — |
| [`local-proxy-detection.md`](local-proxy-detection.md) | Local Proxy 偵測信號 — 主要/次要/排除信號表與判斷流程 | `skills/apk-analysis/techniques/local-proxy/README.md` | — |
| [`media-type-detection.md`](media-type-detection.md) | 媒體類型偵測信號 — Magic Bytes 參考表、靜態 vs 動畫判斷、Container Probe 指令 | `skills/apk-analysis/techniques/media-hls/README.md` | Magic Bytes 參考表已提取到 [`intelligence/engineering/heuristics/magic-bytes-reference.md`](../../heuristics/magic-bytes-reference.md) |
