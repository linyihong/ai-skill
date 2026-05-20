> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-20 - Frida Server Closed Preflight Dry Run

Status: candidate

#### One-line Summary

Android Frida attach 顯示 `unable to connect to remote frida-server: closed` 時，先當成 device 端 server lifecycle 問題處理，不要誤判成 App hook crash。

#### Human Explanation

APK 動態分析中，Frida attach 失敗和 App hook 崩潰是兩種不同問題。若 log 只有 Frida banner 後接 `Failed to attach: unable to connect to remote frida-server: closed`，且 target App PID 仍存在，通常代表 device 端 `frida-server` 沒有持續運作，或 attach 前 server 已被關閉。這時繼續調 hook、縮短 window、或改 UI 流程都不會解決根因。

更穩定的做法是先重啟 `frida-server`，用最小 hook / no-JS dry-run 驗證 attach channel，再進 feature capture。dry-run 應短窗、無 UI destructive action、只確認 hook 初始化訊號與 target App PID 保持存活。

#### Trigger

- Frida log 出現 `unable to connect to remote frida-server: closed`。
- capture log 沒有任何自訂 `[INIT]` 或 hook 訊號。
- App PID / foreground package 仍正常，沒有崩潰跡象。
- 同一 hook 腳本先前可用，但本輪 attach 一開始就失敗。

#### Evidence

- Tool: Android `frida-server`, `frida-ps`, minimal Frida attach, target PID check.
- Sanitized excerpt: attach 失敗只到 Frida transport 層；重啟 device 端 server 後，no-JS minimal attach 可輸出 hook `[INIT]` 訊號，target process 持續存活。
- Evidence path: project raw logs 留在 gitignored capture output；本 lesson 只保留 generalized failure signature。

#### Generalized Lesson

把 attach transport failure 和 hook-induced crash 分開：

1. 先檢查 target App PID 是否仍存在。
2. 檢查 device 上是否有 `frida-server` 行程。
3. 重啟 `frida-server` 後先跑 15-20 秒 minimal dry-run。
4. 只有 dry-run 穩定後才進入 feature capture。
5. 若 dry-run 成功但 feature capture 崩潰，才回頭分析 hook scope / UI timing / JS injection。

#### Agent Action

看到 `frida-server: closed` 時：

- 不要直接修改 hook 邏輯。
- 不要把 failure 歸因為 App 崩潰。
- 先重啟 server，執行 no-JS / minimal attach dry-run。
- 在報告中分開寫：server attach 狀態、hook init 狀態、App PID 存活狀態。

#### Goal / Action / Validation

- Goal: 避免把 Frida transport failure 誤診為 App hook crash。
- Action: 以 server lifecycle preflight gate 擋在 feature capture 前。
- Validation or reference source: dry-run log 有 hook `[INIT]` 訊號，target PID 保持存活，device 上 `frida-server` 行程仍存在。

#### Applies When

- Android root / emulator 以 external `frida-server` attach target App。
- Hook 腳本需要多次短窗重跑。
- 分析目標是穩定 capture，而不是測 server 部署本身。

#### Does Not Apply When

- Frida 已成功 attach 並輸出 hook 訊號後，target App 才崩潰。
- Crash log 明確指向 injected JS、JNI callback、native access violation 或 script exception。
- 使用 Frida Gadget 而非 device `frida-server`。

#### Validation

- failure signature 文件化為 server lifecycle issue。
- 下一次 capture 前先有 minimal dry-run。
- feature capture 的失敗分類不混用 `attach failure` 與 `hook crash`。

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` 的 Frida deployment / capture automation 段落。
- `analysis/apk/tools-and-failures.md` 或對應 Frida troubleshooting 文件。

#### Required Linked Updates

- Project capture checklist 可引用此 pattern，提醒 `frida-server closed` 先做 server preflight。
- 不需同步 project raw evidence；具體 PID、device serial、App 名稱與本機路徑留在 project docs / logs。
