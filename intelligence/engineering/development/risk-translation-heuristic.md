# Risk Translation Heuristic（風險轉譯經驗法則）

**Status**: `candidate-intelligence`
**Source**: [`analysis/development-guidance/risk-translation.md`](../../analysis/development-guidance/risk-translation.md), [`workflow/software-delivery/execution-flow.md`](../../workflow/software-delivery/execution-flow.md)

## 原則

**If an observation describes what an attacker or competitor can do, translate it into what a developer must prevent.**

如果一個觀察描述攻擊者或競爭對手可以做什麼，將其轉譯為開發者必須防止什麼。

## 為什麼

1. **安全分析與開發指引之間存在語言鴻溝** — 分析師說「request can be replayed」，開發者需要知道「implement nonce + timestamp check」。
2. **沒有轉譯的觀察無法變成 actionable 的開發任務** — 「Token is long-lived」不是一個 ticket，但「Implement token rotation with configurable TTL」是。
3. **控制層選擇決定防禦深度** — 同一個風險可以在 client、backend、contract、tooling、build 等不同層級解決，選擇錯誤的層級會產生虛假的安全感。
4. **Client-side hardening 不是 authorization 的替代品** — 客戶端控制可以被繞過，backend 控制才是最後防線。

## 何時適用

- 從 APK analysis、penetration test、code review 中產出開發指引。
- 將安全觀察轉換為 product backlog items。
- 設計新的 control 或 security feature 時。
- 評估現有控制的充分性時。

## 何時不適用

- 純粹的架構討論（沒有具體的攻擊面觀察）。
- 已知且已解決的風險（已有 control 且驗證有效）。
- 專案已明確接受該風險（residual risk 已記錄）。

## 決策流程

```text
有具體觀察？
  ├── 這是攻擊者能做什麼？
  │     └── 轉譯為開發者必須防止什麼
  ├── 這個控制應該在哪一層？
  │     ├── Backend/API → authorization, replay defense, rate limits
  │     ├── Client app → safe storage, secure defaults
  │     ├── Full-stack contract → OpenAPI, typed clients, contract tests
  │     ├── Tooling → rule engine, diagnostics
  │     ├── Third-party → sanitized excerpts, live-test gates
  │     ├── Embedded/firmware → hardware context, driver boundary
  │     ├── Build/release → obfuscation, symbol stripping
  │     └── Monitoring → anomaly detection, abuse patterns
  ├── 定義 control：required control + owner + implementation note + validation + residual risk
  └── 歸檔到對應的 controls/ 目錄
```

## 常見誤用

| 誤用 | 正確 |
|------|------|
| 「這是 client 端的問題，讓 client team 修」 | Client-side 控制可被繞過；backend 必須有獨立防禦 |
| 「這個風險太技術性，開發者不會懂」 | 轉譯為開發者熟悉的語言和框架模式 |
| 「我們已經有 WAF，所以不用擔心」 | WAF 不是應用層控制的替代品 |
| 「把 vendor 的 security doc 直接貼到 controls 裡」 | 提取 sanitized 的整合指引，保留 vendor 原文在專案文件 |

## Token Impact

避免將未轉譯的安全觀察直接放入開發指引，導致開發者無法理解或無法採取行動。一個好的轉譯節省開發者 30-60 分鐘的研究時間。

---

← [回到 engineering/app-development-guidance/](README.md)
