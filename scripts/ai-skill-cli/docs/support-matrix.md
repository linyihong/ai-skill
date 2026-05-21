# 支援矩陣：Ai-skill CLI Runtime

> **上游計畫**：[`2026-05-21-0834-cross-platform-go-script-runtime.md`](../../../plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md)

## 平台支援

| 平台 | 桌面單一 binary | Git 依賴 | Runtime DB | 行動平台定位 |
| --- | --- | --- | --- | --- |
| Windows | 支援目標 | 需要外部 Git | pure Go SQLite 目標 | 不適用 |
| macOS | 支援目標 | 需要外部 Git | pure Go SQLite 目標 | 不適用 |
| Linux | 支援目標 | 需要外部 Git | pure Go SQLite 目標 | 不適用 |
| iOS | 不支援任意 native binary | 只能透過 app 或遠端執行 | 只評估 Browser/WASM 或 app-contained | control plane / inspect UI / remote trigger |
| Android | 可行性評估目標 | Termux / app / remote | Termux 或 app-contained | app sandbox / remote runner |

## 桌面基準

桌面支援目標：

- Windows
- macOS
- Linux

基準假設：

- `ai-skill` 是 runtime toolchain 的單一 Go binary。
- YAML、JSON、SQLite engine、runtime logic、scheduler、migration / repair logic 在可行時應編入 binary。
- Git 不包進 binary。writeback、commit、push、hooks、close-loop 命令仍需要外部 Git。
- Release builder 產生 Windows amd64、macOS amd64/arm64、Linux amd64/arm64 artifacts 與 `SHA256SUMS`；核心 runtime commands 預設不要求 Ruby、Python 或外部 `sqlite3` CLI。

## iOS 邊界

iOS 不是任意 native binary 目標。

可評估路線：

| 路線 | 定位 | 備註 |
| --- | --- | --- |
| App-contained runtime | 可行 | iOS app 必須內建 runtime、Git / file access、SQLite 與 UI |
| Browser/WASM | 可行 | 較適合 inspect UI、replay UI、state validation 與 governance control plane |
| SSH remote runner | 高可行性 | iPhone 作為 control plane；runtime 在桌面、VPS、NAS、Mac mini 或 Linux host 執行 |
| 任意 native binary | 不支援 | iOS security model 不允許一般用途 executable persistence |

## Android 邊界

Android 可行性必須分開評估：

- Termux 可能比 iOS 支援更多本機執行能力。
- App sandbox 的限制類型與 iOS 類似，但細節不同。
- Remote runner 仍是保守選項。

## 不支援 / 阻擋條件

| 條件 | 必要行為 |
| --- | --- |
| 桌面 close-loop 缺 Git | 阻斷並提示安裝 |
| 要求 iOS native binary | 拒絕並提示 App、Browser/WASM 或 SSH remote runner |
| 缺寫入權限 | 阻斷寫入命令 |
| merge / rebase / cherry-pick 狀態 | 阻斷 commit / push |
| 不支援的平台 | 回傳穩定的 `unsupported_platform` exit code |

## Mirror 策略

`sync-cursor-bundle` 的跨平台預設策略是 copy fallback：

- Windows、受限權限或未明確允許 symlink 的環境，一律規劃 copy / replace managed mirror content。
- symlink 只作為未來明確 opt-in；啟用前必須有權限檢查與 managed / unmanaged target fixture。
- dry-run 必須以明確 `--target` 指向 fake Cursor root 或使用者指定 root，不得預設寫入真實 `$HOME`。
