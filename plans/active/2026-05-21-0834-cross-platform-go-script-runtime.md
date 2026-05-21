# Cross-Platform Go Script Runtime

> **狀態**：in-progress
> **建立時間**：2026-05-21 08:34
> **目的**：將 `scripts/` 從依賴特定 shell、Ruby runtime 與本機環境假設，升級為可在 Windows、macOS、Linux 穩定執行的單一 Go binary 工具層；優先內建 pure Go SQLite（`modernc.org/sqlite` candidate）、YAML、JSON 與 runtime logic；desktop Git 是 external dependency，不包進 binary，缺 Git 時由 `doctor` / linked-update flow 阻斷並提示安裝；iOS 以 App sandbox、Browser/WASM 或 SSH remote runner 作為可行方向，明確不支援 native arbitrary binary，Android 另評估 Termux / app sandbox / remote runner，不預設納入桌面同等支援範圍。

## 背景

目前 `scripts/` 同時包含 shell、Ruby、Python 與 git hook script。這些工具已能支援本庫的 close-loop、runtime refresh、SQLite index、hook 安裝與新專案初始化，但執行成功常隱含以下前提：

- 使用者環境有 POSIX shell、Ruby、Python、SQLite CLI、git 與相容的檔案權限模型。
- 腳本以 macOS / Linux 路徑與 shell 行為為主，Windows 需要 Git Bash、WSL 或額外 runtime。
- hook 與 close-loop 行為依賴本機 git 設定、執行權限與環境變數。
- runtime compiler 與 generated surface 驗證依賴 Ruby 與 SQLite CLI，agent 在不同環境容易漏跑或跑不完整。

本計畫不把「改用 Go」視為單純語言替換，而是把 scripts 層重構為可測試、可發佈、可跨平台執行的工具產品。

## 目標

1. 建立一個跨平台 CLI 入口，例如 `ai-skill`，以 Go 實作主要 script 能力。
2. 在 Windows、macOS、Linux 上提供一致的命令、輸出、exit code、dry-run 與錯誤訊息。
3. 產出可發佈的單一 binary，將 YAML parser、JSON parser、SQLite engine、runtime logic、scheduler、migration / repair logic 盡量編入 executable，降低使用者端安裝成本。
4. 將現有 shell / Ruby script 的行為先規格化，再分批遷移，避免一次重寫造成 close-loop 失效。
5. 建立完整測試矩陣：unit test、golden output、fixture repo、跨 OS CI、端到端 dry-run、runtime.db / SQLite assertion。
6. 保留必要的 Ruby compiler 或逐步移植策略，確保 runtime source-of-truth 與 generated artifact 不漂移。
7. 評估 iOS / Android 是否適合作為「App 內建 runtime」、「Browser/WASM control plane」、「SSH / remote runner 控制端」或「不可支援平台」，並寫出明確結論。

## 非目標

- 不在第一階段移除所有現有 script；新 CLI 應先與既有 script 並行，通過驗證後再 deprecate。
- 不把手機平台承諾成完整本機執行環境；iOS / Android 的 sandbox、檔案系統、git、SQLite 與 process model 需先評估。
- 不改變 Ai-skill 的 source-of-truth 分層；CLI 只是執行介面，不能取代 `governance/`、`runtime/`、`knowledge/` 與 `plans/` 的 canonical 文件。
- 不把工具 mirror 當 source repo；reference-first 與 writeback transaction 規則維持不變。

## Architecture Compatibility Preflight

| 欄位 | 內容 |
| --- | --- |
| Trigger | 建立跨平台 script runtime plan，尚未進入 implementation phase |
| Checked sources | `plans/README.md`、`scripts/README.md`、`enforcement/dependency-reading.md`、`enforcement/linked-updates.md`、`enforcement/neutral-language.md`、`governance/lifecycle/knowledge-update-flow.md`、`runtime/compiler/compiler-engine.rb` |
| Current architecture | `scripts/` 是 automation / validation / close-loop 執行層；`runtime/compiler/` 是 runtime.db compiler source；`knowledge/runtime/` 是 generated runtime surfaces；`plans/active/*.md` 會被 compiler 納入 plans index |
| Conflicts | 無直接衝突；但若 implementation 改 `scripts/` 或 `runtime/compiler/`，必須同步 `scripts/README.md`、必要 validator、runtime.db 與 generated reports |
| Decision | proceed with planning only；implementation phase 前必須重新做 preflight，確認 candidate files、source-of-truth 與 test strategy 仍有效 |
| Validation | 本計畫建立後需更新 `plans/README.md`，執行 knowledge runtime refresh、runtime compiler、`validate-runtime-db`、link check、lint、commit / push / readback |

## Development Workflow Alignment

本計畫必須依 [`workflow/software-delivery/`](../../workflow/software-delivery/README.md) 的開發流程執行，不得直接從 Go implementation 開始。

| Software-delivery stage | 本計畫對應產物 | Blocking gate |
| --- | --- | --- |
| Product / impact alignment | Change brief：為什麼要降低 script runtime 的環境依賴、Windows friction、close-loop 漏跑風險 | 未說明使用者 / contributor / agent runtime 的影響前，不得開始 CLI 設計 |
| Requirements cognition / BDD-lite | Command behavior scenarios：`doctor`、`close-loop`、`runtime compile`、missing Git、missing permission、dirty tree | 未建立 Given / When / Then 與 acceptance criteria 前，不得開始 Go implementation |
| Domain architecture cognition | Command contract、舊腳本 parity inventory、dependency detector、Git adapter、runtime compiler adapter、filesystem adapter、reporting / exit-code boundary | 未定義 domain boundary、舊腳本 parity 與 side effects 前，不得建立 production package layout |
| Test effectiveness | Fixture plan、golden output、negative cases、missing dependency、Windows path、runtime.db assertion、legacy script parity fixture | 未定義 fixture、舊腳本覆蓋率與 failure tests 前，不得宣稱 Phase 1 skeleton 可完成 |
| Artifact gates | Command contract、script parity inventory、support matrix、exit-code table、side-effect registry、test fixture plan | 未完成文件 artifact gate 前，不得進入 Phase 1 Go CLI skeleton |

**文件先行 blocking gate**：

在新增 `scripts/ai-skill-cli/go.mod`、`scripts/ai-skill-cli/cmd/ai-skill/` 或任何 production Go implementation 前，必須先完成並 review：

- `scripts/ai-skill-cli/docs/README.md`
- command contract：每個命令的 arguments、輸出、exit code、side effects、dry-run 行為。
- script parity inventory：每個現有 shell / Ruby / Python / git hook / compiler 入口的功能、輸入、輸出、副作用、外部依賴、目標 CLI、parity 狀態與最低測試證據。
- support matrix：Windows、macOS、Linux、iOS、Android 的支援等級與限制。
- BDD-lite scenarios：尤其是 missing Git、missing permission、dirty tree、merge / rebase state、legacy script parity、runtime.db assertion。
- test fixture plan：temporary repo、fake home、PATH isolation、missing Git、Windows path、legacy script parity、runtime DB golden fixture。

若上述任一項未完成，本計畫不得進入 Go implementation。若必須例外，需先在本 plan 的 Open Questions 或 decision record 寫明原因與風險。

## Current Script Inventory

完整逐項 parity 盤點見 [`scripts/ai-skill-cli/docs/script-parity-inventory.md`](../../scripts/ai-skill-cli/docs/script-parity-inventory.md)。本節只保留高層分類與主要風險。

| 類型 | 目前檔案 | 主要風險 |
| --- | --- | --- |
| Shell CLI | `init-new-project.sh`、`sync-cursor-bundle.sh`、`ai-skill-close-loop.sh`、`agent-goals.sh`、`install-hooks.sh` | Windows shell 相容性、執行權限、path separator、symlink、環境變數、git hook 行為 |
| Ruby generator / validator | `validate-knowledge-runtime.rb`、`refresh-knowledge-runtime.rb`、`generate-*.rb`、`query-*.rb`、`validate-runtime-db.rb`、`runtime/compiler/compiler-engine.rb` | Ruby runtime、gem/native extension、encoding、SQLite CLI / library 差異 |
| Python helper | `set-roo-global-custom-instructions.py` | Python availability、path 與使用者設定檔位置 |
| Git hooks | `scripts/git-hooks/pre-commit`、`scripts/git-hooks/post-commit` | hook shell、PATH、環境變數、跨平台 git hook 限制 |

## Proposed Target Architecture

### 1. Go CLI 作為穩定入口

建立 `scripts/ai-skill-cli/cmd/ai-skill`：

```text
ai-skill
  init-project
  sync-cursor-bundle
  close-loop
  goals
  runtime refresh
  runtime validate
  runtime compile
  runtime query
  hooks install
  doctor
```

CLI 必須提供：

- `--dry-run`：所有會寫檔、commit、push、安裝 hook、建立 symlink 的命令都必須支援。
- `--json`：提供 agent / CI 可解析輸出。
- `--plain`：提供人類可讀輸出，避免只依賴彩色終端。
- 穩定 exit code：validation failure、dirty working tree、missing dependency、permission denied、unsupported platform 要有不同 code。
- path abstraction：統一處理 Windows drive、UNC path、POSIX path、symlink 與 path normalization。

### 2. 單一 binary 與內建 dependency 原則

跨平台目標應以「使用者不需要另外安裝 Ruby、Python、sqlite3 CLI、gem、pip 或 C compiler」為方向。Go CLI 應把主要 runtime dependencies 編進 binary：

| Dependency | 優先策略 | 原因 |
| --- | --- | --- |
| YAML parser | 使用 Go library，例如 `gopkg.in/yaml.v3` | 編譯後隨 binary 發佈，不要求使用者安裝外部 YAML 工具 |
| JSON parser | 使用 Go standard library | 無外部 runtime 依賴 |
| SQLite engine | 優先評估 pure Go SQLite，例如 `modernc.org/sqlite` | 避免 CGO、C compiler、SQLite dev headers 與 Windows build friction |
| Runtime / scheduler / migration logic | Go native implementation | 讓 `ai-skill runtime migrate`、`replay`、`repair` 類命令不依賴 shell / Ruby |
| Git 操作 | 不包進 binary；desktop 預設呼叫使用者本機 `git` binary，缺 Git 時由 `doctor` 與 close-loop / linked-update flow 阻斷並提示安裝 | 保留 credential helper、SSH key、hooks、GPG signing、LFS、submodule 與使用者既有 Git 語意 |

SQLite 決策尤其重要：

- **預設不採用 CGO SQLite**：例如 `mattn/go-sqlite3` 需要 CGO、C compiler 與平台相依 build chain，Windows 維護成本較高。
- **優先採用 pure Go SQLite**：例如 `modernc.org/sqlite` 可把 SQLite engine 編入 binary，較符合單一 binary 與低部署摩擦目標。
- **若 performance / compatibility 必須使用 CGO**：必須在 plan 或 ADR 中明確記錄原因、fallback、CI matrix 與 Windows 安裝成本，不可默默引入。

使用者端理想體驗：

```bash
git clone <repo>
./bin/ai-skill runtime migrate
./bin/ai-skill runtime replay
./bin/ai-skill runtime repair
```

上述命令不應要求使用者先安裝 Ruby、Python、sqlite3 CLI、pip、gem 或 C compiler。若某階段仍是 wrapper mode，`doctor` 必須清楚標示哪些能力仍依賴外部 runtime，以及對應遷移計畫。

#### Git external dependency policy

Desktop 平台不應把 Git 包進 `ai-skill` binary。Git 是成熟且使用者環境高度客製化的外部工具，內建 Git 反而容易造成 credential、SSH、hooks、GPG signing、LFS、submodule 與 global config 行為不一致。

策略：

- `ai-skill doctor` 必須檢查 `git` 是否存在、版本是否符合最低需求、是否可執行 `git rev-parse` / `git status`。
- 任何需要 linked updates、writeback transaction、commit、push、runtime sync 或 close-loop 的命令，如果偵測不到 Git，必須阻斷並提醒使用者安裝 Git，不得進入半套更新流程。
- 錯誤訊息必須明確說明：目前 binary 內建 runtime dependencies，但 Git 是 desktop 必要外部依賴。
- 若只需要讀取 repo metadata，可評估 Go git library 作為輔助；但 commit、push、status、hooks、credential 相關行為預設仍呼叫本機 `git`。
- CI / contributor 文件需把 Git 列為必備工具，而 Ruby、Python、sqlite3 CLI 則逐步降為 wrapper / compatibility mode 依賴。

### 3. 以 adapter 包住既有 script，再逐步移植

第一版 Go CLI 不必立刻重寫所有 Ruby compiler / generator。建議分三層：

| 層級 | 策略 | 目的 |
| --- | --- | --- |
| Wrapper mode | Go CLI 呼叫既有 script，統一參數、輸出、exit code | 快速建立跨平台入口與測試 harness |
| Native mode | 將 shell script 與 path/git/file 操作移植到 Go | 優先解決 Windows/macOS/Linux 差異 |
| Compiler migration | 評估 Ruby generator / compiler 是否移植到 Go 或保留 Ruby | 避免 premature rewrite 破壞 runtime semantics |

### 4. Runtime compiler 的保守策略

`runtime/compiler/compiler-engine.rb` 是 `runtime/runtime.db` 的 source-to-generated compiler。遷移策略必須保守：

- Phase 1-2 只用 Go CLI 包裝 compiler，不改變 compiler semantics。
- Phase 3 建立 golden DB / JSON snapshot，比對 Ruby compiler 與 Go prototype 的輸出。
- 只有在 snapshot parity、SQLite schema validation、generated_surfaces content assertion 全部通過後，才可切換 compiler source。
- 若保留 Ruby compiler，Go CLI 的 `doctor` 必須能清楚提示 Ruby / encoding / SQLite 缺失。

## Phase Plan

### Phase 0：Discovery & Contract（P0）

**目標**：先定義 script runtime 的行為契約，不急著改寫。

Current artifacts：

- [`README.md`](../../scripts/ai-skill-cli/docs/README.md)：Phase 0 artifact index。
- [`change-brief.md`](../../scripts/ai-skill-cli/docs/change-brief.md)：change brief 與 scope / blocker。
- [`command-contract.md`](../../scripts/ai-skill-cli/docs/command-contract.md)：command surface、exit codes、side effects。
- [`support-matrix.md`](../../scripts/ai-skill-cli/docs/support-matrix.md)：desktop / mobile 支援矩陣。
- [`bdd-scenarios.md`](../../scripts/ai-skill-cli/docs/bdd-scenarios.md)：BDD-lite scenarios。
- [`test-fixture-plan.md`](../../scripts/ai-skill-cli/docs/test-fixture-plan.md)：fixture plan。

Tasks：

- [x] 盤點所有 `scripts/`、`runtime/compiler/` 與 git hook 的輸入、輸出、寫檔位置、exit code、環境變數、外部命令依賴。
- [x] 建立 command contract 文件，列出每個命令的 arguments、dry-run 行為、side effects、validation signal。
- [x] 建立 platform support matrix：Windows、macOS、Linux、iOS、Android。
- [x] 建立風險清單：symlink、chmod、git hook、SQLite、Ruby gems、encoding、shell quoting、路徑大小寫。
- [x] 決定 repository layout：`scripts/ai-skill-cli/docs/` 放文件先行 artifacts，未來 `scripts/ai-skill-cli/cmd/ai-skill/`、`scripts/ai-skill-cli/internal/`、`scripts/ai-skill-cli/testdata/` 放 Go code 與 fixtures。

Phase 0 風險清單：

| 風險 | 必要控制 |
| --- | --- |
| Symlink / copy fallback | Mirror sync 必須區分受管理檔案與使用者檔案；Windows fallback 不得假設 POSIX symlink 行為。 |
| `chmod` / executable bit | Native Go 命令的核心行為不得依賴 POSIX executable bit；hook installation 必須回報平台限制。 |
| Git hook behavior | `hooks install` 需要 Git，且不得觸發 commit / push；hook mutation 前必須回報 unsafe repo state。 |
| SQLite engine mismatch | 優先採用 pure Go SQLite；wrapper mode 必須回報仍在使用外部 Ruby / SQLite 行為的範圍。 |
| Ruby gems / runtime availability | Wrapper mode 回傳 `missing_dependency`；native migration 必須先有 parity tests，才可取代 compiler 行為。 |
| Encoding / UTF-8 | Runtime wrapper 必須強制 UTF-8-compatible execution，並驗證 generated surfaces，不只看 command exit code。 |
| Shell quoting | Windows tests 必須避免 Git Bash / WSL 假設，並使用 argv-level command construction。 |
| Path case / separator differences | Path abstraction 必須 normalize drive letters、UNC paths、spaces、mixed separators，以及 case-sensitive / insensitive comparisons。 |

Completion criteria：

- 所有現有 scripts 都有 command contract。
- 每個命令都標示是否可 native Go migration、需要 wrapper、或暫不支援。
- iOS / Android 有初步結論：native local run、remote control、或 unsupported。

### Phase 1：Go CLI Skeleton & Doctor（P0）

**目標**：建立可跨平台編譯與自我診斷的 Go CLI。

Tasks：

- [x] 在 `scripts/ai-skill-cli/` 新增 `go.mod` 與 CLI skeleton。
- [x] 建立 dependency policy：pure Go dependency 優先；需要 CGO、外部 binary 或平台 SDK 時必須列為 exception。
- [x] 選型 SQLite library，預設評估 `modernc.org/sqlite`，並記錄是否排除 `mattn/go-sqlite3` 作為預設方案。
- [x] 實作 `ai-skill doctor`：檢查 Git external dependency、SQLite、Ruby、Python、repo root、write permission、hooksPath、PATH；其中缺 Git 必須明確提示安裝。
- [x] 實作 path / OS abstraction，禁止散落 OS-specific string manipulation。
- [x] 建立 `--json` / `--plain` output contract。
- [ ] 建立基本 unit tests 與 GitHub Actions matrix：Windows、macOS、Linux。

Progress notes：

- 已建立最小 `doctor` slice：platform、Git discovery、repo root、write permission 與 runtime DB presence checks；缺 Git 且傳入 `--require-git` 時回傳 `missing_dependency` / `missing_git`，且不產生 mutations。
- 已建立 `doctor` unit tests，覆蓋缺 Git、`--require-git --json`、plain output 與 unknown command；尚未建立 GitHub Actions matrix。
- 已新增 [`dependency-policy.md`](../../scripts/ai-skill-cli/docs/dependency-policy.md)，確認 pure Go dependency 優先、Git 作為 external dependency、`modernc.org/sqlite` 作為預設 SQLite engine，且 `mattn/go-sqlite3` 不作為預設。
- `doctor --check-runtime` 已用 `modernc.org/sqlite` 建立 in-memory query proof，並在找到 `runtime.db` 時執行 integrity check。
- `doctor` 已回報 PATH、Git、repo root、hooksPath、write permission、Ruby / Python wrapper-mode dependency diagnostics；Ruby / Python 只作為 wrapper-mode optional dependency，不列為長期核心依賴。
- 已新增 `internal/pathutil`，集中處理 report path normalization 與 PATH summary；unit tests 覆蓋 Windows drive、UNC path、mixed separators、spaces 與 relative paths。

Completion criteria：

- `go test ./...` 在三大桌面 OS 通過。
- `ai-skill doctor --json` 輸出可被 CI / agent 解析。
- 缺 Git 時，`doctor` 與需要 linked-update / close-loop 的命令會以穩定 exit code 阻斷並顯示安裝提示。
- SQLite、YAML、JSON 的基本讀寫測試不依賴外部 CLI 或本機 C compiler。
- 沒有任何命令在未傳 `--confirm` 或非 dry-run 模式下執行破壞性操作。

### Phase 2：Shell Script Migration（P1）

**目標**：優先移除最容易受平台影響的 shell 腳本依賴。

Candidate commands：

- `init-new-project.sh` → `ai-skill init-project`
- `agent-goals.sh` → `ai-skill goals`
- `install-hooks.sh` → `ai-skill hooks install`
- `sync-cursor-bundle.sh` → `ai-skill sync-cursor-bundle`
- `ai-skill-close-loop.sh` → `ai-skill close-loop`

Tasks：

- [ ] 用 fixtures 模擬新專案初始化，不寫入真實使用者目錄。
- [ ] 用 temporary git repo 測試 close-loop、dirty owner group、lock、merge/rebase state、dry-run。
- [ ] 用 missing-git fixture 測試 linked-update / close-loop 命令會阻斷並提醒使用者安裝 Git。
- [ ] 對 symlink / copy fallback 建立 Windows 行為策略。
- [ ] 更新 `scripts/README.md`，標示舊 script 與新 CLI 的對應關係與 deprecation policy。

Completion criteria：

- 新 CLI 與舊 script 在 fixtures 上輸出一致或差異有文件化理由。
- Windows 不需要 Git Bash / WSL 即可跑 native CLI。
- close-loop 不會混入 unrelated dirty changes。

### Phase 3：Runtime & Knowledge Tooling Strategy（P1）

**目標**：處理 Ruby generator / validator 的跨平台問題。

Tasks：

- [ ] 先建立 Go wrapper：`ai-skill runtime refresh`、`runtime validate`、`runtime compile`。
- [ ] Wrapper 必須固定 UTF-8 環境，並在缺 Ruby / SQLite 時給明確修復建議。
- [ ] 建立 golden fixture：同一組 source 產出固定 `runtime-report.md`、model reports、SQLite index、`runtime.db` assertion。
- [ ] 建立 native SQLite proof-of-concept：用 pure Go SQLite 開啟、查詢、寫入測試 DB，確認 Windows / macOS / Linux 無外部 sqlite3 CLI 依賴。
- [ ] 評估哪些 Ruby validator 適合原生 Go 重寫，哪些應保留 Ruby。
- [ ] 若開始移植 compiler，先建立 Ruby vs Go parity test，不得直接替換 production compiler。

Completion criteria：

- Go CLI 可一鍵跑完整更新流程並輸出 machine-readable summary。
- `runtime.db` 的 `generated_surfaces` 必須可查到 modified source 的內容 assertion。
- native SQLite path 可在三大桌面 OS 執行基本 migration / query / assertion。
- Ruby / Go parity test 未通過前，不切換 runtime compiler source。

### Phase 4：Cross-Platform Release & Distribution（P1）

**目標**：讓使用者不需要先安裝完整開發環境也能跑主要工具。

Tasks：

- [ ] 建立 release artifact：Windows `.exe`、macOS universal / arch-specific、Linux amd64 / arm64。
- [ ] 建立 `bin/` 或 release artifact layout 決策：repository 是否提交 binary、只在 GitHub Releases 發佈，或提供 local build output。
- [ ] 建立 checksum 與版本輸出：`ai-skill version`。
- [ ] 建立 GitHub Actions cross-compile workflow，輸出 Windows、Linux、macOS artifacts。
- [ ] 建立 upgrade / rollback 文件。
- [ ] 評估 Homebrew、Scoop、winget、GitHub Releases、直接下載 binary 的維護成本。

Completion criteria：

- 三大桌面 OS 都能下載 binary 後執行 `doctor` 與 dry-run commands。
- 發佈 artifact 不要求使用者安裝 Ruby、Python、sqlite3 CLI、pip、gem 或 C compiler 才能跑核心命令。
- 文件說明 source build 與 binary install 的差異。
- release 流程有 dry-run 與 artifact verification。

### Phase 5：Mobile Feasibility Evaluation（P2）

**目標**：明確回答 iOS / Android 是否可行，而不是模糊承諾。

#### iOS execution model

iOS 不是 general-purpose executable OS，不能假設使用者可像桌面系統一樣下載 binary 後執行：

```bash
git clone <repo>
./bin/ai-skill runtime migrate
```

iOS 的核心限制：

- 任意 executable 不能持久化後直接執行。
- 可執行邏輯必須存在於 App sandbox、browser sandbox、或遠端機器上。
- Git、terminal、interpreter、SQLite engine 若要在 iOS 本機執行，通常必須由 App 內建，例如 Git client app、terminal app 或專用 wrapper app。
- Safari 下載 binary 不代表可直接執行。

因此 iOS 不應列入「native single binary desktop target」。可行路線應分成：

| 路線 | 定位 | 可行性 | 主要風險 |
| --- | --- | --- | --- |
| App 內建 runtime | iOS app 內建 Git / runtime / terminal / SQLite / config editor | 可行但需要 App 開發與發佈 | App sandbox、檔案存取、App Store policy、版本更新、credential storage |
| Browser/WASM | `runtime.wasm` 在 browser sandbox 執行 governance runtime、YAML / state inspect、部分 replay | 可行，最接近免安裝 native runtime | local repo 存取、持久化、SQLite WASM、效能、離線能力、browser storage 限制 |
| SSH / remote runner | iPhone 作為 control plane，實際 runtime 在 VPS、NAS、Mac mini、Linux mini PC 或桌面機 | 高可行，最符合治理 runtime 的 control-plane 性質 | 遠端授權、金鑰管理、網路可用性、審計、安全邊界 |
| Native arbitrary binary | iOS 直接下載並執行 `ai-skill` binary | 不可作為目標 | iOS security model 不允許一般用途 executable persistence |

建議預設方向：

- 桌面與 CI：以 Go single binary 作為 primary runtime。
- iOS：以 control plane / inspect UI / remote trigger 為主，不承諾本機 native binary。
- Browser/WASM：可作為 governance runtime inspect、replay UI、state validation 的候選方向。
- SSH / remote runner：可作為近期最實用方案，讓 iPhone 管理遠端 Linux / macOS runner。

Evaluation dimensions：

| 平台 | 可行方向 | 主要限制 |
| --- | --- | --- |
| Android | Termux / app sandbox / remote runner client | git、SQLite、檔案權限、背景任務、使用者資料路徑、shell compatibility |
| iOS | App 內建 runtime / Browser-WASM / SSH remote runner / control plane UI | sandbox、任意 binary 不可執行、Git 與 repo 存取、credential storage、browser storage、遠端授權 |

Decision options：

- **App-contained local runner**：由 iOS / Android App 內建 runtime、Git、SQLite 與檔案管理；不等同任意下載 binary 執行。
- **Browser/WASM runner**：把部分 runtime 編成 WASM，在 browser sandbox 執行 state inspect、replay、validation 或 UI。
- **Remote control client**：手機只觸發桌面 / server runner，iOS 上最可能可行。
- **Unsupported**：若成本高於收益，明確標示不支援，不讓 agent 誤判。

Completion criteria：

- 寫出 iOS / Android support decision record。
- iOS decision record 必須明確排除 native arbitrary binary，並在 App-contained、Browser/WASM、SSH remote runner 之間做取捨。
- 若不支援，CLI `doctor`、Browser/WASM UI 或文件要明確顯示 unsupported reason。
- 若支援 remote control client，必須另開安全與授權計畫，不混在本計畫直接實作。

### Phase 6：Deprecation & Closure（P2）

**目標**：完成舊 script 到新 CLI 的治理閉環。所有 replacement 範圍內的舊 script 在新 CLI parity、fixtures、文件與 release gate 全部通過後，必須直接刪除；不得長期保留兩套入口造成 drift。只有 Git hook adapter、明確 `tool-specific` adapter，或使用者另行批准的短期 compatibility wrapper 可以保留，且必須有 owner、移除條件與期限。

Tasks：

- [ ] 更新 `scripts/README.md`：新 CLI 為 primary，列出舊 script 刪除順序、剩餘 adapter 例外與移除條件。
- [ ] 更新 `enforcement/linked-updates.md` 中與 close-loop / scripts 有關的說明。
- [ ] 更新必要的 git hook 文件與 ai-tools 文件。
- [ ] 刪除已被 native CLI 或已驗證 wrapper replacement 覆蓋的舊 shell / Ruby / Python script；不得只標記 deprecated 後長期保留。
- [ ] 若保留 thin wrapper，wrapper 必須只轉呼叫 `ai-skill`，並在同一階段記錄刪除日期或刪除條件。
- [ ] 執行 Plan Completion Closure：validator、linked updates、`plans/README.md` 狀態、搬移 archived、commit / push、readback。

Completion criteria：

- replacement 範圍內的舊 script 已刪除；剩餘 hook / tool-specific adapter 都有明確保留原因、owner 與移除條件。
- 所有文件、runtime generated surfaces、測試與 release artifact 一致。
- `script-parity-inventory.md` 中每個舊入口都有最終 disposition：`deleted`、`hook adapter retained`、`tool-specific adapter retained` 或 `explicitly out of scope`。
- plan 完成後移到 `plans/archived/` 或明確標註持續生效例外。

## Testing Strategy

| 測試類型 | 必須覆蓋 |
| --- | --- |
| Unit tests | path normalization、OS detection、exit code、JSON output、dry-run planner |
| Golden tests | command output、generated reports、runtime.db assertion、hook templates |
| Fixture tests | temporary repo、dirty files、merge/rebase state、missing Git、missing dependency、permission denied |
| Cross-OS CI | Windows、macOS、Linux；至少覆蓋 amd64，arm64 視 CI 能力納入 |
| Compatibility tests | 新 CLI wrapper vs 舊 script 的行為比對 |
| Runtime validation | `ruby scripts/refresh-knowledge-runtime.rb`、`ruby scripts/validate-runtime-db.rb`、SQLite content assertion |
| Release tests | binary checksum、`ai-skill version`、`doctor`、dry-run commands |

測試原則：

- 所有 destructive / write operations 先以 dry-run fixture 測試。
- 涉及 git 的測試一律使用 temporary repo，不碰使用者真實 working tree。
- 缺 Git 情境必須用 PATH isolation 或 fake executable fixture 測試，確認 linked-update / close-loop 會阻斷並提示安裝。
- 涉及 home directory、Cursor bundle、hook install 的測試必須支援 fake home / fake config path。
- Windows 測試不得假設 POSIX shell 存在。
- Runtime compiler migration 必須使用 parity test，不得只比較 exit code。

## Documentation Requirements

Implementation phase 必須同步建立或更新：

- `scripts/README.md`：CLI 使用方式、舊 script 對應表、環境需求、dry-run、錯誤處理。
- `scripts/ai-skill-cli/docs/` 下的 command contract 文件：每個命令的輸入、輸出、side effects、exit code。
- `scripts/ai-skill-cli/README.md`：build、test、release、cross-compile。
- `governance/lifecycle/knowledge-update-flow.md` 或相關 validator 文件：若完整更新流程改由 CLI 統一執行，需同步更新 Step 9-11。
- `ai-tools/`：只有當工具 mirror / Cursor bundle / hook 行為有變更時才更新。

## Affected Files

| 檔案 | 變更類型 | Phase |
| --- | --- | --- |
| `plans/active/2026-05-21-0834-cross-platform-go-script-runtime.md` | 新增計畫 | Phase 0 |
| `scripts/ai-skill-cli/` | 新增 CLI runtime 開發根目錄；`docs/` 放文件先行 artifacts，未來 `cmd/` / `internal/` / `testdata/` 放 Go implementation 與 fixtures | Phase 0 / Phase 1 |
| `plans/README.md` | 新增 active plan 索引 | Phase 0 |
| `scripts/README.md` | 未來更新 CLI mapping / deprecation policy | Phase 2 / Phase 6 |
| `scripts/*.sh` | 未來 wrapper / deprecation / replacement | Phase 2 / Phase 6 |
| `scripts/*.rb` | 未來 wrapper / parity migration | Phase 3 |
| `runtime/compiler/compiler-engine.rb` | 若移植 compiler，需建立 parity gate | Phase 3 |
| `scripts/ai-skill-cli/cmd/ai-skill/` | 未來新增 | Phase 1 |
| `scripts/ai-skill-cli/go.mod` / `scripts/ai-skill-cli/go.sum` | 未來新增 | Phase 1 |
| `bin/` 或 GitHub Releases artifact layout | 未來決策 | Phase 4 |
| `.github/workflows/` 或既有 CI 設定 | 未來新增跨 OS 測試 | Phase 1 / Phase 4 |

## Recommended Execution Order

1. 先執行 Phase 0，完成 command contract 與 support matrix。
2. 再執行 Phase 1，建立 Go CLI skeleton、dependency policy、doctor 與跨 OS CI。
3. 優先執行 Phase 2，因為 shell scripts 是 Windows 相容性風險最高的部分。
4. Phase 3 必須保守推進，先證明 pure Go SQLite 與 runtime assertion 可行；runtime compiler 不可在沒有 parity test 前替換。
5. Phase 4 在 CLI 行為穩定後再做 release；release 目標是單一 binary，避免先發佈仍依賴 Ruby / Python / sqlite3 CLI 的不穩定工具。
6. Phase 5 可以與 Phase 1-2 並行做 feasibility research，但不得阻塞桌面平台支援。
7. Phase 6 只有在新 CLI 覆蓋主要能力且文件、測試、runtime surfaces 都通過後才能執行；執行時以刪除舊 script 為預設，保留 adapter 必須逐項例外化。

## Open Questions

- Go CLI root 已決定放在 `scripts/ai-skill-cli/`；後續 open question 是 package boundary 要如何拆分為 `cmd/`、`internal/` 與 `testdata/`。
- Binary artifact 是否應命名為 `ai-skill`、`runtime`，或拆成 `ai-skill` CLI 與 `runtime` subcommand？目前傾向單一 `ai-skill` binary，避免多工具分裂。
- SQLite library 是否採用 `modernc.org/sqlite` 作為預設？若採用 CGO SQLite，哪些 compatibility / performance 證據足以抵消部署成本？
- 是否允許把 release binary commit 到 `bin/`，或只透過 GitHub Releases / CI artifacts 發佈？若 commit binary，需評估 repo size、review 與安全掃描成本。
- Desktop Git 是否只設定最低版本，或還需要檢查 credential helper、SSH key、LFS 與 submodule 支援？最低要求應從 close-loop 與 linked-update 命令實際需求推導。
- Ruby runtime compiler 是長期保留，還是以 parity test 分階段移植到 Go？
- Release artifact 是否由 GitHub Actions 產生，或先只支援 source build？
- iOS / Android 的主要使用場景是本機執行、遠端觸發，還是只需要讀取文件與狀態？
- 是否需要把 `ai-skill close-loop --commit --push` 變成所有 agent 的預設 commit path，取代手動 git 流程？

## Validation Plan

本計畫本身完成時必須驗證：

- `plans/README.md` 已加入 active plan 索引。
- Markdown link check 通過。
- `ruby scripts/refresh-knowledge-runtime.rb` 通過。
- `ruby runtime/compiler/compiler-engine.rb` 通過，且 `runtime/runtime.db` 已包含 plans index 更新。
- `ruby scripts/validate-runtime-db.rb` 通過。
- SQLite assertion 可查到本計畫的 `plan_id` 或標題。
- `git diff --check` 通過。
- commit / push / readback 後 `git status --short --branch` clean。
