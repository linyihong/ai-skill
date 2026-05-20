# Ambiguity Resolution

## 目的

當需求存在影響 behavior、contract、architecture 或 validation 的 ambiguity 時，先降級 autonomy，不讓 agent 自行補需求。

## 步驟

1. 將不確定項標成 `assumption`、`open question`、`scoped out` 或 `invalidated`。
2. 判斷是否阻擋 implementation。
3. 高風險 ambiguity 需要 human alignment。
4. 更新 behavior contract / acceptance criteria。
