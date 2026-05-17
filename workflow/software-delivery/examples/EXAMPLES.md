# Examples — Simplicity First & Surgical Changes

本文件提供具體的 code-level 範例，展示 Simplicity First 與 Surgical Changes 兩條核心原則在實際開發中的應用。每個範例展示「LLM 常見的錯誤行為」vs「應該怎麼做」。

---

## 1. Think Before Coding（先思考再寫 Code）

### 範例 1：隱含假設

**使用者請求：**「加一個匯出使用者資料的功能」

**❌ LLM 常見錯誤（默默假設）**

```python
def export_users(format='json'):
    """Export all users to a file."""
    users = User.query.all()

    if format == 'json':
        with open('users.json', 'w') as f:
            json.dump([u.to_dict() for u in users], f)
    elif format == 'csv':
        with open('users.csv', 'w') as f:
            writer = csv.DictWriter(f, fieldnames=['id', 'email', 'name'])
            writer.writeheader()
            writer.writerows([u.to_dict() for u in users])

    return f"Exported {len(users)} users"
```

**問題：**
- 假設要匯出**所有**使用者（沒問分頁？隱私？）
- 假設檔案位置（沒問放哪裡）
- 假設要匯出哪些欄位
- 假設 CSV fieldnames 而沒檢查實際資料結構

**✅ 應該做的事（先列出假設再問）**

```
實作之前，我需要先確認：

1. **範圍**：匯出所有使用者還是篩選後的子集？（隱私影響）
2. **格式**：你說的「匯出」是指：
   - 瀏覽器下載檔案？
   - 背景 job + email 通知？
   - API endpoint 回傳資料？
3. **欄位**：哪些使用者欄位？（有些可能是敏感的）
4. **數量**：一般有多少使用者？（影響做法）

最簡單的做法：加一個 API endpoint 回傳分頁 JSON。
需要更多資訊才能決定檔案匯出的做法。

你的偏好是什麼？
```

### 範例 2：多種解讀

**使用者請求：**「讓搜尋更快」

**❌ LLM 常見錯誤（默默選一種）**

```python
# 直接加了 caching、database indexes、async processing
# 完全沒問「哪種快」
@lru_cache(maxsize=1000)
async def search(query: str) -> List[Result]:
    # ... 200 行最佳化 code
```

**✅ 應該做的事**

```
「讓搜尋更快」可能代表：

1. **回應時間更快**（< 100ms 取代 500ms）
   - 加 database indexes
   - 快取常用查詢
   - 預估工時：2 小時

2. **處理更多並發搜尋**（throughput）
   - 加 async/await
   - Connection pooling
   - 預估工時：4 小時

3. **感覺更快**（UX）
   - 立即顯示部分結果
   - Progressive loading
   - 預估工時：3 小時

目前搜尋約 500ms。哪個面向對你的 use case 最重要？
```

---

## 2. Simplicity First（簡潔優先）

### 範例 1：過度抽象

**使用者請求：**「加一個計算折扣的 function」

**❌ LLM 常見錯誤（過度工程）**

```python
from abc import ABC, abstractmethod
from enum import Enum
from typing import Protocol, Union
from dataclasses import dataclass


class DiscountStrategy(ABC):
    @abstractmethod
    def calculate(self, amount: float) -> float:
        pass


class PercentageDiscount(DiscountStrategy):
    def __init__(self, percentage: float):
        self.percentage = percentage

    def calculate(self, amount: float) -> float:
        return amount * (self.percentage / 100)


class FixedDiscount(DiscountStrategy):
    def __init__(self, fixed_amount: float):
        self.fixed_amount = fixed_amount

    def calculate(self, amount: float) -> float:
        return min(self.fixed_amount, amount)


@dataclass
class DiscountConfig:
    strategy: DiscountStrategy
    min_purchase: float = 0.0
    max_discount: float = float('inf')


class DiscountCalculator:
    def __init__(self, config: DiscountConfig):
        self.config = config

    def apply_discount(self, amount: float) -> float:
        if amount < self.config.min_purchase:
            return 0
        discount = self.config.strategy.calculate(amount)
        return min(discount, self.config.max_discount)


# Usage requires 30+ lines of setup for a simple calculation
```

**✅ 應該做的事（簡單）**

```python
def calculate_discount(amount: float, percent: float) -> float:
    """Calculate discount amount. percent should be 0-100."""
    return amount * (percent / 100)


# Usage
discount = calculate_discount(100.0, 10.0)  # $10 off
```

**何時才需要加複雜度：** 當你**真的**需要多種 discount type 的時候。如果那個需求之後才出現，到時候再 refactor。

### 範例 2：Speculative Features

**使用者請求：**「把使用者偏好存到資料庫」

**❌ LLM 常見錯誤（不必要的功能）**

```python
class PreferenceManager:
    def __init__(self, db, cache=None, validator=None):
        self.db = db
        self.cache = cache or InMemoryCache()
        self.validator = validator or DefaultValidator()

    def save(self, user_id: int, prefs: dict,
             merge: bool = True,
             validate: bool = True,
             notify: bool = False) -> bool:
        """Save preferences with optional merging, validation, and notifications."""

        if validate:
            errors = self.validator.validate(prefs)
            if errors:
                raise ValidationError(errors)

        existing = self.get(user_id) if merge else {}

        if merge:
            existing.update(prefs)
            prefs = existing

        self.db.execute(
            "UPDATE users SET preferences = ? WHERE id = ?",
            (json.dumps(prefs), user_id)
        )

        if self.cache:
            self.cache.set(f"prefs:{user_id}", prefs)

        if notify:
            self.notify_preference_change(user_id, prefs)

        return True

    def notify_preference_change(self, user_id: int, prefs: dict):
        # Another 30 lines for a feature nobody asked for
        pass
```

**✅ 應該做的事（只做被要求的）**

```python
def save_preferences(db, user_id: int, preferences: dict):
    """Save user preferences to database."""
    db.execute(
        "UPDATE users SET preferences = ? WHERE id = ?",
        (json.dumps(preferences), user_id)
    )
```

**之後再加：** Caching（當效能有問題時）、validation（當壞資料出現時）、merging（當需求出現時）。

---

## 3. Surgical Changes（外科手術式修改）

### 範例 1：順便 Refactor

**使用者請求：**「修掉 empty email 會讓 validator crash 的 bug」

**❌ LLM 常見錯誤（改太多）**

```diff
  def validate_user(user_data):
-     # Check email format
-     if not user_data.get('email'):
+     """Validate user data."""
+     email = user_data.get('email', '').strip()
+
+     # Validate email
+     if not email:
          raise ValueError("Email required")
-
-     # Basic email validation
-     if '@' not in user_data['email']:
+     if '@' not in email or '.' not in email.split('@')[1]:
          raise ValueError("Invalid email")

-     # Check username
-     if not user_data.get('username'):
+     # Validate username
+     username = user_data.get('username', '').strip()
+     if not username:
          raise ValueError("Username required")
+     if len(username) < 3:
+         raise ValueError("Username too short")
+     if not username.isalnum():
+         raise ValueError("Username must be alphanumeric")

      return True
```

**問題：**
- 「改善」了 email validation（超出 bug fix 範圍）
- 加了沒人要求的 username validation
- 改了 comments
- 加了 docstring

**✅ 應該做的事（只改那幾行）**

```diff
  def validate_user(user_data):
      # Check email format
-     if not user_data.get('email'):
+     email = user_data.get('email', '')
+     if not email or not email.strip():
          raise ValueError("Email required")

      # Basic email validation
-     if '@' not in user_data['email']:
+     if '@' not in email:
          raise ValueError("Invalid email")

      # Check username
      if not user_data.get('username'):
          raise ValueError("Username required")

      return True
```

**只改了：** 處理 empty email 的特定幾行。

### 範例 2：Style Drift

**使用者請求：**「在 upload function 加 logging」

**❌ LLM 常見錯誤（重新格式化全部）**

```diff
- def upload_file(file_path, destination):
+ def upload_file(file_path: str, destination: str) -> bool:
+     """Upload file to destination with logging."""
+     logger.info(f"Uploading {file_path} to {destination}")
+
      try:
-         with open(file_path, 'rb') as f:
-             data = f.read()
-
-         response = requests.post(destination, files={'file': data})
-
-         if response.status_code == 200:
-             return True
-         else:
-             return False
+         with open(file_path, "rb") as f:
+             data = f.read()
+         response = requests.post(destination, files={"file": data})
+         success = response.status_code == 200
+
+         if success:
+             logger.info(f"Upload successful: {file_path}")
+         else:
+             logger.error(f"Upload failed: {file_path}, status={response.status_code}")
+
+         return success
      except Exception as e:
-         print(f"Error: {e}")
+         logger.exception(f"Upload error for {file_path}: {e}")
          return False
```

**問題：**
- 改了 quote style（`''` → `""`）
- 加了沒人要求的 type hints
- 加了 docstring
- 重新格式化 whitespace
- 改了 boolean return logic

**✅ 應該做的事（匹配既有風格）**

```diff
+ import logging
+
+ logger = logging.getLogger(__name__)

  def upload_file(file_path, destination):
+     logger.info(f'Starting upload: {file_path}')
      try:
          with open(file_path, 'rb') as f:
              data = f.read()

          response = requests.post(destination, files={'file': data})

          if response.status_code == 200:
+             logger.info(f'Upload successful: {file_path}')
              return True
          else:
+             logger.error(f'Upload failed: {file_path}, status={response.status_code}')
              return False
      except Exception as e:
-         print(f"Error: {e}")
+         logger.exception(f'Upload error: {file_path}')
          return False
```

**匹配了：** 單引號、沒有 type hints、既有的 boolean pattern、間距風格。

---

## 4. Goal-Driven Execution（目標驅動執行）

### 範例 1：模糊 vs 可驗證

**使用者請求：**「修 authentication 系統」

**❌ LLM 常見錯誤（模糊做法）**

```
我來修 authentication 系統：
1. 審查 code
2. 找出問題
3. 做改善
4. 測試變更

[然後直接開始改 code，沒有明確的成功標準]
```

**✅ 應該做的事（可驗證的目標）**

```
要「修 authentication」，我需要定義 success criteria。
我們在解決什麼具體問題？

如果問題是「使用者改密碼後還能用舊 session 登入」：

Plan:
1. 寫 test：改密碼 → 驗證舊 session 被 invalidated
   驗證：Test fails（reproduces bug）

2. 實作：改密碼時 invalidate sessions
   驗證：Test passes

3. 檢查 edge cases：多個 active sessions、concurrent 改密碼
   驗證：Additional tests pass

4. 驗證沒有 regression：既有 auth tests 仍然通過
   驗證：Full test suite green

目前的 auth test coverage：[顯示 coverage]
你看到的具體 auth 問題是什麼？
```

### 範例 2：Incremental with Verification

**使用者請求：**「在 API 加 rate limiting」

**❌ LLM 常見錯誤（一次做完）**

```python
# 一次 commit 300 行實作完整的 rate limiting：
# Redis、多種策略、config system、monitoring
# 沒有清楚的驗證步驟
```

**✅ 應該做的事（逐步 + 驗證）**

```
Rate limiting 計畫：

1. 加基本的 in-memory rate limiting（單一 endpoint）
   驗證：
   - Test: 100 requests → 前 10 個成功，其他拿 429
   - Manual: curl endpoint 11 次，看到 rate limit error

2. 提取成 middleware（套用到所有 endpoint）
   驗證：
   - Test: Rate limits 套用到 /users 和 /posts
   - 既有 endpoint tests 仍然通過

3. 加 Redis backend（多 server 用）
   驗證：
   - Test: Rate limit 在 app restart 後仍然存在
   - Test: 兩個 app instance 共享 rate limit counter

4. 加 configuration（每個 endpoint 不同 rate）
   驗證：
   - Test: /search 允許 10/min，/users 允許 100/min
   - Config file 解析正確

每個步驟都可以獨立驗證和 deploy。
從 step 1 開始？
```

### 範例 3：Test-First Verification

**使用者請求：**「sorting 在有重複分數時會壞掉」

**❌ LLM 常見錯誤（直接修，沒先重現）**

```python
# 直接改了 sort logic，沒先確認 bug
def sort_scores(scores):
    return sorted(scores, key=lambda x: (-x['score'], x['name']))
```

**✅ 應該做的事（先重現再修）**

```python
# 1. 先寫一個 reproducing test
def test_sort_with_duplicate_scores():
    """Test sorting when multiple items have same score."""
    scores = [
        {'name': 'Alice', 'score': 100},
        {'name': 'Bob', 'score': 100},
        {'name': 'Charlie', 'score': 90},
    ]

    result = sort_scores(scores)

    # The bug: order is non-deterministic for duplicates
    # Run this test multiple times, it should be consistent
    assert result[0]['score'] == 100
    assert result[1]['score'] == 100
    assert result[2]['score'] == 90

# Verify: Run test 10 times → fails with inconsistent ordering

# 2. 現在用 stable sort 修
def sort_scores(scores):
    """Sort by score descending, then name ascending for ties."""
    return sorted(scores, key=lambda x: (-x['score'], x['name']))

# Verify: Test passes consistently
```

---

## Anti-Patterns Summary

| 原則 | Anti-Pattern | 修正方式 |
|------|-------------|---------|
| Think Before Coding | 默默假設檔案格式、欄位、範圍 | 列出假設，先問清楚 |
| Simplicity First | 為一個 discount function 寫 Strategy pattern | 一個 function 直到真的需要複雜度 |
| Surgical Changes | 修 bug 時順便改 quote style、加 type hints | 只改解決問題的那幾行 |
| Goal-Driven | 「我來審查 code 並改善」 | 「先寫 reproducing test → 讓它過 → 驗證沒 regression」 |

## Key Insight

「過度複雜」的範例看起來**不是明顯錯誤**——它們遵循 design patterns 和 best practices。問題在於 **timing**：在需要之前就加了複雜度，導致：

- Code 更難理解
- 引入更多 bug
- 花更長時間實作
- 更難測試

「簡單」的版本：
- 更容易理解
- 更快實作
- 更容易測試
- 可以在之後真正需要複雜度時再 refactor

**好的 code 是解決今天問題的簡單 code，不是解決明天問題的 premature code。**
