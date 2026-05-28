# GitHub Copilot Bootstrap Entry

本檔為 Copilot project-wide custom instructions 的 thin pointer。不要在此複製 bootstrap obligations、格式、enum、examples、goal ledger、close-loop 或 runtime phase 細節。

## Mandatory Startup

在回覆任何使用者請求前，必須先完成下方 Required Reads 並依 canonical bootstrap contract 回報 Bootstrap Receipt。不得因為請求看似只是列檔、簡單查詢、read-only、說明原因或已由對話摘要提供 context 而跳過 bootstrap。

## Required Reads

1. [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md)
2. [`runtime/core-bootstrap.yaml`](../runtime/core-bootstrap.yaml)
3. [`ai-tools/agent/copilot.md`](../ai-tools/agent/copilot.md)

依 canonical bootstrap contract 執行 required reads、Bootstrap Receipt、per-turn Cognitive Mode reporting 與 close-loop checks。若 Copilot 功能無法強制執行某項 gate，回報限制，並讓 repository hooks / CI / `ai-skill runtime validate` 作為 enforcement boundary。
