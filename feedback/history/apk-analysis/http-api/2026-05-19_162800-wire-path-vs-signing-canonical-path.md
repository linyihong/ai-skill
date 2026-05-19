# Wire Path vs Signing Canonical Path

Status: candidate

## Context

在分析加密／簽章 API 時，live request 可能同時存在兩個不同的 path 概念：

- 實際 HTTP wire URL path：request 送到 server 的 path。
- 簽章 canonical path：App 內部餵給 `eh` / signature generator 的 path material。

這兩者不一定完全相同。TATA gossip live smoke 中，HTTP wire path 使用 `/v1/api/public/` 才會回到 encrypted JSON；但 `eh` 產生仍使用 App 內部 `api/public/?...` canonical material。若把兩者混成同一個值，server 可能回 HTML/error page，表面上像是授權、簽章或 decrypt 缺口。

## Rule

遇到 encrypted API 回 HTML 或非 JSON 時，不要先歸咎於 identity/signing/decrypt 缺失。先分別驗證：

1. wire URL path 是否與實際 app request family 一致；
2. signing canonical path 是否與 App crypto helper 取用的 material 一致；
3. 兩者是否需要分開建模與文件化。

## Evidence

- TATA guest login 文件與既有實作使用 `/v1/api/public/` 作為 wire path。
- 既有 `DartEncryptAESProvider` / `GuestLoginClient` 註解顯示 `eh` signing path 使用 `api/public/`，不是 `v1/api/public/`。
- Gossip live smoke 從 `/api/public` 改為 `/v1/api/public/` 後，categories/articles/detail 由 HTML failure 進入可解析 response，並能下載 detail content image。
