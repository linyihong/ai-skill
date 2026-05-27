# Inflated Cognitive Mode Reporting

Status: validated
Class: `governance-drift` / `validation-gap`

## Trigger

當 agent 在 Cognitive Contract 中宣告高於實際任務風險的 mode、cost 或 activation reason，使用此 pattern。

常見訊號：

- 簡單文件或 typo 變更宣告 `DEEP`、`FORENSIC`、`RECOVERY` 或 `LOCKDOWN`。
- `cognitive_cost` 與 `execution_mode × context_mode` lookup 不一致。
- `activation_reason` 使用不在 `runtime/cognitive-modes-discovery.yaml` 的 free-form signal。
- 高風險 mode 沒有 `Capability summary`，reader 無法知道 mode label 對應哪些能力或 gate。

## Failure Mode

Cognitive Mode 報告從 runtime contract 退化成 agent self-description。Agent 可能為了看起來謹慎而膨脹 mode，或為了簡短輸出而低報 mode / cost，導致 reviewer 無法依報告判斷實際 validation depth。

## Risk

- **Label without contract**：`STRICT` / `DEEP` 只剩抽象標籤，沒有 capability semantics。
- **False confidence**：人類誤以為已啟用高風險 validation。
- **Cognitive fatigue**：所有任務都 full report 或高風險 mode，真正重要訊號被稀釋。
- **Audit drift**：commit history 無法反查 mode 為何被啟用。

## Required Agent Action

1. 依 `runtime/cognitive-modes-discovery.yaml` 選擇 `activation_reason`，不得自造 signal。
2. 讓 `cognitive_cost` 由 `runtime/cognitive-modes-cost-class.yaml` derivation 決定。
3. 只有全 6 維 default 時使用 compact form。
4. 高風險 mode 必須附 `Capability summary`。
5. 若 validator 擋下 cognitive block，調整 mode 或拆分 commit，不要用 opt-out 掩蓋真實 task。

## Prevention Gate

Commit-msg hook 必須檢查：

| Check | Required behavior |
| --- | --- |
| `validateCognitiveCost` | declared cost 必須等於 derived cost |
| `validateActivationSignals` | activation signal 必須存在於 discovery contract |
| `validateCapabilitySnippet` | high-risk mode 必須有 capability summary |
| compact form gate | 任一非 default 維度都不可使用 compact form |

## 驗證

符合下列條件時，此 pattern 已被防止：

- `phase6-cognitive-contract-v2-cost-class-v1` PASS。
- `phase6-cognitive-contract-v2-activation-signal-v1` PASS。
- `phase6-cognitive-contract-v2-capability-snippet-v1` PASS。
- `phase6-cognitive-contract-v2-inflated-rejection-v1` PASS。
- `go test ./...` 通過，且 commit-msg hook 會擋下 cost mismatch / unknown signal / missing capability snippet。

## Linked Rules

- [`runtime/cognitive-modes.yaml`](../../runtime/cognitive-modes.yaml)
- [`runtime/cognitive-modes-cost-class.yaml`](../../runtime/cognitive-modes-cost-class.yaml)
- [`runtime/cognitive-modes-discovery.yaml`](../../runtime/cognitive-modes-discovery.yaml)
- [`cognitive-mode-resolution-bypass.md`](cognitive-mode-resolution-bypass.md)
- [`failure-to-validator-closure.md`](failure-to-validator-closure.md)

← [Back to failure patterns](README.md)
