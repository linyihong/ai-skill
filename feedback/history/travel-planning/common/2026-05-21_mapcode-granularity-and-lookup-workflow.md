> 遵守 [共用規則索引](../../../../enforcement/README.md) 與 [feedback-lessons](../../../../enforcement/feedback-lessons.md)；本檔只寫本條 lesson，不重複貼上共用政策全文。

# Mapcode Granularity and Lookup Workflow

### 2026-05-21 - 日本自駕行程 mapcode 粒度規則與查詢工具鏈

Status: active

#### One-line Summary

自駕行程規劃時，同一大景點下若有距離相遠的小景點，每個停車點都必須有獨立 mapcode；查詢工具依 MapFan spots → Mapion phonebook → 旅遊部落格 → 最近バス停代替 的優先順序執行。

#### Human Explanation

在規劃日本自駕行程時，初版只為「大景點」標註一個 mapcode（如奧入瀨溪流整體只放石ヶ戸）。但使用者指出：若大景點下有多個停車點且彼此距離較遠，每個停車點都需要獨立的 mapcode，否則開車時無法逐點設定目的地。

例如：奧入瀨溪流 7.3km 路段中有三個主要停車點（石ヶ戸、阿修羅の流れ路肩、銚子大滝），各距 2～3km，需分開列出。

同日不同目的地（如角館 → 田澤湖）距離達 25km+ 時，也需各自一行 mapcode。

另一個常見問題是「自然景觀類」景點（滝、渓流の流れ等）在主流地図サービスにマップコードが登録されていない場合がある。このときは以下の代替策を使う：
- 最寄りのバス停・駐車場のマップコードを代替として使い、注釈を付ける
- 路肩駐車ポイントの場合は「路肩」と明記

#### Trigger

- 自駕行程中某一天有複數景點且彼此相隔 2km 以上
- 渓流・海岸線・高原など「沿線」タイプの観光地（複数の停車点がある）
- 同一日程に「A 景點 → 同日 B 景點（距離遠）」のパターンがある
- ユーザーが「mapcode を追加して」と言ったとき、既存エントリの粒度が粗すぎないか確認する

#### Evidence

- タスク：東北旅行（奧入瀨溪流 10/31、角館→田澤湖 11/2）のmapcodeを追加
- 初版：奧入瀨溪流を石ヶ戸の1行のみで記載
- ユーザーフィードバック：「奧入瀨溪流には小景點がたくさんある。距離が遠い場合はすべてmapcode補記が必要」
- 修正後：石ヶ戸・阿修羅の流れ（路肩）・銚子大滝 の3行に分割
- 同日：田澤湖（御座石神社）のmapcodeが角館とは別に必要と判明し追加

#### Generalized Lesson

**mapcode 粒度ルール（自駕行程）**

```
同一景點下に停車ポイントが複数かつ相互に 2km 以上離れている
  → 各停車点に独立行を設ける

同日に異なる目的地が 2か所以上ある
  → 各目的地に独立行を設ける（距離問わず）
```

**mapcode 查詢工具鏈（優先順序）**

**Step 1：MapFan spots ページ（最優先）**
```
URL: https://mapfan.com/spots/[spot_code]
```
- Google で `site:mapfan.com [景點名]` で spot URL を取得
- WebFetch でページ取得 → `マップコード` 文字列に続く数字を抽出
- 注意：mapfan.com と mapion.co.jp は Chrome extension でブロックされる（`navigate` 不可）

**Step 2：Mapion 電話帳ページ**
```
URL: https://www.mapion.co.jp/phonebook/[カテゴリ]/[地域コード]/[施設ID]/
```
- Google で `site:mapion.co.jp [景點名] [市区町村]` で URL を取得
- カテゴリ分類例：M06007（観光名所）、M06005（神社寺院）、M12001（バス停）
- WebFetch でページ取得 → `マップコード` を抽出

**Step 3：旅遊部落格（WebSearch → WebFetch）**
- `[景點名] マップコード [予想コード前3桁]` で検索
- 自然景観スポット（滝・渓流）は専門旅行サイトではなく個人ブログに掲載されていることが多い

**Step 4：最近バス停・駐車場を代替として使用**
- `site:mapion.co.jp [最寄りバス停名] バス停` で mapcode 取得
- 表に `[景點名]（[代替ポイント名]付近）` と注釈を明記

**Step 5：未取得の場合**
- 表に `要確認（Tel XXXX-XX-XXXX 索取）` と記載
- カーナビ不可の旨をユーザーに伝える

**自然景観スポット（滝・渓流）の注意点**

| 種類 | mapcode 登録状況 | 対処 |
|------|----------------|------|
| 道の駅・駐車場 | ほぼ必ず登録あり | Step 1-2で取得可 |
| 旅館・ホテル | ほぼ必ず登録あり | Step 1-2で取得可 |
| 神社・仏閣 | ほぼ登録あり | Step 1-2で取得可 |
| 滝・渓流の流れ | **登録なしが多い** | Step 3-4で代替 |
| 路肩駐車ポイント | **未登録** | 最寄りバス停を代替使用 |

#### Anti-patterns

- ❌ 渓流・海岸線などの沿線スポットを「大景點 1行 = 1 mapcode」でまとめる
- ❌ 同日に A→B（25km 離れ）で B の mapcode を記載しない
- ❌ mapcode が見つからないとき、表に行を追加せず無言でスキップする
- ❌ 路肩停車ポイントに mapcode がないからといって記載を省略する（バス停代替で対処）

#### Related

- ツール：`WebFetch`、`WebSearch`、`mcp__Claude_in_Chrome__navigate`（mapfan/mapion はブロックされるため注意）
- 參考：[hotel-availability-check-workflow.md](./2026-05-19_hotel-availability-check-workflow.md)（Chrome extension ブロックドメインの対処パターン共通）
- mapcode 信頼済みソース：`mapfan.com/spots/`、`mapion.co.jp/phonebook/`、個人旅行ブログ（ks2345.sakura.ne.jp 等）
