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

（pilot 階段逐步建立）
