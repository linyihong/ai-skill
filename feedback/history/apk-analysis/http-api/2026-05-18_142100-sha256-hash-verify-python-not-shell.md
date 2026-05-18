# SHA256 Hash Verification: Use Python Instead of Shell to Avoid Quoting Issues

## One-line Summary
When computing SHA256 hashes to verify service names or API identifiers, use Python's `hashlib` instead of shell commands (`shasum`/`sha256sum`) to avoid subtle quoting/encoding issues in subprocess calls.

## Human Explanation
When reverse-engineering API calls, you often need to verify which service name produces a specific SHA256 hash prefix (e.g., matching a `serviceHash` from a Frida capture). Using shell commands like `echo -n "ServiceName" | shasum -a 256` seems straightforward, but when called from a subprocess (Java `Runtime.exec`, build scripts, etc.), shell quoting and encoding differences can produce incorrect results.

Python's `hashlib.sha256(b'ServiceName').hexdigest()[:16]` is more reliable because:
1. No shell interpretation layer — the bytes are passed directly
2. Consistent encoding — no risk of hidden characters or shell escaping
3. Reproducible — same result regardless of locale, shell, or environment

## Trigger
Trying to verify that `Skit.playSkit` produces the service hash `f103361467a8b211` from a Frida capture. The shell command `echo -n "Skit.playSkit" | shasum -a 256 | cut -c1-16` gave `578434e0deb97af2` (wrong), while Python `hashlib.sha256(b'Skit.playSkit').hexdigest()[:16]` gave `f103361467a8b211` (correct).

## Evidence
- Shell (via Java subprocess): `echo -n "Skit.playSkit" | shasum -a 256` → `578434e0deb97af2...`
- Python: `hashlib.sha256(b'Skit.playSkit').hexdigest()[:16]` → `f103361467a8b211`
- The discrepancy was caused by shell quoting issues in the subprocess call (the `echo -n` received different bytes than expected)
- The Python result matched the Frida-captured `serviceHash: f103361467a8b211`

## Generalized Lesson
**Always use Python (or the same language as your analysis toolchain) for hash computation, not shell commands.** Shell quoting, encoding, and subprocess piping introduce too many variables that can silently produce wrong results.

## Agent Action
When you need to compute a hash to verify an identifier:
1. Use Python: `python3 -c "import hashlib; print(hashlib.sha256(b'input').hexdigest()[:16])"`
2. Or use the same language as your test/analysis code (Java: `MessageDigest`, JavaScript: `crypto.createHash`)
3. Avoid shell pipelines for hash computation in automated scripts
4. If you must use shell, verify the result with Python first

## Goal / Action / Validation
- **Goal**: Correct SHA256 hash verification for API service name matching
- **Action**: Use Python `hashlib` instead of shell `shasum`/`sha256sum`
- **Validation**: Both methods should produce the same hash for simple ASCII strings; if they differ, Python is the ground truth

## Applies When
- Computing SHA256 hashes to match against Frida-captured `serviceHash` values
- Any hash computation in automated scripts or subprocess calls
- Debugging hash mismatches between expected and actual values

## Does Not Apply When
- Running the shell command directly in an interactive terminal (no subprocess layer)
- Using a programming language's built-in crypto library directly (not through shell)

## Validation
For any ASCII string, verify that:
```bash
python3 -c "import hashlib; print(hashlib.sha256(b'test').hexdigest()[:16])"
# Should match: 9f86d081884c7d65
echo -n "test" | shasum -a 256 | cut -c1-16
# Should also be: 9f86d081884c7d65
```
If they differ, the shell command has a quoting/encoding issue.

## Promotion Target
`feedback/feedback-lessons.md` (referenced as a toolchain reliability pattern)

## Required Linked Updates
None
