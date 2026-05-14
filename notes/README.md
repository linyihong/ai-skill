# Notes

`notes/` 存放個人觀察筆記，記錄開發過程中的發現、實驗結果與臨時想法。

## 用途

- 記錄尚未正式化的觀察與經驗
- 作為未來建立正式文件或 pipeline 的參考素材
- **不屬於 durable knowledge** — 筆記內容可能過時或不完整，不應作為決策依據

## 目前文件

| 檔案 | 說明 |
|------|------|
| [`intelligence-extraction-observations.md`](intelligence-extraction-observations.md) | Technique → Intelligence 分解過程中的觀察記錄，包含哪些 decision 可 atomize、哪些不易 atomize、格式觀察等 |

## 規則

1. 筆記內容不保證正確性或時效性
2. 當筆記內容被正式化為規則或文件後，應在筆記中標註「已遷移」並連結目標路徑
3. 不從 notes/ 建立 routing registry record（非 knowledge layer）

## 誰會參考這裡（Inbound References）

本層無 routing registry record（非 knowledge layer）。變更時無需通知其他層，但 promotion 時需更新目標層。

## 與既有層的關係

- [`intelligence/`](../intelligence/README.md)：筆記中發現的可重用 pattern 可 promotion 至此
- [`workflow/`](../workflow/README.md)：筆記中發現的流程可正式化至此
