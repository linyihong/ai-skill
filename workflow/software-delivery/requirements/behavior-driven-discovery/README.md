# Behavior-Driven Discovery

## 目的

在 architecture 或 implementation 之前，先理解 expected behavior。這是 BDD-lite 的 discovery stage，不要求 Gherkin runner。

## 步驟

1. 找出 actor / system role。
2. 找出 expected observable behavior。
3. 寫出 in-scope / out-of-scope。
4. 對齊 shared language — 涉及 Ai-skill framework / runtime / cognitive / architecture 詞彙時，**先查 [`knowledge/glossary/ai-skill.md`](../../../../knowledge/glossary/ai-skill.md)** 取 canonical 定義；不在 glossary 的業務詞彙才在 requirements layer 自定。詞義衝突依 [`knowledge/glossary/README.md`](../../../../knowledge/glossary/README.md) §Vocabulary Resolution Priority 解析。
5. 若出現 ambiguity，轉到 `ambiguity-resolution/`。

## 參考

- `intelligence/engineering/requirements/behavior-modeling/`
- [`knowledge/glossary/`](../../../../knowledge/glossary/README.md) — Ai-skill 共享語彙的 canonical 來源
