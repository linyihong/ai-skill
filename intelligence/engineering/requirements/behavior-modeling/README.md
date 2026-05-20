# Behavior Modeling

本目錄保存 behavior understanding 的判斷智慧。它回答「使用者真正期待什麼可觀察行為」以及「行為邊界到哪裡」。

## 目前條目

| 文件 | 用途 |
| --- | --- |
| [`ubiquitous-language-alignment.md`](ubiquitous-language-alignment.md) | 將需求語言與 domain language 對齊。 |
| [`scenario-framing.md`](scenario-framing.md) | 把需求轉成 scenario，而不是自行補功能。 |
| [`actor-intent-modeling.md`](actor-intent-modeling.md) | 分辨 actor、intent、permission 與 expected outcome。 |
| [`acceptance-boundaries.md`](acceptance-boundaries.md) | 定義 acceptance criteria 的範圍與停止點。 |
| [`ambiguity-detection.md`](ambiguity-detection.md) | 辨識需要 human alignment 的 ambiguity。 |

## 原則

Behavior model 先於 domain model。若行為邊界不清，agent 不應直接建立 aggregate、API 或 implementation slice。
