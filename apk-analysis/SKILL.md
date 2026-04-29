---
name: apk-analysis
description: Guides authorized APK traffic analysis, dynamic capture, Flutter/Dart AOT inspection, response decoding, evidence documentation, and reusable learning updates. Use when analyzing APK network behavior, reverse engineering API flows with permission, planning Frida/pcap/Proxyman/blutter workflows, or turning newly discovered APK analysis techniques into reusable docs.
---

# APK Analysis

Use this skill for authorized APK analysis only. The goal is to recover how an app communicates, captures, decodes, and validates data in a reproducible way, then turn findings into documentation and reusable tests without leaking secrets or target-specific private details.

## Quick Start

1. Confirm scope and authorization:
   - Identify the APK, version, device/emulator, architecture, and allowed analysis actions.
   - Do not collect or publish credentials, full tokens, private keys, personal data, or unrelated third-party traffic.
2. Separate method from target facts:
   - Put reusable techniques in this skill folder.
   - Put target-specific endpoints, hosts, schemas, and findings in the project’s API/reference docs.
3. Start with network path triage:
   - Check localhost traffic.
   - Check whole-device pcap for external TLS/SNI/host timing.
   - Check whether system proxy/MITM is actually used.
   - Check Java HTTP hooks only if evidence suggests Java networking.
   - For Flutter apps, inspect Dart AOT/native paths early.
4. Prefer high-semantic hooks:
   - Request object hooks are better than raw socket hooks.
   - Response interceptor/decode hooks are better than trying to reconstruct TLS bytes.
5. Convert dynamic results into durable assets:
   - Redacted request/response samples.
   - Offline decoders or fixtures.
   - API/schema docs.
   - Contract tests where the project has an SDK or client implementation.

## Default Workflow

Read [WORKFLOW.md](WORKFLOW.md) before doing hands-on analysis. Use it as the default decision tree.

Use [TOOLS.md](TOOLS.md) when preparing an environment or choosing between adb, pcap, Proxyman/mitmproxy, Frida, jadx, apktool, blutter, or offline decoding scripts.

Use [DOCUMENTATION.md](DOCUMENTATION.md) when writing human-readable results.

Use [FEEDBACK.md](FEEDBACK.md) when a new technique, failure pattern, or validation rule should be added back into this skill.

Use [RUNBOOK.md](RUNBOOK.md) when starting a new APK project or when the user asks how to apply this skill to another product.

## Required Output Style

When reporting analysis progress, include:

- What was tested.
- What evidence was observed.
- What was ruled out.
- What remains unknown.
- The next lowest-risk action.

When documenting a new finding, include:

- Trigger or UI path.
- Tool and command summary.
- Evidence file path or sanitized excerpt.
- Generalized lesson.
- Follow-up validation.

## Safety and Sanitization

Never write raw secrets into reusable skill docs:

- Full Authorization tokens.
- Session cookies.
- Device identifiers that belong to a real user/device.
- AES/HMAC secrets unless they are synthetic examples.
- Private hostnames not meant for the reusable method guide.
- Personal user data.

Use placeholders:

```text
<package-name>
<device-serial>
<api-host>
<proxy-host>:<proxy-port>
<token-redacted>
<secret-redacted>
```

## Feedback Loop

If analysis discovers a new reusable idea:

1. Add it to `FEEDBACK.md` as a dated candidate lesson.
2. Generalize it so it is not tied to one APK.
3. Add evidence and validation criteria.
4. Promote it into `WORKFLOW.md`, `TOOLS.md`, or `DOCUMENTATION.md` only after it has been validated or is clearly labeled as experimental.

Do not silently overwrite prior lessons. Append, compare, and explain why the new rule supersedes or narrows the old one.
