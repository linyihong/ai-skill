> 遵守 [共用規則索引](../../../../shared-rules/README.md)、[dependency-reading](../../../../shared-rules/dependency-reading.md)、[neutral-language](../../../../shared-rules/neutral-language.md)、[goal-action-validation](../../../../shared-rules/goal-action-validation.md) 與 [feedback-lessons](../../../../shared-rules/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

### 2026-05-12 - Private Adapter Inside Test Module

Status: candidate

#### One-line Summary

當 live smoke 測試需要 private adapter（簽章、解密、不透明參數提供者），但這些實作不能進入公開 SDK reactor 時，應將它們以 test-scoped class 形式放在獨立測試模組內，而非建立新的 Maven 模組。

#### Human Explanation

在 TATA SDK 開發過程中，SDK 核心（typed query、parser、mock transport、BDD 測試）位於公開 reactor 模組。而 private adapter（AES 簽章 mutator、AES 解密 decoder、語言不透明參數提供者、加密材料）需要類似正式環境的機敏值，這些值不能進入公開文件或已發布的 JAR。直覺上可能會想為 adapter 建立一個新的 Maven 模組，但專案規則已經定義了一個測試模組（`tata-sdk-tests`），它具有以下特性：

1. 與 parent reactor 分開建置（未列在 `<modules>` 中）。
2. 已包含 live smoke 測試基礎設施（`LiveExternalGossipFetchTest`、`LiveExternalGossipEnvTest`）。
3. 已透過 `<scope>test</scope>` 或 compile-time dependency 依賴 SDK 核心。
4. 已排除在公開同步邊界之外。

將 adapter 放在 `src/test/java/.../live/adapter/` 下，可保持 test-scoped、不可發布，並與使用它的 smoke 測試放在同一位置。不需要新模組、新 POM、新 reactor 條目。

#### Trigger

- 需要為 live smoke 測試實作 private adapter（簽章、解密、不透明值提供者、加密材料）。
- 專案已有一個獨立測試模組，該模組不在 parent reactor 中，且已排除在公開同步之外。
- 該測試模組已匯入 SDK 核心並擁有 live 測試基礎設施。
- 有人提議為 adapter 建立新的 Maven 模組，而非使用現有測試模組。

#### Evidence

- 工具：Maven 模組結構檢視、parent POM 的 `<modules>` 列表、測試模組的 `pom.xml`、現有 live 測試檔案。
- 去敏摘要：adapter 套件 `com.tata.sdk.tests.live.adapter` 包含 `TataSigningMutator`、`TataDecryptDecoder`、`TataCryptoMaterial`、`TataLanguageProvider`——全部為 test-scoped，程式碼中無正式環境機敏值，由環境變數驅動。
- 證據路徑：`apk-analysis-sdk/pom.xml` 顯示 `tata-sdk-tests` 不在 `<modules>` 中。`apk-analysis-sdk/tata-sdk-tests/pom.xml` 宣告了對 `tata-sdk` 的依賴。Live 測試已存在於 `src/test/java/com/tata/sdk/tests/live/` 下。

#### Generalized Lesson

當需要將 private adapter 程式碼（簽章、解密、不透明參數提供者、加密材料）加入一個已有獨立測試模組（不在 parent reactor 中）的專案時，應將 adapter 放在該測試模組內，以 test-scoped class 形式存在。不要建立新的 Maven 模組。該測試模組已具備：

- SDK 核心作為依賴。
- 被排除在公開同步／發布之外。
- 包含將使用 adapter 的 live 測試基礎設施。
- 與 reactor 分開獨立建置（不會污染公開建置）。

如果沒有這樣的測試模組，則建立一個明確排除在 parent reactor 和公開同步邊界之外的測試模組。

**整合模式**：將 adapter 整合到現有 live test harness 時，應採用**選擇性升級**（optional enhancement）模式——live test harness 自動偵測環境變數中是否有 adapter crypto material，若有則自動升級為 adapter 模式（動態簽章 + 解密），若無則維持既有行為（預先捕獲值注入 + passthrough）。這確保：

- 現有測試配置不受影響（向後相容）。
- 新配置只需設定 adapter 的 env var，無需修改測試程式碼。
- 單一 `requestMutator()` 和 `responseDecoder()` 方法根據 adapter 可用性動態選擇實作。
- Identity header（Authorization、Cookie）仍可獨立於簽章之外注入。

#### Agent Action

1. 檢查專案是否已有獨立測試模組（不在 parent reactor 的 `<modules>` 中）。
2. 如果有，將 private adapter class 放在該模組的 `src/test/java/.../<adapter-package>/` 下。
3. 如果沒有，建立一個新的測試模組，附上自己的 `pom.xml`，將其排除在 parent reactor 之外，並加入公開同步排除列表。
4. 將所有機敏值保留在環境變數或設定檔中——永遠不要將正式環境機敏值寫死在程式碼中。
5. 對 unit test 使用 synthetic test vector，而非正式環境機敏值。
6. 為 adapter 套件加入 `package-info.java`，內含去敏規則。
7. 將 adapter 整合到 live test harness 時，採用選擇性升級模式：
   - Live test harness 的 `fromSystem()` 方法偵測 adapter env var 是否存在。
   - 若存在，建立 adapter 實例（`TataSigningMutator`、`TataDecryptDecoder`）並傳入 harness。
   - `requestMutator()` 在有 adapter 時使用動態簽章，否則使用預先捕獲值注入。
   - `responseDecoder()` 在有 adapter 時使用解密 decoder，否則使用 passthrough。
   - `isRunnable()` 的 `hasPrivateRequestMaterial()` 和 `hasResponseBoundary()` 將 adapter 可用性納入考量。
8. 更新 feedback history 索引（`feedback_history/README.md`），加入新條目。

#### Goal / Action / Validation

- 目標：防止在已有合適測試模組的情況下，為 private adapter 程式碼建立不必要的 Maven 模組。
- 行動：將 private adapter class 放在現有獨立測試模組中，以 test-scoped 程式碼形式存在。
- 驗證或參考來源：adapter 在獨立建置測試模組時可編譯且測試通過（`mvn test -pl tata-sdk-tests` 或 `cd tata-sdk-tests && mvn test`）。Adapter class 不會發布到任何公開儲存庫。

#### Applies When

- 需要將 private adapter 程式碼（簽章、解密、不透明參數提供者、加密材料）加入專案。
- 專案已有獨立測試模組，不在 parent reactor 中。
- 該測試模組已依賴 SDK 核心並擁有 live 測試基礎設施。

#### Does Not Apply When

- 專案沒有獨立測試模組，且建立一個新模組不合理。
- Adapter 程式碼並非 private，可以放在公開 SDK 核心中。
- Adapter 需要作為跨多個專案的可重用函式庫。

#### Validation

在認為本 lesson 已適用之前，確認以下所有項目：

- Adapter 套件位於 `src/test/java/` 下（test-scoped，不會編譯進主要 JAR）。
- 測試模組未列在 parent reactor 的 `<modules>` 中。
- 測試模組可獨立建置（從其目錄執行 `mvn test` 或使用 `-f` 旗標）。
- Unit test 使用 synthetic vector，而非正式環境機敏值。
- `package-info.java` 記錄了套件的去敏規則。
- Feedback history 索引已更新。

#### Promotion Target

- `WORKFLOW.md`
- `checklists/`
- `process/`

#### Required Linked Updates

- 更新 skill feedback 根索引與 `common/README.md`。
- 如果 lesson 來自特定專案交接，將專案特定的 route 名稱、class 名稱和 env-var key 保留在專案文件中。
- 晉升時，從模組結構指引建立 cross-link，讓未來的 SDK 工作在建立新模組前能看到此模式。
