# Java TSV `trim()` Before `split("\t", -1)` Destroys Trailing Empty Fields

## Problem

When parsing tab-separated value (TSV) files in Java, calling `String.trim()` on a line before `split("\t", -1)` silently drops records with trailing empty fields.

## Root Cause

`String.trim()` removes all characters ≤ U+0020, which **includes the tab character (0x09)**. When a TSV line has a trailing empty field (e.g., `id\tfalse\t`), `trim()` converts it to `id\tfalse` (removing the trailing tab). Then `split("\t", -1)` produces `["id", "false"]` with length 2 instead of the expected `["id", "false", ""]` with length 3.

## Example

```java
// Write: produces "id\tfalse\t\n"
line = record.id + "\t" + record.blocked + "\t" + "";

// Read: BROKEN — trim() removes trailing tab
String trimmed = line.trim();                    // "id\tfalse"
String[] parts = trimmed.split("\t", -1);        // ["id", "false"] — length 2!
if (parts.length >= 3) { /* never reached */ }   // Record silently dropped!

// Read: CORRECT — no trim()
String[] parts = line.split("\t", -1);            // ["id", "false", ""] — length 3!
if (parts.length >= 3) { /* works */ }
```

## Fix

Replace:
```java
String trimmed = line.trim();
if (trimmed.isEmpty() || trimmed.startsWith("#")) continue;
String[] parts = trimmed.split("\t", -1);
```

With:
```java
if (line.isEmpty() || line.startsWith("#")) continue;
String[] parts = line.split("\t", -1);
```

## Applies When

- Parsing TSV files with `split("\t", -1)` where trailing empty fields are valid
- Using `trim()` before splitting by tab

## Does Not Apply When

- Using a proper CSV/TSV parser library (e.g., Apache Commons CSV, OpenCSV)
- Trailing empty fields are not meaningful in the data format
- Fields are separated by non-whitespace delimiters (e.g., `,` which is not removed by `trim()`)
