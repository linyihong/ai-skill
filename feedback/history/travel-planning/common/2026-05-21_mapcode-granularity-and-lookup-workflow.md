> 遵守 [共用規則索引](../../../../enforcement/README.md)、[dependency-reading](../../../../enforcement/dependency-reading.md)、[neutral-language](../../../../enforcement/neutral-language.md)、[goal-action-validation](../../../../enforcement/goal-action-validation.md) 與 [feedback-lessons](../../../feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Mapcode Granularity and Lookup Workflow

### 2026-05-21 - 日本自駕行程 mapcode 粒度規則與查詢工具鏈

Status: validated

#### One-line Summary

自駕行程規劃時，同一大景點下若有距離相遠的小景點，每個停車點都必須有獨立 mapcode；查詢工具依 MapFan spots → Mapion phonebook → 旅遊部落格 → 最近巴士站代替 的優先順序執行。

#### Human Explanation

在規劃日本自駕行程時，初版只為「大景點」標註一個 mapcode（如奧入瀨溪流整體只放石ヶ戸）。但使用者指出：若大景點下有多個停車點且彼此距離較遠，每個停車點都需要獨立的 mapcode，否則開車時無法逐點設定目的地。

例如：奧入瀨溪流 7.3km 路段中有三個主要停車點（石ヶ戸、阿修羅の流れ路肩、銚子大滝），各距 2～3km，需分開列出。

同日不同目的地（如角館 → 田澤湖）距離達 25km+ 時，也需各自一行 mapcode。

另一個常見問題是「自然景觀類」景點（例如瀑布、溪流路段）不一定在主流地圖服務中登錄 Mapcode。此時使用以下替代方式：
- 使用最近的巴士站或停車場 Mapcode 作為替代，並加註說明
- 若是路肩停車點，明確標示「路肩」

#### Trigger

- 自駕行程中某一天有複數景點且彼此相隔 2km 以上
- 溪流、海岸線、高原等沿線型觀光地，且有多個停車點
- 同一日程出現「A 景點 → 同日 B 景點（距離遠）」的模式
- 使用者要求加入 mapcode 時，需檢查既有條目的粒度是否過粗

#### Evidence

- 任務：為東北旅行（奧入瀨溪流 10/31、角館→田澤湖 11/2）追加 mapcode
- 初版：奧入瀨溪流只以石ヶ戸 1 行表示
- 使用者回饋：「奧入瀨溪流有很多小景點。距離遠時都需要補記 mapcode」
- 修正後：拆成石ヶ戸、阿修羅の流れ（路肩）、銚子大滝 3 行
- 同日：確認田澤湖（御座石神社）的 mapcode 需與角館分開列出並追加

#### Generalized Lesson

**mapcode 粒度規則（自駕行程）**

```
同一景點下有多個停車點，且彼此相距 2km 以上
  → 各停車點建立獨立行

同日有 2 個以上不同目的地
  → 各目的地建立獨立行（不論距離）
```

**mapcode 查詢工具鏈（優先順序）**

**Step 1：MapFan spots 頁（最優先）**
```
URL: https://mapfan.com/spots/[spot_code]
```
- 用 Google 搜尋 `site:mapfan.com [景點名]` 取得 spot URL
- 用 WebFetch 取得頁面 → 抽出 `マップコード` 後面的數字
- 注意：mapfan.com 與 mapion.co.jp 會被 Chrome extension 封鎖（`navigate` 不可）

**Step 2：Mapion 電話帳頁**
```
URL: https://www.mapion.co.jp/phonebook/[カテゴリ]/[地域コード]/[施設ID]/
```
- 用 Google 搜尋 `site:mapion.co.jp [景點名] [市区町村]` 取得 URL
- 類別分類例：M06007（觀光名所）、M06005（神社寺院）、M12001（巴士站）
- 用 WebFetch 取得頁面 → 抽出 `マップコード`

**Step 3：旅遊部落格（WebSearch → WebFetch）**
- 搜尋 `[景點名] マップコード [預估前三碼]`
- 自然景觀點（瀑布、溪流）常出現在個人部落格，而不是專門旅行網站

**Step 4：使用最近巴士站或停車場作為替代**
- 用 `site:mapion.co.jp [最寄りバス停名] バス停` 搜尋並取得 mapcode
- 在表格中寫成 `[景點名]（[代替點名]附近）`，並明確加註

**Step 5：無法取得時**
- 在表格中記載 `要確認（Tel XXXX-XX-XXXX 索取）`
- 告知使用者該點無法直接用車用導航輸入

**自然景觀點（瀑布、溪流）的注意事項**

| 種類 | mapcode 登錄狀況 | 對應方式 |
|------|----------------|------|
| 道の駅、停車場 | 幾乎都有登錄 | Step 1-2 可取得 |
| 旅館、飯店 | 幾乎都有登錄 | Step 1-2 可取得 |
| 神社、寺院 | 多數有登錄 | Step 1-2 可取得 |
| 瀑布、溪流路段 | **常未登錄** | Step 3-4 替代 |
| 路肩停車點 | **通常未登錄** | 使用最近巴士站替代 |

#### Anti-patterns

- ❌ 將溪流、海岸線等沿線景點用「大景點 1 行 = 1 mapcode」概括
- ❌ 同日 A→B（相距 25km）時未記載 B 的 mapcode
- ❌ 找不到 mapcode 時，不在表格中加入行並沉默跳過
- ❌ 因路肩停車點沒有 mapcode 就省略記載；應使用巴士站替代處理

#### Agent Action

- 建立或審查自駕行程的 mapcode 表時，**必須**逐日確認各停車點，檢查同一景點內是否有相距 2km+ 的 sub-spot
- 找不到 mapcode 時也**必須**在表格中加入行，記載替代碼或「要確認（Tel XXXX）」
- 不得把「大景點 1 行概括」視為最終成果物

#### Goal / Action / Validation

- Goal: 讓所有自駕停車點都能在車用導航中個別輸入
- Action: 套用粒度規則（2km+ → 分割），並用工具鏈取得各 mapcode
- Validation or reference source: 行程表的 mapcode 表中「要確認」為零，或已記載替代碼與註解

#### Applies When

- 建立或更新日本自駕行程的 mapcode 表時
- 行程包含溪流、海岸線、高原等沿線型觀光地，且有多個停車點時
- 同日設定 2 個以上不同目的地時

#### Does Not Apply When

- 只有步行或公共交通的行程（不需要車用導航）
- 使用者要求的是 Google Maps 或地址導航，而不是 mapcode

#### Validation

- 行程表的所有停車點都有 mapcode 行
- 沿線景點（例如奧入瀨溪流）已拆成多行
- 使用替代碼時，已記載註解（例如「〇〇巴士站附近」）

#### Promotion Target

- `workflow/travel-planning/execution-flow.md`（Step 11 國家與地區特定檢查中的日本自駕段落，已加入 mapcode 粒度規則）
- `analysis/travel/sources-and-tools.md`（mapcode 查詢工具鏈）

#### Required Linked Updates

- `workflow/travel-planning/execution-flow.md`：Step 11 的日本自駕段落加入「沿線景點需依停車點建立獨立 mapcode 行」規則
- `analysis/travel/sources-and-tools.md`：補充 mapcode 查詢工具鏈（MapFan spots → Mapion 電話帳 → 旅行部落格 → 巴士站替代）
- `knowledge/summaries/travel-planning.md`：補充 mapcode 粒度規則的一行摘要
- 已依 [reusable-guidance-boundary.md](../../../../enforcement/reusable-guidance-boundary.md) 確認：具體旅館名稱、地名已抽象化，僅保留通用規則與工具鏈

#### Related

- 工具：`WebFetch`、`WebSearch`、`mcp__Claude_in_Chrome__navigate`（注意 mapfan/mapion 會被 Chrome extension 封鎖）
- 參考：[hotel-availability-check-workflow.md](./2026-05-19_hotel-availability-check-workflow.md)（Chrome extension 封鎖域名的處理模式共通）
- mapcode 可信來源：`mapfan.com/spots/`、`mapion.co.jp/phonebook/`、個人旅行部落格（例如 ks2345.sakura.ne.jp）
