# Evidence-First Analysis Routing（以證據驅動分析路線選擇）

**Status**: `validated-intelligence`
**Source**: Multiple `feedback/history/apk-analysis/common/` lessons

## 問題

APK analysis session 經常浪費 token，原因如下：

1. 在證據縮小路線之前，就先讀取所有 technique categories。
2. 在確認流量是否到達 proxy 之前，就先假設是 certificate pinning。
3. 當 stack 是 Flutter/Dart native 時，卻先加 broad Java hooks。
4. 把「proxy 顯示無流量」解讀為「沒有網路活動」，而非「client 繞過系統 proxy」。

## 原則

**讓證據驅動路線選擇，而非假設。**

在選擇工具或 technique category 之前，依序回答這些問題：

1. **是否有 localhost 流量？** → pcap loopback interface。
2. **是否有外部 TLS？** → whole-device pcap 看 SNI/IP/timing。
3. **流量是否到達系統 proxy？** → 檢查 CONNECT/proxy log。
4. **使用哪個 HTTP stack？** → Java hook → native connect trace → Flutter/Dart AOT。
5. **只有在證據指向某個 category 之後**，才載入對應的 technique docs。

## 決策表

| 證據 | 路線 | 不要 |
| --- | --- | --- |
| 沒有 CONNECT 到 proxy | 檢查 proxy config、injection timing 或 native stack | 假設是 pinning |
| Java hook 沒有顯示 target host | 切換到 native/Flutter trace | 繼續加更多 Java hooks |
| 偵測到 Flutter libapp.so | 使用 blutter/unflutter + Dart AOT hooks | 讀取 HTTP API technique docs |
| Proxy CONNECT 成功但 TLS 失敗 | 檢查 CA trust、network config、pinning | 假設「proxy 不能用」 |
| 只有 pcap SNI 可見 | 在 decompiled code 中靜態搜尋 host/path | 跳過 request shape 直接解密 |
| Frida constructor chain 顯示 `PBC.ctor`（PaddedBlockCipher）但 mode（CBC vs SIC vs GCM）無法從 `processBlock` count 區分 | 執行 live proxy test：分別用不同 mode 加密相同明文，比對 HTTP response status code | 僅依賴 block count 猜測 mode（43 blocks 可對應多種 mode） |

## Token 影響

遵循 evidence-first routing 可以將初始 context 負載從讀取所有 technique folders（每個約 4500 tokens）減少到只讀取匹配的 category（先讀 summary 約 500 tokens）。

---

← [回到 intelligence/engineering/analysis/](README.md)
