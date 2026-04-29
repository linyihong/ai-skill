# APK 分析文件寫法

分析文件的目標是讓人和 AI 都能重現推理，而不是只留下最後答案。

## 文件分層

建議分成三層：

| 層級 | 內容 | 是否放入 skill |
| --- | --- | --- |
| 方法論 | 流量路徑判斷、工具選擇、hook 策略、去敏規則。 | 可以，放在本資料夾。 |
| 專案結論 | 目標 App 的 API、host、endpoint、schema、媒體規則。 | 不放 skill，放專案 API / reference docs。 |
| 原始證據 | pcap、MITM export、Frida log、raw response、decrypted fixture。 | 不放 skill；放 gitignored 或專案指定位置，文件只引用去敏摘要。 |

## 單次分析筆記模板

```markdown
# [APK / 功能] 分析紀錄

## Scope

- APK:
- Version:
- Package:
- Device / emulator:
- Authorization:
- Goal:

## Environment

- OS:
- adb:
- Frida:
- Proxy tool:
- Static tools:

## Hypotheses

| Hypothesis | Test | Result |
| --- | --- | --- |
| localhost bridge | lo pcap | |
| system proxy / MITM | proxy capture | |
| Java HTTP stack | Java hook | |
| Flutter / native | connect backtrace / AOT strings | |

## Evidence

| Evidence | Path / excerpt | Interpretation |
| --- | --- | --- |
| pcap | `<path>` | |
| hook log | `<path>` | |
| static search | `<path or command>` | |

## Findings

- Finding 1.
- Finding 2.

## Unknowns

- Unknown 1.

## Next Steps

1. Next validation.
2. Next fixture or test.

## Sanitization

- Tokens redacted:
- Device identifiers redacted:
- User data removed:
```

## API / Schema 文件模板

```markdown
## Endpoint Name

| Field | Value |
| --- | --- |
| Method | `GET` / `POST` |
| Path | `/path` |
| Auth | Required / Optional |
| Source | pcap / MITM / hook / replay |

### Request

| Parameter / Header | Meaning | Required | Notes |
| --- | --- | --- | --- |

### Response Wrapper

| Field | Type | Notes |
| --- | --- | --- |

### Decrypted / Inner Payload

| Field | Type | Notes |
| --- | --- | --- |

### Evidence

- Sanitized log:
- Fixture:
- UI path:

### Validation

- Replay:
- Contract test:
- Manual verification:
```

## 去敏規則

必須遮蔽：

- `Authorization`、cookie、session token。
- device id、install id、advertising id。
- 真實帳號、電話、email、邀請碼。
- AES/HMAC key material。
- 能直接重放付費內容或個人內容的 URL。
- 本機絕對路徑、使用者名稱、私有工作目錄、clone 位置。請改用 `<AI_SKILL_REPO>`、`<PROJECT_ROOT>`、`<WORKSPACE>` 等 placeholder。

可以保留：

- header 名稱。
- path shape。
- query key 名稱。
- response top-level schema。
- 已去敏的 fixture。
- magic bytes、容器格式、演算法步驟。

## 證據鏈要求

好文件不只寫「成功」，還要寫為什麼相信它成功：

- pcap 證明對外 TLS host 存在。
- proxy CONNECT 證明導流成功。
- hook log 證明 request object 在 TLS 前可見。
- decrypt hook 或離線 decoder 證明 inner JSON 正確。
- fixture / test 證明規則可重跑。

## 失敗也要記錄

失敗紀錄應包含：

- 嘗試了什麼。
- 期望看到什麼。
- 實際看到什麼。
- 排除了什麼假設。
- 是否要重試，或是否停止投入。

例：

```text
Java OkHttp hook installed successfully, but no target host/path appeared while pcap showed TLS traffic to the API host. This rules out the Java OkHttp path for the tested flow and shifts the next step to native/Flutter analysis.
```

## 技巧回饋文件要給人讀

寫入 `FEEDBACK.md` 時，請避免只有工具名與短結論。每條技巧都應包含：

- `One-line Summary`：一句話講重點。
- `Human Explanation`：給人看的背景與誤判風險。
- `Trigger`：什麼現象會觸發這條技巧。
- `Evidence`：去敏證據或觀察。
- `Generalized Lesson`：抽象後的通用規則。
- `Agent Action`：下次 AI 要採取的具體行動。
- `Applies When` / `Does Not Apply When`：適用邊界。
- `Validation`：怎麼確認這條技巧有效。

好的 lesson 應該像這樣：

```markdown
### Proxy failure 要先拆成導流與 TLS 兩層

One-line Summary:
代理看不到明文時，先確認「有沒有進代理」，再談憑證或 pinning。

Human Explanation:
很多人看到 MITM 沒有明文就直接判斷是 pinning。更可靠的順序是先看 App 是否真的連到 proxy。如果仍直連目標 host，問題在導流或初始化時機；如果已經進 proxy 才 TLS failed，才查 CA / pinning。

Agent Action:
先檢查 CONNECT 或 connect target，不要先寫 pinning 結論。
```

## 回填規則

每次分析完成後：

- 目標 API 結論回填專案 API 文件。
- 解碼規則回填協議/解密文件。
- SDK 或 client 行為回填 BDD / tests。
- 通用技巧回填 `FEEDBACK.md`，驗證後再整理到本 skill 的主文件。
