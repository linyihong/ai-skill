# Self-Generation Audits Slice（SDK live + Authorized identity material）

> **Cognitive Slice**：`apk-self-generation-audits`（從 [`../artifact-gates.md`](../artifact-gates.md) §8+§9 抽出的 focused slice，對應 [`governance/cognitive-slice-taxonomy.md`](../../../governance/cognitive-slice-taxonomy.md) §7.5）。

| slice 欄位 | 值 |
|---|---|
| `id` | `apk-self-generation-audits` |
| `purpose` | 強制 audit「除授權身分材料外，SDK 能否自行生成」與「身分材料能否自生成」；防止錯誤宣稱 standalone SDK 可達 |
| `type` | **`failure`**（type: failure 因其本質是 red-line / caveat 規則；不照表填會宣稱錯誤的 self-generation verdict） |
| `tags` | artifact-gate, security |
| `load_when` | 涉及 SDK live self-generation 或 authorized identity material（device id、install id、account、session seed、attestation）操作 |
| `do_not_load_when` | 任務無 self-generation / signing / identity 風險、純被動 observation |
| `owner_layer` | workflow |
| `layer_justification` | 規定「runtime factor 必填分類表、verdict 判讀規則」的 caveat / red-line gate；通過 workflow membership test，是失敗預防型 gate 而非新觀察方法（非 analysis）或長期 pattern（非 intelligence） |
| `canonical_source` | 本檔（原 `artifact-gates.md` §8 SDK Live Self-Generation Audit + §9 Authorized Identity Material Self-Generation Audit） |
| `dependencies` | `apk-domain-runtime-baseline`、`apk-sanitization`（evidence 不得保留 raw token / device id） |
| `dependency_budget` | default `max_depth:2` / `max_runtime_dependencies:4` |
| `validation_signal` | Scenario AG-C（security focus 任務應 PRIMARY 載入本 slice） |

## 8. SDK Live Self-Generation Audit

當使用者的目標是「像某些既有 SDK 一樣，除了呼叫方合法提供的 **授權身分材料** 外，其餘 host、路由、簽章、session、decrypt 都能由 SDK 自行生成」時，必須在專案 baseline 或 SDK-readiness 文件加一張 **runtime factor classification** 表：

| Classification | 意義 | 可開始 live SDK self-generation? |
| --- | --- | --- |
| `sdk-generatable` | 可由公開 SDK 程式、常數、演算法、穩定 public config 或已去敏規則自行生成；不需要私有 runtime bridge。 | 是 |
| `identity-material-bound` | 需要授權方提供或初始化的身分材料，例如 device id、install id、授權帳號、session seed、合法裝置初始化結果。 | 可，若這是唯一剩餘未知或唯一使用者提供項 |
| `private-adapter-required` | 需要 raw service、私有 host 選擇、簽章 key、decrypt key、in-app bridge、未公開 provider，或只能靠 app runtime 生成。 | 否 |
| `unknown` | 還不知道來源、時效、錯誤行為或是否可重建。 | 否 |
| `scoped-out` | 不屬於本 SDK live scope（例如 media download、write actions）。 | 不阻塞該 scope，但必須明寫 |

建議至少列：

| Runtime factor | 必問問題 |
| --- | --- |
| Base endpoint / host | SDK 能否從固定 fallback、public config、DNS/config API 自行選擇？還是必須 private host table / app storage？ |
| Route binding / service | raw route id/service 是否可由 SDK deterministic 生成？ |
| Authorized identity material | 是否只剩 device id / install id / 授權帳號 / session seed 等呼叫方身分材料需要注入？每個 key group 能否由 SDK/tool 自行生成或初始化？ |
| Session/bootstrap | guest/device login 是否可由 SDK 生成？是否仍需要 app-only token、captcha、human login 或私有 WebView state？ |
| Opaque query/header | 每個 opaque 欄位是 app 常數、locale、device/session 派生，還是 response/session 私有值？ |
| Signing/gateway | canonicalization、排序、hash/HMAC/AES、timestamp/random 來源是否可重建且已 fixture 驗證？ |
| Response decrypt/unwrap | SDK 是否能自行把 wire response 解成 JSON？key/IV/KDF 是否還在 private app helper？ |
| Error/session recovery | token 過期、bad signature、bad device、bad opaque 的 code 與 refresh/writeback 是否已 live matrix 驗證？ |
| Pagination/data truth | 是否已知如何終止分頁、辨識空資料、避免把錯誤 envelope 當空列表？ |
| Media/download | 若 scope 包含媒體，signed URL、key unwrap、decrypt、package 是否可重建；不含媒體時標 scoped-out。 |

完成後給出 verdict：

```text
Live SDK self-generation verdict:
- ready except authorized identity material: yes/no
- remaining non-device blockers: <factor list>
- allowed next work: live SDK implementation / private adapter only / offline parser only
```

**判讀規則：**只要仍有 `private-adapter-required` 或 `unknown` 的 base host、route service、signing、decrypt、session bootstrap、opaque provider，就不能說「只剩授權身分材料」。

## 9. Authorized Identity Material Self-Generation Audit

當 runtime 需要 device / install / account / session seed / vendor attestation / server-issued session 類材料時，必須逐 key group 回答「能否自生成」與「怎麼生成」。

最低表格：

| Field | 必填內容 |
| --- | --- |
| Key group / surface | 欄位名稱群、storage key 名、request key 名或 provider function boundary；只寫 name/shape，不寫 raw value。 |
| Role in live access | 它是 app/build constant、device/install material、account material、guest/session seed、vendor attestation、server-issued session，還是其它 runtime factor。 |
| Self-generation verdict | `sdk-generatable` / `caller-provided` / `server-issued` / `trusted-bridge` / `private-adapter-required` / `unknown` / `scoped-out`。 |
| Generation recipe or provider boundary | 若 `sdk-generatable`，寫 sanitized recipe：inputs、algorithm family、canonical order、storage key name、refresh trigger、validation fixture；若不能，寫由 caller、server response、trusted bridge 或 private adapter 提供。 |
| Lifecycle and reset behavior | first install、cold start、guest/login、preserved session、logout、`clear app data`、reinstall、token expiry 時如何建立、重用、更新或清除。 |
| Cooldown / risk controls | 是否會觸發 rate limit、device health、attestation check、account lock、captcha/human step；只寫 status/error class，不寫可濫用細節。 |
| Error / negative matrix | missing、empty、stale、bad-fixed、bad-signature、expired-session 等情況的 wrapper/UI/recovery class；若未驗證，標 `pending`。 |
| Validation evidence | 去敏 hook summary、static provider trace、fixture、replay parity、unit/contract test；不得保存 raw token、device id、account、vendor payload、signature 或 host。 |

判斷規則：

- `sdk-generatable` 需要可重跑的生成 recipe 或測試，不能只因為值看起來像 UUID、hash、locale、random 就宣稱可生成。
- `caller-provided` 仍需定義 lifecycle、reset/cooldown、health/error 行為；否則是 `unknown` 或 `private-adapter-required`。
- `server-issued` 可以由 SDK 建模 storage/refresh boundary，但 raw material 來自授權 server response。
- `trusted-bridge` / `private-adapter-required` 可以讓 private live smoke 成立，但不能支撐「standalone self-generating SDK」宣稱。
- 若任何 live-required identity key group 的 generation recipe、provider boundary、reset/cooldown 或 error matrix 是 `unknown`，live-facing development 只能繼續在 private adapter / bridge scope。

**Finish gate：** 若本輪目標包含「可程式化拉取真實資料」「接 SDK transport」「寫 integration test」之一，而專案尚無 baseline 或僅有 API 條目：必須在同一工作單建立 **skeleton baseline**，並把 **open** 項寫成可驗證問題。若本輪要開始開發 SDK、client、app tool、live integration，baseline 不能只停在 skeleton。若 live flow 需要 device/install/account/session/vendor/server-issued material，必須先補 authorized identity material self-generation audit。

---

← [回到 artifact-gates 索引](../artifact-gates.md) | [workflow/apk-analysis/](../README.md)
