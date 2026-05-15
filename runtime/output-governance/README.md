# Output Governance

## 目的

將語言偏好、文件輸出規則從分散的 prose 檔案升級為 **declarative YAML**，讓 runtime 可以直接檢查、強制執行、並在 validation phase 自動驗證輸出品質。

## 設計原則

1. **Canonical Source 仍在 prose**：`enforcement/neutral-language.md`、`enforcement/sanitization.md`、`enforcement/tool-neutral-documentation.md` 仍為原始規則定義。本目錄的 YAML 為 compiled version，由 compiler 同步更新。
2. **Phase-aware 檢查**：Output governance gates 掛在 `validation` 與 `finalize` phase，確保每輪輸出都經過語言/格式/去敏檢查。
3. **工具中立**：語言政策定義核心規則，各工具的具體設定方式（Roo Code SQLite、Cursor `.cursor/rules/`、Claude `CLAUDE.md`）留在 `ai-tools/agent/*.md`。
4. **Compiler 整合**：compiler 在編譯 generated YAML 時同時檢查 output rules。

## 檔案結構

```
runtime/output-governance/
├── README.md                          # 設計原則與使用說明（本檔）
├── language-policy.yaml               # 語言強制規則
├── output-rules.yaml                  # 文件輸出規則
└── governance-gates.yaml              # Output governance blocking gates
```

## 與既有層的關係

| 元件 | 關係 |
|------|------|
| `runtime/phases/phase-machine.yaml` | Governance gates 掛在 `validation` 與 `finalize` phase |
| `runtime/gates/blocking-gates.yaml` | Governance gates 為 blocking gates 的子集 |
| `runtime/compiler/compiler-rules.yaml` | Compiler 在編譯時應檢查 output rules |
| `runtime/intelligence/intelligence-routing.yaml` | Intelligence 知識的輸出也受 governance 規範 |
| `enforcement/neutral-language.md` | 語言規則的 prose source |
| `enforcement/sanitization.md` | 去敏規則的 prose source |
| `enforcement/tool-neutral-documentation.md` | 工具中立性規則的 prose source |
| `ai-tools/agent/*.md` | 各工具的語言設定方式應 reference governance YAML |

## 使用方式

Agent 在 `validation` phase 應：

1. 讀取 `language-policy.yaml` 確認語言一致性
2. 讀取 `output-rules.yaml` 確認輸出格式與去敏
3. 通過 `governance-gates.yaml` 中定義的所有 gate
4. 若任一 gate 未通過，進入 `recovery` phase 修正

## 誰會參考這裡（Inbound References）

- `runtime/phases/phase-machine.yaml` — validation phase 的 blocking gates
- `runtime/compiler/compiler-rules.yaml` — compiler 的 output governance check
- `enforcement/goal-action-validation.md` — validation flow 參考
- `ai-tools/agent/roo.md` — 語言設定方式 reference
- `ai-tools/agent/claude.md` — 語言設定方式 reference
- `ai-tools/agent/cursor.md` — 語言設定方式 reference
