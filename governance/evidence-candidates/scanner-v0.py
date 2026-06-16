#!/usr/bin/env python3
"""Evidence Candidate Scanner v0 — ASSEMBLER (stateless).

Design: plans/active/2026-06-16-1131-evidence-candidate-system.md §Phase 1C.

This is NOT the ai-skill CLI and is NOT wired into runtime / hooks. It is a
standalone, stateless assembler: given an artifact + EXPLICIT criteria_hits
(annotated outside the scanner) + the rule registry, it assembles a well-formed
candidate and persists it to the gitignored inbox.

scanner DOES:   schema validate / pointer resolve / dedupe / invariant check / persist inbox
scanner DOES NOT: infer / match / classify / score / rank / accept / expire

Invariants enforced:
  - source.artifact MUST reference an original artifact, NOT another candidate.
  - criteria_hits MUST originate outside the scanner (criteria_source.actor present
    and not the scanner itself).
Output ordering is undefined (emit only). No confidence is produced (Q1 frozen).

Usage:  scanner-v0.py < input.json      (reads one annotated artifact as JSON)
Stateless: same input -> same candidate id (content hash) -> idempotent write.
"""
import sys, json, re, hashlib, pathlib

ROOT = pathlib.Path(__file__).resolve().parent          # governance/evidence-candidates/
REGISTRY = ROOT / "evidence-rules"
INBOX = ROOT / "inbox"
SCANNER_ACTORS = {"scanner", "scanner-v0", "self"}       # forbidden as criteria_source.actor


def fail(msg):
    print(f"REJECT: {msg}", file=sys.stderr)
    sys.exit(1)


def main():
    try:
        data = json.load(sys.stdin)
    except Exception as e:
        fail(f"input is not valid JSON: {e}")

    # --- schema validate ---
    src = data.get("source") or {}
    for f in ("repo", "artifact"):
        if not src.get(f):
            fail(f"source.{f} required")
    matched_plans = data.get("matched_plans") or []
    criteria_hits = data.get("criteria_hits") or []
    if not matched_plans:
        fail("matched_plans must be non-empty")
    if not criteria_hits:
        fail("criteria_hits must be non-empty")
    actor = (data.get("criteria_source") or {}).get("actor")
    if not actor:
        fail("criteria_source.actor required (criteria_hits MUST originate outside scanner)")

    # --- invariant: criteria_hits originate outside scanner ---
    if actor.lower() in SCANNER_ACTORS:
        fail(f"criteria_source.actor='{actor}' is the scanner itself; criteria_hits must come from outside (human / matcher-v2=Phase 2)")

    # --- invariant: source must be an original artifact, not another candidate ---
    art = str(src["artifact"])
    if re.match(r"^C-[0-9a-fA-F]{6,}$", art) or "evidence-candidates/inbox" in art:
        fail(f"source.artifact '{art}' looks like a candidate; candidate MUST NOT reference another candidate")

    # --- pointer resolve: every matched_plan must have a registry pointer ---
    for plan in matched_plans:
        ptr = REGISTRY / f"{plan}.pointer.yaml"
        if not ptr.exists():
            fail(f"no registry pointer for matched_plan '{plan}' ({ptr})")

    # --- deterministic id (dedupe / idempotent) ---
    basis = json.dumps({
        "source": {k: src.get(k) for k in ("repo", "artifact", "commit")},
        "matched_plans": sorted(matched_plans),
        "criteria_hits": sorted(criteria_hits),
    }, sort_keys=True, ensure_ascii=False)
    cid = "C-" + hashlib.sha1(basis.encode("utf-8")).hexdigest()[:8]

    candidate = {
        "id": cid,
        "source": {k: src.get(k) for k in ("repo", "artifact", "commit")},
        "matched_plans": matched_plans,          # order preserved as given; scanner does not rank
        "criteria_hits": criteria_hits,
        "criteria_source": {"actor": actor},
        "status": "create",                       # scanner never sets accept/discard/expire
    }

    INBOX.mkdir(parents=True, exist_ok=True)
    out = INBOX / f"{cid}.json"
    if out.exists():
        print(f"IDEMPOTENT: {cid} already in inbox (no duplicate written)")
    else:
        out.write_text(json.dumps(candidate, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
        print(f"EMIT: {cid} -> inbox/{cid}.json")
    # emit-only; ordering undefined; no ranking/scoring/confidence.


if __name__ == "__main__":
    main()
