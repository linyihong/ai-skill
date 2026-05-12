# Claude Code Auto-Bootstrap

啟動時自動載入以下流程：

1. 讀 [CORE_BOOTSTRAP.md](CORE_BOOTSTRAP.md) — 3 條必讀核心規則（~800 tokens）
2. 讀 [README.md](README.md) — 超短入口，了解 OS layout
3. 依 task intent 查 [skills-index.yaml](skills-index.yaml) 找到對應 skill
4. 依 [runtime/router/activation-rules.yaml](runtime/router/activation-rules.yaml) 決定 lazy-load rules
5. 先讀 `knowledge/summaries/` 對應 summary（300-500 tokens），需要時才展開全文

詳細使用說明見 [ai-tools/claude.md](ai-tools/claude.md)。
