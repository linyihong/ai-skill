### 2026-04-30 - Frida 只有 banner 時先做最小 hook 健康檢查

Status: validated

#### One-line Summary

Frida 沒輸出不一定是 hook 點錯，可能是 client、權限、sandbox、attach 時機或 App 沒觸發流程。

#### Human Explanation

實務上常看到 Frida 啟動後只有 banner，沒有任何自訂 log。這時不要立刻重寫完整 hook。先拆問題：Frida client 能不能列 process？最小 `console.log` hook 能不能載入？目標 PID 對不對？App 是否真的觸發該流程？執行環境是否限制了 Frida？這樣可以避免把工具連線問題誤判成分析結論。

#### Trigger

- Frida log 只有 banner。
- 完整 hook 沒有 `[INIT]` 或自訂 log。
- App 操作時 pcap 有流量，但 hook 沒任何事件。

#### Evidence

授權分析中曾遇過完整 hook 無輸出，但最小 attach 可輸出；調整執行環境或 attach/spawn 方式後，完整 hook 才正常。

#### Generalized Lesson

先用最小 hook 驗證 Frida client/server/目標 PID/權限，再檢查 hook offset 或業務流程。

#### Agent Action

遇到 Frida 無輸出時，先建議最小測試：列 process、attach PID、輸出 `HOOK_LOADED`、再逐步加 hook。不要直接擴大 hook 範圍。

#### Promotion Target

- `TOOLS.md`
- `WORKFLOW.md`
