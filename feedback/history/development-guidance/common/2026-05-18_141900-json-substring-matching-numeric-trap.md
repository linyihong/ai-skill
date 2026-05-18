# JSON Substring Matching Trap in API Response Validation

## One-line Summary
When validating JSON API responses with substring matching (e.g., `json.contains("\"code\":1")`), numeric values like `1002` will falsely match because they contain the substring `"code":1`. Always parse the actual numeric value for comparison.

## Human Explanation
A common pattern in live tests is to check if an API response indicates success by looking for specific field values. Using `String.contains()` for this is dangerous because:

- `"code":1002` contains the substring `"code":1` (because `1002` starts with `1`)
- `"code":10` contains the substring `"code":1`
- Any code starting with the target digit will falsely match

This is especially problematic when the API returns `ret:200` (HTTP-level success) but `data.code` indicates business-logic failure (e.g., `1002` = "resource not found").

## Trigger
Writing a `Predicate<String>` to check if a PLAY API response succeeded, where the initial implementation used:
```java
json.contains("\"code\":0") || json.contains("\"code\":1")
```
This incorrectly matched `"code":1002` as success.

## Evidence
- `Skit.playSkit` returned `{"ret":200,"data":{"code":1002,"msg":"短劇不存在"}}`
- The predicate `json.contains("\"code\":1")` returned `true` because `"code":1002` contains the substring `"code":1`
- Fixed by using `extractJsonValue(json, "code")` to parse the actual numeric value

## Generalized Lesson
**Always parse JSON numeric fields to their actual numeric type before comparison.** Never use substring matching (`contains()`, `indexOf()`) to check numeric values in JSON, as any value starting with the target digits will produce false positives.

## Agent Action
When writing API response validation logic:
1. Extract the field value using a proper JSON parser or a dedicated extraction function that returns the exact value
2. Parse it to the appropriate numeric type (int, long, etc.)
3. Compare using numeric operators (`==`, `>`, `<`), not string operations
4. If a full JSON parser is unavailable, write an extraction function that reads until the next delimiter (comma, brace, bracket)

## Goal / Action / Validation
- **Goal**: Reliable API response validation without false positives from substring matching
- **Action**: Use numeric parsing instead of substring matching for JSON numeric fields
- **Validation**: Test with edge cases: `code:1002`, `code:10`, `code:1`, `code:0` — only the exact match should pass

## Applies When
- Writing live API tests that check response status codes
- Parsing JSON responses without a full JSON parser library
- Any scenario where numeric field values are checked via string operations

## Does Not Apply When
- Using a proper JSON parser (Jackson, Gson, org.json) that handles type coercion correctly
- Checking string fields (where substring matching is appropriate)
- The numeric values are guaranteed to be single digits

## Validation
Test the predicate against known edge cases:
- `"code":1002` → should NOT match `code == 1`
- `"code":10` → should NOT match `code == 1`
- `"code":1` → should match `code == 1`
- `"code":0` → should match `code == 0`

## Promotion Target
`feedback/feedback-lessons.md` (referenced as a validation pattern for API response checking)

## Required Linked Updates
None
