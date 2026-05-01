# Local Proxy / Loopback Feedback Lessons

Lessons for loopback routing, local proxy handlers, Netty request extraction, and embedded routing layers.

| File | Status | Topic | Notes |
| --- | --- | --- | --- |
| `2026-04-30_120009-內建-sing-box-tun-類通道可能繞過-wi-fi-系統代理.md` | candidate | 內建 sing-box／TUN 類通道可能繞過 Wi‑Fi 系統代理 | Embedded route/TUN-like traffic bypass. |
| `2026-04-30_120010-本機-loopback-proxyserver-轉發會讓-wi-fi-http-mitm-看不到業務-connect.md` | candidate | 本機 loopback「ProxyServer」轉發會讓 Wi‑Fi HTTP MITM 看不到業務 CONNECT | Local loopback proxy routing. |
| `2026-05-01_114300-local-proxy-handler-uri-hook.md` | validated | Hook local proxy handler URI, not just OkHttp | Hook handler URI to recover upstream route. |
| `2026-05-01_131000-cast-netty-request-for-handler-route.md` | validated | Cast Netty request interfaces for handler routes | Cast Netty interfaces before reading method/URI. |
| `2026-05-01_132400-netty-aggregated-request-tostring-headers.md` | validated | Netty aggregated request toString can expose headers | Use `toString()` as private fallback for structure. |
