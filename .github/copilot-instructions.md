# GitHub Copilot Bootstrap Entry

本檔為 Copilot project-wide custom instructions 的 thin pointer。不要在此複製 bootstrap obligations、格式、enum、examples、goal ledger、close-loop 或 runtime phase 細節。

## Required Reads

1. [`CORE_BOOTSTRAP.md`](../CORE_BOOTSTRAP.md)
2. [`runtime/core-bootstrap.yaml`](../runtime/core-bootstrap.yaml)
3. [`ai-tools/agent/copilot.md`](../ai-tools/agent/copilot.md)

依 canonical bootstrap contract 執行 required reads、Bootstrap Receipt、per-turn Cognitive Mode reporting 與 close-loop checks。若 Copilot 功能無法強制執行某項 gate，回報限制，並讓 repository hooks / CI / `ai-skill runtime validate` 作為 enforcement boundary。
