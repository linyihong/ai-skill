### 2026-05-14 - Bash Regex Bracket Inside `[[ =~ ]]` Syntax Error

#### One-line Summary
In bash `[[ ... ]]` conditionals, using `=~` with regex patterns containing literal `]` causes syntax errors because `]` is interpreted as closing the `[[` bracket; the fix is to store the regex in a variable.

#### Human Explanation
When writing bash scripts, it's common to use `[[ "${var}" =~ pattern ]]` for regex matching. However, if the regex pattern contains a literal `]` character (e.g., `\[.\]` to match `[x]`), bash's parser gets confused — it sees the `]` inside the pattern and thinks it's the closing `]]` of the conditional expression, causing a syntax error like `syntax error in conditional expression: unexpected token '\]'`.

The root cause is that bash's `[[` compound command has special parsing rules for the `]]` terminator, and `]` inside `=~` patterns is not protected from this parsing.

#### Trigger
Any bash script using `[[ "${var}" =~ ...\]... ]]` where the regex contains a literal `]` character.

#### Evidence
```bash
# ❌ Syntax error: unexpected token '\]'
if [[ "${line}" =~ ^[[:space:]]*\[.\] ]]; then
  ...
fi

# ✅ Works correctly: regex stored in variable
local re='^[[:space:]]*\[.\]'
if [[ "${line}" =~ $re ]]; then
  ...
fi
```

#### Generalized Lesson
When using `=~` inside `[[ ... ]]` in bash, **always store regex patterns containing `]` in a variable** before using them. This avoids the parser conflict between the regex's `]` and the `[[` terminator `]]`.

This applies to any `]` in the pattern, including:
- `\[.\]` — matching a single character inside brackets
- `\[[xX]\]` — matching `[x]` or `[X]`
- `\[ \]` — matching `[ ]`

#### Agent Action
When writing bash `[[ "${var}" =~ ... ]]` conditionals:
1. Check if the regex contains any literal `]` character
2. If yes, extract the regex into a local variable first
3. Use `[[ "${var}" =~ $variable_name ]]` instead of inline pattern

#### Goal / Action / Validation
- **Goal**: Avoid bash syntax errors when using regex with `]` inside `[[ =~ ]]`
- **Action**: Store regex in variable before use
- **Validation**: Run `bash -n script.sh` to check for syntax errors

#### Applies When
- Writing bash scripts with `[[ ... =~ ... ]]` conditionals
- Regex patterns need to match literal `[` or `]` characters
- Any pattern containing `]` that would be confused with `[[` terminator

#### Does Not Apply When
- Using `grep -E` or `sed` instead of `[[ =~ ]]`
- Regex patterns without `]` character
- Using `case` statements instead of `[[ ]]`

#### Validation
Run `bash -n script.sh` to verify no syntax errors. If the script passes, the fix is correct.

#### Promotion Target
`enforcement/failure-patterns/` — This is a general bash scripting gotcha that could be a failure pattern if it recurs.

#### Required Linked Updates
None.
