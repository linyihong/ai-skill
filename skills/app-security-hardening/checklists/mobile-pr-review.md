# Mobile PR Review Checklist

- No tokens, authorization headers, raw responses, or personal data are logged.
- Debug flags, test endpoints, and verbose interceptors are not enabled for release.
- API calls do not trust client-only state for authorization or pricing/business integrity.
- Storage changes document what data is cached, where, and for how long.
- Platform channels, deeplinks, intents, and native modules validate inputs and avoid broad capabilities.
- New dependencies or SDKs have permissions and telemetry reviewed.
- Tests or review evidence cover high-risk control changes.
