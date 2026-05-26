#!/usr/bin/env bash
# Stop hook: ensure last assistant response contains the per-turn
# `### Cognitive Mode 報告` block before Claude is allowed to stop.
#
# Per runtime/core-bootstrap.yaml §per_turn_obligations[obligation.cognitive.mode_report].
# Implements D-strategy (strengthen per-turn obligation visibility) from the
# bootstrap-contract-yaml-migration archived plan: UserPromptSubmit pre-injects
# reminder; Stop hook post-validates the response.
#
# Behavior:
#   - Block present in last assistant message  → exit 0 (allow stop)
#   - Block missing → exit 2 with stderr reason → Claude continues, appends, retries stop
#   - stop_hook_active=true → exit 0 (break infinite loop)
#   - transcript unreadable → exit 0 (fail open)
#
# Uses Python (universally available on macOS/Linux) instead of jq for
# portability — many systems do not have jq pre-installed.

set -uo pipefail

input_json="$(cat)"

result=$(python3 - "$input_json" <<'PYEOF'
import json
import os
import sys

raw = sys.argv[1] if len(sys.argv) > 1 else ""
try:
    payload = json.loads(raw)
except Exception:
    # Bad input — fail open.
    print("ALLOW")
    sys.exit(0)

if payload.get("stop_hook_active") is True:
    print("ALLOW_LOOP_GUARD")
    sys.exit(0)

transcript_path = payload.get("transcript_path", "")
if not transcript_path or not os.path.isfile(transcript_path):
    print("ALLOW_NO_TRANSCRIPT")
    sys.exit(0)

marker = "### Cognitive Mode 報告"
last_assistant_text = None
try:
    with open(transcript_path, "r", encoding="utf-8") as fh:
        for line in fh:
            line = line.strip()
            if not line:
                continue
            try:
                entry = json.loads(line)
            except Exception:
                continue
            role = entry.get("type") or entry.get("role")
            if role != "assistant":
                continue
            # Content can be: a string OR a list of content blocks.
            # Walk all string-bearing fields.
            chunks = []
            msg = entry.get("message")
            if isinstance(msg, dict):
                c = msg.get("content")
                if isinstance(c, str):
                    chunks.append(c)
                elif isinstance(c, list):
                    for item in c:
                        if isinstance(item, dict) and isinstance(item.get("text"), str):
                            chunks.append(item["text"])
                        elif isinstance(item, str):
                            chunks.append(item)
            else:
                c = entry.get("content")
                if isinstance(c, str):
                    chunks.append(c)
            if chunks:
                last_assistant_text = "\n".join(chunks)
except Exception:
    print("ALLOW_READ_ERROR")
    sys.exit(0)

if last_assistant_text is None:
    print("ALLOW_NO_ASSISTANT_MSG")
    sys.exit(0)

if marker in last_assistant_text:
    print("ALLOW_BLOCK_PRESENT")
    sys.exit(0)

print("BLOCK_MISSING")
sys.exit(2)
PYEOF
)
exit_code=$?

if [[ $exit_code -eq 0 ]]; then
  exit 0
fi

cat >&2 <<'REASON'
[ai-skill Stop hook] Missing obligation: your final response did not include the `### Cognitive Mode 報告` block.

Per runtime/core-bootstrap.yaml §per_turn_obligations[obligation.cognitive.mode_report], every final user-facing response MUST end with a Cognitive Mode block (compact 1-line for trivial all-default tasks: `Cognitive: <e>·<c>·<g>·<m> / V:<v> / Cost:<cost> / Sig:<signal>`; full 6-row markdown table otherwise).

Please append the block to your response now, then stop again. Canonical format spec: runtime/core-bootstrap.yaml. Query active obligations: `ai-skill runtime obligations`.
REASON
exit 2
