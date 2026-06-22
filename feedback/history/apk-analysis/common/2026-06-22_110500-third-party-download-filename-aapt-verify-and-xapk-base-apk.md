> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md)、[sanitization](../../../../enforcement/sanitization.md)、[reusable-guidance-boundary](../../../../enforcement/reusable-guidance-boundary.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-06-22 - Third-party download filename vs actual package; XAPK analysis target

Status: candidate

#### One-line Summary

第三方鏡像檔名常含目標 package 字串，但 APK 內容可能是**商店 client 本身**；靜態分析前必須 `aapt dump badging` 驗 package。XAPK 的分析對象是內層 **base APK**，不是 XAPK 容器或 density split。

#### Human Explanation

從第三方下載站取得的檔名（例如 `<mirror>-<expected-package>.apk`）容易被誤讀為「就是目標 App」。實務上可能下載到**下載器／商店 App**（badging 顯示的 package 與預期不符）。若未驗證就開始 API 分析，會浪費整輪 triage。另：XAPK 是 zip 容器（base + split + manifest.json）；業務邏輯、DEX、native、Frida、jadx 都在 **base APK**；`config.*dpi` 等 split 通常只有資源。

#### Trigger

- 檔名含預期 package 或 mirror 站名，但尚未跑 badging。
- 收到 `.xapk` 或 split bundle，不確定要解哪個檔。
- 靜態結果（MainActivity package 不符、stack 與 Play 類型描述不一致）與任務預期不一致。

#### Evidence

- Tool: `aapt dump badging`, `unzip -l`, checksum 比對 base vs 容器內檔、可選 `adb pull` base vs 本地 base。
- Sanitized excerpt: 檔名暗示 package A，badging 顯示 package B（第三方商店 client）；XAPK 內 base 與另存 base 檔 checksum 一致；已安裝 pull 與本地 base checksum 一致。
- Evidence path: 具體 incident（檔名、badging 輸出、checksum）留在 `<PROJECT_ROOT>/docs/` 或 `<PROJECT_ROOT>/capture/`；**本 lesson 不含** App 專名、package 真值、mirror URL 或本機路徑。

#### Generalized Lesson

1. **Identity gate（任何 APK 分析第一步）**：`aapt dump badging <file> | grep package` — package、version、launchable-activity 必須與**任務指定 package**一致。
2. **XAPK**：解包 → 讀 `manifest.json` 的 base entry → 對 **base `.apk`** 做 badging / jadx / strings / 安裝（split 僅在 `install-multiple` 時需要）。
3. **保留 artifact**：canonical base APK + 可選完整 XAPK；安裝後可用 `adb pull` base 與本地 checksum 交叉驗證。
4. **檔名不可當 identity 證據**；badging 才是 gate。

#### Agent Action

1. 收到 mirror 下載檔 → 立即 badging；package 不符則**停止分析**並重新取得 artifact。
2. XAPK → 解包 base；所有靜態/動態分析指向 base；split 只列在 install 清單。
3. 寫入 Ai-skill lesson 前依 [sanitization.md](../../../../enforcement/sanitization.md) 檢查：不得含真實 package 對照表、mirror 連結、checksum 原文、`<PROJECT_ROOT>` 下具體子目錄名。

#### Goal / Action / Validation

- Goal: 避免對錯誤 artifact 做 API/stack triage。
- Action: 在 traffic triage 之前加 identity gate（badging + 可選 checksum vs 已安裝 pull）。
- Validation or reference source: `workflow/apk-analysis/execution-flow.md` §1；`analysis/apk/traffic-triage.md` 主線選擇前 package 已確認。

#### Applies When

- 第三方 mirror / split bundle / 非 Play 直連來源。
- 多 App 分析 workspace 中檔名與 Play listing 混用。

#### Does Not Apply When

- Artifact 來自已驗 package 的上游 pipeline（例如已 badging 的 CI 產物）。
- 純 `adb pull` 且 package 已由 `pm list packages` 確認。

#### Validation

- badging package 與任務指定 package 一致才進 stack triage。
- 報告明確寫「analysis target = base APK」，不是 XAPK 外殼。

#### Promotion Target

- `workflow/apk-analysis/execution-flow.md` §2 Quick Start
- `runtime/onboarding/apk-analysis-setup.md`

#### Required Linked Updates

- 依 [`linked-updates.md`](../../../../enforcement/linked-updates.md)：`feedback/history/apk-analysis/README.md` 索引待追加。
- 已依 [`reusable-guidance-boundary.md`](../../../../enforcement/reusable-guidance-boundary.md) 與 [`sanitization.md`](../../../../enforcement/sanitization.md) 檢查：正文無 App 專名、真實 host、token、本機路徑、device id。
