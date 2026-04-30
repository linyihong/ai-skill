### 2026-04-30 - Proxy failure 要先拆成導流與 TLS 兩層

Status: promoted

#### One-line Summary

代理看不到明文時，先確認「有沒有進代理」，再談憑證或 pinning。

#### Human Explanation

很多人看到 Proxyman / Burp / mitmproxy 沒有明文，就直接判斷是 certificate pinning。這常常太早下結論。更可靠的判斷順序是先看 App 的連線目標是否已經變成 proxy；如果還是直連目標 host，問題是導流或初始化時機，不是 TLS。只有已經進 proxy 且 handshake 失敗時，才應該查 CA trust、network security config、custom trust 或 pinning。

#### Trigger

MITM 工具沒有看到明文，或顯示 `SSL Handshake Failed`。

#### Evidence

在授權 APK 分析中，曾觀察到兩種完全不同的 failure：

- 流量沒有進 proxy，裝置仍直接連目標 host `:443`。
- 流量已進 proxy，但 TLS handshake 因 CA trust / pinning 失敗。

#### Generalized Lesson

不要把「代理工具看不到明文」直接等同於 pinning。先看是否有 CONNECT / connect target 到 proxy；只有導流成功後，才進入 CA / pinning 排查。

#### Agent Action

下次遇到 MITM 失敗時，先要求或執行導流驗證：檢查 proxy 是否收到 CONNECT，或用 connect trace 觀察目標是否為 `<proxy-host>:<proxy-port>`。不要先寫 pinning 結論。

#### Promotion Target

已整理到 `WORKFLOW.md` 與 `TOOLS.md`。
