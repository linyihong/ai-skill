> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[sanitization](../../../../enforcement/sanitization.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-20 - Runtime H5 External Replay Gate

Status: candidate

#### One-line Summary

嵌入式 H5 API 即使在 WebView runtime 內成功，也應先用 opt-in external replay gate 驗證 gateway / session 邊界，再決定是否落成 SDK public client。

#### Human Explanation

Hybrid App 的 H5 WebView 可能在 runtime 內持有額外 gateway token、origin / referer、UA shape、host selection 或 TLS / CORS context。分析者若只看到 WebView 內 response 成功，就直接把 endpoint、codec 和 schema 做成 SDK client，容易把「runtime-bound」能力誤判成「external SDK 可重放」能力。

較安全的做法是建立一個 disabled-by-default 的外部 replay gate：只有在本機或 CI secret 明確提供授權材料、H5 request codec material、identity material 與 endpoint override 時才送出 request。測試只輸出去敏 summary，例如 host/path、HTTP status、business code、是否有 encrypted data；不保存 token、uid、完整 URL、ciphertext、raw response、key / IV 或特定使用者資料。

#### Trigger

- H5 response schema 已在 WebView runtime、MITM、DevTools 或 hook 中確認。
- 外部 replay 與 runtime replay 結果不一致，或 business code 顯示 gateway / session 邊界仍未通過。
- 團隊想把 H5 endpoint 包成 SDK public API，但尚未證明 SDK 外部 client 能穩定通過。
- H5 request 需要 bundle-derived codec、runtime token、origin / referer 或 special user-agent。

#### Evidence

- Tooling pattern: JUnit / integration-test style opt-in probe, environment-only secrets, sanitized JSON summary.
- Sanitized behavior: external replay 可達 HTTP gateway，但 business result 可能仍是 boundary failure；這時只能記為 readiness blocker，不可宣稱 SDK client ready。
- Evidence path: raw request / response 與 secrets 留在 project-local ignored evidence 或 secret store；reusable lesson 只保存 pattern。

#### Generalized Lesson

Runtime H5 API 的 SDK readiness 應分成三層：

1. **Runtime schema evidence**：WebView runtime 內成功並解出 payload shape。
2. **External replay gate**：SDK 外部 client 使用可重建的 request flow 得到 success business code。
3. **Public client contract**：只有 gate 通過後，才補 BDD / executable feature / production client。

如果第 2 層仍停在 business boundary failure，完成狀態應是「H5 URL / schema documented」，不是「SDK list client implemented」。

#### Agent Action

下次分析嵌入式 H5 API 並準備轉成 SDK 時：

- 不要把 WebView runtime success 直接等同 external SDK readiness。
- 先加 disabled-by-default external replay gate。
- Gate input 全部走 env / secret store，不落 tracked fixture。
- Gate output 只寫 sanitized summary，不寫 raw URL、token、uid、ciphertext、key / IV。
- Public client 前必須有 gate pass + docs-first BDD / executable feature。

#### Goal / Action / Validation

- Goal: 防止 runtime-only H5 API 被過早包成 SDK public client。
- Action: 加 opt-in external replay gate，將 business success 與 boundary failure 都可稽核化。
- Validation or reference source: gate 預設 skip；有授權 env 時輸出去敏 summary；只有 gateway pass 才允許進 production client planning。

#### Applies When

- Android WebView / Flutter InAppWebView / hybrid App H5。
- H5 endpoint 有 request / response codec、runtime token、origin / referer 或 gateway host selection。
- 任務要把 APK / H5 analysis 轉成 SDK 或 app-development guidance。

#### Does Not Apply When

- 目標只是 UI / schema 文件化，不打算做 external client。
- Endpoint 是普通 App API，已由既有 SDK session HTTP client 穩定重放。
- 授權邊界不允許外部 replay，或 secrets 無法安全注入測試。

#### Validation

- Gate 預設不送網路 request。
- Tracked docs 不含 endpoint secret、raw token、raw user identity、完整 URL、ciphertext、key / IV。
- Project contract 明確寫出「runtime schema evidence」與「external replay pass」不是同一件事。

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` 的 H5 → SDK readiness gate。
- `workflow/software-delivery/execution-flow.md` 的 Change Intake / docs-first gate：runtime evidence 不等於 public client contract。

#### Required Linked Updates

- Project development contract 應連到 opt-in replay gate。
- Integration matrix 應標註 gate 是 opt-in 且需要 secret / local evidence。
- 若 promotion 到 workflow，需同步 sanitization guidance 與 SDK docs-first closure。
