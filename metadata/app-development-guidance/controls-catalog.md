# Security Controls Catalog

本文件列出跨平台安全控制的種類與核心原則。承接 [`skills/app-development-guidance/controls/`](../../skills/app-development-guidance/controls/) 的內容，提取為 tool-neutral 的分析參考。

> **遷移狀態**：此文件為新分層的 reference target，`skills/app-development-guidance/controls/` 已不再作為 active entrypoint。新內容請直接寫入此文件。

## 控制類型

| 控制 | 核心原則 | 原始來源 |
|------|----------|----------|
| **API Transport** | Enforce HTTPS, pinning as risk-based control, server-side freshness controls (timestamp/nonce/idempotency), backend auth must not trust client flags | [`controls/api-transport.md`](../../skills/app-development-guidance/controls/api-transport.md) |
| **Auth & Session** | Token scope, refresh, revocation, logout, session invalidation, account binding | [`controls/auth-session.md`](../../skills/app-development-guidance/controls/auth-session.md) |
| **Local Storage** | Secure storage, cache, backups, screenshots, offline data | [`controls/local-storage.md`](../../skills/app-development-guidance/controls/local-storage.md) |
| **Logging & Telemetry** | Log redaction, crash reports, analytics, security observability | [`controls/logging-telemetry.md`](../../skills/app-development-guidance/controls/logging-telemetry.md) |
| **Anti-Tamper** | Root/hook/emulator signals, anti-tamper limits, risk scoring | [`controls/anti-tamper-risk.md`](../../skills/app-development-guidance/controls/anti-tamper-risk.md) |
| **Release Build** | Obfuscation, debug flag removal, symbol stripping, dependency and secret checks | [`controls/release-build.md`](../../skills/app-development-guidance/controls/release-build.md) |

## 使用原則

1. **Prefer `controls/` before `platforms/`, `languages/`, or `implementation/`** — 當 lesson 主要是關於安全屬性而非框架特定或可建置步驟時，先放在 controls。
2. **Controls 是跨平台的** — 核心原則在此定義，平台差異在 `platforms/`，可建置步驟在 `implementation/`。
3. **Client-side controls 不是 authorization 的替代品** — 客戶端控制可被繞過，backend 控制才是最後防線。

## 常見過度宣稱

- Pinning 不能防止所有逆向工程。
- Request signing 在 app 被分析後無法保護靜態 client secret。
- Encrypted request headers 不能對可在加密前 instrument client 的一方隱藏明文。
- Client checks 不能取代 server-side authorization。

## 與其他層的關係

- `analysis/app-development-guidance/risk-translation.md` 提供如何將觀察轉譯為控制的流程。
- `intelligence/engineering/development/risk-translation-heuristic.md` 提供控制選擇的背後原則。
- `skills/app-development-guidance/controls/` 是原始來源，已不再作為 active entrypoint。
