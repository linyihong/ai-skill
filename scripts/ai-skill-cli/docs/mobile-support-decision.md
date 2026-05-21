# Mobile Support Decision

## Decision

`ai-skill` 支援邊界如下：

| Platform | Decision | Reason |
| --- | --- | --- |
| Windows / macOS / Linux | Supported as native single binary | Desktop OS 可執行 repo-local / release artifact binaries |
| iOS | No arbitrary native binary support | iOS 不允許一般用途 executable persistence |
| Android | Not a first-class native target yet | Termux 可行但路徑、Git、權限與背景任務仍需專門驗證 |
| Mobile control plane | Supported direction | 手機觸發桌面 / server runner 最符合安全與平台限制 |

## iOS

iOS 不列入「下載 binary 後直接執行」目標。不可承諾：

```bash
./bin/ai-skill runtime validate --repo .
```

可行方向：

| Option | Decision | Notes |
| --- | --- | --- |
| App-contained runtime | Deferred | 需要 App 開發、檔案權限、Git/SQLite/runtime bundling、App Store policy 評估 |
| Browser/WASM | Research candidate | 適合 inspect UI、state validation、replay UI；local repo access 是主要限制 |
| SSH / remote runner | Preferred near-term route | iPhone 作 control plane，runtime 在 macOS/Linux/Windows runner |
| Arbitrary native binary | Unsupported | 不符合 iOS security model |

## Android

Android 不與 desktop single-binary target 混在一起承諾。

可行方向：

- Termux runner：可能可行，但需驗證 Git、SQLite、PATH、repo filesystem、background execution。
- App-contained runtime：可行但需 Android app scope。
- Remote runner client：保守且接近 iOS control-plane 模式。

## Required Follow-Up If Mobile Is Pursued

若要支援 mobile control client，另開安全與授權計畫，至少涵蓋：

- Runner authentication.
- SSH key / token storage.
- Audit log.
- Remote command allowlist.
- Lost-device revocation.
- Network failure and retry semantics.

## CLI Behavior

Current CLI should not claim iOS native binary support. Documentation and support matrix must point users to App-contained, Browser/WASM, or remote runner options.
