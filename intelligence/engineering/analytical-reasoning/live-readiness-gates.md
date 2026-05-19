# Live SDK/Client Readiness Gates（實作就緒門檻）

**Status**: `validated-intelligence`
**Source**: `workflow/apk-analysis/artifact-gates.md`（development readiness gate、SDK live self-generation audit）

## 問題

APK analysis 經常產出詳細的 API schemas，但遺漏了 live SDK/client 實作所需的 runtime factors。團隊開始對正式服務寫 code 之後，才發現缺少 opaque parameters、未記載的 session dependencies 或 private signing paths。

## 原則

**僅有 API shape 不足以開始 live SDK/client 工作。** 在開始 live-facing 實作之前，專案必須具備 Domain/Runtime Baseline，回答以下問題：

| Factor | 必要項目 |
| --- | --- |
| Endpoint/path family | 已知 host family、path prefix、multi-CDN/gateway |
| Route/service binding | Raw route id 或 hash → deterministic 或 private-adapter |
| Session/bootstrap dependency | Guest/device login 可產生？Session seed 來源？ |
| Opaque parameter source & lifetime | 每個 opaque field：app constant、device-derived 或 response-derived？ |
| Signing/gateway prerequisites | Canonicalization、sorting、hash/HMAC algorithm 已 fixture-verified？ |
| Response decrypt/unwrap boundary | Wire → JSON：key/IV/KDF 已知且 fixture-verified？ |
| Pagination truth | `has_next` flag 或 heuristic？已知 anti-pattern 風險？ |
| Error/session recovery | Token expiry、bad signature、bad device codes 已 live-verified？ |
| Replay checklist | 重現 list API call 的最小步驟？ |

## Identity Material Self-Generation Audit

當 live access 需要 device/install/account/session/vendor material 時，將每個 key group 分類：

| 分類 | 意義 | 可以開始 live SDK？ |
| --- | --- | --- |
| `sdk-generatable` | Public SDK code、constants 或 algorithms 可產生。 | 是 |
| `identity-material-bound` | 需要授權的身分材料（device id、account、session seed）。 | 是，如果這是唯一剩餘的未知項 |
| `private-adapter-required` | 需要 raw service、private key、in-app bridge 或未公開的 provider。 | 否 |
| `unknown` | 來源、lifetime 或 error behavior 尚未知道。 | 否 |
| `scoped-out` | 不在 scope 內（例如 media download）。 | 不阻擋，但必須明確標示 |

## Verdict 格式

```text
Live SDK self-generation verdict:
- ready except authorized identity material: yes/no
- remaining non-device blockers: <factor list>
- allowed next work: live SDK implementation / private adapter only / offline parser only
```

## Token 影響

避免浪費實作週期。單一遺漏的 runtime factor 可能阻擋數週的 SDK 工作。此 gate 確保在寫 code 之前就先問對問題。

---

← [回到 intelligence/engineering/analytical-reasoning/](README.md)
