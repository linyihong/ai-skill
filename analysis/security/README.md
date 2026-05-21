# Security Analysis

`analysis/security/` 放**安全相關的可重用分析方法**：如何觀察、拆解、取得證據以判斷一個系統的特定安全屬性。不放具體 incident、不放修補步驟（修補屬 `workflow/`），不放架構決策（屬 `intelligence/engineering/`）。

## 目前入口

| 文件 | 用途 |
| --- | --- |
| [`dual-token-audit.md`](dual-token-audit.md) | 系統同時使用兩套 token 簽章/加密機制（如 JWT + JWE、HMAC + 對稱簽章、平台 token + 廠商 token）時的審計方法 |

## 放什麼

- 安全屬性的觀察方法（如何看出此系統的 token 流向、加密邊界、權限模型）。
- Audit checklist、judgement heuristics、failure signals。
- Cross-cutting 的證據蒐集路線（log、code、config、network 各層怎麼配對）。

## 不放什麼

- 具體 CVE 與修補步驟（屬 `workflow/security/` 或專案 owner）。
- 安全架構決策（屬 `intelligence/engineering/architecture/system-boundaries/`）。
- Raw findings、real tokens、specific hosts、payload；違反 `enforcement/sanitization.md`。
- AI agent 自己的安全規則（屬 `enforcement/`）。

## 與其他層的關係

- 安全 architecture decision → [`intelligence/engineering/architecture/system-boundaries/`](../../intelligence/engineering/architecture/system-boundaries/README.md)
- 安全 anti-patterns → [`intelligence/engineering/anti-patterns/`](../../intelligence/engineering/anti-patterns/README.md)
- 修補執行流程 → `workflow/`（尚未建立）

← [回到 analysis/](../README.md)
