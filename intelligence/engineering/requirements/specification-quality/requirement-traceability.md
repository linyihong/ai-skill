# Requirement Traceability

**Status**: `candidate-intelligence`

## 模型

```text
requirement
→ behavior contract
→ acceptance criteria
→ validation target
→ execution artifact
```

## 目的

防止：

```text
requirement
→ inferred implementation
→ unverifiable success claim
```

## 規則

每個 requirement 若會影響 observable behavior，至少要能追到一個 validation target。若 validation target 缺失，不能宣稱完成。
