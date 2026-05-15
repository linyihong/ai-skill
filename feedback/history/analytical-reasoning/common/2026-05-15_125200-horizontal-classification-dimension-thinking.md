# Horizontal Classification Dimension Thinking

## 觀察

在將 `java-tsv-trim-split-trailing-empty.md` 分類時，最初直接放入 `analytical-reasoning/failure/`，因為該目錄已有 failure 子層。但這個分類忽略了知識的「語言特定性」——Java `String.trim()` 行為是 Java 標準庫的語言特定知識，不屬於跨語言的分析技術失敗模式。

## 教訓

分類知識時，不應只檢查現有子層是否能容納，而應先思考：

1. **橫向思考**：這份知識是否屬於全新的分類維度？
   - 語言特定知識 → `language-specific/<lang>/`
   - 框架特定知識 → `framework-specific/<framework>/`
   - 平台特定知識 → `platform-specific/<platform>/`
2. **只有當不屬於新維度時**，才檢查現有子層是否能容納。

## 已套用的改善

- `knowledge-update-flow.md` Step 2.4 已加入橫向維度決策樹
- `intelligence/engineering/` 下新增 `language-specific/` 維度
- `intelligence/engineering/language-specific/java/failure/` 已建立

## 適用範圍

任何需要將知識文件分類到 `intelligence/engineering/` 下的場景。

## 觸發信號

- 知識內容涉及特定語言的標準庫行為
- 知識內容涉及特定框架的 API 行為
- 知識內容涉及特定平台的 runtime 特性
- 直覺上「放在現有目錄好像不太對」但說不出原因
